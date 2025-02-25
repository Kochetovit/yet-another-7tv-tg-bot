package dto

type Update struct {
	UpdateID int           `json:"update_id"`
	Message  *InputMessage `json:"message,omitempty"`
}

type InputMessage struct {
	Chat *Chat  `json:"chat,omitempty"`
	Text string `json:"text"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type SetWebhook struct {
	URL string `json:"url"`
}
