package main

import "context"

type DnsHost interface {
	checkAndUpdate(ctx context.Context, recordType string, value string) error
	init(config *Configuration) error
}
