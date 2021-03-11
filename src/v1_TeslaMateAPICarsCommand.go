package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// TeslaMateAPICarsCommandV1 func
func TeslaMateAPICarsCommandV1(c *gin.Context) {

	// creating required vars
	var TeslaAccessToken, TeslaVehicleID string
	var jsonData map[string]interface{}
	var err error
	var command string

	// check if commands are enabled.. if not we need to abort
	if getEnvAsBool("ENABLE_COMMANDS", false) == false {
		log.Println("[warning] TeslaMateAPICarsCommandV1 ENABLE_COMMANDS is not true.. returning 403 forbidden.")
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to access commands"})
		return
	}

	// if request method is GET return list of commands
	if c.Request.Method == "GET" {
		c.JSON(http.StatusOK, gin.H{"enabled_commands": allowList})
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
		log.Println("[error] TeslaMateAPICarsCommandV1 CarID is invalid (zero)!")
		c.JSON(http.StatusBadRequest, gin.H{"error": "CarID invalid"})
		return
	}

	// getting request body to pass to Tesla
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("[error] TeslaMateAPICarsCommandV1 error in first ioutil.ReadAll", err)
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

	// get TeslaVehicleID and TeslaAccessToken
	query := `
		SELECT
			eid as TeslaVehicleID,
			(SELECT access FROM tokens LIMIT 1) as TeslaAccessToken
		FROM cars
		WHERE id = $1
		LIMIT 1;`
	rows, err := db.Query(query, CarID)

	// checking for errors in query
	if err != nil {
		log.Fatal(err)
	}

	// defer closing rows
	defer rows.Close()

	// looping through all results (even if it's only one..)
	for rows.Next() {
		// scanning row and putting values into the drive
		err = rows.Scan(
			&TeslaVehicleID,
			&TeslaAccessToken,
		)
	}

	// checking for errors in query when doing scan action
	if err != nil {
		log.Println("[error] TeslaMateAPICarsCommandV1 error in sql query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://owner-api.teslamotors.com/api/1/vehicles/"+TeslaVehicleID+command, strings.NewReader(string(reqBody)))
	req.Header.Add("Authorization", "Bearer "+TeslaAccessToken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	// check response error
	if err != nil {
		log.Println("[error] TeslaMateAPICarsCommandV1 error in http request to https://owner-api.teslamotors.com:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	defer resp.Body.Close()
	defer client.CloseIdleConnections()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[error] TeslaMateAPICarsCommandV1 error in second ioutil.ReadAll:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	json.Unmarshal([]byte(respBody), &jsonData)

	// print to log about request
	if gin.IsDebugging() {
		log.Println("[debug] TeslaMateAPICarsCommandV1 " + c.Request.RequestURI + " returned data:")
		js, _ := json.Marshal(jsonData)
		log.Printf("[debug] %s\n", js)
	}

	if resp.StatusCode == http.StatusOK {
		log.Println("[info] TeslaMateAPICarsCommandV1 " + c.Request.RequestURI + " executed successful.")
	} else {
		log.Println("[error] TeslaMateAPICarsCommandV1 " + c.Request.RequestURI + " error in execution!")
	}
	c.JSON(resp.StatusCode, jsonData)
}
