package smsgate

type MessageState uint8

const (
	MessageStateScheduled     MessageState = iota // The message is scheduled
	MessageStateEnroute                           // The message is in enroute state
	MessageStateDelivered                         // Message is delivered to destination
	MessageStateExpired                           // Message validity period has expired
	MessageStateDeleted                           // Message has been deleted
	MessageStateUndeliverable                     // Message is undeliverable
)

type SubmitRequest struct {
	Source         string
	Destination    string
	Content        string
	DeliveryReport bool
}

type SubmitResponse struct {
	MessageID string
}

type QueryResponse struct {
	MessageID string
	State     MessageState
}
