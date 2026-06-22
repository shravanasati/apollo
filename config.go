package main

import (
	"fmt"
	"path/filepath"
	"time"
)

type Configuration struct {
	Timeouts   map[string]int `json:"timeouts"`
	PlayBeep   bool           `json:"play_beep"`
	Notify     bool           `json:"notify"`
	PlaySpeech bool           `json:"play_speech"`
}

// getConfig reads the config file and returns a Configuration struct.
func getConfig() *Configuration {
	configFile := filepath.Join(getApolloDir(), "config.json")
	configContent, e := readFile(configFile)
	if e != nil {
		panic(e)
	}

	config, e := readJSON(configContent)
	if e != nil {
		panic(e)
	}

	return config
}

// writeConfig writes the config file.
func writeConfig(c *Configuration) {
	configFile := filepath.Join(getApolloDir(), "config.json")
	jsonStr, e := jsonify(c)
	if e != nil {
		panic(e)
	}
	writeFile(configFile, jsonStr)
}

type ParsedDurations map[string]time.Duration

func parseDurationsFromConfig(c *Configuration) ParsedDurations {
	pt := make(ParsedDurations)
	for k, v := range c.Timeouts {
		timeoutDuration, err := time.ParseDuration(fmt.Sprintf("%vs", v))
		if err != nil {
			continue
		}
		pt[k] = timeoutDuration
	}
	return pt
}
