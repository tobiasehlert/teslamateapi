package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// defining db var
var db *sql.DB

// main function
func main() {

	// init of API with connection to database
	initAPI()
	defer db.Close()

	// kicking off Gin in value r
	r := gin.Default()

	// /cars endpoint
	r.GET("/cars", func(c *gin.Context) {
		// TeslaMateAPICars to get cars
		result, ValidResponse := TeslaMateAPICars()

		c.Header("Content-Type", "application/json")
		if ValidResponse {
			c.String(http.StatusOK, result)
		} else {
			c.String(http.StatusNotFound, result)
		}
	})

	// /ping endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// run this! :)
	r.Run()
}

// initAPI func
func initAPI() {

	// declare error var for use insite initAPI
	var err error

	// creating connection string towards postgres
	dbhost := getEnv("TM_DB_HOST", "database")
	dbport := getEnvAsInt("TM_DB_PORT", 5432)
	dbuser := getEnv("TM_DB_USER", "teslamate")
	dbpass := getEnv("TM_DB_PASS", "secret")
	dbname := getEnv("TM_DB_NAME", "teslamate")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpass, dbname)

	// opening connection to postgres
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic(err)
	}

	// doing ping to database to test connection
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	// showing database successfully connected
	log.Println("successfully connected to postgres.")
}

// getEnv func - read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// getEnvAsInt func - read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// getEnvAsBool func - read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

// convertStringToBool func
func convertStringToBool(data string) bool {
	if value, err := strconv.ParseBool(data); err == nil {
		return value
	}
	// else..
	log.Println("Error: convertStringToBool could not return value correct.. returning false")
	return false
}

// convertStringToFloat func
func convertStringToFloat(data string) float64 {
	if value, err := strconv.ParseFloat(data, 64); err == nil {
		return value
	}
	// else..
	log.Println("Error: convertStringToFloat could not return value correct.. returning 0.0")
	return 0.0
}

// convertStringToInteger func
func convertStringToInteger(data string) int {
	if value, err := strconv.Atoi(data); err == nil {
		return value
	}
	// else..
	log.Println("Error: convertStringToInteger could not return value correct.. returning 0")
	return 0
}

// kilometersToMiles func
func kilometersToMiles(km float64) float64 {
	return (km * 0.62137119223733)
}

// milesToKilometers func
func milesToKilometers(mi float64) float64 {
	return (mi * 1.609344)
}

// celsiusToFahrenheit func
func celsiusToFahrenheit(c float64) float64 {
	return (c*9/5 + 32)
}

// fahrenheitToCelsius func
func fahrenheitToCelsius(f float64) float64 {
	return ((f - 32) * 5 / 9)
}
