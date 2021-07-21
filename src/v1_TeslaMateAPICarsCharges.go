package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// TeslaMateAPICarsChargesV1 func
func TeslaMateAPICarsChargesV1(c *gin.Context) {

	// getting CarID param from URL
	CarID := convertStringToInteger(c.Param("CarID"))
	// query options to modify query when collecting data
	ResultPage := convertStringToInteger(c.DefaultQuery("page", "1"))
	ResultShow := convertStringToInteger(c.DefaultQuery("show", "100"))

	// creating structs for /cars/<CarID>/charges
	// Car struct - child of Data
	type Car struct {
		CarID   int    `json:"car_id"`   // smallint
		CarName string `json:"car_name"` // text
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
	}
	// Incomplete Charges struct - child of Data
	type IncompleteCharges struct {
		ChargeID          int            `json:"charge_id"`           // int
		StartDate         string         `json:"start_date"`          // string
		Address           string         `json:"address"`             // string
	}
	// TeslaMateUnits struct - child of Data
	type TeslaMateUnits struct {
		UnitsLength      string `json:"unit_of_length"`      // string
		UnitsTemperature string `json:"unit_of_temperature"` // string
	}
	// Data struct - child of JSONData
	type Data struct {
		Car                Car                    `json:"car"`
		Charges            []Charges              `json:"charges"`
		IncompleteCharges  []IncompleteCharges    `json:"incomplete_charges"`
		TeslaMateUnits     TeslaMateUnits         `json:"units"`
	}
	// JSONData struct - main
	type JSONData struct {
		Data Data `json:"data"`
	}

	// creating required vars
	var ChargesData []Charges
	var IncompleteChargesData []IncompleteCharges
	var UnitsLength, UnitsTemperature, CarName string
	var ValidResponse bool // default is false

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
			(SELECT unit_of_length FROM settings LIMIT 1) as unit_of_length,
			(SELECT unit_of_temperature FROM settings LIMIT 1) as unit_of_temperature,
			cars.name
		FROM charging_processes
		LEFT JOIN cars ON car_id = cars.id
		LEFT JOIN addresses address ON address_id = address.id
		LEFT JOIN positions position ON position_id = position.id
		LEFT JOIN geofences geofence ON geofence_id = geofence.id
		WHERE charging_processes.car_id=$1
		ORDER BY start_date DESC
		LIMIT $2 OFFSET $3;`
	rows, err := db.Query(query, CarID, ResultShow, ResultPage)

	// checking for errors in query
	if err != nil {
		log.Fatal(err)
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
			&UnitsLength,
			&UnitsTemperature,
			&CarName,
		)
		
		// checking for errors after scanning
		// Incomplete charges may still be in progress, or a previous charge that failed to finalize
		if err != nil {
			// Check if charge is available with incomplete data
			dest := []interface{}{ // Standard MySQL columns
				new(int), // ChargeID
				new(string), // StartDate
				new(sql.RawBytes), // EndDate
				new(string), // Address
				new(sql.RawBytes), // ChargeEnergyAdded
				new(sql.RawBytes), // ChargeEnergyUsed
				new(sql.RawBytes), // Cost
				new(sql.RawBytes), // StartRange
				new(sql.RawBytes), // EndRange
				new(sql.RawBytes), // StartRange
				new(sql.RawBytes), // EndRange
				new(sql.RawBytes), // StartBatteryLevel
				new(sql.RawBytes), // EndBatteryLevel
				new(sql.RawBytes), // DurationMin
				new(sql.RawBytes), // DurationStr
				new(sql.RawBytes), // OutsideTempAvg
				new(sql.RawBytes), // UnitsLength
				new(sql.RawBytes), // UnitsTemperature
				new(sql.RawBytes), // CarName
			}
			err2 = rows.Scan(dest...)
			
			if dest[0] != nil {
				if err2 != nil {
					log.Fatal(err2)
				}
				// There is a charge available, it's just incomplete
				incompleteCharge := IncompleteCharges{
					ChargeID: dest[0], 
					StartDate: dest[1],
					Address: dest[3],
				}
			
				// adjusting to timezone differences from UTC to be userspecific
				incompleteCharge.StartDate = getTimeInTimeZone(incompleteCharge.StartDate)

				// appending charge to ChargesData
				IncompleteChargesData = append(IncompleteChargesData, incompleteCharge)

				continue
				
			} else {
				log.Fatal(err)
			}
		}

		// converting values based of settings UnitsLength
		if UnitsLength == "mi" {
			charge.RangeIdeal.StartRange = kilometersToMiles(charge.RangeIdeal.StartRange)
			charge.RangeIdeal.EndRange = kilometersToMiles(charge.RangeIdeal.EndRange)
			charge.RangeRated.StartRange = kilometersToMiles(charge.RangeRated.StartRange)
			charge.RangeRated.EndRange = kilometersToMiles(charge.RangeRated.EndRange)
		}
		// converting values based of settings UnitsTemperature
		if UnitsTemperature == "F" {
			charge.OutsideTempAvg = celsiusToFahrenheit(charge.OutsideTempAvg)
		}

		// adjusting to timezone differences from UTC to be userspecific
		charge.StartDate = getTimeInTimeZone(charge.StartDate)
		charge.EndDate = getTimeInTimeZone(charge.EndDate)

		// appending charge to ChargesData
		ChargesData = append(ChargesData, charge)
		ValidResponse = true
	}

	// checking for errors in the rows result
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	
	// if no errors, but ChargeData is empty, return a valid response with an empty set
	if len(ChargesData) == 0 {
		ValidResponse = true
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
			IncompleteCharges: IncompleteChargesData,
			TeslaMateUnits: TeslaMateUnits{
				UnitsLength:      UnitsLength,
				UnitsTemperature: UnitsTemperature,
			},
		},
	}

	// print to log about request
	if gin.IsDebugging() {
		log.Println("[debug] TeslaMateAPICarsChargesV1 " + c.Request.RequestURI + " returned data:")
		js, _ := json.Marshal(jsonData)
		log.Printf("[debug] %s\n", js)
	}

	// return jsonData
	if ValidResponse {
		log.Println("[info] TeslaMateAPICarsChargesV1 " + c.Request.RequestURI + " executed successful.")
		c.JSON(http.StatusOK, jsonData)
	} else {
		log.Println("[error] TeslaMateAPICarsChargesV1 " + c.Request.RequestURI + " error in execution!")
		c.JSON(http.StatusNotFound, gin.H{"error": "something went wrong in TeslaMateAPICarsChargesV1.."})
	}
}
