package heysender

// MessageBuilder provides a fluent interface for building messages
type MessageBuilder struct {
	message Message
}

// NewMessageBuilder creates a new MessageBuilder
func NewMessageBuilder(fromEmail, fromName, subject, html string) *MessageBuilder {
	return &MessageBuilder{
		message: Message{
			FromEmail: fromEmail,
			FromName:  fromName,
			Subject:   subject,
			HTML:      html,
			To:        []Recipient{},
		},
	}
}

// SetText sets the plain text body
func (b *MessageBuilder) SetText(text string) *MessageBuilder {
	b.message.Text = text
	return b
}

// AddTo adds a TO recipient
func (b *MessageBuilder) AddTo(email, name string) *MessageBuilder {
	b.message.To = append(b.message.To, Recipient{Email: email, Name: name})
	return b
}

// AddCC adds a CC recipient
func (b *MessageBuilder) AddCC(email, name string) *MessageBuilder {
	if b.message.CC == nil {
		b.message.CC = []Recipient{}
	}
	b.message.CC = append(b.message.CC, Recipient{Email: email, Name: name})
	return b
}

// AddBCC adds a BCC recipient
func (b *MessageBuilder) AddBCC(email string) *MessageBuilder {
	if b.message.BCC == nil {
		b.message.BCC = []string{}
	}
	b.message.BCC = append(b.message.BCC, email)
	return b
}

// SetReplyTo sets the reply-to address
func (b *MessageBuilder) SetReplyTo(replyTo interface{}) *MessageBuilder {
	b.message.ReplyTo = replyTo
	return b
}

// AddAttachment adds an attachment
func (b *MessageBuilder) AddAttachment(name, base64Content string) *MessageBuilder {
	if b.message.Attachments == nil {
		b.message.Attachments = []Attachment{}
	}
	b.message.Attachments = append(b.message.Attachments, Attachment{
		Name:    name,
		Content: base64Content,
	})
	return b
}

// AddTag adds a custom tag
func (b *MessageBuilder) AddTag(key, value string) *MessageBuilder {
	if b.message.Tags == nil {
		b.message.Tags = []Tag{}
	}
	b.message.Tags = append(b.message.Tags, Tag{Key: key, Value: value})
	return b
}

// AddHeader adds a custom header
func (b *MessageBuilder) AddHeader(key, value string) *MessageBuilder {
	if b.message.Headers == nil {
		b.message.Headers = []Header{}
	}
	b.message.Headers = append(b.message.Headers, Header{Key: key, Value: value})
	return b
}

// SetCustomContent sets custom content for bulk messaging
func (b *MessageBuilder) SetCustomContent(content map[string]map[string]string) *MessageBuilder {
	b.message.CustomContent = content
	return b
}

// SetTracking enables or disables tracking
func (b *MessageBuilder) SetTracking(enabled bool) *MessageBuilder {
	b.message.Tracking = &enabled
	return b
}

// SetListUnsubscribe enables or disables list-unsubscribe header
func (b *MessageBuilder) SetListUnsubscribe(enabled bool) *MessageBuilder {
	b.message.ListUnsubscribe = &enabled
	return b
}

// SetRetentionTime sets message retention time in days
func (b *MessageBuilder) SetRetentionTime(days int) *MessageBuilder {
	b.message.RetentionTime = &days
	return b
}

func (b *MessageBuilder) SetAnonymizeOptions(options []AnonymizeOption) *MessageBuilder {
	b.message.AnonymizeOptions = options
	return b
}

// SetWebhook sets custom webhook for this message
func (b *MessageBuilder) SetWebhook(url string, events []EventType) *MessageBuilder {
	b.message.Webhook = &WebhookConfig{
		URL:    url,
		Events: events,
	}
	return b
}

// Build returns the built message
func (b *MessageBuilder) Build() Message {
	return b.message
}
