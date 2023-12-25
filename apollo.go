package main

import (
	"fmt"
	"github.com/inancgumus/screen"
	"github.com/thatisuday/commando"
	"path/filepath"
	"time"
)

const (
	NAME        = "apollo"
	VERSION     = "0.2.0"
	DESCRIPTION = "apollo is a command line utility which helps you being healthy by reminding you to take breaks at fixed intervals of time."
)

func init() {
	if !exists(filepath.Join(getApolloDir(), "config.json")) {
		c := Configuration{
			Timeouts: map[string]int{
				"eyes_timeout":     900,
				"water_timeout":    1200,
				"exercise_timeout": 1500,
			},
			PlayBeep: true,
			Notify:   true,
		}
		writeConfig(&c)
	}
}

func initTimeoutsFromConfig(config Configuration) map[string]time.Time {
	timeoutMap := map[string]time.Time{}
	for k := range config.Timeouts {
		timeoutMap[k] = time.Now()
	}
	return timeoutMap
}

func startLoop(jobs chan notification) {
	config := getConfig()

	timeouts_inits := initTimeoutsFromConfig(*config)

	jobs <- notification{"Apollo Alarm", "Apollo has been started."}

	for {
		screen.MoveTopLeft()

		for k, v := range config.Timeouts {
			if time.Since(timeouts_inits[k]).Seconds() > float64(v) {
				jobs <- notification{"Apollo Alarm", fmt.Sprintf("%s: Take a break.", k)}
				timeouts_inits[k] = time.Now()
			}
		}

		for k, v := range timeouts_inits {
			passed := time.Since(v).Round(time.Second)
			timeoutDuration, err := time.ParseDuration(fmt.Sprintf("%vs", config.Timeouts[k]))
			if err != nil {
				panic("unable to parse timeouts!")
			}
			remaining := timeoutDuration - passed

			fmt.Printf("Time passed, remaining since the last %s: %v, %v.\n", k, passed, remaining)
		}

		time.Sleep(time.Second)
		screen.Clear()
	}
}

func main() {
	fmt.Println(NAME, VERSION)
	go deletePreviousInstallation()

	// Create a command line application.
	commando.
		SetExecutableName(NAME).
		SetVersion(VERSION).
		SetDescription(DESCRIPTION)

	// the root command
	commando.
		Register(nil).
		SetShortDescription("Start the apollo loop.").
		SetDescription("Start the main apollo loop, i.e., reminding about taking breaks.").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			jobs := make(chan notification, 5)
			go notifier(jobs)

			startLoop(jobs)
		})

	// config command
	// config := getConfig()
	// commando.
	// 	Register("config").
	// 	SetShortDescription("Configure the apollo settings.").
	// 	SetDescription("Configure the apollo settings, i.e., the timeouts for the different tasks and whether to play beeps and send notifications.").
	// 	AddFlag("play-beep,p", "Whether to play a beep when the timeout is reached.", commando.Bool, config.PlayBeep).
	// 	AddFlag("notify,n", "Whether to send a notification when the timeout is reached.", commando.Bool, config.Notify).
	// 	SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	// 		config := Configuration{
	// 			EyesTimeout:     flags["eyes-timeout"].Value.(int),
	// 			WaterTimeout:    flags["water-timeout"].Value.(int),
	// 			ExerciseTimeout: flags["exercise-timeout"].Value.(int),
	// 			PlayBeep:        flags["play-beep"].Value.(bool),
	// 			Notify:          flags["notify"].Value.(bool),
	// 		}
	// 		writeConfig(&config)
	// 	})

	// update command
	commando.
		Register("update").
		SetShortDescription("Update the apollo to the latest version.").
		SetDescription("Update the apollo to the latest version.").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			update()
		})

	// run the command line application
	commando.Parse(nil)
}
