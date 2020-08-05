package routers

import (
	"net/http"

	"github.com/vivshankar/verifypass-go/internal/dpcm"
	"github.com/vivshankar/verifypass-go/internal/oidc"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/vivshankar/verifypass-go/internal/config"
)

func Profile(c *gin.Context) {

	ok, token := checkAuthenticated(c)
	if !ok {
		return
	}

	// Call DUA to check if consent needed
	res, err := checkForDUA(c, token, []*dpcm.DUAItem{
		&dpcm.DUAItem{
			PurposeID:    config.ProfilePurposeID,
			AccessTypeID: config.ReadAccessType,
		},
	}, "/profile", "Consent for Profile Display")

	if err != nil {
		log.Errorf("DUA check could not complete; %v", err)
	}

	if !res {
		return
	}

	// Introspect the token
	payload, err := oidc.Introspect(token)
	if err != nil {
		log.Errorf("Unable to introspect the token; err = %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := &map[string]string{
		"displayName": payload["preferred_username"].(string),
		"givenName":   "",
		"familyName":  "",
		"mobile":      "",
		"email":       "",
	}

	if val, ok := payload["given_name"]; ok {
		(*user)["givenName"], _ = val.(string)
	}

	if val, ok := payload["family_name"]; ok {
		(*user)["familyName"], _ = val.(string)
	}

	if val, ok := payload["mobile_number"]; ok {
		(*user)["mobile"], _ = val.(string)
	}

	if val, ok := payload["email"]; ok {
		(*user)["email"], _ = val.(string)
	}

	c.HTML(http.StatusOK, "profile.tmpl", gin.H{
		"title":      "User Profile",
		"isLoggedIn": true,
		"user":       *user,
	})
}
