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
		log.Println("[error] TeslaMateAPICarsCommand missing X-Command-Token header.. throwing error!")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthroized"})
		return
	}

	// if request method is GET return list of commands
	if c.Request.Method == "GET" {
		c.JSON(http.StatusOK, allowList)
		return
	}

	// authentication for the endpoint
	validToken, errorMessage := validateAuthToken(c)
	if !validToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorMessage})
		return
	}

	// getting CarID param from URL
	ParamCarID := c.Param("CarID")
	var CarID int
	if ParamCarID != "" {
		CarID = convertStringToInteger(ParamCarID)
	}

	// validating that CarID is not zero
	if CarID == 0 {
		log.Println("[error] TeslaMateAPICarsCommand CarID is invalid (zero)!")
		c.JSON(http.StatusBadRequest, gin.H{"error": "CarID invalid"})
		return
	}

	// getting request body to pass to Tesla
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("[error] TeslaMateAPICarsCommand error in first ioutil.ReadAll", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	// I am not a fan of hardcoding "/api/v1/"
	//   it would be nice to find a way to retrieve api.Group
	command = (c.Request.RequestURI[len("/api/v1/cars/"+ParamCarID):])

		log.Print("[warning] TeslaMateAPICarsCommand command: " + command + " not allowed")
	if !checkArrayContainsString(allowList, command) {
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
		log.Println("[error] TeslaMateAPICarsCommand error in token sql query ", err)
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
		log.Println("[error] TeslaMateAPICarsCommand error in cars sql query ", err)
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
		log.Println("[error] TeslaMateAPICarsCommand error in http request to owner-api.teslamotors.com ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	defer resp.Body.Close()
	defer client.CloseIdleConnections()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[error] TeslaMateAPICarsCommand error in second ioutil.ReadAll ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	json.Unmarshal([]byte(respBody), &jsonData)

	// print to log about request
	if gin.IsDebugging() {
		log.Println("[debug] TeslaMateAPICarsCommand " + command + " returned data:")
		js, _ := json.Marshal(jsonData)
		log.Printf("[debug] %s\n", js)
	}

	c.JSON(resp.StatusCode, jsonData)
}

func getCommandToken() string {
	// get token from environment variable COMMAND_TOKEN
	token := getEnv("COMMAND_TOKEN", "")
	if token == "" || len(token) < 32 {
		log.Println("[warning] getCommandToken environment variable COMMAND_TOKEN not set, is empty, or too short. All POST commands will return unauthroized.")
		token = ""
	}
	return token
}

	} else {
		log.Print("[info] getAllowList COMMANDS from environment variables set, " + commandAllowListLocation + " will be ignored.")
	}
}
