package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// setting TeslaMateApi version number
// TODO: get the value from git-tag later..
var apiVersion = "1.8.0"

// defining db var
var db *sql.DB

// defining envToken that contains API_TOKEN value
var envToken string

// list of allowed commands
var allowList []string

// main function
func main() {

	// setting log parameters
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	// setting application to ReleaseMode if DEBUG_MODE is false
	if !getEnvAsBool("DEBUG_MODE", false) {
		// setting GIN_MODE to ReleaseMode
		gin.SetMode(gin.ReleaseMode)
		log.Printf("[info] TeslaMateApi running in release mode.")
	} else {
		// setting GIN_MODE to DebugMode
		gin.SetMode(gin.DebugMode)
		log.Printf("[info] TeslaMateApi running in debug mode.")
	}

	// init of API with connection to database
	initDBconnection()
	defer db.Close()

	// run initAuthToken to validate environment vars
	initAuthToken()
	// initialize allowList stored for /command section
	initCommandAllowList()

	// Connect to the MQTT broker
	statusCache, err := startMQTT()
	if err != nil {
		log.Fatalf("[error] TeslaMateApi MQTT connection failed: %s", err)
	}

	if getEnvAsBool("API_TOKEN_DISABLE", false) == true {
		log.Println("[warning] validateAuthToken - header authorization bearer token disabled. Authorizaiton: Bearer token will not be required for commands.")
	}

	// kicking off Gin in value r
	r := gin.Default()

	// root endpoint telling API is running
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "TeslaMateApi container runnnig..", "path": "/"})
	})

	// TeslaMateApi /api endpoints
	api := r.Group("/api")
	{
		// TeslaMateApi /api root
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "TeslaMateApi container runnnig..", "path": "/api"})
		})

		// TeslaMateApi /api/v1 endpoints
		v1 := api.Group("/v1")
		{
			// TeslaMateApi /api/v1 root
			v1.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "TeslaMateApi v1 runnnig..", "path": "/api/v1"})
			})

			// v1 /api/v1/cars endpoints
			v1.GET("/cars", TeslaMateAPICarsV1)
			v1.GET("/cars/:CarID", TeslaMateAPICarsV1)

			// v1 /api/v1/cars/:CarID/charges endpoints
			v1.GET("/cars/:CarID/charges", TeslaMateAPICarsChargesV1)
			v1.GET("/cars/:CarID/charges/:ChargeID", TeslaMateAPICarsChargesDetailsV1)

			// v1 /api/v1/cars/:CarID/command endpoints
			v1.GET("/cars/:CarID/command", TeslaMateAPICarsCommandV1)
			v1.GET("/cars/:CarID/commands", TeslaMateAPICarsCommandV1)
			v1.POST("/cars/:CarID/command/:Command", TeslaMateAPICarsCommandV1)

			// v1 /api/v1/cars/:CarID/drives endpoints
			v1.GET("/cars/:CarID/drives", TeslaMateAPICarsDrivesV1)
			v1.GET("/cars/:CarID/drives/:DriveID", TeslaMateAPICarsDrivesDetailsV1)

			// v1 /api/v1/cars/:CarID/logging endpoints
			v1.GET("/cars/:CarID/logging", TeslaMateAPICarsLoggingV1)
			v1.PUT("/cars/:CarID/logging/:Command", TeslaMateAPICarsLoggingV1)

			// v1 /api/v1/cars/:CarID/status endpoints
			v1.GET("/cars/:CarID/status", statusCache.TeslaMateAPICarsStatusV1)

			// v1 /api/v1/cars/:CarID/updates endpoints
			v1.GET("/cars/:CarID/updates", TeslaMateAPICarsUpdatesV1)

			// v1 /api/v1/cars/:CarID/wake_up endpoints
			v1.POST("/cars/:CarID/wake_up", TeslaMateAPICarsCommandV1)

			// v1 /api/v1/globalsettings endpoints
			v1.GET("/globalsettings", TeslaMateAPIGlobalsettingsV1)
		}

		// /api/ping endpoint
		api.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
	}

	// TeslaMateApi endpoints (before versioning)
	r.GET("/cars", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/api/v1"+c.Request.RequestURI) })
	r.GET("/cars/:CarID", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/api/v1"+c.Request.RequestURI) })
	r.GET("/cars/:CarID/charges", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/api/v1"+c.Request.RequestURI) })
	r.GET("/cars/:CarID/charges/:ChargeID", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/api/v1"+c.Request.RequestURI) })
	r.GET("/cars/:CarID/drives", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/api/v1"+c.Request.RequestURI) })
	r.GET("/cars/:CarID/drives/:DriveID", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/api/v1"+c.Request.RequestURI) })
	r.GET("/cars/:CarID/status", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/api/v1"+c.Request.RequestURI) })
	r.GET("/cars/:CarID/updates", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/api/v1"+c.Request.RequestURI) })
	r.GET("/globalsettings", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/api/v1"+c.Request.RequestURI) })

	// run this and serve on 0.0.0.0:8080
	r.Run(":8080")
}

// initDBconnection func
func initDBconnection() {

	// declare error var for use insite initAPI
	var err error
	dbsslmode := "disable"

	// creating connection string towards postgres
	dbhost := getEnv("DATABASE_HOST", "database")
	dbport := getEnvAsInt("DATABASE_PORT", 5432)
	dbuser := getEnv("DATABASE_USER", "teslamate")
	dbpass := getEnv("DATABASE_PASS", "secret")
	dbname := getEnv("DATABASE_NAME", "teslamate")
	// dbpool := getEnvAsInt("DATABASE_POOL_SIZE", 10)
	dbtimeout := (getEnvAsInt("DATABASE_TIMEOUT", 60000) / 1000)
	dbssl := getEnvAsBool("DATABASE_SSL", false)
	// dbipv6 := getEnvAsBool("DATABASE_IPV6", false)
	if dbssl {
		dbsslmode = "prefer"
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d", dbhost, dbport, dbuser, dbpass, dbname, dbsslmode, dbtimeout)

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
	if gin.IsDebugging() {
		log.Println("[debug] initDBconnection - successfully completed (connected to postgres).")
	}
}

func getTimeInTimeZone(datestring string) string {

	// getting timezone from environment
	UsersTimezone := getEnv("TZ", "Europe/Berlin")

	// format the dates are stored in postgres
	TeslaMateDateFormat := "2006-01-02T15:04:05Z"

	// parsing datestring into TeslaMateDateFormat
	t, _ := time.Parse(TeslaMateDateFormat, datestring)

	// loading timezone location
	UserLocation, _ := time.LoadLocation(UsersTimezone)

	// formatting in users location in RFC3339 format
	ReturnDate := t.In(UserLocation).Format(time.RFC3339)

	// logging time conversion to log
	if gin.IsDebugging() {
		log.Println("[debug] getTimeInTimeZone - UTC", t.Format(time.RFC3339), "time converted to", UsersTimezone, "is", ReturnDate)
	}

	return ReturnDate
}

/*
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
*/

// getEnv func - read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
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

/*
// getEnvAsFloat func - read an environment variable into a float64 or return default value
func getEnvAsFloat(name string, defaultVal float64) float64 {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseFloat(valStr, 64); err == nil {
		return val
	}
	return defaultVal
}
*/

// getEnvAsInt func - read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// convertStringToBool func
func convertStringToBool(data string) bool {
	if value, err := strconv.ParseBool(data); err == nil {
		return value
	}
	// else..
	log.Println("[warning] convertStringToBool error could not return value correct.. returning false")
	return false
}

// convertStringToFloat func
func convertStringToFloat(data string) float64 {
	if value, err := strconv.ParseFloat(data, 64); err == nil {
		return value
	}
	// else..
	log.Println("[warning] convertStringToFloat error could not return value correct.. returning 0.0")
	return 0.0
}

// convertStringToInteger func
func convertStringToInteger(data string) int {
	if value, err := strconv.Atoi(data); err == nil {
		return value
	}
	// else..
	log.Println("[warning] convertStringToInteger error could not return value correct.. returning 0")
	return 0
}

// kilometersToMiles func
func kilometersToMiles(km float64) float64 {
	return (km * 0.62137119223733)
}

// kilometersToMilesNilSupport func
func kilometersToMilesNilSupport(km NullFloat64) NullFloat64 {
	km.Float64 = (km.Float64 * 0.62137119223733)
	return (km)
}

/*
// milesToKilometers func
func milesToKilometers(mi float64) float64 {
	return (mi * 1.609344)
}

// milesToKilometersNilSupport func
func milesToKilometersNilSupport(mi NullFloat64) NullFloat64 {
	mi.Float64 = (mi.Float64 * 1.609344)
	return (mi)
}
*/

// celsiusToFahrenheit func
func celsiusToFahrenheit(c float64) float64 {
	return (c*9/5 + 32)
}

// celsiusToFahrenheitNilSupport func
func celsiusToFahrenheitNilSupport(c NullFloat64) NullFloat64 {
	c.Float64 = (c.Float64*9/5 + 32)
	return (c)
}

/*
// fahrenheitToCelsius func
func fahrenheitToCelsius(f float64) float64 {
	return ((f - 32) * 5 / 9)
}

// fahrenheitToCelsiusNilSupport func
func fahrenheitToCelsiusNilSupport(f NullFloat64) NullFloat64 {
	f.Float64 = ((f.Float64 - 32) * 5 / 9)
	return (f)
}
*/

// checkArrayContainsString func - check if string is inside stringarray
func checkArrayContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
