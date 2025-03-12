package main

import (
	"net/http"

	"github.com/denniscmcom/pacobot/cmd"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	var authRes cmd.AuthRes

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello world",
		})
	})

	// TODO: Pass username in parameters
	r.GET("/id", func(c *gin.Context) {
		cmd.GetBroadcasterUserId("denniscmartin", authRes.AccessToken)

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.GET("/auth", func(c *gin.Context) {
		authUrl := cmd.GetAuthUrl()

		c.Redirect(http.StatusMovedPermanently, authUrl)
	})

	r.GET("/auth-validate", func(c *gin.Context) {
		msg := "failed"

		if cmd.IsAuthTokenValid(authRes.AccessToken) {
			msg = "ok"
		}

		c.JSON(http.StatusOK, gin.H{
			"message": msg,
		})
	})

	r.GET("/auth-refresh", func(c *gin.Context) {
		authRes = cmd.RefreshAuthToken(authRes.AccessToken, authRes.RefreshToken)

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.GET("/auth-revoke", func(c *gin.Context) {
		cmd.RevokeAuthToken(authRes.AccessToken)

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.GET("/twitch-auth-code-callback", func(c *gin.Context) {
		authCode := c.Query("code")
		authRes = cmd.GetAuthToken(authCode)

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.GET("/connect", func(c *gin.Context) {
		go cmd.ConnSocket(authRes.AccessToken)

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.Run()
}
