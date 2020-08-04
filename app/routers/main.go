package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Main(c *gin.Context) {

	ok, _ := checkAuthenticated(c)
	if !ok {
		return
	}

	c.HTML(http.StatusOK, "main.tmpl", gin.H{
		"title":      "Home",
		"isLoggedIn": true,
	})
}
