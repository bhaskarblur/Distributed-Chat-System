package models

type ChatMessage struct {
	EventID        string `json:"event_id"`
	ChatID         string `json:"chat_id"`
	SenderUserID   string `json:"sender_user_id"`
	ReceiverUserID string `json:"receiver_user_id"`
	MessageType    string `json:"message_type"`
	Message        string `json:"message"`
}
