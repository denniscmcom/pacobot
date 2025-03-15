package api

import (
	"log"
	"sync"
)

var authStore sync.Map

func setAccessToken(accessToken string) {
	authStore.Store("accessToken", accessToken)
}

func setRefreshToken(refreshToken string) {
	authStore.Store("refreshToken", refreshToken)
}

func getAccessToken() string {
	value, exists := authStore.Load("accessToken")

	if !exists {
		log.Fatal("api: access token not found")
	}

	return value.(string)
}

func getRefreshToken() string {
	value, exists := authStore.Load("refreshToken")

	if !exists {
		log.Fatal("api: refresh token not found")
	}

	return value.(string)
}
