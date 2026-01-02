package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// TeslaMateAPICarsDrivesV1 func
func TeslaMateAPICarsDrivesV1(c *gin.Context) {

	// define error messages
	var CarsDrivesError1 = "Unable to load drives."
	var CarsDrivesError2 = "Invalid date format."

	// getting CarID param from URL
	CarID := convertStringToInteger(c.Param("CarID"))
	// query options to modify query when collecting data
	ResultPage := convertStringToInteger(c.DefaultQuery("page", "1"))
	ResultShow := convertStringToInteger(c.DefaultQuery("show", "100"))

	// get startDate and endDate from query parameters
	parsedStartDate, err := parseDateParam(c.Query("startDate"))
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsDrivesV1", CarsDrivesError2, err.Error())
		return
	}
	parsedEndDate, err := parseDateParam(c.Query("endDate"))
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsDrivesV1", CarsDrivesError2, err.Error())
		return
	}
	// get optional minDistance and maxDistance filters from query parameters
	minDistanceParam := c.Query("minDistance")
	minDistance := 0.0
	if minDistanceParam != "" {
		minDistance = convertStringToFloat(minDistanceParam)
		if minDistance < 0 {
			minDistance = 0
		}
	}
	maxDistanceParam := c.Query("maxDistance")
	maxDistance := 0.0
	if maxDistanceParam != "" {
		maxDistance = convertStringToFloat(maxDistanceParam)
		if maxDistance < 0 {
			maxDistance = 0
		}
	}

	// creating structs for /cars/<CarID>/drives
	// Car struct - child of Data
	type Car struct {
		CarID   int        `json:"car_id"`   // smallint
		CarName NullString `json:"car_name"` // text (nullable)
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
	// Drives struct - child of Data
	type Drives struct {
		DriveID           int             `json:"drive_id"`            // int
		StartDate         string          `json:"start_date"`          // string
		EndDate           string          `json:"end_date"`            // string
		StartAddress      string          `json:"start_address"`       // string
		EndAddress        string          `json:"end_address"`         // string
		OdometerDetails   OdometerDetails `json:"odometer_details"`    // OdometerDetails
		DurationMin       int             `json:"duration_min"`        // int
		DurationStr       string          `json:"duration_str"`        // string
		SpeedMax          int             `json:"speed_max"`           // int
		SpeedAvg          float64         `json:"speed_avg"`           // float64
		PowerMax          int             `json:"power_max"`           // int
		PowerMin          int             `json:"power_min"`           // int
		BatteryDetails    BatteryDetails  `json:"battery_details"`     // BatteryDetails
		RangeIdeal        PreferredRange  `json:"range_ideal"`         // PreferredRange
		RangeRated        PreferredRange  `json:"range_rated"`         // PreferredRange
		OutsideTempAvg    float64         `json:"outside_temp_avg"`    // float64
		InsideTempAvg     float64         `json:"inside_temp_avg"`     // float64
		EnergyConsumedNet *float64        `json:"energy_consumed_net"` // Energy consumed (net) in kWh
		ConsumptionNet    *float64        `json:"consumption_net"`     // Ø Consumption (net) per distance unit
	}
	// DrivesSummary struct - child of Data
	type DrivesSummary struct {
		Drives       int      `json:"drives"`        // int - count of drives
		Distance     float64  `json:"distance"`      // float64 - total distance
		DurationMin  int      `json:"duration_min"`  // int - total duration in minutes
		ConsumedNet  *float64 `json:"consumed_net"`  // float64 - total energy consumed (nullable)
	}
	// TeslaMateUnits struct - child of Data
	type TeslaMateUnits struct {
		UnitsLength      string `json:"unit_of_length"`      // string
		UnitsTemperature string `json:"unit_of_temperature"` // string
	}
	// Data struct - child of JSONData
	type Data struct {
		Car            Car            `json:"car"`
		Summary        DrivesSummary  `json:"summary"`
		Drives         []Drives       `json:"drives"`
		TeslaMateUnits TeslaMateUnits `json:"units"`
	}
	// JSONData struct - main
	type JSONData struct {
		Data Data `json:"data"`
	}

	// creating required vars
	var (
		CarName                       NullString
		DrivesData                    []Drives
		UnitsLength, UnitsTemperature string
		SummaryData                   DrivesSummary
		totalDistanceKm               float64
		totalDurationMin              int
		totalEnergyConsumed           *float64
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
			CASE 
				WHEN (start_rated_range_km - end_rated_range_km) > 0 
				THEN (start_rated_range_km - end_rated_range_km) * cars.efficiency 
				ELSE NULL 
			END as energy_consumed_net,
			CASE 
				WHEN (duration_min > 1 AND distance > 1 AND ( start_position.usable_battery_level IS NULL OR end_position.usable_battery_level IS NULL OR ( end_position.battery_level - end_position.usable_battery_level ) = 0 )) AND NULLIF(distance, 0) IS NOT NULL
				THEN (start_rated_range_km - end_rated_range_km) * cars.efficiency / NULLIF(distance, 0) * 1000
				ELSE NULL 
			END as consumption_net,
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
		WHERE drives.car_id=$1 AND end_date IS NOT NULL`

	// Parameters to be passed to the query
	var queryParams []any
	queryParams = append(queryParams, CarID)
	paramIndex := 2

	// Add date filtering if provided
	if parsedStartDate != "" {
		query += fmt.Sprintf(" AND drives.start_date >= $%d", paramIndex)
		queryParams = append(queryParams, parsedStartDate)
		paramIndex++
	}
	if parsedEndDate != "" {
		query += fmt.Sprintf(" AND drives.end_date <= $%d", paramIndex)
		queryParams = append(queryParams, parsedEndDate)
		paramIndex++
	}

	// Add minimum/maximum distance filtering if provided
	if minDistance > 0 || maxDistance > 0 {
		var unitsLength string
		err = db.QueryRow("SELECT unit_of_length FROM settings LIMIT 1").Scan(&unitsLength)
		if err != nil {
			TeslaMateAPIHandleErrorResponse(
				c,
				"TeslaMateAPICarsDrivesV1",
				CarsDrivesError1,
				fmt.Sprintf("unable to retrieve unit_of_length from settings table: %v", err),
			)
			return
		}
		if unitsLength == "mi" {
			if minDistance > 0 {
				minDistance = milesToKilometers(minDistance)
			}
			if maxDistance > 0 {
				maxDistance = milesToKilometers(maxDistance)
			}
		}

		if minDistance > 0 {
			query += fmt.Sprintf(" AND distance >= $%d", paramIndex)
			queryParams = append(queryParams, minDistance)
			paramIndex++
		}
		if maxDistance > 0 {
			query += fmt.Sprintf(" AND distance <= $%d", paramIndex)
			queryParams = append(queryParams, maxDistance)
			paramIndex++
		}
	}

	// Build summary query with the same filters (but without pagination)
	summaryQuery := `
		SELECT
			COUNT(drives.id) as total_drives,
			COALESCE(SUM(drives.distance), 0) as total_distance_km,
			COALESCE(SUM(drives.duration_min), 0) as total_duration_min,
			SUM(
				CASE
					WHEN (start_rated_range_km - end_rated_range_km) > 0
					THEN (start_rated_range_km - end_rated_range_km) * cars.efficiency
					ELSE NULL
				END
			) as total_energy_consumed
		FROM drives
		LEFT JOIN cars ON car_id = cars.id
		LEFT JOIN positions start_position ON start_position_id = start_position.id
		LEFT JOIN positions end_position ON end_position_id = end_position.id
		WHERE drives.car_id=$1 AND end_date IS NOT NULL`

	// Build summary query params (same as main query params, but without pagination)
	summaryParams := []any{CarID}
	summaryParamIndex := 2

	if parsedStartDate != "" {
		summaryQuery += fmt.Sprintf(" AND drives.start_date >= $%d", summaryParamIndex)
		summaryParams = append(summaryParams, parsedStartDate)
		summaryParamIndex++
	}
	if parsedEndDate != "" {
		summaryQuery += fmt.Sprintf(" AND drives.end_date <= $%d", summaryParamIndex)
		summaryParams = append(summaryParams, parsedEndDate)
		summaryParamIndex++
	}

	// Add distance filtering to summary query if provided
	if minDistance > 0 || maxDistance > 0 {
		// minDistance and maxDistance are already converted to km above
		if minDistance > 0 {
			summaryQuery += fmt.Sprintf(" AND distance >= $%d", summaryParamIndex)
			summaryParams = append(summaryParams, minDistance)
			summaryParamIndex++
		}
		if maxDistance > 0 {
			summaryQuery += fmt.Sprintf(" AND distance <= $%d", summaryParamIndex)
			summaryParams = append(summaryParams, maxDistance)
			summaryParamIndex++
		}
	}

	summaryQuery += ";"

	// Execute summary query
	err = db.QueryRow(summaryQuery, summaryParams...).Scan(
		&SummaryData.Drives,
		&totalDistanceKm,
		&totalDurationMin,
		&totalEnergyConsumed,
	)

	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsDrivesV1", CarsDrivesError1, err.Error())
		return
	}

	// Convert distance based on unit settings
	if UnitsLength == "mi" {
		SummaryData.Distance = kilometersToMiles(totalDistanceKm)
	} else {
		SummaryData.Distance = totalDistanceKm
	}

	SummaryData.DurationMin = totalDurationMin
	SummaryData.ConsumedNet = totalEnergyConsumed

	query += fmt.Sprintf(`
        ORDER BY start_date DESC
        LIMIT $%d OFFSET $%d;`, paramIndex, paramIndex+1)

	queryParams = append(queryParams, ResultShow, ResultPage)

	rows, err := db.Query(query, queryParams...)

	// checking for errors in query
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsDrivesV1", CarsDrivesError1, err.Error())
		return
	}

	// defer closing rows
	defer rows.Close()

	// looping through all results
	for rows.Next() {

		// creating drive object based on struct
		drive := Drives{}

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
			&drive.EnergyConsumedNet,
			&drive.ConsumptionNet,
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
			if drive.ConsumptionNet != nil {
				*drive.ConsumptionNet = kilometersToMiles(*drive.ConsumptionNet)
			}
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
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsDrivesV1", CarsDrivesError1, err.Error())
			return
		}

		// appending drive to DrivesData
		DrivesData = append(DrivesData, drive)
	}

	// checking for errors in the rows result
	err = rows.Err()
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsDrivesV1", CarsDrivesError1, err.Error())
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
			Summary: SummaryData,
			Drives:  DrivesData,
			TeslaMateUnits: TeslaMateUnits{
				UnitsLength:      UnitsLength,
				UnitsTemperature: UnitsTemperature,
			},
		},
	}

	// return jsonData
	TeslaMateAPIHandleSuccessResponse(c, "TeslaMateAPICarsDrivesV1", jsonData)
}
