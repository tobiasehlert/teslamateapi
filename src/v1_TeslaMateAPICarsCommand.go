package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// TeslaMateAPICarsCommandV1 func
func TeslaMateAPICarsCommandV1(c *gin.Context) {

	// creating required vars
	var accessToken string
	var vehicleID string
	var jsonData map[string]interface{}
	var err error
	var command string
	var requestToken string

	// verify X-Command-Token
	requestToken = c.Request.Header.Get("X-Command-Token")
	if requestToken != commandToken || requestToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthroized"})
		return
	}

	// if request method is GET return list of commands
	if c.Request.Method == "GET" {
		c.JSON(http.StatusOK, allowList)
		return
	}

	// getting CarID param from URL
	ParamCarID := c.Param("CarID")
	var CarID int
	if ParamCarID != "" {
		CarID = convertStringToInteger(ParamCarID)
	}

	if CarID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CarID invalid"})
		return
	}

	// getting request body to pass to Tesla
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	// I am not a fan of hardcoding "/api/v1/"
	//   it would be nice to find a way to retrieve api.Group
	command = (c.Request.RequestURI[len("/api/v1/cars/"+ParamCarID):])

	if !contains(allowList, command) {
		log.Print("command: " + command + " not allowed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthroized"})
		return
	}

	// allow /command/wake_up but silently `redirect` it to /wake_up
	if command == "/command/wake_up" {
		command = "/wake_up"
	}

	// getting access token
	err = db.QueryRow(`
		SELECT
			access
		FROM tokens
		LIMIT 1;
	`).Scan(&accessToken)

	// checking for errors in query
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	// get vehicle ID
	err = db.QueryRow(`
		SELECT
			eid
		FROM cars
		WHERE id = $1
		LIMIT 1;
	`, CarID).Scan(&vehicleID)

	// ToDo: ?cleanup DB connections? -- I can't find an example of closing db.QueryRow() ¯\_(ツ)_/¯

	// checking for errors in query
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://owner-api.teslamotors.com/api/1/vehicles/"+vehicleID+command, strings.NewReader(string(reqBody)))
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	// check response error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	defer resp.Body.Close()
	defer client.CloseIdleConnections()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	json.Unmarshal([]byte(respBody), &jsonData)

	// print to log about request
	if gin.IsDebugging() {
		log.Println("[TeslaMateApi] TeslaMateAPICarsCommand " + command + " returned data:")
		js, _ := json.Marshal(jsonData)
		log.Printf("%s\n", js)
	}

	c.JSON(resp.StatusCode, jsonData)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getCommandToken() string {
	// get token from environment variable COMMAND_TOKEN
	token := getEnv("COMMAND_TOKEN", "")
	if token == "" || len(token) < 32 {
		log.Println("Environment variable COMMAND_TOKEN not set, is empty, or too short. All commands will return unauthroized.")
		token = ""
	}
	return token
}

func getAllowList() []string {

	var allowAll bool = getEnvAsBool("COMMANDS_ALL", false)

	// https://tesla-api.timdorr.com/vehicle/commands/wake
	if getEnvAsBool("COMMANDS_WAKE", false) || allowAll {
		allowList = append(allowList,
			"/command/wake_up",
			"/wake_up")
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
	if getEnvAsBool("COMMAND_SUNROOF", false) || allowAll {
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
		if err != nil {
			log.Println("COMMANDS_ALLOWLIST: " + commandAllowListLocation + " not found and will be ignored")
		} else {
			byteValue, err := ioutil.ReadAll(commandAllowListFile)
			if err != nil {
				log.Println("error reading COMMANDS_ALLOWLIST: " + commandAllowListLocation + " it will be ignored")
			} else {
				err = json.Unmarshal(byteValue, &allowListFile)
				if err != nil {
					log.Println("error parsing JSON.. COMMANDS_ALLOWLIST: " + commandAllowListLocation + " it will be ignored")
				} else {
					allowList = append(allowList, allowListFile...)
					commandAllowListFile.Close()
				}
			}
		}
	} else {
		log.Print("COMMANDS from environemnt variables set, " + commandAllowListLocation + " will be ignored.")
	}

	log.Println("List of allowed Commands: " + strings.Join(allowList, ", "))

	return allowList

}
