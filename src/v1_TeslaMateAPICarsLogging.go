package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// TeslaMateAPICarsLoggingV1 func
func TeslaMateAPICarsLoggingV1(c *gin.Context) {

	// creating required vars
	var (
		jsonData map[string]any
		err      error
	)

	// check if commands are enabled.. if not we need to abort
	if !getEnvAsBool("ENABLE_COMMANDS", false) {
		log.Println("[warning] TeslaMateAPICarsLoggingV1 ENABLE_COMMANDS is not true.. returning 403 forbidden.")
		TeslaMateAPIHandleOtherResponse(c, http.StatusForbidden, "TeslaMateAPICarsLoggingV1", gin.H{"error": "You are not allowed to access logging commands"})
		return
	}

	// if request method is GET return list of commands
	if c.Request.Method == http.MethodGet {
		TeslaMateAPIHandleSuccessResponse(c, "TeslaMateAPICarsLoggingV1", gin.H{"enabled_commands": allowList})
		return
	}

	// authentication for the endpoint
	validToken, errorMessage := validateAuthToken(c)
	if !validToken {
		TeslaMateAPIHandleOtherResponse(c, http.StatusUnauthorized, "TeslaMateAPICarsLoggingV1", gin.H{"error": errorMessage})
		return
	}

	// getting CarID param from URL and validating that it's not zero
	CarID := convertStringToInteger(c.Param("CarID"))
	if CarID == 0 {
		log.Println("[error] TeslaMateAPICarsLoggingV1 CarID is invalid (zero)!")
		TeslaMateAPIHandleOtherResponse(c, http.StatusBadRequest, "TeslaMateAPICarsLoggingV1", gin.H{"error": "CarID invalid"})
		return
	}

	// getting request body to pass to Tesla
	reqBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("[error] TeslaMateAPICarsLoggingV1 error in first io.ReadAll", err)
		TeslaMateAPIHandleOtherResponse(c, http.StatusInternalServerError, "TeslaMateAPICarsLoggingV1", gin.H{"error": "internal io reading error"})
		return
	}

	// getting :Command
	command := ("/logging/" + c.Param("Command"))

	if !checkArrayContainsString(allowList, command) {
		log.Printf("[warning] TeslaMateAPICarsLoggingV1 command not allowed!")
		TeslaMateAPIHandleOtherResponse(c, http.StatusUnauthorized, "TeslaMateAPICarsLoggingV1", gin.H{"error": "unauthorized"})
		return
	}

	client := &http.Client{}
	putURL := ""
	if getEnvAsBool("TESLAMATE_SSL", false) {
		putURL = "https://"
	} else {
		putURL = "http://"
	}
	putURL = putURL + getEnv("TESLAMATE_HOST", "teslamate") + ":" + getEnv("TESLAMATE_PORT", "4000") + "/api/car/" + strconv.Itoa(CarID) + command
	req, _ := http.NewRequest(http.MethodPut, putURL, strings.NewReader(string(reqBody)))
	req.Header.Set("User-Agent", "TeslaMateApi/"+apiVersion+" https://github.com/tobiasehlert/teslamateapi")
	resp, err := client.Do(req)

	// check response error
	if err != nil {
		log.Println("[error] TeslaMateAPICarsLoggingV1 error in http request to http://teslamate:", err)
		TeslaMateAPIHandleOtherResponse(c, http.StatusInternalServerError, "TeslaMateAPICarsLoggingV1", gin.H{"error": "internal http request error"})
		return
	}

	defer resp.Body.Close()
	defer client.CloseIdleConnections()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[error] TeslaMateAPICarsLoggingV1 error in second io.ReadAll:", err)
		TeslaMateAPIHandleOtherResponse(c, http.StatusInternalServerError, "TeslaMateAPICarsLoggingV1", gin.H{"error": "internal io reading error"})
		return
	}
	json.Unmarshal([]byte(respBody), &jsonData)

	// return jsonData
	// use TeslaMateAPIHandleOtherResponse since we use the statusCode from Tesla API
	TeslaMateAPIHandleOtherResponse(c, resp.StatusCode, "TeslaMateAPICarsLoggingV1", jsonData)
}
