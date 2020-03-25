package gcp

import (
	"context"
	"os"

	"github.com/sjljrvis/deploynow/log"
	"golang.org/x/oauth2/google"
	dns "google.golang.org/api/dns/v1"
)

// CreateDNS will create A name record in cloud provider
func CreateDNS(repository_name string) (string, error) {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, dns.CloudPlatformScope)
	if err != nil {
		log.Error().Msgf("[GCP DNS] %s", err.Error())
		return "", err
	}

	dnsService, err := dns.New(c)
	if err != nil {
		log.Error().Msgf("[GCP DNS] %s", err.Error())
		return "", err
	}

	project := os.Getenv("GCP_PROJECT_ID")
	managedZone := os.Getenv("GCP_DNS_ZONE")

	rb := &dns.Change{
		Kind:      "dns#change",
		IsServing: true,
		Additions: []*dns.ResourceRecordSet{
			{
				Kind:    "dns#resourceRecordSet",
				Name:    repository_name + ".upweb.io.",
				Type:    "A",
				Ttl:     300,
				Rrdatas: []string{os.Getenv("01_UPWEB_IP")},
			},
		},
	}

	resp, err := dnsService.Changes.Create(project, managedZone, rb).Context(ctx).Do()
	if err != nil {
		log.Error().Msgf("[GCP DNS] %s", err.Error())
		return "", err
	}
	return resp.Id, nil
}
