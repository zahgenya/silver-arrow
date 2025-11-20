package telegram

type TgUser struct {
	Ok     bool `json:"ok"`
	Result struct {
		ID                      int64  `json:"id"`
		IsBot                   bool   `json:"is_bot"`
		FirstName               string `json:"first_name"`
		Username                string `json:"username"`
		CanJoinGroups           bool   `json:"can_join_groups"`
		CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
		SupportsInlineQueries   bool   `json:"supports_inline_queries"`
		CanConnectToBusiness    bool   `json:"can_connect_to_business"`
		HasMainWebApp           bool   `json:"has_main_web_app"`
	} `json:"result"`
}

type TgSendMessage struct {
	ChatId		int64	`json:"chat_id"`
	Text		string	`json:"text"`
}

type LatestPriceResponse struct {
    BaseVolume  float64 `json:"base-volume"`
    LatestPrice float64 `json:"latest-price"`
    Symbol      string  `json:"symbol"`
}

type TgGetUpdate struct {
	Ok		bool		`json:"ok"`
	Result	[]TgUpdate	`json:"result"`
}

type TgUpdate struct {
	UpdateId	int			`json:"update_id"`
	Message		*TgMessage	`json:"message"`
}

type TgMessage struct {
	MessageId	int64	`json:"message_id"`
	Chat		TgChat	`json:"chat"`
	Text		string	`json:"text"`
}

type TgChat struct {
	Id		int64	`json:"id"`
	Type	string	`json:"type"`
}