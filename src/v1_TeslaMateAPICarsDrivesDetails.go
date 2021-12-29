package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// TeslaMateAPICarsDrivesDetailsV1 func
func TeslaMateAPICarsDrivesDetailsV1(c *gin.Context) {

	// getting CarID and DriveID param from URL
	CarID := convertStringToInteger(c.Param("CarID"))
	DriveID := convertStringToInteger(c.Param("DriveID"))

	// creating structs for /cars/<CarID>/drives/<DriveID>
	// Car struct - child of Data
	type Car struct {
		CarID   int    `json:"car_id"`   // smallint
		CarName string `json:"car_name"` // text
	}
	// OdometerDetails struct - child of Drives
	type OdometerDetails struct {
		OdometerStart    float64 `json:"odometer_start"`    // float64
		OdometerEnd      float64 `json:"odometer_end"`      // float64
		OdometerDistance float64 `json:"odometer_distance"` // float64
	}
	// BatteryDetails struct - child of Drives
	type BatteryDetails struct {
		StartUsableBatteryLevel int  `json:"start_usable_battery_level"` // int
		StartBatteryLevel       int  `json:"start_battery_level"`        // int
		EndUsableBatteryLevel   int  `json:"end_usable_battery_level"`   // int
		EndBatteryLevel         int  `json:"end_battery_level"`          // int
		ReducedRange            bool `json:"reduced_range"`              // bool
		IsSufficientlyPrecise   bool `json:"is_sufficiently_precise"`    // bool
	}
	// PreferredRange struct - child of Drives
	type PreferredRange struct {
		StartRange float64 `json:"start_range"` // float64
		EndRange   float64 `json:"end_range"`   // float64
		RangeDiff  float64 `json:"range_diff"`  // float64
	}
	// ClimateInfo struct - child of DriveDetails
	type ClimateInfo struct {
		InsideTemp           NullFloat64 `json:"inside_temp"`            // numeric(4,1)
		OutsideTemp          NullFloat64 `json:"outside_temp"`           // numeric(4,1)
		IsClimateOn          NullBool    `json:"is_climate_on"`          // boolean
		FanStatus            NullInt64   `json:"fan_status"`             // integer
		DriverTempSetting    NullFloat64 `json:"driver_temp_setting"`    // numeric(4,1)
		PassengerTempSetting NullFloat64 `json:"passenger_temp_setting"` // numeric(4,1)
		IsRearDefrosterOn    NullBool    `json:"is_rear_defroster_on"`   // boolean
		IsFrontDefrosterOn   NullBool    `json:"is_front_defroster_on"`  // boolean
	}
	// BatteryInfo struct - child of DriveDetails
	type BatteryInfo struct {
		EstBatteryRange      NullFloat64 `json:"est_battery_range"`       // numeric(6,2)
		IdealBatteryRange    NullFloat64 `json:"ideal_battery_range"`     // numeric(6,2)
		RatedBatteryRange    NullFloat64 `json:"rated_battery_range"`     // numeric(6,2)
		BatteryHeater        NullBool    `json:"battery_heater"`          // boolean
		BatteryHeaterOn      NullBool    `json:"battery_heater_on"`       // boolean
		BatteryHeaterNoPower NullBool    `json:"battery_heater_no_power"` // boolean
	}
	// DriveDetails struct - child of Drive
	type DriveDetails struct {
		DetailID           int         `json:"detail_id"`            // integer
		Date               string      `json:"date"`                 // timestamp without time zone
		Latitude           float64     `json:"latitude"`             // numeric(8,6)
		Longitude          float64     `json:"longitude"`            // numeric(9,6)
		Speed              int         `json:"speed"`                // smallint
		Power              int         `json:"power"`                // smallint
		Odometer           float64     `json:"odometer"`             // double precision
		BatteryLevel       int         `json:"battery_level"`        // smallint
		UsableBatteryLevel NullInt64   `json:"usable_battery_level"` // smallint
		Elevation          NullInt64   `json:"elevation"`            // smallint
		ClimateInfo        ClimateInfo `json:"climate_info"`         // struct
		BatteryInfo        BatteryInfo `json:"battery_info"`         // struct
	}
	// Drive struct - child of Data
	type Drive struct {
		DriveID         int             `json:"drive_id"`         // int
		StartDate       string          `json:"start_date"`       // string
		EndDate         string          `json:"end_date"`         // string
		StartAddress    string          `json:"start_address"`    // string
		EndAddress      string          `json:"end_address"`      // string
		OdometerDetails OdometerDetails `json:"odometer_details"` // OdometerDetails
		DurationMin     int             `json:"duration_min"`     // int
		DurationStr     string          `json:"duration_str"`     // string
		SpeedMax        int             `json:"speed_max"`        // int
		SpeedAvg        float64         `json:"speed_avg"`        // float64
		PowerMax        int             `json:"power_max"`        // int
		PowerMin        int             `json:"power_min"`        // int
		BatteryDetails  BatteryDetails  `json:"battery_details"`  // BatteryDetails
		RangeIdeal      PreferredRange  `json:"range_ideal"`      // PreferredRange
		RangeRated      PreferredRange  `json:"range_rated"`      // PreferredRange
		OutsideTempAvg  float64         `json:"outside_temp_avg"` // float64
		InsideTempAvg   float64         `json:"inside_temp_avg"`  // float64
		DriveDetails    []DriveDetails  `json:"drive_details"`    // struct
	}
	// TeslaMateUnits struct - child of Data
	type TeslaMateUnits struct {
		UnitsLength      string `json:"unit_of_length"`      // string
		UnitsTemperature string `json:"unit_of_temperature"` // string
	}
	// Data struct - child of JSONData
	type Data struct {
		Car            Car            `json:"car"`
		Drive          Drive          `json:"drive"`
		TeslaMateUnits TeslaMateUnits `json:"units"`
	}
	// JSONData struct - main
	type JSONData struct {
		Data Data `json:"data"`
	}

	// creating required vars
	var DriveData Drive
	var DriveDetailsData []DriveDetails
	var UnitsLength, UnitsTemperature, CarName string

	// getting data from database
	query := `
		SELECT
			drives.id AS drive_id,
			start_date,
			end_date,
			COALESCE(start_geofence.name, CONCAT_WS(', ', COALESCE(start_address.name, nullif(CONCAT_WS(' ', start_address.road, start_address.house_number), '')), start_address.city)) AS start_address,
			COALESCE(end_geofence.name, CONCAT_WS(', ', COALESCE(end_address.name, nullif(CONCAT_WS(' ', end_address.road, end_address.house_number), '')), end_address.city)) AS end_address,
			start_km,
			end_km,
			distance,
			duration_min,
			TO_CHAR((duration_min * INTERVAL '1 minute'), 'HH24:MI') as duration_str,
			speed_max,
			COALESCE(distance / NULLIF(duration_min, 0) * 60, 0) AS speed_avg,
			power_max,
			power_min,
			COALESCE(start_position.usable_battery_level, start_position.battery_level) as start_usable_battery_level,
			start_position.battery_level as start_battery_level,
			COALESCE(end_position.usable_battery_level, end_position.battery_level) as end_usable_battery_level,
			end_position.battery_level as end_battery_level,
			case when ( start_position.battery_level != start_position.usable_battery_level OR end_position.battery_level != end_position.usable_battery_level ) = true then true else false end  as reduced_range,
			duration_min > 1 AND distance > 1 AND ( start_position.usable_battery_level IS NULL OR end_position.usable_battery_level IS NULL OR ( end_position.battery_level - end_position.usable_battery_level ) = 0 ) as is_sufficiently_precise,
			start_ideal_range_km,
			end_ideal_range_km,
			COALESCE( NULLIF ( GREATEST ( start_ideal_range_km - end_ideal_range_km, 0 ), 0 ),0 ) as range_diff_ideal_km,
			start_rated_range_km,
			end_rated_range_km,
			COALESCE( NULLIF ( GREATEST ( start_rated_range_km - end_rated_range_km, 0 ), 0 ),0 ) as range_diff_rated_km,
			outside_temp_avg,
			inside_temp_avg,
			(SELECT unit_of_length FROM settings LIMIT 1) as unit_of_length,
			(SELECT unit_of_temperature FROM settings LIMIT 1) as unit_of_temperature,
			cars.name
		FROM drives
		LEFT JOIN cars ON car_id = cars.id
		LEFT JOIN addresses start_address ON start_address_id = start_address.id
		LEFT JOIN addresses end_address ON end_address_id = end_address.id
		LEFT JOIN positions start_position ON start_position_id = start_position.id
		LEFT JOIN positions end_position ON end_position_id = end_position.id
		LEFT JOIN geofences start_geofence ON start_geofence_id = start_geofence.id
		LEFT JOIN geofences end_geofence ON end_geofence_id = end_geofence.id
		WHERE drives.car_id=$1 AND end_date IS NOT NULL AND drives.id = $2;`
	rows, err := db.Query(query, CarID, DriveID)

	// defer closing rows
	defer rows.Close()

	// checking for errors in query
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsDrivesDetailsV1", "Unable to load drive.", err.Error())
		return
	}

	// looping through all results
	for rows.Next() {

		// creating drive object based on struct
		drive := Drive{}

		// scanning row and putting values into the drive
		err = rows.Scan(
			&drive.DriveID,
			&drive.StartDate,
			&drive.EndDate,
			&drive.StartAddress,
			&drive.EndAddress,
			&drive.OdometerDetails.OdometerStart,
			&drive.OdometerDetails.OdometerEnd,
			&drive.OdometerDetails.OdometerDistance,
			&drive.DurationMin,
			&drive.DurationStr,
			&drive.SpeedMax,
			&drive.SpeedAvg,
			&drive.PowerMax,
			&drive.PowerMin,
			&drive.BatteryDetails.StartUsableBatteryLevel,
			&drive.BatteryDetails.StartBatteryLevel,
			&drive.BatteryDetails.EndUsableBatteryLevel,
			&drive.BatteryDetails.EndBatteryLevel,
			&drive.BatteryDetails.ReducedRange,
			&drive.BatteryDetails.IsSufficientlyPrecise,
			&drive.RangeIdeal.StartRange,
			&drive.RangeIdeal.EndRange,
			&drive.RangeIdeal.RangeDiff,
			&drive.RangeRated.StartRange,
			&drive.RangeRated.EndRange,
			&drive.RangeRated.RangeDiff,
			&drive.OutsideTempAvg,
			&drive.InsideTempAvg,
			&UnitsLength,
			&UnitsTemperature,
			&CarName,
		)

		// converting values based of settings UnitsLength
		if UnitsLength == "mi" {
			drive.OdometerDetails.OdometerStart = kilometersToMiles(drive.OdometerDetails.OdometerStart)
			drive.OdometerDetails.OdometerEnd = kilometersToMiles(drive.OdometerDetails.OdometerEnd)
			drive.OdometerDetails.OdometerDistance = kilometersToMiles(drive.OdometerDetails.OdometerDistance)
			drive.SpeedMax = int(kilometersToMiles(float64(drive.SpeedMax)))
			drive.SpeedAvg = kilometersToMiles(drive.SpeedAvg)
			drive.RangeIdeal.StartRange = kilometersToMiles(drive.RangeIdeal.StartRange)
			drive.RangeIdeal.EndRange = kilometersToMiles(drive.RangeIdeal.EndRange)
			drive.RangeIdeal.RangeDiff = kilometersToMiles(drive.RangeIdeal.RangeDiff)
			drive.RangeRated.StartRange = kilometersToMiles(drive.RangeRated.StartRange)
			drive.RangeRated.EndRange = kilometersToMiles(drive.RangeRated.EndRange)
			drive.RangeRated.RangeDiff = kilometersToMiles(drive.RangeRated.RangeDiff)
		}
		// converting values based of settings UnitsTemperature
		if UnitsTemperature == "F" {
			drive.OutsideTempAvg = celsiusToFahrenheit(drive.OutsideTempAvg)
			drive.InsideTempAvg = celsiusToFahrenheit(drive.InsideTempAvg)
		}
		// adjusting to timezone differences from UTC to be userspecific
		drive.StartDate = getTimeInTimeZone(drive.StartDate)
		drive.EndDate = getTimeInTimeZone(drive.EndDate)

		// checking for errors after scanning
		if err != nil {
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsDrivesDetailsV1", "Unable to load drive.", err.Error())
			return
		}

		// appending drive to DriveData
		DriveData = drive

		// getting detailed drive data from database
		query = `
		 			SELECT
						id AS detail_id,
						date,
						latitude,
						longitude,
						COALESCE(speed, 0) AS speed,
						power,
						odometer,
						battery_level,
						usable_battery_level,
						elevation,
						inside_temp,
						outside_temp,
						is_climate_on,
						fan_status,
						driver_temp_setting,
						passenger_temp_setting,
						is_rear_defroster_on,
						is_front_defroster_on,
						est_battery_range_km,
						ideal_battery_range_km,
						rated_battery_range_km,
						battery_heater,
						battery_heater_on,
						battery_heater_no_power
		 			FROM positions
		 			WHERE drive_id = $1
		 			ORDER BY id ASC;`
		rows, err = db.Query(query, DriveID)

		// defer closing rows
		defer rows.Close()

		// checking for errors in query
		if err != nil {
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsDrivesDetailsV1", "Unable to load drive details.", err.Error())
			return
		}

		// looping through all results
		for rows.Next() {

			// creating drivedetails object based on struct
			drivedetails := DriveDetails{}

			// scanning row and putting values into the drive
			err = rows.Scan(
				&drivedetails.DetailID,
				&drivedetails.Date,
				&drivedetails.Latitude,
				&drivedetails.Longitude,
				&drivedetails.Speed,
				&drivedetails.Power,
				&drivedetails.Odometer,
				&drivedetails.BatteryLevel,
				&drivedetails.UsableBatteryLevel,
				&drivedetails.Elevation,
				&drivedetails.ClimateInfo.InsideTemp,
				&drivedetails.ClimateInfo.OutsideTemp,
				&drivedetails.ClimateInfo.IsClimateOn,
				&drivedetails.ClimateInfo.FanStatus,
				&drivedetails.ClimateInfo.DriverTempSetting,
				&drivedetails.ClimateInfo.PassengerTempSetting,
				&drivedetails.ClimateInfo.IsRearDefrosterOn,
				&drivedetails.ClimateInfo.IsFrontDefrosterOn,
				&drivedetails.BatteryInfo.EstBatteryRange,
				&drivedetails.BatteryInfo.IdealBatteryRange,
				&drivedetails.BatteryInfo.RatedBatteryRange,
				&drivedetails.BatteryInfo.BatteryHeater,
				&drivedetails.BatteryInfo.BatteryHeaterOn,
				&drivedetails.BatteryInfo.BatteryHeaterNoPower,
			)

			// converting values based of settings UnitsLength
			if UnitsLength == "mi" {
				drivedetails.Odometer = kilometersToMiles(drivedetails.Odometer)
				drivedetails.Speed = int(kilometersToMiles(float64(drivedetails.Speed)))
				drivedetails.BatteryInfo.EstBatteryRange = kilometersToMilesNilSupport(drivedetails.BatteryInfo.EstBatteryRange)
				drivedetails.BatteryInfo.IdealBatteryRange = kilometersToMilesNilSupport(drivedetails.BatteryInfo.IdealBatteryRange)
				drivedetails.BatteryInfo.RatedBatteryRange = kilometersToMilesNilSupport(drivedetails.BatteryInfo.RatedBatteryRange)
			}
			// converting values based of settings UnitsTemperature
			if UnitsTemperature == "F" {
				drivedetails.ClimateInfo.InsideTemp = celsiusToFahrenheitNilSupport(drivedetails.ClimateInfo.InsideTemp)
				drivedetails.ClimateInfo.OutsideTemp = celsiusToFahrenheitNilSupport(drivedetails.ClimateInfo.OutsideTemp)
				drivedetails.ClimateInfo.DriverTempSetting = celsiusToFahrenheitNilSupport(drivedetails.ClimateInfo.DriverTempSetting)
				drivedetails.ClimateInfo.PassengerTempSetting = celsiusToFahrenheitNilSupport(drivedetails.ClimateInfo.PassengerTempSetting)
			}
			// adjusting to timezone differences from UTC to be userspecific
			drivedetails.Date = getTimeInTimeZone(drivedetails.Date)

			// checking for errors after scanning
			if err != nil {
				TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsDrivesDetailsV1", "Unable to load drive details.", err.Error())
				return
			}

			// appending drive to DriveData
			DriveDetailsData = append(DriveDetailsData, drivedetails)
			DriveData.DriveDetails = DriveDetailsData
		}

		// checking for errors in the rows result
		err = rows.Err()
		if err != nil {
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsDrivesDetailsV1", "Unable to load drive details.", err.Error())
			return
		}

	}

	// checking for errors in the rows result
	err = rows.Err()
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsDrivesDetailsV1", "Unable to load drive details.", err.Error())
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
			Drive: DriveData,
			TeslaMateUnits: TeslaMateUnits{
				UnitsLength:      UnitsLength,
				UnitsTemperature: UnitsTemperature,
			},
		},
	}

	// return jsonData
	TeslaMateAPIHandleSuccessResponse(c, "TeslaMateAPICarsDrivesDetailsV1", jsonData)
}
