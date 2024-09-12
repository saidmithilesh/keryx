package event

type Event struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	SenderID  string `json:"sender_id"`
	Recipient string `json:"recipient"`

	Data interface{} `json:"data"`

	SenderTimestamp int64 `json:"sender_timestamp"`
	ServerTimestamp int64 `json:"server_timestamp"`

	DeliveryStatus int8 `json:"delivery_status"`
}
