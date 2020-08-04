package routers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vivshankar/verifypass-go/internal/config"
	"github.com/vivshankar/verifypass-go/internal/dpcm"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ConsentView struct {
	ID             string
	PurposeID      string
	PurposeName    string
	AttributeID    string
	AttributeName  string
	AccessTypeID   string
	AccessTypeName string
	Expires        string
	ConsentID      string
	Consented      bool
	Version        int
}

func Consents(c *gin.Context) {

	ok, token := checkAuthenticated(c)
	if !ok {
		return
	}

	log.Infof("I am logged in! Token = %s", token)

	scope := c.DefaultQuery("scope", "")
	if scope == "" {
		c.HTML(http.StatusOK, "consents.tmpl", gin.H{
			"title":      c.DefaultQuery("title", "My consents"),
			"isLoggedIn": true,
			"consents":   nil,
			"callback":   c.DefaultQuery("callback", "/"),
		})

		return
	}

	purposes := make([]string, 0)

	if scope != "all" {
		for _, item := range strings.Split(scope, ",") {
			purposeID, _, _ := splitScope(item)
			purposes = append(purposes, purposeID)
		}
	} else {
		purposes = append(purposes, config.EULAPurposeID)
		purposes = append(purposes, config.ProfilePurposeID)
		purposes = append(purposes, config.MFAPurposeID)
	}

	dspReq := &dpcm.DSPRequest{
		PurposeID: purposes,
	}

	dspResp, err := dpcm.DSP(token, dspReq)
	if err != nil {
		log.Errorf("DSP call resulted in error: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	consents := preprocessDSP(dspResp, strings.Split(c.DefaultQuery("scope", ""), ","))
	c.HTML(http.StatusOK, "consents.tmpl", gin.H{
		"title":      c.DefaultQuery("title", "My consents"),
		"isLoggedIn": true,
		"consents":   consents,
		"callback":   c.DefaultQuery("callback", "/"),
	})
}

func RecordConsents(c *gin.Context) {

	ok, token := checkAuthenticated(c)
	if !ok {
		return
	}

	log.Infof("I am logged in! Token = %s", token)

	var body map[string]bool
	err := c.ShouldBindJSON(&body)

	if err != nil {
		b, _ := c.GetRawData()
		log.Errorf("Error occurred in binding JSON; err = %v, body = %v", err, string(b))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var consents []*dpcm.Consent = make([]*dpcm.Consent, 0)
	for k, v := range body {
		purposeID, attributeID, accessTypeID := splitScope(k)
		state := dpcm.ConsentStateAllow
		if !v {
			state = dpcm.ConsentStateDeny
		}

		log.Infof("[DEBUG] v = %v, k = %v, state = %v", v, k, state)
		consents = append(consents, &dpcm.Consent{
			PurposeID:    purposeID,
			AttributeID:  attributeID,
			AccessTypeID: accessTypeID,
			State:        state,
		})
	}

	err = dpcm.BulkCreateConsents(token, consents)
	if err != nil {
		log.Errorf("Error occurred in bulk creation of consents; err = %v, body = %v", err, body)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

// Scope uses the structure {{purpose}}:{{attribute}}/{{accessTypeId}}
func preprocessDSP(r *dpcm.DSPResponse, scope []string) []*ConsentView {

	var consentViewMap map[string]*ConsentView = make(map[string]*ConsentView)

	if scope[0] != "all" {
		// Extract the fragments
		for _, item := range scope {

			purpose, attributeID, accessTypeID := splitScope(item)
			view := &ConsentView{
				ID:             item,
				PurposeID:      purpose,
				AttributeID:    attributeID,
				AccessTypeID:   accessTypeID,
				PurposeName:    r.Purposes[purpose].Name,
				AccessTypeName: r.AccessTypes[accessTypeID].Name,
			}

			if attributeID != "" {
				view.AttributeName = r.Attributes[attributeID].Name
			}

			consentViewMap[item] = view
		}
	} else {
		// Just iterate through all purposes, attributes and access types
		for purposeID, purpose := range r.Purposes {

			if len(purpose.Attributes) == 0 {
				attributeID := ""
				for _, accessType := range purpose.AccessTypes {

					key := fmt.Sprintf("%s:%s/%s", purposeID, attributeID, accessType.AccessTypeID)
					view := &ConsentView{
						ID:             key,
						PurposeID:      purposeID,
						AttributeID:    attributeID,
						AccessTypeID:   accessType.AccessTypeID,
						PurposeName:    r.Purposes[purposeID].Name,
						AccessTypeName: r.AccessTypes[accessType.AccessTypeID].Name,
					}

					consentViewMap[key] = view
				}
			} else {

				for _, attribute := range purpose.Attributes {

					attributeID := attribute.AttributeID
					attributeName := r.Attributes[attributeID].Name
					for _, accessType := range attribute.AccessTypes {
						key := fmt.Sprintf("%s:%s/%s", purposeID, attributeID, accessType.AccessTypeID)
						view := &ConsentView{
							ID:             key,
							PurposeID:      purposeID,
							AttributeID:    attributeID,
							AccessTypeID:   accessType.AccessTypeID,
							PurposeName:    r.Purposes[purposeID].Name,
							AccessTypeName: r.AccessTypes[accessType.AccessTypeID].Name,
							AttributeName:  attributeName,
						}

						consentViewMap[key] = view
					}
				}
			}
		}
	}

	// Iterate the consents
	for k, v := range r.Consents {
		// Build the consent view key
		key := fmt.Sprintf("%s:%s/%s", v.PurposeID, v.AttributeID, v.AccessTypeID)
		view, ok := consentViewMap[key]
		if !ok {
			// Consent may no longer be relevant because the purpose has been updated and
			// the attribute, for example, might have been removed.
			continue
		}

		// Check if the consent has expired
		now := time.Now().Unix()
		if now <= v.EndTime && (v.State == dpcm.ConsentStateAllow || v.State == dpcm.ConsentStateOptIn) {
			view.Consented = true
			view.Expires = time.Unix(v.EndTime, 0).Format("Jan 2 15:04:05 2006")
			view.ConsentID = k
			view.Version = v.Version
		}
	}

	var consentViews []*ConsentView
	for _, v := range consentViewMap {
		consentViews = append(consentViews, v)
	}

	return consentViews
}

func splitScope(item string) (string, string, string) {
	toks := strings.Split(item, ":")
	purpose := toks[0]
	toks = strings.Split(toks[1], "/")
	attributeID := toks[0]
	accessTypeID := toks[1]

	return purpose, attributeID, accessTypeID
}

func checkForDUA(c *gin.Context, token string, items []*dpcm.DUAItem, callback string) (bool, error) {

	req := &dpcm.DUARequest{
		Items: items,
	}

	resp, err := dpcm.DUA(token, req)
	if err != nil {
		log.Errorf("DUA call resulted in error: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return false, err
	}

	log.Infof("[DEBUG] DUA Response: %v", resp)
	var scopes []string = make([]string, 0)
	for _, purposeItem := range resp {

		for _, attrItem := range purposeItem.Result {

			if !attrItem.Approved {
				attributeID := attrItem.AttributeID
				if attributeID == "" {
					attributeID = purposeItem.AttributeID
				}

				scopes = append(scopes, fmt.Sprintf("%s:%s/%s", purposeItem.PurposeID, attributeID, purposeItem.AccessTypeID))
			}
		}
	}

	if len(scopes) > 0 {
		qs := strings.Join(scopes, ",")
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/consents?callback=%s&scope=%s", callback, url.QueryEscape(qs)))
		return false, nil
	}

	return true, nil
}
