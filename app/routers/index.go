package routers

import (
	"net/http"

	"github.com/vivshankar/verifypass-go/internal/config"
	"github.com/vivshankar/verifypass-go/internal/dpcm"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Index(c *gin.Context) {

	ok, token := checkAuthenticated(c)
	if !ok {
		return
	}

	log.Infof("I am logged in! Token = %s", token)

	// Check for EULA
	eulaRequest := &dpcm.DUARequest{
		Items: []*dpcm.DUAItem{
			&dpcm.DUAItem{
				PurposeID:    config.EULAPurposeID,
				AccessTypeID: config.EULAAccessType,
			},
		},
	}

	eulaResp, err := dpcm.DUA(token, eulaRequest)
	if err != nil {
		log.Errorf("DUA call resulted in error: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.Infof("DUA Response: %v", eulaResp)
	for _, purposeItem := range eulaResp {

		for _, attrItem := range purposeItem.Result {

			if !attrItem.Approved {
				// Prompt EULA
				c.Redirect(http.StatusTemporaryRedirect, "/eula")
				return
			}
		}
	}

	// Prompt EULA
	c.Redirect(http.StatusTemporaryRedirect, "/main")
}
