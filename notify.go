package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	htgotts "github.com/hegedustibor/htgo-tts"
	handlers "github.com/hegedustibor/htgo-tts/handlers"
	voices "github.com/hegedustibor/htgo-tts/voices"
)

type notification struct {
	title string
	body  string
}

var speech = htgotts.Speech{Folder: filepath.Join(getApolloDir(), "audio"), Language: voices.English, Handler: &handlers.Native{}}

func speak(s string) {
	speech.Speak(strings.ReplaceAll(s, "_", " "))
}

func notifier(jobs <-chan notification) {
	config := getConfig()
	var wg sync.WaitGroup

	for entry := range jobs {
		wg.Add(1)

		go func() {
			defer wg.Done()
			if config.PlayBeep {
				if err := beeep.Beep(beeep.DefaultFreq, 1500); err != nil {
					fmt.Println("unable to play beep")
					fmt.Println(err)
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			if config.Notify {
				if err := beeep.Notify(entry.title, entry.body, filepath.Join(getApolloDir(), "logo.png")); err != nil {
					fmt.Println("unable to send a notification")
					fmt.Println(err)
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			if config.PlaySpeech {
				speak(entry.body)
			}
		}()

		wg.Wait()
		time.Sleep(time.Second * 5) // delay between alerts
	}
}
