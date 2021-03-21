package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	intitools "github.com/0xJeti/intitools/pkg/intigo"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	conf := &config{}

	// Create channel for accepting signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		cancel()
	}()

	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
					conf.init(os.Args)
				case os.Interrupt:
					cancel()
					os.Exit(1)
				}
			case <-ctx.Done():
				log.Printf("Done.")
				os.Exit(1)
			}
		}
	}()

	if err := run(ctx, conf, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, conf *config, out io.Writer) error {
	conf.init(os.Args)

	c := intitools.NewClient(conf.username, conf.password)
	c.SlackWebhookURL = conf.webhookurl

	log.SetOutput(os.Stdout)

	log.Printf("Starting monitoring with tick %s", conf.tick)
	httpctx := context.Background()
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-time.Tick(conf.tick):
			err := c.Authenticate()
			if err != nil {
				log.Printf("Authentication error: %s\n", err)
				continue

			}

			numActivities, err := c.CheckActivity(httpctx)
			if err != nil {
				log.Printf("CheckActivity error: %s\n", err)
				continue
			}

			if numActivities == 0 {
				continue
			}

			res, err := c.GetActivities(httpctx)

			if err != nil {
				log.Printf("GetActivities error: %s\n", err)
				continue
			}

			for idx, activity := range res.Activities {
				if idx > numActivities-1 {
					break
				}

				message, err := c.FormatActivityMessage(activity)
				if err != nil {
					log.Printf("FormatActivityMessage error: %s\n", err)
					continue
				}

				err = c.SlackSend(message)
				if err != nil {
					log.Printf("SlackSend error: %s\n", err)
					continue
				}

			}

			c.LastViewed = time.Now().Unix()
		}
	}

}
