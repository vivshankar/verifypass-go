package main

import (
	"html/template"
	"strconv"
	"time"

	"github.com/vivshankar/verifypass-go/app/routers"
	"github.com/vivshankar/verifypass-go/internal/config"

	// "goginapp/plugins" if you create your own plugins import them here
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func port() string {
	port := config.Port
	return ":" + port
}

func main() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	funcMap := template.FuncMap{
		"formatTime": func(raw string) string {
			if raw == "" {
				return "-"
			}

			val, _ := strconv.ParseInt(raw, 10, 64)
			t := time.Unix(val, 0)

			return t.Format("Jan 2 15:04:05 2006")
		},
	}

	router := gin.Default()
	router.RedirectTrailingSlash = false
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// Static webpage content endpoints
	router.SetFuncMap(funcMap)
	router.LoadHTMLGlob("web/template/*")
	router.Use(static.Serve("/static", static.LocalFile("./web/static", false)))
	router.Use(static.Serve("/img", static.LocalFile("./web/img", false)))
	router.GET("/", routers.Index)
	router.GET("/main", routers.Main)
	router.GET("/login", routers.Login)
	router.GET("/session/logout", routers.Logout)
	router.GET("/oauth/callback", routers.Callback)
	router.GET("/eula", routers.PromptEULA)
	router.GET("/consents", routers.Consents)
	router.GET("/profile", routers.Profile)
	router.GET("/mfa", routers.MFA)

	// API endpoints
	router.POST("/api/eula/consent", routers.RecordEULAConsent)
	router.POST("/api/consents", routers.RecordConsents)

	log.Info("Starting verifypass on port " + port())

	router.Run(port())
}
