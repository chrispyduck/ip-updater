package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

type Cloudflare struct {
	api *cloudflare.API

	zoneId   string
	zoneName string
	hostname string

	initialized bool
}

func (c *Cloudflare) init(config *Configuration) error {
	log.Debug("Logging into Cloudflare API")
	api, err := cloudflare.NewWithAPIToken(config.Cloudflare.ApiToken, cloudflare.Debug(config.Debug))
	if err != nil {
		return fmt.Errorf("failed to initialize Cloudflare API: %w", err)
	}
	c.api = api

	id, err := api.ZoneIDByName(config.Cloudflare.Zone)
	if err != nil {
		return fmt.Errorf("failed to obtain zone ID for domain %s: %w", config.Cloudflare.Zone, err)
	}
	c.zoneId = id
	c.zoneName = config.Cloudflare.Zone
	c.hostname = config.Hostname
	c.initialized = true
	return nil
}

func (c *Cloudflare) checkAndUpdate(ctx context.Context, recordType string, value string) error {
	if !c.initialized {
		return fmt.Errorf("the cloudflare provider has not been initialized")
	}

	query := cloudflare.DNSRecord{
		Name: c.hostname,
		Type: recordType,
	}
	records, err := c.api.DNSRecords(ctx, c.zoneId, query)
	if err != nil {
		return fmt.Errorf("failed to get current value for %s record %s: %w", recordType, c.hostname, err)
	}
	log.Debugf("Found %d matching record(s) for %s %s", len(records), c.hostname, recordType)
	if len(records) == 0 {
		// no match - make new record
		_, err := c.api.CreateDNSRecord(ctx, c.zoneId, query)
		if err != nil {
			return fmt.Errorf("unable to create new %s record %s: %w", recordType, c.hostname, err)
		}
		log.Infof("Created new record %s %s pointing to %s", c.hostname, recordType, value)
	} else if len(records) == 1 {
		// one match - update it
		record := records[0]
		if record.Content == value {
			log.Infof("No change needed to %s record %s", recordType, c.hostname)
			return nil
		}
		record.Content = value
		err := c.api.UpdateDNSRecord(ctx, c.zoneId, record.ID, record)
		if err != nil {
			return fmt.Errorf("unable to update existing %s record %s: %w", recordType, c.hostname, err)
		}
		log.Infof("Updated %s %s to %s", c.hostname, recordType, value)
	} else {
		return fmt.Errorf("found an unexpected number of matching records for %s %s. Unable to update. Matches = %+v", recordType, c.hostname, records)
	}

	return nil
}
