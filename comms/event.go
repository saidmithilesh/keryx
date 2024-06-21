// Package comms provides communication-related structures and functions.
package comms

import (
	"encoding/json"
	"fmt"
	"time"

	// Package for generating and working with UUIDs
	"github.com/google/uuid"
)

// Constants representing different delivery statuses.
const (
	StatusWait         = 0 // Waiting status
	StatusSent         = 1 // Sent status
	StatusDelivered    = 2 // Delivered status
	StatusAcknowledged = 3 // Acknowledged status

	EventTypeMessage  = "message"
	EventTypeJoinRoom = "join_room"
)

// DeliveryStatus represents the delivery status of a message.
type DeliveryStatus struct {
	ReceiverID uuid.UUID // ID of the message receiver
	Status     int       // Delivery status
}

// Event represents a communication event.
type Event struct {
	ID   uuid.UUID `json:"id"`   // Unique identifier for the event
	Type string    `json:"type"` // Type of event

	// ID of the sender of the event
	SenderID uuid.UUID `json:"senderId"`

	// Payload of the event
	Payload []byte `json:"payload"`

	// Time when the event was sent by the sender
	SenderTime time.Time `json:"senderTime"`
	// Time when the event was received by the server
	ServerTime time.Time `json:"serverTime"`
}

// populate sets the ServerTime of the event to the current time.
func (event *Event) populate() {
	event.ID = uuid.New()
	event.ServerTime = time.Now()
}

func (event *Event) validate() error {
	switch event.Type {
	case EventTypeMessage, EventTypeJoinRoom:
		return nil
	default:
		return fmt.Errorf("invalid event type %s", event.Type)
	}
}

// toBytes serializes the Event object to a JSON byte slice.
func (event *Event) toBytes() ([]byte, error) {
	return json.Marshal(event)
}

// EventFromBytes deserializes a JSON byte slice into an Event object
// and populates its ServerTime.
func EventFromBytes(packet []byte) (event *Event, err error) {
	// Deserialize the JSON byte slice into the Event object
	event = &Event{}
	err = json.Unmarshal(packet, event)
	if err != nil {
		// Return error if deserialization fails
		return nil, err
	}

	// Set the ServerTime to the current time
	event.populate()
	return event, nil
}
