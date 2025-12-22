package createDomain

import (
	"fmt"
	"heysender"
	"log"
)

func main() {
	client := heysender.NewClient("your-api-key", "your-api-secret")
	domainReq := heysender.CreateDomainRequest{
		URL: "example.heysender.com",
	}
	response, err := client.CreateDomain(domainReq)

	if err != nil {
		log.Fatalf("Failed to create domain: %v", err)
	}

	fmt.Printf("ID: %d, URL: %s, Validated: %t\n",
		response.ID, response.URL, response.Validated)

}
