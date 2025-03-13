package socket

type Res_Metadata struct {
	Metadata Metadata `json:"metadata"`
}

type Res_Metadata_Notif struct {
	Metadata Metadata_Notif `json:"metadata"`
}

type Res_Welcome struct {
	Metadata Metadata        `json:"metadata"`
	Payload  Payload_Welcome `json:"payload"`
}

type Res_Keepalive struct {
	Metadata Metadata          `json:"metadata"`
	Payload  Payload_Keepalive `json:"payload"`
}

type Res_Notif_ChannelChatMsg struct {
	Metadata Metadata_Notif               `json:"metadata,omitempty"`
	Payload  Payload_Notif_ChannelChatMsg `json:"payload,omitempty"`
}

type Metadata struct {
	MsgId        string `json:"message_id"`
	MsgType      string `json:"message_type"`
	MsgTimestamp string `json:"message_timestamp"`
}

type Metadata_Notif struct {
	Metadata
	SubType    string `json:"subscription_type"`
	SubVersion string `json:"subscription_version"`
}

type Payload_Welcome struct {
	Session struct {
		Id               string `json:"id"`
		Status           string `json:"status"`
		ConnectedAt      string `json:"connected_at"`
		KeepaliveTimeout int    `json:"keepalive_timeout_seconds"`
	} `json:"session"`
}

type Payload_Keepalive struct {
}

type Payload_Notif_ChannelChatMsg struct {
	Subscription Payload_Sub_Notif_ChannelChatMsg   `json:"subscription"`
	Event        Payload_Event_Notif_ChannelChatMsg `json:"event"`
}

type Payload_Sub_Notif struct {
	Id      string `json:"id"`
	Status  string `json:"status"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Cost    int    `json:"cost"`

	Transport struct {
		Method    string `json:"method"`
		SessionId string `json:"session_id"`
	} `json:"transport"`

	CreatedAt string `json:"created_at"`
}

type Payload_Sub_Notif_ChannelChatMsg struct {
	Payload_Sub_Notif

	Condition struct {
		BroadcasterUserId string `json:"broadcaster_user_id"`
		UserId            string `json:"user_id"`
	} `json:"condition"`
}

type Payload_Event_Notif_ChannelChatMsg struct {
	BroadcasterUserId    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	ChatterUserId        string `json:"chatter_user_id"`
	ChatterUserLogin     string `json:"chatter_user_login"`
	ChatterUserName      string `json:"chatter_user_name"`
	MsgId                string `json:"message_id"`

	Msg struct {
		Text string `json:"text"`

		Fragments []struct {
			Type      string `json:"type"`
			Text      string `json:"text"`
			Cheermote string `json:"cheermote"`
			Emote     string `json:"emote"`
			Mention   string `json:"mention"`
		} `json:"fragments"`
	} `json:"message"`

	Color string `json:"color"`

	Badges []struct {
		SetId string `json:"set_id"`
		Id    string `json:"id"`
		Info  string `json:"info"`
	} `json:"badges"`

	MsgType                     string `json:"message_type"`
	Cheer                       string `json:"cheer"`
	Reply                       string `json:"reply"`
	ChannelPointsCustomRewardId string `json:"channel_points_custom_reward_id"`
	SourceBroadcasterUserId     string `json:"source_broadcaster_user_id"`
	SourceBroadcasterUserLogin  string `json:"source_broadcaster_user_login"`
	SourceBroadcasterUserName   string `json:"source_broadcaster_user_name"`
	SourceMessageId             string `json:"source_message_id"`
	SourceBadges                string `json:"source_badges"`
}
