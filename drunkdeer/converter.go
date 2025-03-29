package main

import (
	"encoding/json"
)

type DDConfigKey struct {
	Keyname    string  `json:"keyname"`
	Actuation  float32 `json:"action_point"`
	Downstroke float32 `json:"downstroke"`
	Upstroke   float32 `json:"upstroke"`
}

type DDConfig struct {
	StorageName string        `json:"storagename"`
	Showname    string        `json:"showname"`
	Keys        []DDConfigKey `json:"keys_array"`
}

func (c *DDConfig) getMostUsedActuation() float32 {
	var distinctActuations map[float32]int = make(map[float32]int)
	for _, key := range c.Keys {
		distinctActuations[key.Actuation]++
	}

	var mostUsedActuation float32 = 0
	var mostUsedCount int = 0
	for actuation, count := range distinctActuations {
		if count > mostUsedCount {
			mostUsedCount = count
			mostUsedActuation = actuation
		}
	}

	return mostUsedActuation
}

func (c *DDConfig) getMostUsedDownstroke() float32 {
	var distinctDownstrokes map[float32]int = make(map[float32]int)
	for _, key := range c.Keys {
		distinctDownstrokes[key.Downstroke]++
	}

	var mostUsedDownstroke float32 = 0
	var mostUsedCount int = 0
	for downstroke, count := range distinctDownstrokes {
		if count > mostUsedCount {
			mostUsedCount = count
			mostUsedDownstroke = downstroke
		}
	}

	return mostUsedDownstroke
}

func (c *DDConfig) getMostUsedUpstroke() float32 {
	var distinctUpstrokes map[float32]int = make(map[float32]int)
	for _, key := range c.Keys {
		distinctUpstrokes[key.Upstroke]++
	}

	var mostUsedUpstroke float32 = 0
	var mostUsedCount int = 0
	for upstroke, count := range distinctUpstrokes {
		if count > mostUsedCount {
			mostUsedCount = count
			mostUsedUpstroke = upstroke
		}
	}

	return mostUsedUpstroke
}

func (c *DDConfig) getModelFromStorageName() string {
	sub := c.StorageName[5:8]

	return sub
}

func (c *DDConfig) convertToCLIConfig() *Config {
	var config Config

	config.Model = c.getModelFromStorageName()
	config.DefaultActuation = c.getMostUsedActuation()
	config.RapidTrigger.DefaultDownstroke = c.getMostUsedDownstroke()
	config.RapidTrigger.DefaultUpstroke = c.getMostUsedUpstroke()
	config.RapidTrigger.Enabled = true

	config.Turbo = false

	config.ActuationPoints = make(map[string]float32)
	config.RapidTriggers = make(map[string][2]float32)

	config.Light = LightSettings{}
	config.Light.Enabled = false

	for _, key := range c.Keys {
		var createRtEntry bool = false

		if key.Actuation != 0 && key.Actuation != config.DefaultActuation {
			config.ActuationPoints[key.Keyname] = key.Actuation
		}

		if key.Downstroke != 0 && key.Downstroke != config.RapidTrigger.DefaultDownstroke {
			createRtEntry = true
		}

		if key.Upstroke != 0 && key.Upstroke != config.RapidTrigger.DefaultUpstroke {
			createRtEntry = true
		}

		if createRtEntry {
			config.RapidTriggers[key.Keyname] = [2]float32{key.Downstroke, key.Upstroke}
		}
	}

	return &config
}

func parseDrunkDeerConfig(config []byte) *DDConfig {
	var ddConfig DDConfig
	err := json.Unmarshal(config, &ddConfig)
	handleError("Error parsing JSON", err)

	return &ddConfig
}
