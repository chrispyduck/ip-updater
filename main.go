package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	config, err := bootstrap()
	if err != nil {
		log.Fatal(err)
	}

	var provider = new(Cloudflare)
	err = provider.init(config)
	if err != nil {
		log.Fatal("Failed to initialize Cloudflare provider: ", err)
	}

	err = checkAndUpdateIP(ctx, config.QueryUrl.IPv4, "A", provider)
	if err != nil {
		log.Error("Failed to update A record: ", err)
	}

	err = checkAndUpdateIP(ctx, config.QueryUrl.IPv6, "AAAA", provider)
	if err != nil {
		log.Error("Failed to update AAAA record: ", err)
	}
}

func checkAndUpdateIP(ctx context.Context, queryUrl string, recordType string, provider DnsHost) error {
	ip, err := getIP(queryUrl)
	if err != nil {
		return fmt.Errorf("failed to obtain current WAN IP address from %s: %w", queryUrl, err)
	}

	childContext, cancelFunc := context.WithTimeout(ctx, 5000*time.Millisecond)
	err = provider.checkAndUpdate(childContext, recordType, ip)
	cancelFunc()

	if err != nil {
		return fmt.Errorf("failed to update WAN IP address: %w", err)
	}

	return nil
}

func getIP(requestUrl string) (string, error) {
	response, err := http.Get(requestUrl)
	if err != nil {
		return "", fmt.Errorf("failed to obtain IP from remote API at %s: %w", requestUrl, err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response from remote API %s: %w", requestUrl, err)
	}

	rawIP := string(body)
	parsedIP := net.ParseIP(rawIP)
	if parsedIP == nil {
		return "", fmt.Errorf("invalid IP address returned from remote API: %s", rawIP)
	}

	return string(body), nil
}
