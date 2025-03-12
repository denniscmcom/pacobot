package cmd

import (
	"bytes"
	"encoding/json"
	"io"
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

	log.Println(authRes.Scope)

	return authRes
}

func IsAuthTokenValid(authToken string) bool {
	endpoint := "https://id.twitch.tv/oauth2/validate"
	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		log.Fatalf("Error creating request: %v\n", err)
	}

	req.Header.Set("Authorization", "OAuth "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error sending request: %v\n", err)
	}

	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func RevokeAuthToken(authToken string) {
	config := readConfig()

	data := url.Values{}
	data.Set("client_id", config.ClientId)
	data.Set("token", authToken)

	endpoint := "https://id.twitch.tv/oauth2/revoke"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(data.Encode()))

	if err != nil {
		log.Fatalf("Error creating request: %v\n", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error sending request: %v\n", err)
	}

	defer resp.Body.Close()
}

func RefreshAuthToken(authToken, refreshToken string) AuthRes {
	config := readConfig()

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", config.ClientId)
	data.Set("client_secret", config.ClientSecret)

	endpoint := "https://id.twitch.tv/oauth2/token"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(data.Encode()))

	if err != nil {
		log.Fatalf("Error creating request: %v\n", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error sending request: %v\n", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var authResponse AuthRes

	if err := json.Unmarshal(body, &authResponse); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return authResponse
}

// TODO: Return broadcaste user id
func GetBroadcasterUserId(userName, authToken string) {
	config := readConfig()

	endpoint := "https://api.twitch.tv/helix/users?login=" + userName
	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		log.Fatalf("Error creating request: %v\n", err)
	}

	req.Header.Set("Client-ID", config.ClientId)
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error sending request: %v\n", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	log.Println(string(body))
}
