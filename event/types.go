package event

type ChannelChatMsgSubPayload struct {
	Payload struct {
		Event struct {
			ChatterUserId    string `json:"chatter_user_id"`
			ChatterUserLogin string `json:"chatter_user_login"`
			ChatterUserName  string `json:"chatter_user_name"`
			Msg              struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"message"`
		} `json:"event"`
	} `json:"payload"`
}
