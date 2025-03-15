package api

import (
	"net/http"

	"github.com/denniscmcom/pacobot/auth"
	"github.com/denniscmcom/pacobot/socket"
	"github.com/gin-gonic/gin"
)

func GetUserHandler(c *gin.Context) {
	userName := c.Query("username")
	user := auth.GetUser(userName, getAccessToken())

	c.JSON(http.StatusOK, gin.H{
		"message": user.Data[len(user.Data)-1].Id,
	})
}

func AuthHandler(c *gin.Context) {
	authUrl := auth.GetAuthUrl()
	c.Redirect(http.StatusMovedPermanently, authUrl)
}

func AuthValidateHandler(c *gin.Context) {
	msg := "failed"

	if auth.IsAuthTokenValid(getAccessToken()) {
		msg = "ok"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": msg,
	})
}

func AuthRefreshHandler(c *gin.Context) {
	authRes := auth.RefreshAuthToken(getAccessToken(), getRefreshToken())
	setAccessToken(authRes.AccessToken)
	setRefreshToken(authRes.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func AuthRevokeHandler(c *gin.Context) {
	auth.RevokeAuthToken(getAccessToken())

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func TwitchCallbackHandler(c *gin.Context) {
	authCode := c.Query("code")
	authRes := auth.GetAuthToken(authCode)
	authStore.Store("accessToken", authRes.AccessToken)
	authStore.Store("refreshToken", authRes.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func ConnectHandler(c *gin.Context) {
	go socket.Connect(getAccessToken())

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
