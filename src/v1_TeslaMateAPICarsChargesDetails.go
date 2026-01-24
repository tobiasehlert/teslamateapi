package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// TeslaMateAPICarsChargesDetailsV1 func
func TeslaMateAPICarsChargesDetailsV1(c *gin.Context) {

	// define error messages
	var (
		CarsChargesDetailsError1 = "Unable to load charge."
		CarsChargesDetailsError2 = "Unable to load charge details."
	)

	// getting CarID and ChargeID param from URL
	CarID := convertStringToInteger(c.Param("CarID"))
	ChargeID := convertStringToInteger(c.Param("ChargeID"))

	// creating structs for /cars/<CarID>/charges/<ChargeID>
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
		Odometer          float64         `json:"odometer"`            // float64
		Latitude          float64         `json:"latitude"`            // float64
		Longitude         float64         `json:"longitude"`           // float64
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
	)

	// getting data from database
	query := `
		SELECT
			charging_processes.id AS charge_id,
			charging_processes.start_date,
			charging_processes.end_date,
			COALESCE(geofence.name, CONCAT_WS(', ', COALESCE(address.name, nullif(CONCAT_WS(' ', address.road, address.house_number), '')), address.city)) AS address,
			COALESCE(charging_processes.charge_energy_added, 0) AS charge_energy_added,
			COALESCE(GREATEST(charge_energy_used, charging_processes.charge_energy_added), 0) AS charge_energy_used,
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
			position.latitude,
			position.longitude,
			(SELECT unit_of_length FROM settings LIMIT 1) as unit_of_length,
			(SELECT unit_of_temperature FROM settings LIMIT 1) as unit_of_temperature,
			cars.name
		FROM charging_processes
		LEFT JOIN cars ON car_id = cars.id
		LEFT JOIN addresses address ON address_id = address.id
		LEFT JOIN positions position ON position_id = position.id
		LEFT JOIN geofences geofence ON geofence_id = geofence.id
		LEFT JOIN charges ON charging_processes.id = charges.id
		WHERE charging_processes.car_id=$1 AND charging_processes.id=$2 AND charging_processes.end_date IS NOT NULL
		ORDER BY start_date DESC;`
	row := db.QueryRow(query, CarID, ChargeID)

	// scanning row and putting values into the charge
	err := row.Scan(
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
		&charge.Latitude,
		&charge.Longitude,
		&UnitsLength,
		&UnitsTemperature,
		&CarName,
	)

	switch err {
	case sql.ErrNoRows:
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesDetailsV1", "No rows were returned!", err.Error())
		return
	case nil:
		// nothing wrong.. continuing
		break
	default:
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesDetailsV1", CarsChargesDetailsError1, err.Error())
		return
	}

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
	rows, err := db.Query(query, ChargeID)

	// checking for errors in query
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesDetailsV1", CarsChargesDetailsError2, err.Error())
		return
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
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesDetailsV1", CarsChargesDetailsError2, err.Error())
			return
		}

		// appending drive to ChargeData
		ChargeDetailsData = append(ChargeDetailsData, chargedetails)
		charge.ChargeDetails = ChargeDetailsData
	}

	// checking for errors in the rows result
	err = rows.Err()
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsChargesDetailsV1", CarsChargesDetailsError2, err.Error())
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
			Charge: charge,
			TeslaMateUnits: TeslaMateUnits{
				UnitsLength:      UnitsLength,
				UnitsTemperature: UnitsTemperature,
			},
		},
	}

	// return jsonData
	TeslaMateAPIHandleSuccessResponse(c, "TeslaMateAPICarsChargesDetailsV1", jsonData)
}
