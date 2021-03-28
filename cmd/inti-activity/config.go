package main

import (
	"fmt"
	"os"
	"time"

	"github.com/namsral/flag"
)

const defaultTick = 60 * time.Second

type config struct {
	tick        time.Duration
	username    string
	password    string
	webhookurl  string
	webhooktype string
	sendlast    int
}

func (c *config) init(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.String(flag.DefaultConfigFlagname, "", "Path to config file")

	var (
		tick        = flags.Duration("tick", defaultTick, "Ticking interval")
		username    = flags.String("username", "", "Intigriti username (e-mail)")
		password    = flags.String("password", "", "Intigriti password")
		webhookurl  = flags.String("webhook", "", "Webhook URL")
		webhooktype = flags.String("type", "slack", "Webhook type [slack|discord]")
		sendlast    = flags.Int("last", 0, "Number of activity entries sent on start (for debugging)")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	if *username == "" || *password == "" || *webhookurl == "" {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flags.PrintDefaults()
		os.Exit(1)
	}

	c.username = *username
	c.tick = *tick
	c.password = *password
	c.webhookurl = *webhookurl
	c.webhooktype = *webhooktype
	c.sendlast = *sendlast

	return nil
}
