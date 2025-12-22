package heysender

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// NewClient creates a new Heysender client
func NewClient(apiKey, apiSecret string) *HeysenderClient {
	return &HeysenderClient{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   "https://app.heysender.com",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// request makes an HTTP request to the API
func (c *HeysenderClient) request(method, endpoint string, body interface{}) ([]byte, error) {
	url := c.BaseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication
	credentials := base64.StdEncoding.EncodeToString([]byte(c.APIKey + ":" + c.APISecret))
	req.Header.Set("Authorization", "Basic "+credentials)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, &HeysenderError{
			Message:    string(responseBody),
			StatusCode: resp.StatusCode,
		}
	}

	return responseBody, nil
}

// ==================== DOMAIN TYPES AND METHODS ====================

// GetDomains retrieves a list of domains
func (c *HeysenderClient) GetDomains() ([]Domain, error) {
	data, err := c.request("GET", "/api/domains", nil)
	if err != nil {
		return nil, err
	}

	var domains []Domain
	if err := json.Unmarshal(data, &domains); err != nil {
		return nil, fmt.Errorf("failed to unmarshal domains: %w", err)
	}

	return domains, nil
}

// CreateDomain creates a new domain
func (c *HeysenderClient) CreateDomain(req CreateDomainRequest) (*Domain, error) {
	data, err := c.request("POST", "/api/domains", req)
	if err != nil {
		return nil, err
	}

	var domain Domain
	if err := json.Unmarshal(data, &domain); err != nil {
		return nil, fmt.Errorf("failed to unmarshal domain: %w", err)
	}

	return &domain, nil
}

// UpdateDomain updates a domain with a new DKIM key
func (c *HeysenderClient) UpdateDomain(domain, dkimKey string) (string, error) {
	response, err := c.request("PUT", fmt.Sprintf("/api/domains/%s", domain), map[string]string{
		"dkim_key": dkimKey,
	})
	if err != nil {
		return "", err
	}
	var message string
	if err := json.Unmarshal(response, &message); err != nil {
		return "", fmt.Errorf("failed to unmarshal domain: %w", err)
	}
	return message, nil
}

// DeleteDomain deletes a domain
func (c *HeysenderClient) DeleteDomain(domain string) error {
	_, err := c.request("DELETE", fmt.Sprintf("/api/domains/%s", domain), nil)
	return err
}

// ValidateDomain validates domain SPF and DKIM
func (c *HeysenderClient) ValidateDomain(domain string) (*DomainValidation, error) {
	data, err := c.request("GET", fmt.Sprintf("/api/domains/%s/validate", domain), nil)
	if err != nil {
		return nil, err
	}

	var validation DomainValidation
	if err := json.Unmarshal(data, &validation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal validation: %w", err)
	}

	return &validation, nil
}

// ==================== SMTP USER TYPES AND METHODS ====================

// GetSMTPUsers retrieves SMTP users for a domain
func (c *HeysenderClient) GetSMTPUsers(domainID int) ([]SMTPUser, error) {
	data, err := c.request("GET", fmt.Sprintf("/api/smtp/%d", domainID), nil)
	if err != nil {
		return nil, err
	}

	var smtpUsers AllSMTPUserResponse
	if err := json.Unmarshal(data, &smtpUsers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SMTP users: %w", err)
	}
	users := make([]SMTPUser, len(smtpUsers.Data))
	for i, config := range smtpUsers.Data {
		users[i] = config.ToSMTPUser()
	}
	return users, nil
}

// CreateSMTPUser creates a new SMTP user
func (c *HeysenderClient) CreateSMTPUser(domainID int, req CreateSMTPUserRequest) (*SMTPUser, error) {
	if req.AnonymizeOptions == nil {
		req.AnonymizeOptions = []AnonymizeOption{AnonymizeNone}
	}

	data, err := c.request("POST", fmt.Sprintf("/api/smtp/%d", domainID), req)

	if err != nil {
		return nil, err
	}

	var user SMTPUser
	fmt.Printf("heer hits")
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SMTP user: %w", err)
	}

	return &user, nil

}

// DeleteSMTPUser deletes an SMTP user
func (c *HeysenderClient) DeleteSMTPUser(domainID, userID int) error {
	_, err := c.request("DELETE", fmt.Sprintf("/api/smtp/%d/%d", domainID, userID), nil)
	return err
}

// ResetSMTPPassword generates a new password for an SMTP user
func (c *HeysenderClient) ResetSMTPPassword(domainID, userID int) (*SMTPUser, error) {
	data, err := c.request("GET", fmt.Sprintf("/api/smtp/%d/%d/newpassword", domainID, userID), nil)
	if err != nil {
		return nil, err
	}

	var user SMTPUser
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SMTP user: %w", err)
	}

	return &user, nil
}

// ==================== WEBHOOK TYPES AND METHODS ====================

// NewWebhookRequest creates a webhook request with specified events
func NewWebhookRequest(url string, events []EventType) WebhookRequest {
	req := WebhookRequest{URL: url}

	for _, event := range events {
		switch event {
		case EventQueued:
			req.Queued = true
		case EventSent:
			req.Sent = true
		case EventAttempt:
			req.Attempt = true
		case EventSoftBounce:
			req.SoftBounce = true
		case EventHardBounce:
			req.HardBounce = true
		case EventComplaint:
			req.Complaint = true
		case EventUnsubscribe:
			req.Unsubscribe = true
		case EventOpen:
			req.Open = true
		case EventClick:
			req.Click = true
		}
	}

	return req
}

// GetWebhooks retrieves webhooks for a domain
func (c *HeysenderClient) GetWebhooks(domain string) ([]Webhook, error) {
	data, err := c.request("GET", fmt.Sprintf("/api/webhooks/%s", domain), nil)
	if err != nil {
		return nil, err
	}

	var webhooks []Webhook
	if err := json.Unmarshal(data, &webhooks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhooks: %w", err)
	}

	return webhooks, nil
}

// GetWebhook retrieves a specific webhook
func (c *HeysenderClient) GetWebhook(domain string, webhookID int) (*Webhook, error) {
	data, err := c.request("GET", fmt.Sprintf("/api/webhooks/%s/%d", domain, webhookID), nil)
	if err != nil {
		return nil, err
	}

	var webhook Webhook
	if err := json.Unmarshal(data, &webhook); err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhook: %w", err)
	}

	return &webhook, nil
}

// CreateWebhook creates a new webhook
func (c *HeysenderClient) CreateWebhook(domain string, req WebhookRequest) (map[string]interface{}, error) {
	data, err := c.request("POST", fmt.Sprintf("/api/webhooks/%s", domain), req)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}

// UpdateWebhook updates a webhook
func (c *HeysenderClient) UpdateWebhook(domain string, webhookID int, req WebhookRequest) error {
	_, err := c.request("PUT", fmt.Sprintf("/api/webhooks/%s/%d", domain, webhookID), req)
	return err
}

// DeleteWebhook deletes a webhook
func (c *HeysenderClient) DeleteWebhook(domain string, webhookID int) error {
	_, err := c.request("DELETE", fmt.Sprintf("/api/webhooks/%s/%d", domain, webhookID), nil)
	return err
}

// ==================== MESSAGE TYPES AND METHODS ====================

// SendMessage sends an email message
func (c *HeysenderClient) SendMessage(message Message) ([]MessageResponse, error) {
	data, err := c.request("POST", "/api/message", message)
	if err != nil {
		return nil, err
	}

	var responses []MessageResponse
	if err := json.Unmarshal(data, &responses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message response: %w", err)
	}

	return responses, nil
}

// GetMessage retrieves message information
func (c *HeysenderClient) GetMessage(messageID string) (*MessageInfo, error) {
	data, err := c.request("GET", fmt.Sprintf("/api/message/%s", messageID), nil)
	if err != nil {
		return nil, err
	}

	var info MessageInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message info: %w", err)
	}

	return &info, nil
}

// GetMessageByRecipient retrieves message information for a specific recipient
func (c *HeysenderClient) GetMessageByRecipient(messageID, recipient string) (*MessageInfo, error) {
	data, err := c.request("GET", fmt.Sprintf("/api/message/%s/%s", messageID, recipient), nil)
	if err != nil {
		return nil, err
	}

	var info MessageInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message info: %w", err)
	}

	return &info, nil
}

// ==================== SUPPRESSION TYPES AND METHODS ====================

// GetSuppressions retrieves suppressions by domain and type
func (c *HeysenderClient) GetSuppressions(domain string, suppressionType SuppressionType) (*SuppressionList, error) {
	data, err := c.request("GET", fmt.Sprintf("/api/suppressions/%s/%s", domain, suppressionType), nil)
	if err != nil {
		return nil, err
	}

	var list SuppressionList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("failed to unmarshal suppression list: %w", err)
	}

	return &list, nil
}

// RemoveBounce removes an email from bounce suppressions
func (c *HeysenderClient) RemoveBounce(domain, email string) error {
	_, err := c.request("DELETE", fmt.Sprintf("/api/suppressions/%s/bounce/%s", domain, email), nil)
	return err
}
