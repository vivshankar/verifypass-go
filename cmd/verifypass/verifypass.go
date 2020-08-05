package main

import (
	"github.com/vivshankar/verifypass-go/app"
	"github.com/vivshankar/verifypass-go/internal/config"

	// "goginapp/plugins" if you create your own plugins import them here
	"os"

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

	router := gin.Default()

	app.Register(router)

	log.Info("Starting verifypass on port " + port())
	router.Run(port())
}
