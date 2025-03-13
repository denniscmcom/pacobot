package auth

type Config struct {
	ClientId          string `json:"client_id"`
	ClientSecret      string `json:"client_secret"`
	BroadcasterUserId string `json:"broadcaster_user_id"`
}

type UserRes struct {
	Data []struct {
		Id string `json:"id"`
	} `json:"data"`
}

type AuthRes struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
}
