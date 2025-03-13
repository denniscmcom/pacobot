package main

import (
	"log"
	"net/http"

	"github.com/denniscmcom/pacobot/auth"
	"github.com/denniscmcom/pacobot/bot"
	"github.com/denniscmcom/pacobot/socket"
	"github.com/gin-gonic/gin"
)

type PageData struct {
	Title string
}

func main() {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()
	r.LoadHTMLGlob("./www/*.html")

	var authRes auth.AuthRes

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title": "Index",
		})
	})

	r.GET("/user", func(c *gin.Context) {
		userName := c.Query("username")
		user := auth.GetUser(userName, authRes.AccessToken)

		c.JSON(http.StatusOK, gin.H{
			"message": user.Data[len(user.Data)-1].Id,
		})
	})

	r.GET("/auth", func(c *gin.Context) {
		authUrl := auth.GetAuthUrl()

		c.Redirect(http.StatusMovedPermanently, authUrl)
	})

	r.GET("/auth-validate", func(c *gin.Context) {
		msg := "failed"

		if auth.IsAuthTokenValid(authRes.AccessToken) {
			msg = "ok"
		}

		c.JSON(http.StatusOK, gin.H{
			"message": msg,
		})
	})

	r.GET("/auth-refresh", func(c *gin.Context) {
		authRes = auth.RefreshAuthToken(authRes.AccessToken, authRes.RefreshToken)

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.GET("/auth-revoke", func(c *gin.Context) {
		auth.RevokeAuthToken(authRes.AccessToken)

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.GET("/twitch-auth-code-callback", func(c *gin.Context) {
		authCode := c.Query("code")
		authRes = auth.GetAuthToken(authCode)

		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.GET("/connect", func(c *gin.Context) {
		go socket.Connect(authRes.AccessToken)

		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.POST("/timer", func(c *gin.Context) {
		form := c.Request.PostForm
		log.Println(form)
		timesec := form.Get("tiempo-oculto")
		log.Println(timesec)
		args := []string{"timer", timesec}
		bot.HandleCmd(args)
	})

	r.Run()
}
