// env_config example: NewClientFromEnv.
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	bearer, err := client.IssueAccessToken(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Token issued from env:", bearer[:20]+"...")
}
