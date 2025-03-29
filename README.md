# Custom DrunkDeer Driver CLI
This CLI is a custom driver for the DrunkDeer keyboard. It allows you to configure the keyboard's settings, including actuation points, light settings, turbo mode and rapid trigger.

## Reason for this project
This project was created out of frustration, DrunkDeer webdriver's servers are so awful that getting into the WebDriver can sometimes take up to 10 minutes (especially uncached, I have tendencies to reload using CTRL+SHIFT+R). This project is a workaround for that, it allows you to configure the keyboard without the need for the web driver.

## REMAPPING IS NOT SUPPORTED YET

## Installation
```bash
go install github.com/2xxn/cli-drunkdeer/drunkdeer@latest
```

## Usage
```bash
drunkdeer [command] [value?] [options?]
```

### You will generally need those 4 commands
```bash
drunkdeer list - list all available devices
drunkdeer profiles - list all available profiles
drunkdeer import [path-to-config-file] - import a DRUNKDEER ANTLER CONFIG FILE
drunkdeer load [profile-name] - load a profile into the keyboard
```
#### To import someone's CLI config file you can do `drunkdeer load [url/relative or absolute path]`


### To learn more
```bash
drunkdeer
drunkdeer -h
```

## Config File structure
### This entry is for myself and the more advanced users
### The config file is a JSON file that contains the following structure:
<!-- <p>List of character names can be found in [this file](https://github.com/2xxn/cli-drunkdeer/pkg/consts.go)</p> -->
<p>Actuation point should be between 0.1mm and 3.9mm (although both are unadvised, you should do 0.2mm at lowest)</p>

```json
{
    "model": "A75",
    "turbo": false,
    "defaultActuation": 2.0,
    "rapidTrigger": {
        "enabled": false,
        "defaultDownstroke": 0.0,
        "defaultUpstroke": 0.0
    },
    "light": {
        "enabled": true,
        "direction": 0,
        "speed": 5,
        "brightness": 9,
        "sequence": 5
    },
    "actuationPoints": {
        "W": 0.2,
        "A": 0.2,
        "S": 0.2,
        "D": 0.2,
        "TAB": 3.8
    },
    "rapidTriggers": {
        "A": [0.2, 0.2],
        "S": [0.2, 0.2]
    }
}
```