package createSmptUser

import (
	"fmt"
	"heysender"
	"log"
)

func main() {
	client := heysender.NewClient("your-api-key", "your-api-secret")
	domains, err := client.GetDomains()
	if err != nil {
		log.Fatalf("Failed to get domains: %v", err)
	}

	if len(domains) > 0 {

		for _, resp := range domains {
			fmt.Printf("ID: %d, DomainUrl: %s\n",
				resp.ID, resp.URL)
		}

		domain := &domains[0]
		email := fmt.Sprintf("example@%s", domain.URL)
		smtpReq := heysender.CreateSMTPUserRequest{
			SMTPEmail:        email,
			AnonymizeOptions: []heysender.AnonymizeOption{heysender.AnonymizeSubject, heysender.AnonymizeContent},
		}
		response, err := client.CreateSMTPUser(domain.ID, smtpReq)

		if err != nil {
			log.Fatalf("Failed to create SMTP user: %v", err)
		}

		fmt.Printf("ID: %d, DomainId: %d, Password: %s\n",
			response.ID, response.DomainID, response.SMTPPassword)

		responses, err := client.GetSMTPUsers(response.DomainID)
		if err != nil {
			log.Fatalf("Failed to get SMTP users: %v", err)
		}

		for _, resp := range responses {
			fmt.Printf("ID: %d, DomainId: %d, Options: %v\n",
				resp.ID, resp.DomainID, resp.AnonymizeOptions)
		}
	} else {
		fmt.Printf("No domains was returned")
	}
}
