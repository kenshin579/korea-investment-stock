// yaml_config example: NewClientFromYAML.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
)

func main() {
	path := flag.String("config", "./config.yaml", "path to config.yaml")
	flag.Parse()

	client, err := kis.NewClientFromYAML(*path)
	if err != nil {
		log.Fatal(err)
	}
	bearer, err := client.IssueAccessToken(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Token issued from yaml:", bearer[:20]+"...")
}
