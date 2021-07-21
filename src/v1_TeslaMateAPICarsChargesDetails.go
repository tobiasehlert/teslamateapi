package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// TeslaMateAPICarsChargesDetailsV1 func
func TeslaMateAPICarsChargesDetailsV1(c *gin.Context) {

	// getting CarID and ChargeID param from URL
	CarID := convertStringToInteger(c.Param("CarID"))
	ChargeID := convertStringToInteger(c.Param("ChargeID"))

	// creating structs for /cars/<CarID>/charges/<ChargeID>
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
	// ChargerDetails struct - child of ChargeDetails
	type ChargerDetails struct {
		ChargerActualCurrent int `json:"charger_actual_current"` // int
		ChargerPhases        int `json:"charger_phases"`         // int
		ChargerPilotCurrent  int `json:"charger_pilot_current"`  // int
		ChargerPower         int `json:"charger_power"`          // int
		ChargerVoltage       int `json:"charger_voltage"`        // int
	}
	// FastChargerInfo struct - child of ChargeDetails
	type FastChargerInfo struct {
		FastChargerPresent bool       `json:"fast_charger_present"` // bool
		FastChargerBrand   NullString `json:"fast_charger_brand"`   // string
		FastChargerType    string     `json:"fast_charger_type"`    // string

	}
	// BatteryInfo struct - child of ChargeDetails
	type BatteryInfo struct {
		IdealBatteryRange    float64  `json:"ideal_battery_range"`     // float64
		RatedBatteryRange    float64  `json:"rated_battery_range"`     // float64
		BatteryHeater        bool     `json:"battery_heater"`          // bool
		BatteryHeaterOn      bool     `json:"battery_heater_on"`       // bool
		BatteryHeaterNoPower NullBool `json:"battery_heater_no_power"` // bool
	}
	// ChargeDetails struct - child of Charge
	type ChargeDetails struct {
		DetailID             int             `json:"detail_id"`                // integer
		Date                 string          `json:"date"`                     // string
		BatteryLevel         int             `json:"battery_level"`            // int
		UsableBatteryLevel   int             `json:"usable_battery_level"`     // int
		ChargeEnergyAdded    float64         `json:"charge_energy_added"`      // float64
		NotEnoughPowerToHeat NullBool        `json:"not_enough_power_to_heat"` // bool
		ChargerDetails       ChargerDetails  `json:"charger_details"`          // struct
		BatteryInfo          BatteryInfo     `json:"battery_info"`             // struct
		ConnChargeCable      string          `json:"conn_charge_cable"`        // string
		FastChargerInfo      FastChargerInfo `json:"fast_charger_info"`        // struct
		OutsideTemp          float64         `json:"outside_temp"`             // float64
	}
	// Charge struct - child of Data
	type Charge struct {
		ChargeID          int             `json:"charge_id"`           // int
		StartDate         string          `json:"start_date"`          // string
		EndDate           string          `json:"end_date"`            // string
		Address           string          `json:"address"`             // string
		ChargeEnergyAdded float64         `json:"charge_energy_added"` // float64
		ChargeEnergyUsed  float64         `json:"charge_energy_used"`  // float64
		Cost              float64         `json:"cost"`                // float64
		DurationMin       int             `json:"duration_min"`        // int
		DurationStr       string          `json:"duration_str"`        // string
		BatteryDetails    BatteryDetails  `json:"battery_details"`     // BatteryDetails
		RangeIdeal        PreferredRange  `json:"range_ideal"`         // PreferredRange
		RangeRated        PreferredRange  `json:"range_rated"`         // PreferredRange
		OutsideTempAvg    float64         `json:"outside_temp_avg"`    // float64
		ChargeDetails     []ChargeDetails `json:"charge_details"`      // struct
	}
	// Incomplete charge struct - child of Data
	type IncompleteCharge struct {
		ChargeID      int             `json:"charge_id"`      // int
		StartDate     string          `json:"start_date"`     // string
		Address       string          `json:"address"`        // string
		ChargeDetails []ChargeDetails `json:"charge_details"` // struct
	}
	// TeslaMateUnits struct - child of Data
	type TeslaMateUnits struct {
		UnitsLength      string `json:"unit_of_length"`      // string
		UnitsTemperature string `json:"unit_of_temperature"` // string
	}
	// Data struct - child of JSONData
	type Data struct {
		Car               Car               `json:"car"`
		IsChargeCompleted bool              `json:"is_charge_completed"`
		IncompleteCharge  *IncompleteCharge `json:"incomplete_charge,omitempty"`
		Charge            *Charge           `json:"charge,omitempty"`
		TeslaMateUnits    TeslaMateUnits    `json:"units"`
	}
	// JSONData struct - main
	type JSONData struct {
		Data Data `json:"data"`
	}

	// creating required vars
	var ChargeData Charge
	var IncompleteChargeData IncompleteCharge
	var ChargeDetailsData []ChargeDetails
	var IsChargeCompleted bool = true
	var UnitsLength, UnitsTemperature, CarName string
	var ValidResponse bool // default is false

	// getting data from database
	query := `
		SELECT
			charging_processes.id AS charge_id,
			start_date,
			end_date,
			COALESCE(geofence.name, CONCAT_WS(', ', COALESCE(address.name, nullif(CONCAT_WS(' ', address.road, address.house_number), '')), address.city)) AS address,
			COALESCE(charging_processes.charge_energy_added, 0) AS charge_energy_added,
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
		LEFT JOIN charges ON charging_processes.id = charges.id
		WHERE charging_processes.car_id=$1 AND charging_processes.id=$2
		ORDER BY start_date DESC;`
	rows, err := db.Query(query, CarID, ChargeID)

	// checking for errors in query
	if err != nil {
		log.Fatal(err)
	}

	// defer closing rows
	defer rows.Close()

	// looping through all results
	for rows.Next() {

		// creating charge object based on struct
		charge := Charge{}

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
				new(int),          // ChargeID
				new(string),       // StartDate
				new(sql.RawBytes), // EndDate
				new(string),       // Address
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

			var err2 error = rows.Scan(dest...)

			if dest[0] != nil {
				IsChargeCompleted = false

				if err2 != nil {
					log.Fatal(err2)
				}
				// There is a charge available, it's just incomplete
				incompleteCharge := IncompleteCharge{}
				incompleteCharge.ChargeID = *(dest)[0].(*int)
				incompleteCharge.StartDate = *(dest)[1].(*string)
				incompleteCharge.Address = *(dest)[3].(*string)

				// adjusting to timezone differences from UTC to be userspecific
				incompleteCharge.StartDate = getTimeInTimeZone(incompleteCharge.StartDate)

				// appending charge to ChargesData
				IncompleteChargeData = incompleteCharge

				ValidResponse = true
			} else {
				log.Fatal(err)
			}
		} else {
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

			// appending charge to ChargeData
			ChargeData = charge

			ValidResponse = true
		}

		// getting detailed charge data from database
		query = `
 			SELECT
				id AS detail_id,
				date,
				battery_level,
				usable_battery_level,
				charge_energy_added,
				not_enough_power_to_heat,
				COALESCE(charger_actual_current, 0) as charger_actual_current,
				COALESCE(charger_phases, 0) AS charger_phases,
				COALESCE(charger_pilot_current, 0) as charger_pilot_current,
				COALESCE(charger_power, 0) as charger_power,
				COALESCE(charger_voltage, 0) as charger_voltage,
				ideal_battery_range_km AS ideal_battery_range,
				rated_battery_range_km AS rated_battery_range,
				battery_heater,
				battery_heater_on,
				battery_heater_no_power,
				conn_charge_cable,
				fast_charger_present,
				fast_charger_brand,
				fast_charger_type,
				outside_temp
			FROM charges
			WHERE charging_process_id=$1
			ORDER BY id ASC;`
		rows, err = db.Query(query, ChargeID)

		// checking for errors in query
		if err != nil {
			log.Fatal(err)
		}

		// defer closing rows
		defer rows.Close()

		// looping through all results
		for rows.Next() {

			// creating chargedetails object based on struct
			chargedetails := ChargeDetails{}

			// scanning row and putting values into the drive
			err = rows.Scan(
				&chargedetails.DetailID,
				&chargedetails.Date,
				&chargedetails.BatteryLevel,
				&chargedetails.UsableBatteryLevel,
				&chargedetails.ChargeEnergyAdded,
				&chargedetails.NotEnoughPowerToHeat,
				&chargedetails.ChargerDetails.ChargerActualCurrent,
				&chargedetails.ChargerDetails.ChargerPhases,
				&chargedetails.ChargerDetails.ChargerPilotCurrent,
				&chargedetails.ChargerDetails.ChargerPower,
				&chargedetails.ChargerDetails.ChargerVoltage,
				&chargedetails.BatteryInfo.IdealBatteryRange,
				&chargedetails.BatteryInfo.RatedBatteryRange,
				&chargedetails.BatteryInfo.BatteryHeater,
				&chargedetails.BatteryInfo.BatteryHeaterOn,
				&chargedetails.BatteryInfo.BatteryHeaterNoPower,
				&chargedetails.ConnChargeCable,
				&chargedetails.FastChargerInfo.FastChargerPresent,
				&chargedetails.FastChargerInfo.FastChargerBrand,
				&chargedetails.FastChargerInfo.FastChargerType,
				&chargedetails.OutsideTemp,
			)

			// converting values based of settings UnitsLength
			if UnitsLength == "mi" {
				chargedetails.BatteryInfo.IdealBatteryRange = kilometersToMiles(chargedetails.BatteryInfo.IdealBatteryRange)
				chargedetails.BatteryInfo.RatedBatteryRange = kilometersToMiles(chargedetails.BatteryInfo.RatedBatteryRange)

			}
			// converting values based of settings UnitsTemperature
			if UnitsTemperature == "F" {
				chargedetails.OutsideTemp = celsiusToFahrenheit(chargedetails.OutsideTemp)
			}
			// adjusting to timezone differences from UTC to be userspecific
			chargedetails.Date = getTimeInTimeZone(chargedetails.Date)

			// checking for errors after scanning
			if err != nil {
				log.Fatal(err)
			}

			// appending drive to ChargeData
			ChargeDetailsData = append(ChargeDetailsData, chargedetails)
		}

		if IsChargeCompleted {
			ChargeData.ChargeDetails = ChargeDetailsData
		} else {
			IncompleteChargeData.ChargeDetails = ChargeDetailsData
		}
	}

	// checking for errors in the rows result
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	//
	// build the data-blob
	jsonData := JSONData{
		Data{
			Car: Car{
				CarID:   CarID,
				CarName: CarName,
			},
			IsChargeCompleted: IsChargeCompleted,
			TeslaMateUnits: TeslaMateUnits{
				UnitsLength:      UnitsLength,
				UnitsTemperature: UnitsTemperature,
			},
		},
	}

	if IsChargeCompleted {
		jsonData.Data.Charge = &ChargeData
	} else {
		jsonData.Data.IncompleteCharge = &IncompleteChargeData
	}

	// print to log about request
	if gin.IsDebugging() {
		log.Println("[debug] TeslaMateAPICarsChargesDetailsV1 " + c.Request.RequestURI + " returned data:")
		js, _ := json.Marshal(jsonData)
		log.Printf("[debug] %s\n", js)
	}

	// return jsonData
	if ValidResponse {
		log.Println("[info] TeslaMateAPICarsChargesDetailsV1 " + c.Request.RequestURI + " executed successful.")
		c.JSON(http.StatusOK, jsonData)
	} else {
		log.Println("[error] TeslaMateAPICarsChargesDetailsV1 " + c.Request.RequestURI + " error in execution!")
		c.JSON(http.StatusNotFound, gin.H{"error": "something went wrong in TeslaMateAPICarsChargesDetailsV1.."})
	}
}
