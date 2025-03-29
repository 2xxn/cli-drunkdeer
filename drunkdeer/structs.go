package main

type RapidTriggerSettings struct {
	Enabled           bool    `json:"enabled"`
	DefaultDownstroke float32 `json:"defaultDownstroke"`
	DefaultUpstroke   float32 `json:"defaultUpstroke"`
}

type LightSettings struct {
	Enabled    bool `json:"enabled"`
	Direction  int  `json:"direction"`
	Sequence   int  `json:"sequence"`
	Speed      int  `json:"speed"`
	Brightness int  `json:"brightness"`
}

type Config struct {
	Model            string                `json:"model"`
	RapidTrigger     RapidTriggerSettings  `json:"rapidTrigger"`
	Turbo            bool                  `json:"turbo"`
	DefaultActuation float32               `json:"defaultActuation"`
	ActuationPoints  map[string]float32    `json:"actuationPoints"`
	RapidTriggers    map[string][2]float32 `json:"rapidTriggers"`
	Light            LightSettings         `json:"light"`
}

type Args struct {
	Command  string `arg:"positional"`
	CmdValue string `arg:"positional"`
	Import   string `arg:"-i,--import" help:"Import a drunkdeer webdriver profile from the specified file/url"`
	Debug    bool   `arg:"-d,--debug" help:"Enable debug mode"`
	Index    int    `arg:"-i,--index" help:"Keyboard index to use (0 for first device, 1 for second, etc.)"`
	Profiles bool   `arg:"-p,--profiles" help:"Show all available profiles"`
	Reset    bool   `arg:"-r,--reset" help:"Reset the keyboard to default settings"`
	Load     string `arg:"-L,--load" help:"Load a profile from the specified file/url"`
	Save     string `arg:"-S,--save" help:"Save URL as profile"`
	Version  bool   `arg:"-v,--version" help:"Show version information"`
	List     bool   `arg:"-l,--list" help:"List all connected devices"`
}
