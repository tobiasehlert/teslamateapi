package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// initCommandAllowList func
func initCommandAllowList() {

	// allow all commands available below
	allowAll := getEnvAsBool("COMMANDS_ALL", false)

	// https://github.com/adriankumpf/teslamate/discussions/1433
	if getEnvAsBool("COMMANDS_LOGGING", false) || allowAll {
		allowList = append(allowList,
			"/logging/resume",
			"/logging/suspend")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/wake
	if getEnvAsBool("COMMANDS_WAKE", false) || allowAll {
		allowList = append(allowList, "/wake_up")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/alerts
	if getEnvAsBool("COMMANDS_ALERT", false) || allowAll {
		allowList = append(allowList,
			"/command/honk_horn",
			"/command/flash_lights")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/remotestart
	if getEnvAsBool("COMMANDS_REMOTESTART", false) || allowAll {
		allowList = append(allowList, "/command/remote_start_drive")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/homelink
	if getEnvAsBool("COMMANDS_HOMELINK", false) || allowAll {
		allowList = append(allowList, "/command/trigger_homelink")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/speedlimit
	if getEnvAsBool("COMMANDS_SPEEDLIMIT", false) || allowAll {
		allowList = append(allowList,
			"/command/speed_limit_set_limit",
			"/command/speed_limit_activate",
			"/command/speed_limit_deactivate",
			"/command/speed_limit_clear_pin")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/valet
	if getEnvAsBool("COMMANDS_VALET", false) || allowAll {
		allowList = append(allowList,
			"/command/set_valet_mode",
			"/command/reset_valet_pin")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/sentrymode
	if getEnvAsBool("COMMANDS_SENTRYMODE", false) || allowAll {
		allowList = append(allowList, "/command/set_sentry_mode")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/doors
	if getEnvAsBool("COMMANDS_DOORS", false) || allowAll {
		allowList = append(allowList,
			"/command/door_unlock",
			"/command/door_lock")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/trunk
	if getEnvAsBool("COMMANDS_TRUNK", false) || allowAll {
		allowList = append(allowList, "/command/actuate_trunk")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/windows
	if getEnvAsBool("COMMANDS_WINDOWS", false) || allowAll {
		allowList = append(allowList, "/command/window_control")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/sunroof
	if getEnvAsBool("COMMANDS_SUNROOF", false) || allowAll {
		allowList = append(allowList, "/command/sun_roof_control")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/charging
	if getEnvAsBool("COMMANDS_CHARGING", false) || allowAll {
		allowList = append(allowList,
			"/command/charge_port_door_open",
			"/command/charge_port_door_close",
			"/command/charge_start",
			"/command/charge_stop",
			"/command/charge_standard",
			"/command/charge_max_range",
			"/command/set_charge_limit")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/climate
	if getEnvAsBool("COMMANDS_CLIMATE", false) || allowAll {
		allowList = append(allowList,
			"/command/auto_conditioning_start",
			"/command/auto_conditioning_stop",
			"/command/set_temps",
			"/command/set_preconditioning_max",
			"/command/remote_seat_heater_request",
			"/command/remote_steering_wheel_heater_request")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/media
	if getEnvAsBool("COMMANDS_MEDIA", false) || allowAll {
		allowList = append(allowList,
			"/command/media_toggle_playback",
			"/command/media_next_track",
			"/command/media_prev_track",
			"/command/media_next_fav",
			"/command/media_prev_fav",
			"/command/media_volume_up",
			"/command/media_volume_down")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/sharing
	if getEnvAsBool("COMMANDS_SHARING", false) || allowAll {
		allowList = append(allowList, "/command/share")
	}

	// https://tesla-api.timdorr.com/vehicle/commands/softwareupdate
	if getEnvAsBool("COMMANDS_SOFTWAREUPDATE", false) || allowAll {
		allowList = append(allowList,
			"/command/schedule_software_update",
			"/command/cancel_software_update")
	}

	// if allowList is empty, read COMMANDS_ALLOWLIST and append to allowList
	commandAllowListLocation := getEnv("COMMANDS_ALLOWLIST", "allow_list.json")
	if len(allowList) == 0 {
		var allowListFile []string
		commandAllowListFile, err := os.Open(commandAllowListLocation)
		defer commandAllowListFile.Close()
		if err != nil {
			log.Println("[error] getAllowList error with COMMANDS_ALLOWLIST: " + commandAllowListLocation + " not found and will be ignored")
		} else {
			byteValue, err := ioutil.ReadAll(commandAllowListFile)
			if err != nil {
				log.Println("[error] getAllowList error while reading COMMANDS_ALLOWLIST: " + commandAllowListLocation + " it will be ignored")
			} else {
				err = json.Unmarshal(byteValue, &allowListFile)
				if err != nil {
					log.Println("[error] getAllowList error while parsing JSON.. COMMANDS_ALLOWLIST: " + commandAllowListLocation + " it will be ignored")
				} else {
					allowList = append(allowList, allowListFile...)
				}
			}
		}
	} else {
		log.Print("[info] getAllowList COMMANDS from environment variables set, " + commandAllowListLocation + " will be ignored.")
	}

	if gin.IsDebugging() {
		log.Println("[info] initCommandAllowList - generated following list of allowed commands: " + strings.Join(allowList, ", "))
	}
}
