package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

// TODO: Change unmarshall to JSON DECODE

func GetAuthUrl() string {
	config := readConfig()

	baseUrl := &url.URL{
		Scheme: "https",
		Host:   "id.twitch.tv",
		Path:   "/oauth2/authorize",
	}

	params := url.Values{}
	params.Add("client_id", config.ClientId)
	params.Add("force_verify", "true")
	params.Add("redirect_uri", "http://localhost:8080/twitch-auth-code-callback")
	params.Add("response_type", "code")
	params.Add("scope", "channel:bot user:read:chat")
	params.Add("state", "c3ab8aa609ea11e793ae92361f002671")

	baseUrl.RawQuery = params.Encode()

	return baseUrl.String()
}

type AuthRes struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
}

func GetAuthToken(authCode string) AuthRes {
	config := readConfig()

	baseUrl := &url.URL{
		Scheme: "https",
		Host:   "id.twitch.tv",
		Path:   "/oauth2/token",
	}

	formData := url.Values{}
	formData.Add("client_id", config.ClientId)
	formData.Add("client_secret", config.ClientSecret)
	formData.Add("code", authCode)
	formData.Add("grant_type", "authorization_code")
	formData.Add("redirect_uri", "http://localhost:8080/twitch-auth-code-callback")

	res, err := http.PostForm(baseUrl.String(), formData)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal("GetAuthToken")
	}

	var authRes AuthRes

	err = json.NewDecoder(res.Body).Decode(&authRes)

	if err != nil {
		log.Fatal(err)
	}

	return authRes
}

func IsAuthTokenValid(authToken string) bool {
	baseUrl := &url.URL{
		Scheme: "https",
		Host:   "id.twitch.tv",
		Path:   "oauth2/validate",
	}

	req, err := http.NewRequest("GET", baseUrl.String(), nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "OAuth "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func RevokeAuthToken(authToken string) {
	config := readConfig()

	baseUrl := &url.URL{
		Scheme: "https",
		Host:   "id.twitch.tv",
		Path:   "oauth2/revoke",
	}

	data := url.Values{}
	data.Add("client_id", config.ClientId)
	data.Add("token", authToken)

	res, err := http.PostForm(baseUrl.String(), data)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
}

func RefreshAuthToken(authToken, refreshToken string) AuthRes {
	config := readConfig()

	baseUrl := &url.URL{
		Scheme: "https",
		Host:   "id.twitch.tv",
		Path:   "oauth2/token",
	}

	data := url.Values{}
	data.Add("grant_type", "refresh_token")
	data.Add("refresh_token", refreshToken)
	data.Add("client_id", config.ClientId)
	data.Add("client_secret", config.ClientSecret)

	res, err := http.PostForm(baseUrl.String(), data)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var authRes AuthRes

	err = json.NewDecoder(res.Body).Decode(&authRes)

	if err != nil {
		log.Fatal(err)
	}

	return authRes
}

type UserRes struct {
	Data []struct {
		Id string `json:"id"`
	} `json:"data"`
}

func GetUser(userName, authToken string) UserRes {
	config := readConfig()

	baseUrl := &url.URL{
		Scheme: "https",
		Host:   "api.twitch.tv",
		Path:   "helix/users",
	}

	params := url.Values{}
	params.Add("login", userName)

	req, err := http.NewRequest("GET", baseUrl.String(), nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Client-ID", config.ClientId)
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var userRes UserRes

	err = json.NewDecoder(res.Body).Decode(&userRes)

	if err != nil {
		log.Fatal(err)
	}

	return userRes
}
