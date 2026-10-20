package push

import (
	"context"
	"time"
)

// Outbound is one notification for one device, in a form that does not depend on the provider.
type Outbound struct {
	Token string
	// Platform is "ios", "android" or "unknown" (a token registered by an app version that did not say).
	Platform string
	Title    string
	Body     string
	Data     map[string]any
	// ImageURL is an optional https image; empty for a text notification.
	ImageURL string
	// ActionCategory is the interactive notification category (action buttons); empty for none.
	ActionCategory string
	// TTL is how long the provider may keep the message while the device is unreachable.
	TTL time.Duration
	// DedupeTag makes a duplicated send replace the notification already shown instead of stacking a second one,
	// where the platform allows it.
	DedupeTag string
}

// ResultKind is what the provider said about one message.
type ResultKind int

const (
	// ResultAccepted: the provider took the message. This does not say that the device got or showed it.
	ResultAccepted ResultKind = iota
	// ResultRetry: a failure that may go away, try again later.
	ResultRetry
	// ResultInvalidDevice: the provider confirms that the token is dead for good.
	ResultInvalidDevice
	// ResultRejected: this message will never be accepted (bad payload, bad credentials).
	ResultRejected
)

// SendResult is the outcome of one message of a request.
type SendResult struct {
	Kind ResultKind
	// TicketID is what the provider's receipt is looked up with; empty when it gave none.
	TicketID string
	Code     string
	Detail   string
}

// ReceiptKind is what a receipt says.
type ReceiptKind int

const (
	// ReceiptHandedOff: the provider passed the message to APNs or FCM, which took it. It is still not a delivery
	// confirmation: the phone may be off, offline, or have notifications disabled.
	ReceiptHandedOff ReceiptKind = iota
	ReceiptRetry
	ReceiptInvalidDevice
	ReceiptFailed
)

// ReceiptResult is one receipt.
type ReceiptResult struct {
	Kind   ReceiptKind
	Code   string
	Detail string
}

// RequestError means a whole request failed, not one message of it.
type RequestError struct {
	Code   string
	Status int
	// Retryable: the same request may succeed later (network, rate limit, outage, credentials to be fixed).
	Retryable bool
	// Isolate: the request was refused permanently but the cause may be one of its messages, so it is worth
	// splitting it to find out which.
	Isolate bool
	Detail  string
}

func (e *RequestError) Error() string {
	if e.Detail == "" {
		return "push provider request failed: " + e.Code
	}
	return "push provider request failed: " + e.Code + ": " + e.Detail
}

// Provider is a push service. The pipeline knows nothing else about Expo, APNs or FCM.
type Provider interface {
	Name() string
	// MaxBatch is the most messages one Send call accepts.
	MaxBatch() int
	// Send returns one result per message, in order, or a RequestError when the whole request failed.
	Send(ctx context.Context, messages []Outbound) ([]SendResult, error)
	// Receipts returns the receipts that are available; a ticket missing from the map has none yet.
	Receipts(ctx context.Context, ticketIDs []string) (map[string]ReceiptResult, error)
}
