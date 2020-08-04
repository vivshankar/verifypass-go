package routers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/vivshankar/verifypass-go/internal/dpcm"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/vivshankar/verifypass-go/internal/config"
)

func PromptEULA(c *gin.Context) {

	ok, token := checkAuthenticated(c)
	if !ok {
		return
	}

	// Get the presentation cues
	dspRequest := &dpcm.DSPRequest{
		PurposeID: []string{config.EULAPurposeID},
	}
	dsp, err := dpcm.DSP(token, dspRequest)
	if err != nil {
		log.Errorf("Unable to get the DSP data to present to the user; err = %v", err)
		// TODO: Add error template
		return
	}

	eulaPurpose := dsp.Purposes[config.EULAPurposeID]

	c.HTML(http.StatusOK, "eula.tmpl", gin.H{
		"isLoggedIn": true,
		"purposeId":  config.EULAPurposeID,
		"eulaLink":   eulaPurpose.Terms.Ref,
		"eulaName":   eulaPurpose.Name,
		"title":      "Terms of service",
	})
}

func RecordEULAConsent(c *gin.Context) {

	ok, token := checkAuthenticated(c)
	if !ok {
		return
	}

	var body map[string]interface{}
	err := c.ShouldBindJSON(&body)

	if err != nil {
		b, _ := c.GetRawData()
		log.Errorf("Error occurred in binding JSON; err = %v, body = %v", err, string(b))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	consent := &dpcm.Consent{
		PurposeID:    body["purposeId"].(string),
		AccessTypeID: config.EULAAccessType,
		State:        dpcm.ConsentStateAllow,
		EndTime:      strconv.FormatInt(time.Now().AddDate(0, 0, 1).Unix(), 10),
	}

	err = dpcm.CreateConsent(token, consent)
	if err != nil {
		log.Errorf("Error occurred while creating consent; err = %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
