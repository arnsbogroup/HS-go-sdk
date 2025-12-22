package heysender

// AnonymizeOption represents anonymization options for email privacy compliance
type AnonymizeOption string

const (
	AnonymizeNone      AnonymizeOption = "none"
	AnonymizeAll       AnonymizeOption = "all"
	AnonymizeRecipient AnonymizeOption = "recipient"
	AnonymizeSubject   AnonymizeOption = "subject"
	AnonymizeContent   AnonymizeOption = "content"
)

// AnonymizeOptions returns all valid anonymization options
func AnonymizeOptions() []AnonymizeOption {
	return []AnonymizeOption{
		AnonymizeNone,
		AnonymizeAll,
		AnonymizeRecipient,
		AnonymizeSubject,
		AnonymizeContent,
	}
}

// String returns the string representation of the anonymization option
func (a AnonymizeOption) String() string {
	return string(a)
}

// EventType represents webhook event types
type EventType string

const (
	EventQueued      EventType = "queued"
	EventSent        EventType = "sent"
	EventAttempt     EventType = "attempt"
	EventSoftBounce  EventType = "soft_bounce"
	EventHardBounce  EventType = "hard_bounce"
	EventComplaint   EventType = "complaint"
	EventUnsubscribe EventType = "unsubscribe"
	EventOpen        EventType = "open"
	EventClick       EventType = "click"
)

// EventTypes returns all valid event types
func EventTypes() []EventType {
	return []EventType{
		EventQueued,
		EventSent,
		EventAttempt,
		EventSoftBounce,
		EventHardBounce,
		EventComplaint,
		EventUnsubscribe,
		EventOpen,
		EventClick,
	}
}

// String returns the string representation of the event type
func (e EventType) String() string {
	return string(e)
}

// SuppressionType represents suppression list types
type SuppressionType string

const (
	SuppressionBounces      SuppressionType = "bounce"
	SuppressionUnsubscribes SuppressionType = "unsubscribe"
	SuppressionComplaints   SuppressionType = "complaint"
)

// SuppressionTypes returns all valid suppression types
func SuppressionTypes() []SuppressionType {
	return []SuppressionType{
		SuppressionBounces,
		SuppressionUnsubscribes,
		SuppressionComplaints,
	}
}

// String returns the string representation of the suppression type
func (s SuppressionType) String() string {
	return string(s)
}
