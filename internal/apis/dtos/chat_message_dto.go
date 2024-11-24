package dtos

type ChatMessageDto struct {
	ChatID         string `json:"chat_id"`
	ReceiverUserID string `json:"receiver_user_id"`
	MessageType    string `json:"message_type"`
	Message        string `json:"message"`
}
