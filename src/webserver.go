package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

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

	// setting application to ReleaseMode if DEBUG_MODE is false
	if getEnvAsBool("DEBUG_MODE", false) == false {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// kicking off Gin in value r
	r := gin.Default()

	// root endpoint telling API is running
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "TeslaMateApi runnnig..",
		})
	})

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

	// /cars/<CarID> endpoint
	r.GET("/cars/:CarID", func(c *gin.Context) {
		// getting CarID param from URL
		CarID := convertStringToInteger(c.Param("CarID"))
		// TeslaMateAPICars to get cars
		ValidCarID := TeslaMateAPICarsID(CarID)

		if ValidCarID {
			// redirecting to /cars/<CarID>/status page..
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/cars/%d/status", CarID))
		} else {
			// redirecting to /cars since we got nothing to show for only one car..
			c.Redirect(http.StatusFound, "/cars")
		}
	})

	// /cars/<CarID>/charges endpoint
	r.GET("/cars/:CarID/charges", func(c *gin.Context) {
		// getting CarID param from URL
		CarID := convertStringToInteger(c.Param("CarID"))
		// query options to modify query when collecting data
		ResultPage := convertStringToInteger(c.DefaultQuery("page", "1"))
		ResultShow := convertStringToInteger(c.DefaultQuery("show", "100"))
		// TeslaMateAPICarsCharges to get data
		result, ValidResponse := TeslaMateAPICarsCharges(CarID, ResultPage, ResultShow)

		c.Header("Content-Type", "application/json")
		if ValidResponse {
			c.String(http.StatusOK, result)
		} else {
			c.String(http.StatusNotFound, result)
		}
	})

	// /cars/<CarID>/charges/<ChargeID> endpoint
	r.GET("/cars/:CarID/charges/:ChargeID", func(c *gin.Context) {
		// getting CarID and ChargeID param from URL
		CarID := convertStringToInteger(c.Param("CarID"))
		ChargeID := convertStringToInteger(c.Param("ChargeID"))
		// TeslaMateAPICarsChargesDetails to get data
		result, ValidResponse := TeslaMateAPICarsChargesDetails(CarID, ChargeID)

		c.Header("Content-Type", "application/json")
		if ValidResponse {
			c.String(http.StatusOK, result)
		} else {
			c.String(http.StatusNotFound, result)
		}
	})

	// /cars/<CarID>/drives endpoint
	r.GET("/cars/:CarID/drives", func(c *gin.Context) {
		// getting CarID param from URL
		CarID := convertStringToInteger(c.Param("CarID"))
		// query options to modify query when collecting data
		ResultPage := convertStringToInteger(c.DefaultQuery("page", "1"))
		ResultShow := convertStringToInteger(c.DefaultQuery("show", "100"))
		// TeslaMateAPICarsDrives to get data
		result, ValidResponse := TeslaMateAPICarsDrives(CarID, ResultPage, ResultShow)

		c.Header("Content-Type", "application/json")
		if ValidResponse {
			c.String(http.StatusOK, result)
		} else {
			c.String(http.StatusNotFound, result)
		}
	})

	// /cars/<CarID>/drives/<DriveID> endpoint
	r.GET("/cars/:CarID/drives/:DriveID", func(c *gin.Context) {
		// getting CarID and DriveID param from URL
		CarID := convertStringToInteger(c.Param("CarID"))
		DriveID := convertStringToInteger(c.Param("DriveID"))
		// TeslaMateAPICarsDrivesDetails to get data
		result, ValidResponse := TeslaMateAPICarsDrivesDetails(CarID, DriveID)

		c.Header("Content-Type", "application/json")
		if ValidResponse {
			c.String(http.StatusOK, result)
		} else {
			c.String(http.StatusNotFound, result)
		}
	})

	// /cars/<CarID>/status endpoint
	r.GET("/cars/:CarID/status", func(c *gin.Context) {
		// getting mqtt flag
		mqttdisabled := getEnvAsBool("DISABLE_MQTT", false)
		if mqttdisabled {
			c.JSON(http.StatusNotImplemented, gin.H{
				"error": "mqtt disabled.. status not accessible!",
			})
		} else {
			// getting CarID param from URL
			CarID := convertStringToInteger(c.Param("CarID"))
			// TeslaMateAPICarsStatus to get data
			result, ValidResponse := TeslaMateAPICarsStatus(CarID)

			c.Header("Content-Type", "application/json")
			if ValidResponse {
				c.String(http.StatusOK, result)
			} else {
				c.String(http.StatusNotFound, result)
			}
		}
	})

	// /cars/<CarID>/updates endpoint
	r.GET("/cars/:CarID/updates", func(c *gin.Context) {
		// getting CarID param from URL
		CarID := convertStringToInteger(c.Param("CarID"))
		// query options to modify query when collecting data
		ResultPage := convertStringToInteger(c.DefaultQuery("page", "1"))
		ResultShow := convertStringToInteger(c.DefaultQuery("show", "100"))
		// TeslaMateAPICarsUpdates to get data
		result, ValidResponse := TeslaMateAPICarsUpdates(CarID, ResultPage, ResultShow)

		c.Header("Content-Type", "application/json")
		if ValidResponse {
			c.String(http.StatusOK, result)
		} else {
			c.String(http.StatusNotFound, result)
		}
	})

	// /globalsettings endpoint
	r.GET("/globalsettings", func(c *gin.Context) {
		// TeslaMateAPIGlobalSettings to get data
		result, ValidResponse := TeslaMateAPIGlobalSettings()

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

	// run this and serve on 0.0.0.0:8080
	r.Run(":8080")
}

// initAPI func
func initAPI() {

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
	log.Println("successfully connected to postgres.")
}

func getTimeInTimeZone(datestring string) string {

	// getting timezone from environment
	UsersTimezone := getEnv("TZ", "Europe/Berlin")

	// format the dates are stored in postgres
	TeslaMateDateFormat := "2006-01-02T15:04:05Z"

	// parsing datestring into TeslaMateDateFormat
	t, _ := time.Parse(TeslaMateDateFormat, datestring)

	// logging UTC time to log
	// log.Println("UTC: ", t.Format(time.RFC3339))

	// loading timezone location
	UserLocation, _ := time.LoadLocation(UsersTimezone)

	// formatting in users location in RFC3339 format
	ReturnDate := t.In(UserLocation).Format(time.RFC3339)

	// logging Users location-converted time to log
	// log.Println("Location at User: " + ReturnDate)

	return ReturnDate
}

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

// kilometersToMilesNilSupport func
func kilometersToMilesNilSupport(km NullFloat64) NullFloat64 {
	km.Float64 = (km.Float64 * 0.62137119223733)
	return (km)
}

// milesToKilometers func
func milesToKilometers(mi float64) float64 {
	return (mi * 1.609344)
}

// milesToKilometersNilSupport func
func milesToKilometersNilSupport(mi NullFloat64) NullFloat64 {
	mi.Float64 = (mi.Float64 * 1.609344)
	return (mi)
}

// celsiusToFahrenheit func
func celsiusToFahrenheit(c float64) float64 {
	return (c*9/5 + 32)
}

// celsiusToFahrenheitNilSupport func
func celsiusToFahrenheitNilSupport(c NullFloat64) NullFloat64 {
	c.Float64 = (c.Float64*9/5 + 32)
	return (c)
}

// fahrenheitToCelsius func
func fahrenheitToCelsius(f float64) float64 {
	return ((f - 32) * 5 / 9)
}

// fahrenheitToCelsiusNilSupport func
func fahrenheitToCelsiusNilSupport(f NullFloat64) NullFloat64 {
	f.Float64 = ((f.Float64 - 32) * 5 / 9)
	return (f)
}
