package main

import (
	"github.com/denniscmcom/pacobot/api"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	r.GET("/user", api.GetUserHandler)
	r.GET("/auth", api.AuthHandler)
	r.GET("/auth-validate", api.AuthValidateHandler)
	r.GET("/auth-refresh", api.AuthRefreshHandler)
	r.GET("/auth-revoke", api.AuthRevokeHandler)
	r.GET("/twitch-auth-code-callback", api.TwitchCallbackHandler)
	r.GET("/connect", api.ConnectHandler)

	r.Run()
}
