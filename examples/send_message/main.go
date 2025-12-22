package sendMessage

import (
	"fmt"
	"heysender"
	"log"
)

func main() {
	client := heysender.NewClient("your-api-key", "your-api-secret")
	message := heysender.NewMessageBuilder(
		"sender@yourdomain.com",
		"Your Name",
		"Test Email",
		"<h1>Hello</h1><p>World!</p>",
	).
		AddTo("example@heysender.com", "Mr. Sender").
		AddBCC("example2@heysender.com").
		SetTracking(true).
		AddTag("tagTest", "some tag").
		SetAnonymizeOptions([]heysender.AnonymizeOption{
			heysender.AnonymizeRecipient,
			heysender.AnonymizeContent,
		}).
		Build()

	responses, err := client.SendMessage(message)

	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	for _, resp := range responses {
		fmt.Printf("Status: %s, MessageID: %s, Recipient: %s\n",
			resp.Status, resp.MessageID, resp.Recipient)
	}
}
