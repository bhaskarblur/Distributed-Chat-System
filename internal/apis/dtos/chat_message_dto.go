package dtos

type ChatMessageDto struct {
	ChatId      string `json:"chat_id"`
	ReceiverId  string `json:"receiver_id"`
	MessageType string `json:"message_type"`
	Message     string `json:"message"`
}
