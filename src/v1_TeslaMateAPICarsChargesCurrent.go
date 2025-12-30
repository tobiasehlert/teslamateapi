package main

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	maxChargeInactivityThresholdMinutes = 15
)

// TeslaMateAPICarsChargesCurrentV1 func
func TeslaMateAPICarsChargesCurrentV1(c *gin.Context) {

	// define error messages
	var (
		CarsChargesCurrentError1 = "Unable to load current charge."
		CarsChargesCurrentError2 = "Unable to load current charge details."
		CarsChargesCurrentError3 = "No active charging in progress."
	)

	// getting CarID param from URL
	CarID := convertStringToInteger(c.Param("CarID"))

	// creating structs for /cars/<CarID>/charges/current
	// Car struct - child of Data
	type Car struct {
		CarID   int        `json:"car_id"`   // smallint
		CarName NullString `json:"car_name"` // text (nullable)
	}
	// BatteryDetails struct - child of Charge
	type BatteryDetails struct {
		StartBatteryLevel   int `json:"start_battery_level"`   // int
		CurrentBatteryLevel int `json:"current_battery_level"` // int
	}
	// PreferredRange struct - child of Charge
	type PreferredRange struct {
		StartRange   float64 `json:"start_range"`   // float64
		CurrentRange float64 `json:"current_range"` // float64
		AddedRange   float64 `json:"added_range"`   // float64
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
		FastChargerPresent bool    `json:"fast_charger_present"`         // bool
		FastChargerBrand   *string `json:"fast_charger_brand,omitempty"` // string or null
		FastChargerType    *string `json:"fast_charger_type,omitempty"`  // string or null
	}
	// BatteryInfo struct - child of ChargeDetails
	type BatteryInfo struct {
		RatedBatteryRange    float64  `json:"rated_battery_range"`     // float64
		BatteryHeater        bool     `json:"battery_heater"`          // bool
		BatteryHeaterOn      bool     `json:"battery_heater_on"`       // bool
		BatteryHeaterNoPower NullBool `json:"battery_heater_no_power"` // bool
	}
	// ChargeDetails struct - child of Charge
	type ChargeDetails struct {
		DetailID             int             `json:"detail_id"`                   // integer
		Date                 string          `json:"date"`                        // string
		BatteryLevel         int             `json:"battery_level"`               // int
		UsableBatteryLevel   int             `json:"usable_battery_level"`        // int
		ChargeEnergyAdded    float64         `json:"charge_energy_added"`         // float64
		NotEnoughPowerToHeat NullBool        `json:"not_enough_power_to_heat"`    // bool
		ChargerDetails       ChargerDetails  `json:"charger_details"`             // struct
		BatteryInfo          BatteryInfo     `json:"battery_info"`                // struct
		ConnChargeCable      interface{}     `json:"conn_charge_cable,omitempty"` // string or null
		FastChargerInfo      FastChargerInfo `json:"fast_charger_info"`           // struct
		OutsideTemp          float64         `json:"outside_temp"`                // float64
	}
	// Charge struct - child of Data
	type Charge struct {
		ChargeID          int             `json:"charge_id"`           // int
		StartDate         string          `json:"start_date"`          // string
		IsCharging        bool            `json:"is_charging"`         // bool
		Address           string          `json:"address"`             // string
		ChargeEnergyAdded float64         `json:"charge_energy_added"` // float64
		Cost              float64         `json:"cost"`                // float64
		DurationMin       int             `json:"duration_min"`        // int
		DurationStr       string          `json:"duration_str"`        // string
		BatteryDetails    BatteryDetails  `json:"battery_details"`     // BatteryDetails
		RatedRange        PreferredRange  `json:"rated_range"`         // PreferredRange
		OutsideTempAvg    float64         `json:"outside_temp_avg"`    // float64
		Odometer          float64         `json:"odometer"`            // float64
		ChargeDetails     []ChargeDetails `json:"charge_details"`      // struct
	}
	// TeslaMateUnits struct - child of Data
	type TeslaMateUnits struct {
		UnitsLength      string `json:"unit_of_length"`      // string
		UnitsTemperature string `json:"unit_of_temperature"` // string
	}
	// Data struct - child of JSONData
	type Data struct {
		Car            Car            `json:"car"`
		Charge         Charge         `json:"charge"`
		TeslaMateUnits TeslaMateUnits `json:"units"`
	}
	// JSONData struct - main
	type JSONData struct {
		Data Data `json:"data"`
	}

	// creating required vars
	var (
		CarName                       NullString
		charge                        Charge
		ChargeDetailsData             []ChargeDetails
		UnitsLength, UnitsTemperature string
		isCharging                    bool
	)

	// Create temp vars to handle NULL values in the database
	var (
		startRatedRange, currentRatedRange     sql.NullFloat64
		startBatteryLevel, currentBatteryLevel sql.NullInt64
		chargeEnergyAdded, cost                sql.NullFloat64
		outsideTempAvg                         sql.NullFloat64
		odometer                               sql.NullFloat64
		durationMin                            sql.NullFloat64
		durationStr, address                   sql.NullString
	)

	// Construct the query with the preferred range setting
	query := `
		SELECT
			charging_processes.id AS charge_id,
			start_date,
			COALESCE(geofence.name, CONCAT_WS(', ', COALESCE(address.name, nullif(CONCAT_WS(' ', address.road, address.house_number), '')), address.city)) AS address,
			(SELECT charge_energy_added FROM charges WHERE charging_process_id = charging_processes.id ORDER BY id DESC LIMIT 1) AS charge_energy_added,
			COALESCE(cost, 0) AS cost,
	        (SELECT rated_battery_range_km FROM charges WHERE charging_process_id = charging_processes.id ORDER BY id ASC LIMIT 1) AS start_rated_range,
			(SELECT rated_battery_range_km FROM charges WHERE charging_process_id = charging_processes.id ORDER BY id DESC LIMIT 1) AS current_rated_range,
			(SELECT battery_level FROM charges WHERE charging_process_id = charging_processes.id ORDER BY date ASC LIMIT 1) AS start_battery_level,
			(SELECT battery_level FROM charges WHERE charging_process_id = charging_processes.id ORDER BY id DESC LIMIT 1) AS current_battery_level,
			EXTRACT(EPOCH FROM (COALESCE(end_date, NOW()) - start_date))/60 AS duration_min,
			TO_CHAR((EXTRACT(EPOCH FROM (COALESCE(end_date, NOW()) - start_date))/60 * INTERVAL '1 minute'), 'HH24:MI') as duration_str,
			(SELECT outside_temp FROM charges WHERE charging_process_id = charging_processes.id ORDER BY id DESC LIMIT 1) AS outside_temp_avg,
			position.odometer as odometer,
			(SELECT unit_of_length FROM settings LIMIT 1) as unit_of_length,
			(SELECT unit_of_temperature FROM settings LIMIT 1) as unit_of_temperature,
			cars.name,
			end_date IS NULL AS is_charging
		FROM charging_processes
		LEFT JOIN cars ON car_id = cars.id
		LEFT JOIN addresses address ON address_id = address.id
		LEFT JOIN positions position ON position_id = position.id
		LEFT JOIN geofences geofence ON geofence_id = geofence.id
		WHERE charging_processes.car_id=$1
		ORDER BY end_date IS NULL DESC, start_date DESC
		LIMIT 1;`

	row := db.QueryRow(query, CarID)

	// Scanning row and putting values into the temp vars to handle NULLs
	err := row.Scan(
		&charge.ChargeID,
		&charge.StartDate,
		&address,
		&chargeEnergyAdded,
		&cost,
		&startRatedRange,
		&currentRatedRange,
		&startBatteryLevel,
		&currentBatteryLevel,
		&durationMin,
		&durationStr,
		&outsideTempAvg,
		&odometer,
		&UnitsLength,
		&UnitsTemperature,
		&CarName,
		&isCharging,
	)

	switch err {
	case sql.ErrNoRows:
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesCurrentV1", "No current charge found.", "No rows were returned")
		return
	case nil:
		// nothing wrong.. continuing
		break
	default:
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesCurrentV1", CarsChargesCurrentError1, err.Error())
		return
	}

	charge.IsCharging = isCharging
	if address.Valid {
		charge.Address = address.String
	} else {
		charge.Address = "Unknown"
	}

	if chargeEnergyAdded.Valid {
		charge.ChargeEnergyAdded = chargeEnergyAdded.Float64
	}

	if cost.Valid {
		charge.Cost = cost.Float64
	}

	if startRatedRange.Valid {
		charge.RatedRange.StartRange = startRatedRange.Float64
	}

	if currentRatedRange.Valid {
		charge.RatedRange.CurrentRange = currentRatedRange.Float64
	}

	if addedRange := charge.RatedRange.CurrentRange - charge.RatedRange.StartRange; addedRange > 0 {
		charge.RatedRange.AddedRange = addedRange
	}

	if startBatteryLevel.Valid {
		charge.BatteryDetails.StartBatteryLevel = int(startBatteryLevel.Int64)
	}

	if currentBatteryLevel.Valid {
		charge.BatteryDetails.CurrentBatteryLevel = int(currentBatteryLevel.Int64)
	}

	if durationMin.Valid {
		charge.DurationMin = int(durationMin.Float64) // Convert float64 to int
	}

	if durationStr.Valid {
		charge.DurationStr = durationStr.String
	}

	if outsideTempAvg.Valid {
		charge.OutsideTempAvg = outsideTempAvg.Float64
	}

	if odometer.Valid {
		charge.Odometer = odometer.Float64
	}

	// Converting values based on settings UnitsLength
	if UnitsLength == "mi" {
		charge.RatedRange.StartRange = kilometersToMiles(charge.RatedRange.StartRange)
		charge.RatedRange.CurrentRange = kilometersToMiles(charge.RatedRange.CurrentRange)
		charge.RatedRange.AddedRange = kilometersToMiles(charge.RatedRange.AddedRange)
		charge.Odometer = kilometersToMiles(charge.Odometer)
	}
	// Converting values based on settings UnitsTemperature
	if UnitsTemperature == "F" && outsideTempAvg.Valid {
		charge.OutsideTempAvg = celsiusToFahrenheit(charge.OutsideTempAvg)
	}

	// Adjusting to timezone differences from UTC to be user-specific
	charge.StartDate = getTimeInTimeZone(charge.StartDate)

	// Getting detailed charge data from database
	detailsQuery := `
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
		ORDER BY id DESC;`
	rows, err := db.Query(detailsQuery, charge.ChargeID)

	// Checking for errors in query
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesCurrentV1", CarsChargesCurrentError2, err.Error())
		return
	}

	// Defer closing rows
	defer rows.Close()

	// Looping through all results
	for rows.Next() {
		// Create temp variables to handle NULL values
		var (
			detailBatteryLevel, detailUsableBatteryLevel                        sql.NullInt64
			detailChargeEnergyAdded, detailRatedBatteryRange, detailOutsideTemp sql.NullFloat64
			detailConnChargeCable, detailFastChargerType                        sql.NullString
			detailFastChargerBrand                                              sql.NullString
		)

		// Creating chargedetails object based on struct
		chargedetails := ChargeDetails{}

		// Scanning row and putting values into temporary variables
		err = rows.Scan(
			&chargedetails.DetailID,
			&chargedetails.Date,
			&detailBatteryLevel,
			&detailUsableBatteryLevel,
			&detailChargeEnergyAdded,
			&chargedetails.NotEnoughPowerToHeat,
			&chargedetails.ChargerDetails.ChargerActualCurrent,
			&chargedetails.ChargerDetails.ChargerPhases,
			&chargedetails.ChargerDetails.ChargerPilotCurrent,
			&chargedetails.ChargerDetails.ChargerPower,
			&chargedetails.ChargerDetails.ChargerVoltage,
			&detailRatedBatteryRange,
			&chargedetails.BatteryInfo.BatteryHeater,
			&chargedetails.BatteryInfo.BatteryHeaterOn,
			&chargedetails.BatteryInfo.BatteryHeaterNoPower,
			&detailConnChargeCable,
			&chargedetails.FastChargerInfo.FastChargerPresent,
			&detailFastChargerBrand,
			&detailFastChargerType,
			&detailOutsideTemp,
		)

		// Handle NULL values
		if detailBatteryLevel.Valid {
			chargedetails.BatteryLevel = int(detailBatteryLevel.Int64)
		}

		if detailUsableBatteryLevel.Valid {
			chargedetails.UsableBatteryLevel = int(detailUsableBatteryLevel.Int64)
		}

		if detailChargeEnergyAdded.Valid {
			chargedetails.ChargeEnergyAdded = detailChargeEnergyAdded.Float64
		}

		if detailRatedBatteryRange.Valid {
			chargedetails.BatteryInfo.RatedBatteryRange = detailRatedBatteryRange.Float64
		}

		// Properly handle NULL values for string fields by using interface{}
		if detailConnChargeCable.Valid {
			chargedetails.ConnChargeCable = detailConnChargeCable.String
		} else {
			chargedetails.ConnChargeCable = nil
		}

		// Fix for fast_charger_brand and fast_charger_type when "<invalid>"
		if detailFastChargerBrand.Valid && detailFastChargerBrand.String != "<invalid>" {
			chargedetails.FastChargerInfo.FastChargerBrand = &detailFastChargerBrand.String
		} else {
			chargedetails.FastChargerInfo.FastChargerBrand = nil
		}

		if detailFastChargerType.Valid && detailFastChargerType.String != "<invalid>" {
			chargedetails.FastChargerInfo.FastChargerType = &detailFastChargerType.String
		} else {
			chargedetails.FastChargerInfo.FastChargerType = nil
		}

		if detailOutsideTemp.Valid {
			chargedetails.OutsideTemp = detailOutsideTemp.Float64
		}

		// Converting values based on settings UnitsLength
		if UnitsLength == "mi" && detailRatedBatteryRange.Valid {
			chargedetails.BatteryInfo.RatedBatteryRange = kilometersToMiles(chargedetails.BatteryInfo.RatedBatteryRange)
		}

		// Converting values based on settings UnitsTemperature
		if UnitsTemperature == "F" && detailOutsideTemp.Valid {
			chargedetails.OutsideTemp = celsiusToFahrenheit(chargedetails.OutsideTemp)
		}

		// Adjusting to timezone differences from UTC to be user-specific
		chargedetails.Date = getTimeInTimeZone(chargedetails.Date)

		chargedetails.ChargerDetails.ChargerPhases = normalizeChargerPhases(chargedetails.ChargerDetails.ChargerPhases)

		// Checking for errors after scanning
		if err != nil {
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesCurrentV1", CarsChargesCurrentError2, err.Error())
			return
		}

		// Appending chargedetails to ChargeDetailsData
		ChargeDetailsData = append(ChargeDetailsData, chargedetails)
	}

	// Checking for errors in the rows result
	err = rows.Err()
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesCurrentV1", CarsChargesCurrentError2, err.Error())
		return
	}

	// Check if the charge details array contains any entries
	if len(ChargeDetailsData) > 0 {
		// Parse the date of the most recent charge detail
		latestDetailDate, err := time.Parse(time.RFC3339, ChargeDetailsData[0].Date)
		if err != nil {
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesCurrentV1", CarsChargesCurrentError2, "Error parsing charge detail date")
			return
		}

		// Calculate time elapsed since the most recent detail
		timeElapsed := time.Since(latestDetailDate)

		// If the most recent detail is more than 15 minutes old, consider it incomplete/not current
		if timeElapsed.Minutes() > maxChargeInactivityThresholdMinutes {
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesCurrentV1", CarsChargesCurrentError3, "No active charging in progress. There are incomplete charges but last update was more than 15 minutes ago.")
			return
		}
	}

	// Set the ChargeDetails in the charge
	charge.ChargeDetails = ChargeDetailsData

	if charge.RatedRange.StartRange == 0 && len(ChargeDetailsData) > 0 {
		charge.RatedRange.StartRange = ChargeDetailsData[len(ChargeDetailsData)-1].BatteryInfo.RatedBatteryRange
	}

	if addedRange := charge.RatedRange.CurrentRange - charge.RatedRange.StartRange; addedRange > 0 {
		charge.RatedRange.AddedRange = addedRange
	} else {
		charge.RatedRange.AddedRange = 0
	}

	// Build the data-blob
	jsonData := JSONData{
		Data{
			Car: Car{
				CarID:   CarID,
				CarName: CarName,
			},
			Charge: charge,
			TeslaMateUnits: TeslaMateUnits{
				UnitsLength:      UnitsLength,
				UnitsTemperature: UnitsTemperature,
			},
		},
	}

	// Return jsonData
	TeslaMateAPIHandleSuccessResponse(c, "TeslaMateAPICarsChargesCurrentV1", jsonData)
}

// normalizeChargerPhases converts phase values to valid configurations.
// Both phase 2 and 3 are normalized to 3.
// All other values default to 1 (single-phase).
func normalizeChargerPhases(phases int) int {
	switch phases {
	case 2, 3:
		return 3
	default:
		return 1
	}
}
