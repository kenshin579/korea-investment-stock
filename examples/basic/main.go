// Basic example: NewClient + IssueAccessToken.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	kis "github.com/kenshin579/korea-investment-stock"
)

func main() {
	client, err := kis.NewClient(
		os.Getenv("KOREA_INVESTMENT_API_KEY"),
		os.Getenv("KOREA_INVESTMENT_API_SECRET"),
		os.Getenv("KOREA_INVESTMENT_ACCOUNT_NO"),
	)
	if err != nil {
		log.Fatal(err)
	}
	bearer, err := client.IssueAccessToken(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Token issued:", bearer[:20]+"...")
}
