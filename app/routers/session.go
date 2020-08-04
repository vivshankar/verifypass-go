package routers

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/vivshankar/verifypass-go/internal/config"
)

var oauth2Config *oauth2.Config = &oauth2.Config{
	ClientID:     config.ClientID,
	ClientSecret: config.ClientSecret,
	Scopes:       []string{"openid"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  fmt.Sprintf("%s/v1.0/endpoint/default/authorize", config.Tenant),
		TokenURL: fmt.Sprintf("%s/v1.0/endpoint/default/token", config.Tenant),
	},
	RedirectURL: config.RedirectURI,
}

const (
	tokenKey string = "session_token"
)

func Login(c *gin.Context) {

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := oauth2Config.AuthCodeURL("state")
	log.Infof("Visit the URL for the auth dialog: %v", url)

	c.Redirect(http.StatusTemporaryRedirect, url)
	c.Abort()
}

func Logout(c *gin.Context) {

	fmt.Println("In logout...")

	session := sessions.Default(c)
	session.Set("dummy", "val")
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1}) // this sets the cookie with a MaxAge of 0
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/")
	c.Abort()
}

func Callback(c *gin.Context) {

	oauth2Token, err := oauth2Config.Exchange(c, c.Request.URL.Query().Get("code"))
	if err != nil {
		log.Errorf("Error occurred at Exchange: err = %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
	}

	log.Infof("Access Token is = %s", oauth2Token.AccessToken)
	session := sessions.Default(c)
	session.Set(tokenKey, oauth2Token.AccessToken)
	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: int(oauth2Token.Expiry.Sub(time.Now()).Seconds()),
	})
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/")
	c.Abort()
}

func checkAuthenticated(c *gin.Context) (bool, string) {

	session := sessions.Default(c)
	token := session.Get(tokenKey)

	if token == nil {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":      "Welcome to VerifyPass",
			"isLoggedIn": false,
		})

		return false, ""
	}

	return true, token.(string)
}
