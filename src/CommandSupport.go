package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// initCommandAllowList func
func initCommandAllowList() {

	// generate map of all available commands
	CommandList := make(map[string][]string)

	// https://github.com/adriankumpf/teslamate/discussions/1433
	CommandList["COMMANDS_LOGGING"] = []string{
		"/logging/resume",
		"/logging/suspend",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/wake
	CommandList["COMMANDS_WAKE"] = []string{
		"/wake_up",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/alerts
	CommandList["COMMANDS_ALERT"] = []string{
		"/command/honk_horn",
		"/command/flash_lights",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/remotestart
	CommandList["COMMANDS_REMOTESTART"] = []string{
		"/command/remote_start_drive",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/homelink
	CommandList["COMMANDS_HOMELINK"] = []string{
		"/command/trigger_homelink",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/speedlimit
	CommandList["COMMANDS_SPEEDLIMIT"] = []string{
		"/command/speed_limit_set_limit",
		"/command/speed_limit_activate",
		"/command/speed_limit_deactivate",
		"/command/speed_limit_clear_pin",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/valet
	CommandList["COMMANDS_VALET"] = []string{
		"/command/set_valet_mode",
		"/command/reset_valet_pin",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/sentrymode
	CommandList["COMMANDS_SENTRYMODE"] = []string{
		"/command/set_sentry_mode",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/doors
	CommandList["COMMANDS_DOORS"] = []string{
		"/command/door_unlock",
		"/command/door_lock",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/trunk
	CommandList["COMMANDS_TRUNK"] = []string{
		"/command/actuate_trunk",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/windows
	CommandList["COMMANDS_WINDOWS"] = []string{
		"/command/window_control",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/sunroof
	CommandList["COMMANDS_SUNROOF"] = []string{
		"/command/sun_roof_control",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/charging
	CommandList["COMMANDS_CHARGING"] = []string{
		"/command/charge_port_door_open",
		"/command/charge_port_door_close",
		"/command/charge_start",
		"/command/charge_stop",
		"/command/set_charge_limit",
		"/command/set_charging_amps",
		"/command/set_scheduled_charging",
		"/command/set_scheduled_departure",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/climate
	CommandList["COMMANDS_CLIMATE"] = []string{
		"/command/auto_conditioning_start",
		"/command/auto_conditioning_stop",
		"/command/set_temps",
		"/command/set_preconditioning_max",
		"/command/remote_seat_heater_request",
		"/command/remote_steering_wheel_heater_request",
		"/command/set_bioweapon_mode",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/media
	CommandList["COMMANDS_MEDIA"] = []string{
		"/command/media_toggle_playback",
		"/command/media_next_track",
		"/command/media_prev_track",
		"/command/media_next_fav",
		"/command/media_prev_fav",
		"/command/media_volume_up",
		"/command/media_volume_down",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/sharing
	CommandList["COMMANDS_SHARING"] = []string{
		"/command/share",
		"/command/navigation_sc_request",
	}

	// https://tesla-api.timdorr.com/vehicle/commands/softwareupdate
	CommandList["COMMANDS_SOFTWAREUPDATE"] = []string{
		"/command/schedule_software_update",
		"/command/cancel_software_update",
	}

	// not documentet and unsorted new endpoints
	CommandList["COMMANDS_UNKNOWN"] = []string{
		"/command/upcoming_calendar_entries",
		"/command/dashcam_save_clip",
	}

	// allow all commands available below
	allowAll := getEnvAsBool("COMMANDS_ALL", false)

	// looping over CommandList to generate allowList
	for key := range CommandList {
		// checking if env is set from key or if all should be allowed
		if getEnvAsBool(key, false) || allowAll {
			// appending to allowList
			allowList = append(allowList, CommandList[key]...)
		}
	}

	// if allowList is empty, read COMMANDS_ALLOWLIST and append to allowList
	commandAllowListLocation := getEnv("COMMANDS_ALLOWLIST", "allow_list.json")
	if len(allowList) == 0 {
		var allowListFile []string
		commandAllowListFile, err := os.Open(commandAllowListLocation)
		if err != nil {
			log.Println("[error] getAllowList error with COMMANDS_ALLOWLIST: " + commandAllowListLocation + " not found and will be ignored")
		} else {
			defer commandAllowListFile.Close()
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

// encryptAccessToken func
func encryptAccessToken(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	return fmt.Sprintf("%x", ciphertext)
}

// decryptAccessToken func
func decryptAccessToken(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}
