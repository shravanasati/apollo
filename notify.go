package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gen2brain/beeep"
)

type notification struct {
	title string
	body  string
}

func notifier(jobs <-chan notification) {
	config := getConfig()

	for entry := range jobs {
		if config.PlayBeep {
			if err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration); err != nil {
				fmt.Println("unable to play beep")
				fmt.Println(err)
			}
		}
		if config.Notify {
			if err := beeep.Notify(entry.title, entry.body, filepath.Join(getApolloDir(), "logo.png")); err != nil {
				fmt.Println("unable to send a notification")
				fmt.Println(err)
			}
		}
		time.Sleep(time.Second * 5)
	}
}
