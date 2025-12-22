package heysender

import (
	"encoding/json"
	"time"
)

// Domain represents a domain in Heysender
type Domain struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Validated bool   `json:"-"`
}

// Handles the Validated field being returned as both int and bool depending on the endpoint
func (d *Domain) UnmarshalJSON(data []byte) error {
	type Alias Domain

	aux := &struct {
		Validated interface{} `json:"validated"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch v := aux.Validated.(type) {
	case bool:
		d.Validated = v
	case int:
		d.Validated = v != 0
	case string:
		d.Validated = v == "1" || v == "true"
	default:
		d.Validated = false
	}

	return nil
}

// DomainValidation represents domain validation status
type DomainValidation struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	Validated int       `json:"validated"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DKIM      struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"dkim"`
	SPF     string `json:"spf"`
	Current struct {
		SPF  string `json:"spf"`
		DKIM string `json:"dkim"`
	} `json:"current"`
	SPFValid  bool `json:"spfValid"`
	DKIMValid bool `json:"dkimValid"`
}

// CreateDomainRequest represents a request to create a domain
type CreateDomainRequest struct {
	URL            string  `json:"url"`
	CustomSelector *string `json:"custom_selector,omitempty"`
	DKIMKey        *string `json:"dkim_key,omitempty"`
}

// SMTPUser represents an SMTP user
type SMTPUser struct {
	ID               int               `json:"id"`
	DomainID         int               `json:"domain_id"`
	SMTPEmail        string            `json:"smtp_email"`
	AnonymizeOptions []AnonymizeOption `json:"anonymize_options"`
	SMTPPassword     string            `json:"smtp_password,omitempty"` // Only returned on creation
}

// AllSMTPUserResponse represents the response returned by the get all endpoint
type AllSMTPUserResponse struct {
	Data []AllSMTPUser `json:"data"`
}

// AllSMTPUser representes the slightly different SMTPUser returned by the get all endpoint
type AllSMTPUser struct {
	ID                 int    `json:"id"`
	DomainID           int    `json:"domain_id"`
	SMTPEmail          string `json:"smtp_email"`
	AnonymizeAll       int    `json:"anonymize_all"`
	AnonymizeNone      int    `json:"anonymize_none"`
	AnonymizeSubject   int    `json:"anonymize_subject"`
	AnonymizeContent   int    `json:"anonymize_content"`
	AnonymizeRecipient int    `json:"anonymize_recipient"`
}

// Converts AllSMTPUser to normal SMTPUser struct
func (e AllSMTPUser) ToSMTPUser() SMTPUser {
	var options []AnonymizeOption

	// Check if AnonymizeAll is set, or all individual options are 1
	if e.AnonymizeAll == 1 || (e.AnonymizeSubject == 1 && e.AnonymizeContent == 1 && e.AnonymizeRecipient == 1) {
		options = []AnonymizeOption{AnonymizeAll}
	} else if e.AnonymizeNone == 1 || (e.AnonymizeSubject == 0 && e.AnonymizeContent == 0 && e.AnonymizeRecipient == 0) {
		// AnonymizeNone is explicitly set, or all individual options are 0
		options = []AnonymizeOption{AnonymizeNone}
	} else {
		// Build array based on individual flags
		if e.AnonymizeSubject != 0 {
			options = append(options, AnonymizeSubject)
		}
		if e.AnonymizeContent != 0 {
			options = append(options, AnonymizeContent)
		}
		if e.AnonymizeRecipient != 0 {
			options = append(options, AnonymizeRecipient)
		}
	}

	return SMTPUser{
		ID:               e.ID,
		DomainID:         e.DomainID,
		SMTPEmail:        e.SMTPEmail,
		AnonymizeOptions: options,
	}
}

// CreateSMTPUserRequest represents a request to create an SMTP user
type CreateSMTPUserRequest struct {
	SMTPEmail        string            `json:"smtp_email"`
	AnonymizeOptions []AnonymizeOption `json:"anonymize_options,omitempty"`
}

// Webhook represents a webhook configuration
type Webhook struct {
	ID          int       `json:"id"`
	DomainID    int       `json:"domain_id"`
	URL         string    `json:"url"`
	LatestCode  int       `json:"latest_code"`
	LastSuccess string    `json:"last_success"`
	Queued      bool      `json:"queued"`
	Sent        bool      `json:"sent"`
	Attempt     bool      `json:"attempt"`
	SoftBounce  bool      `json:"soft_bounce"`
	HardBounce  bool      `json:"hard_bounce"`
	Complaint   bool      `json:"complaint"`
	Unsubscribe bool      `json:"unsubscribe"`
	Open        bool      `json:"open"`
	Click       bool      `json:"click"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// WebhookRequest represents a request to create or update a webhook
type WebhookRequest struct {
	URL         string `json:"url"`
	Queued      bool   `json:"queued"`
	Sent        bool   `json:"sent"`
	Attempt     bool   `json:"attempt"`
	SoftBounce  bool   `json:"soft_bounce"`
	HardBounce  bool   `json:"hard_bounce"`
	Complaint   bool   `json:"complaint"`
	Unsubscribe bool   `json:"unsubscribe"`
	Open        bool   `json:"open"`
	Click       bool   `json:"click"`
}

// Recipient represents an email recipient
type Recipient struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// Attachment represents an email attachment
type Attachment struct {
	Name    string `json:"name"`
	Content string `json:"content"` // Base64 encoded
}

// Tag represents a custom tag
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Header represents a custom header
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// WebhookConfig represents webhook configuration for a message
type WebhookConfig struct {
	URL    string      `json:"url"`
	Events []EventType `json:"events"`
}

// Message represents an email message
type Message struct {
	FromEmail        string                       `json:"from_email"`
	FromName         string                       `json:"from_name"`
	HTML             string                       `json:"html"`
	Text             string                       `json:"text,omitempty"`
	Subject          string                       `json:"subject"`
	To               []Recipient                  `json:"to"`
	CC               []Recipient                  `json:"cc,omitempty"`
	BCC              []string                     `json:"bcc,omitempty"`
	ReplyTo          interface{}                  `json:"reply_to,omitempty"` // Can be string or []Recipient
	CustomContent    map[string]map[string]string `json:"custom_content,omitempty"`
	Attachments      []Attachment                 `json:"attachments,omitempty"`
	Tags             []Tag                        `json:"tags,omitempty"`
	Headers          []Header                     `json:"headers,omitempty"`
	RetentionTime    *int                         `json:"retention_time,omitempty"`
	ListUnsubscribe  *bool                        `json:"list_unsubscribe,omitempty"`
	Tracking         *bool                        `json:"tracking,omitempty"`
	AnonymizeOptions []AnonymizeOption            `json:"anonymize_options,omitempty"`
	Webhook          *WebhookConfig               `json:"webhook,omitempty"`
}

// MessageResponse represents the response after sending a message
type MessageResponse struct {
	Status    string `json:"status"`
	MessageID string `json:"message_id"`
	Recipient string `json:"recipient"`
	Error     string `json:"error,omitempty"`
}

// MessageInfo represents detailed message information
type MessageInfo struct {
	Protocol      string              `json:"protocol"`
	FromEmail     string              `json:"from_email"`
	Subject       string              `json:"subject"`
	Domain        string              `json:"domain"`
	SMTPUser      int                 `json:"smtp_user,omitempty"`
	MessageID     string              `json:"message_id"`
	Status        string              `json:"status"`
	DeleteAt      time.Time           `json:"delete_at"`
	HTMLContent   string              `json:"html_content"`
	TextContent   string              `json:"text_content"`
	Attachments   []map[string]string `json:"attachments"`
	Events        []map[string]string `json:"events"`
	Headers       []map[string]string `json:"headers"`
	Tags          []map[string]string `json:"tags"`
	RecipientType string              `json:"recipient_type"`
	ToEmail       string              `json:"to_email"`
	CreatedAt     time.Time           `json:"created_at"`
}

// Suppression represents a suppression entry
type Suppression struct {
	Address   string    `json:"address"`
	Code      string    `json:"code"`
	Error     string    `json:"error"`
	CreatedAt time.Time `json:"created_at"`
	Type      string    `json:"type"`
}

// SuppressionList represents a paginated list of suppressions
type SuppressionList struct {
	CurrentPage  int           `json:"current_page"`
	Data         []Suppression `json:"data"`
	FirstPageURL string        `json:"first_page_url"`
	From         int           `json:"from"`
	LastPage     int           `json:"last_page"`
	LastPageURL  string        `json:"last_page_url"`
	NextPageURL  *string       `json:"next_page_url"`
	Path         string        `json:"path"`
	PerPage      int           `json:"per_page"`
	PrevPageURL  *string       `json:"prev_page_url"`
	To           int           `json:"to"`
	Total        int           `json:"total"`
}
