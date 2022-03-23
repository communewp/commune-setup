package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type SetupOptions struct {
	ApexDomain             string
	NginxDefaultConfig     bool
	NginxInitializeCerts   bool
	SimonAddDefaultDomains bool
}

var setupOpts = new(SetupOptions)

func simonAddDefaultDomains() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redisint:6379",
		Password: "",
		DB:       0,
	})

	err := rdb.SAdd(ctx, "simon:adminer."+setupOpts.ApexDomain, "adminer:8080").Err()
	if err != nil {
		log.Fatal(err)
	}
}

func getBoolOrTryAgain(question string, res *bool) {
	fmt.Println(question)
	answer := ""
	for {
		fmt.Scanln(&answer)
		if answer == "yes" {
			*res = true
			break
		} else if answer == "no" {
			*res = false
			break
		}
		fmt.Println("Required: yes or no")
	}
}

func main() {
	fmt.Println("Enter the apex (top-level) domain for your Commune instance: ")
	fmt.Scanln(&setupOpts.ApexDomain)

	getBoolOrTryAgain("Should the Commune default OpenResty configuration be applied?", &setupOpts.NginxDefaultConfig)
	getBoolOrTryAgain("Should default services subdomains be configured in openresty?", &setupOpts.SimonAddDefaultDomains)
	getBoolOrTryAgain("Should a self-signed SSL certificate be generated for SSL fallback?", &setupOpts.NginxInitializeCerts)
	out, err := json.Marshal(&setupOpts)
	if err != nil {
		log.Fatal("Failed to marshal")
	}
	fmt.Println(string(out))
}
