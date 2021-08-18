package main

import (
	"path/filepath"
)

type Configuration struct {
	EyesTimeout     int  `json:"eyes_timeout"`
	WaterTimeout    int  `json:"water_timeout"`
	ExerciseTimeout int  `json:"exercise_timeout"`
	PlayBeep        bool `json:"play_beep"`
	Notify          bool `json:"notify"`
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
