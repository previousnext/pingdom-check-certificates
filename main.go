package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/certifi/gocertifi"
	"github.com/nlopes/slack"
	"github.com/russellcardullo/go-pingdom/pingdom"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const portHTTPS = "443"

var (
	cliSlackReady    = kingpin.Flag("slack-channel-ready", "Slack channel to post timeouts to").Required().OverrideDefaultFromEnvar("SLACK_CHANNEL_READY").String()
	cliSlackNotReady = kingpin.Flag("slack-channel-not-ready", "Slack channel to post good certificates to").Required().OverrideDefaultFromEnvar("SLACK_CHANNEL_NOT_READY").String()
	cliSlackToken    = kingpin.Flag("slack-token", "Slack token for authentication").Required().OverrideDefaultFromEnvar("SLACK_TOKEN").String()
	cliPingdomUser   = kingpin.Flag("pingdom-user", "Pingdom username for authentication").Required().OverrideDefaultFromEnvar("PINGDOM_USER").String()
	cliPingdomPass   = kingpin.Flag("pingdom-pass", "Pingdom password for authentication").Required().OverrideDefaultFromEnvar("PINGDOM_PASS").String()
	cliPingdomToken  = kingpin.Flag("pingdom-token", "Pingdom token for authentication").Required().OverrideDefaultFromEnvar("PINGDOM_TOKEN").String()
	cliDays          = kingpin.Flag("days", "Days until you get a slack alert").Default("30").OverrideDefaultFromEnvar("DAYS").Int()
)

func main() {
	kingpin.Parse()

	var (
		client = pingdom.NewClient(*cliPingdomUser, *cliPingdomPass, *cliPingdomToken)
		api    = slack.New(*cliSlackToken)
		params = slack.PostMessageParameters{
			Username:  "Certificate Bot",
			IconEmoji: ":bot_certificate:",
		}
	)

	checks, err := client.Checks.List()
	if err != nil {
		panic(err)
	}

	// Check the certificates for each host.
	for _, check := range checks {
		timeout, err := getTimeout(check.Hostname, portHTTPS)
		if err != nil {
			log.Println(err)
			continue
		}

		if timeout > *cliDays {
			_, _, err = api.PostMessage(*cliSlackNotReady, fmt.Sprintf("Certificate is not about to expire: *%s* (*%d days*)", check.Hostname, timeout), params)
			if err != nil {
				panic(err)
			}

			continue
		}

		_, _, err = api.PostMessage(*cliSlackReady, fmt.Sprintf("Certificate is about to expire: *%s* (*%d  days*)", check.Hostname, timeout), params)
		if err != nil {
			panic(err)
		}
	}
}

// Helper function to get the timeout of a certificate of a host.
func getTimeout(host, port string) (int, error) {
	cas, err := gocertifi.CACerts()
	if err != nil {
		return 0, err
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", host, port), &tls.Config{
		RootCAs: cas,
	})
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	now := time.Now()

	for _, chain := range conn.ConnectionState().VerifiedChains {
		for _, cert := range chain {
			if len(cert.DNSNames) < 1 {
				continue
			}

			return days(cert.NotAfter.Sub(now)), nil
		}
	}

	return 0, nil
}

// Helper function to convert hours to days.
func days(diff time.Duration) int {
	return int(diff.Hours() / 24)
}
