package main

import (
	"fmt"
	"github.com/thatisuday/commando"
	"path/filepath"
	"time"
	"github.com/inancgumus/screen"
)

const (
	NAME    = "apollo"
	VERSION = "0.1.1"
	DESCRIPTION = "apollo is a command line utility which helps you being healthy by reminding you to take breaks at fixed intervals of time."
)

func init() {
	if !exists(filepath.Join(getApolloDir(), "config.json")) {
		c := Configuration{
			EyesTimeout:     900,
			WaterTimeout:    1200,
			ExerciseTimeout: 1500,
			PlayBeep:        true,
			Notify:          true,
		}
		writeConfig(&c)
	}
}

func startLoop(jobs chan notification) {
	config := getConfig()

	eyesInit := time.Now()

	waterInit := time.Now()

	exerciseInit := time.Now()

	jobs <- notification{"Apollo Alarm", "Apollo has been started."}

	for {
		screen.MoveTopLeft()

		if (time.Since(eyesInit)).Seconds() > float64(config.EyesTimeout) {
			jobs <- notification{"Apollo Alarm", "Eyes timeout. Take a break."}
			eyesInit = time.Now()
		}
		if (time.Since(waterInit)).Seconds() > float64(config.WaterTimeout) {
			jobs <- notification{"Apollo Alarm", "Water timeout. Drink some water."}
			waterInit = time.Now()
		}
		if (time.Since(exerciseInit)).Seconds() > float64(config.ExerciseTimeout) {
			jobs <- notification{"Apollo Alarm", "Exercise timeout. Do some exercise."}
			exerciseInit = time.Now()
		}
		
		fmt.Printf("Time passed since the last eyes timeout: %v.\n", (time.Since(eyesInit)).Round(time.Second))
		fmt.Printf("Time passed since the last water timeout: %v.\n", (time.Since(waterInit)).Round(time.Second))
		fmt.Printf("Time passed since the last exercise timeout: %v.\n", (time.Since(exerciseInit)).Round(time.Second))

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
	config := getConfig()
	commando.
		Register("config").
		SetShortDescription("Configure the apollo settings.").
		SetDescription("Configure the apollo settings, i.e., the timeouts for the different tasks and whether to play beeps and send notifications.").
		AddFlag("eyes-timeout,e", "The time in seconds to wait before the eyes timeout.", commando.Int, config.EyesTimeout).
		AddFlag("water-timeout,w", "The time in seconds to wait before the water timeout.", commando.Int, config.WaterTimeout).
		AddFlag("exercise-timeout,x", "The time in seconds to wait before the exercise timeout.", commando.Int, config.ExerciseTimeout).
		AddFlag("play-beep,p", "Whether to play a beep when the timeout is reached.", commando.Bool, config.PlayBeep).
		AddFlag("notify,n", "Whether to send a notification when the timeout is reached.", commando.Bool, config.Notify).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			config := Configuration{
				EyesTimeout:     flags["eyes-timeout"].Value.(int),
				WaterTimeout:    flags["water-timeout"].Value.(int),
				ExerciseTimeout: flags["exercise-timeout"].Value.(int),
				PlayBeep:        flags["play-beep"].Value.(bool),
				Notify:          flags["notify"].Value.(bool),
			}
			writeConfig(&config)
		})

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
