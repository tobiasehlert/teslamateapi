package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// TeslaMateAPICarsChargesV1 func
func TeslaMateAPICarsChargesV1(c *gin.Context) {

	// define error messages
	var CarsChargesError1 = "Unable to load charges."
	var CarsChargesError2 = "Invalid date format. Please use ISO 8601 format (e.g., 2025-06-25T13:26:34+02:00)."

	// getting CarID param from URL
	CarID := convertStringToInteger(c.Param("CarID"))
	// query options to modify query when collecting data
	ResultPage := convertStringToInteger(c.DefaultQuery("page", "1"))
	ResultShow := convertStringToInteger(c.DefaultQuery("show", "100"))
	StartDate := c.DefaultQuery("startDate", "")
	EndDate := c.DefaultQuery("endDate", "")

	// Parse and validate date parameters
	var ParsedStartDate, ParsedEndDate string

	if StartDate != "" {
		// URL decode the date parameter to handle + signs in timezone offsets
		decodedStartDate, err := url.QueryUnescape(StartDate)
		if err != nil {
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesV1", CarsChargesError2, fmt.Sprintf("Failed to decode startDate parameter: %s", err.Error()))
			return
		}

		// Fix the timezone format by replacing spaces with + for timezone offsets
		// This handles cases where URL decoding converted + to space
		if len(decodedStartDate) > 19 && decodedStartDate[19] == ' ' {
			decodedStartDate = decodedStartDate[:19] + "+" + decodedStartDate[20:]
		}

		// Try to parse the date in various common formats
		var parsedTime time.Time
		formats := []string{
			time.RFC3339,                // 2006-01-02T15:04:05Z07:00
			"2006-01-02T15:04:05-07:00", // Alternative timezone format
			"2006-01-02T15:04:05+07:00", // Alternative timezone format
			"2006-01-02T15:04:05+07:05", // Handle minute offsets like +02:05
			"2006-01-02T15:04:05-07:05", // Handle negative minute offsets
			"2006-01-02T15:04:05Z",      // UTC format
			"2006-01-02 15:04:05",       // Basic format without timezone
			"2006-01-02",                // Date only
		}

		var parseErr error
		for _, format := range formats {
			parsedTime, parseErr = time.Parse(format, decodedStartDate)
			if parseErr == nil {
				break
			}
		}

		if parseErr != nil {
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesV1", CarsChargesError2, fmt.Sprintf("Invalid startDate format '%s': %s", decodedStartDate, parseErr.Error()))
			return
		}

		// Convert to UTC and format for PostgreSQL
		ParsedStartDate = parsedTime.UTC().Format("2006-01-02 15:04:05")
	}

	if EndDate != "" {
		// URL decode the date parameter to handle + signs in timezone offsets
		decodedEndDate, err := url.QueryUnescape(EndDate)
		if err != nil {
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesV1", CarsChargesError2, fmt.Sprintf("Failed to decode endDate parameter: %s", err.Error()))
			return
		}

		// Fix the timezone format by replacing spaces with + for timezone offsets
		// This handles cases where URL decoding converted + to space
		if len(decodedEndDate) > 19 && decodedEndDate[19] == ' ' {
			decodedEndDate = decodedEndDate[:19] + "+" + decodedEndDate[20:]
		}

		// Try to parse the date in various common formats
		var parsedTime time.Time
		formats := []string{
			time.RFC3339,                // 2006-01-02T15:04:05Z07:00
			"2006-01-02T15:04:05-07:00", // Alternative timezone format
			"2006-01-02T15:04:05+07:00", // Alternative timezone format
			"2006-01-02T15:04:05+07:05", // Handle minute offsets like +02:05
			"2006-01-02T15:04:05-07:05", // Handle negative minute offsets
			"2006-01-02T15:04:05Z",      // UTC format
			"2006-01-02 15:04:05",       // Basic format without timezone
			"2006-01-02",                // Date only
		}

		var parseErr error
		for _, format := range formats {
			parsedTime, parseErr = time.Parse(format, decodedEndDate)
			if parseErr == nil {
				break
			}
		}

		if parseErr != nil {
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesV1", CarsChargesError2, fmt.Sprintf("Invalid endDate format '%s': %s", decodedEndDate, parseErr.Error()))
			return
		}

		// Convert to UTC and format for PostgreSQL
		ParsedEndDate = parsedTime.UTC().Format("2006-01-02 15:04:05")
	}

	// creating structs for /cars/<CarID>/charges
	// Car struct - child of Data
	type Car struct {
		CarID   int        `json:"car_id"`   // smallint
		CarName NullString `json:"car_name"` // text (nullable)
	}
	// BatteryDetails struct - child of Charges
	type BatteryDetails struct {
		StartBatteryLevel int `json:"start_battery_level"` // int
		EndBatteryLevel   int `json:"end_battery_level"`   // int
	}
	// PreferredRange struct - child of Charges
	type PreferredRange struct {
		StartRange float64 `json:"start_range"` // float64
		EndRange   float64 `json:"end_range"`   // float64
	}
	// Charges struct - child of Data
	type Charges struct {
		ChargeID          int            `json:"charge_id"`           // int
		StartDate         string         `json:"start_date"`          // string
		EndDate           string         `json:"end_date"`            // string
		Address           string         `json:"address"`             // string
		ChargeEnergyAdded float64        `json:"charge_energy_added"` // float64
		ChargeEnergyUsed  float64        `json:"charge_energy_used"`  // float64
		Cost              float64        `json:"cost"`                // float64
		DurationMin       int            `json:"duration_min"`        // int
		DurationStr       string         `json:"duration_str"`        // string
		BatteryDetails    BatteryDetails `json:"battery_details"`     // BatteryDetails
		RangeIdeal        PreferredRange `json:"range_ideal"`         // PreferredRange
		RangeRated        PreferredRange `json:"range_rated"`         // PreferredRange
		OutsideTempAvg    float64        `json:"outside_temp_avg"`    // float64
		Odometer          float64        `json:"odometer"`            // float64
	}
	// TeslaMateUnits struct - child of Data
	type TeslaMateUnits struct {
		UnitsLength      string `json:"unit_of_length"`      // string
		UnitsTemperature string `json:"unit_of_temperature"` // string
	}
	// Data struct - child of JSONData
	type Data struct {
		Car            Car            `json:"car"`
		Charges        []Charges      `json:"charges"`
		TeslaMateUnits TeslaMateUnits `json:"units"`
	}
	// JSONData struct - main
	type JSONData struct {
		Data Data `json:"data"`
	}

	// creating required vars
	var (
		CarName                       NullString
		ChargesData                   []Charges
		UnitsLength, UnitsTemperature string
	)

	// calculate offset based on page (page 0 is not possible, since first page is minimum 1)
	if ResultPage > 0 {
		ResultPage--
	} else {
		ResultPage = 0
	}
	ResultPage = (ResultPage * ResultShow)

	// getting data from database
	query := `
		SELECT
			charging_processes.id AS charge_id,
			start_date,
			end_date,
			COALESCE(geofence.name, CONCAT_WS(', ', COALESCE(address.name, nullif(CONCAT_WS(' ', address.road, address.house_number), '')), address.city)) AS address,
			COALESCE(charge_energy_added, 0) AS charge_energy_added,
			COALESCE(charge_energy_used, 0) AS charge_energy_used,
			COALESCE(cost, 0) AS cost,
			start_ideal_range_km AS start_ideal_range,
			end_ideal_range_km AS end_ideal_range,
			start_rated_range_km AS start_rated_range,
			end_rated_range_km AS end_rated_range,
			start_battery_level,
			end_battery_level,
			duration_min,
			TO_CHAR((duration_min * INTERVAL '1 minute'), 'HH24:MI') as duration_str,
			outside_temp_avg,
			position.odometer as odometer,
			(SELECT unit_of_length FROM settings LIMIT 1) as unit_of_length,
			(SELECT unit_of_temperature FROM settings LIMIT 1) as unit_of_temperature,
			cars.name
		FROM charging_processes
		LEFT JOIN cars ON car_id = cars.id
		LEFT JOIN addresses address ON address_id = address.id
		LEFT JOIN positions position ON position_id = position.id
		LEFT JOIN geofences geofence ON geofence_id = geofence.id
		WHERE charging_processes.car_id=$1 AND charging_processes.end_date IS NOT NULL`

	// Parameters to be passed to the query
	var queryParams []interface{}
	queryParams = append(queryParams, CarID)
	paramIndex := 2

	// Add date filtering if provided
	if ParsedStartDate != "" {
		query += fmt.Sprintf(" AND charging_processes.start_date >= $%d", paramIndex)
		queryParams = append(queryParams, ParsedStartDate)
		paramIndex++
	}
	if ParsedEndDate != "" {
		query += fmt.Sprintf(" AND charging_processes.start_date <= $%d", paramIndex)
		queryParams = append(queryParams, ParsedEndDate)
		paramIndex++
	}

	query += `
		ORDER BY start_date DESC
		LIMIT $` + fmt.Sprintf("%d", paramIndex) + ` OFFSET $` + fmt.Sprintf("%d", paramIndex+1) + `;`

	queryParams = append(queryParams, ResultShow, ResultPage)

	rows, err := db.Query(query, queryParams...)

	// checking for errors in query
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesV1", CarsChargesError1, err.Error())
		return
	}

	// defer closing rows
	defer rows.Close()

	// looping through all results
	for rows.Next() {

		// creating charge object based on struct
		charge := Charges{}

		// scanning row and putting values into the charge
		err = rows.Scan(
			&charge.ChargeID,
			&charge.StartDate,
			&charge.EndDate,
			&charge.Address,
			&charge.ChargeEnergyAdded,
			&charge.ChargeEnergyUsed,
			&charge.Cost,
			&charge.RangeIdeal.StartRange,
			&charge.RangeIdeal.EndRange,
			&charge.RangeRated.StartRange,
			&charge.RangeRated.EndRange,
			&charge.BatteryDetails.StartBatteryLevel,
			&charge.BatteryDetails.EndBatteryLevel,
			&charge.DurationMin,
			&charge.DurationStr,
			&charge.OutsideTempAvg,
			&charge.Odometer,
			&UnitsLength,
			&UnitsTemperature,
			&CarName,
		)

		// converting values based of settings UnitsLength
		if UnitsLength == "mi" {
			charge.RangeIdeal.StartRange = kilometersToMiles(charge.RangeIdeal.StartRange)
			charge.RangeIdeal.EndRange = kilometersToMiles(charge.RangeIdeal.EndRange)
			charge.RangeRated.StartRange = kilometersToMiles(charge.RangeRated.StartRange)
			charge.RangeRated.EndRange = kilometersToMiles(charge.RangeRated.EndRange)
			charge.Odometer = kilometersToMiles(charge.Odometer)
		}
		// converting values based of settings UnitsTemperature
		if UnitsTemperature == "F" {
			charge.OutsideTempAvg = celsiusToFahrenheit(charge.OutsideTempAvg)
		}

		// adjusting to timezone differences from UTC to be userspecific
		charge.StartDate = getTimeInTimeZone(charge.StartDate)
		charge.EndDate = getTimeInTimeZone(charge.EndDate)

		// checking for errors after scanning
		if err != nil {
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesV1", CarsChargesError1, err.Error())
			return
		}

		// appending charge to ChargesData
		ChargesData = append(ChargesData, charge)
	}

	// checking for errors in the rows result
	err = rows.Err()
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesV1", CarsChargesError1, err.Error())
		return
	}

	//
	// build the data-blob
	jsonData := JSONData{
		Data{
			Car: Car{
				CarID:   CarID,
				CarName: CarName,
			},
			Charges: ChargesData,
			TeslaMateUnits: TeslaMateUnits{
				UnitsLength:      UnitsLength,
				UnitsTemperature: UnitsTemperature,
			},
		},
	}

	// return jsonData
	TeslaMateAPIHandleSuccessResponse(c, "TeslaMateAPICarsChargesV1", jsonData)
}
