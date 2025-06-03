/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/stianeikeland/go-rpio"
	"gopkg.in/telebot.v4"
)

var (
	// TeleToken bot
	TeleToken = os.Getenv("TELE_TOKEN")
)

type TrafficSignal struct {
	Pin int8
	On  bool
}

// cbotCmd represents the cbot command
var cbotCmd = &cobra.Command{
	Use:     "cbot",
	Aliases: []string{"start"},
	Short:   "Telegram bot for controlling traffic light signals",
	Long: `A Telegram bot that allows controlling traffic light signals through GPIO pins.
The bot accepts commands to toggle red, amber, and green lights on/off.

Usage:
  /s red    - Toggle red light
  /s amber  - Toggle amber light
  /s green  - Toggle green light
  hello     - Get a greeting from the bot`,
	Run: func(cmd *cobra.Command, args []string) {
		if TeleToken == "" {
			log.Fatal("TELE_TOKEN environment variable is not set")
		}

		fmt.Printf("cbot %s started\n", appVersion)

		cbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable. %s", err)
		}

		if err := rpio.Open(); err != nil {
			log.Fatalf("Unable to open GPIO: %s", err)
		}
		defer rpio.Close()

		// Initialize traffic signals
		trafficSignals := map[string]*TrafficSignal{
			"red":   {Pin: 12, On: false},
			"amber": {Pin: 27, On: false},
			"green": {Pin: 22, On: false},
		}

		// Initialize all pins as input
		for _, signal := range trafficSignals {
			pin := rpio.Pin(signal.Pin)
			pin.Input()
		}

		cbot.Handle(telebot.OnText, func(m telebot.Context) error {
			log.Printf("Received message: %s", m.Text())
			payload := m.Message().Payload

			switch payload {
			case "hello":
				return m.Send(fmt.Sprintf("Hello I'm cbot %s!", appVersion))

			case "red", "amber", "green":
				signal := trafficSignals[payload]
				pin := rpio.Pin(signal.Pin)

				if !signal.On {
					pin.Output()
					pin.High()
					signal.On = true
				} else {
					pin.Low()
					pin.Input()
					signal.On = false
				}

				return m.Send(fmt.Sprintf("Switched %s light %s", payload, map[bool]string{true: "on", false: "off"}[signal.On]))

			default:
				return m.Send("Usage: /s red|amber|green")
			}
		})

		cbot.Start()
	},
}

func init() {
	rootCmd.AddCommand(cbotCmd)
}
