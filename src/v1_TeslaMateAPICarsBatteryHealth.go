package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// TeslaMateAPICarsBatteryHealthV1 func
func TeslaMateAPICarsBatteryHealthV1(c *gin.Context) {
	var CarsBatteryHealthError1 = "Unable to load battery health data."
	CarID := convertStringToInteger(c.Param("CarID"))

	// creating structs for /cars/<CarID>/battery-health
	// Car struct - child of Data
	type Car struct {
		CarID   int        `json:"car_id"`   // smallint
		CarName NullString `json:"car_name"` // text (nullable)
	}
	// BatteryHealth struct - child of Data
	type BatteryHealth struct {
		MaxRange                float64 `json:"max_range"`                 // float64
		CurrentRange            float64 `json:"current_range"`             // float64
		MaxCapacity             float64 `json:"max_capacity"`              // float64
		CurrentCapacity         float64 `json:"current_capacity"`          // float64
		RatedEfficiency         float64 `json:"rated_efficiency"`          // float64
		BatteryHealthPercentage float64 `json:"battery_health_percentage"` // float64
	}
	// TeslaMateUnits struct - child of Data
	type TeslaMateUnits struct {
		UnitsLength      string `json:"unit_of_length"`      // string
		UnitsTemperature string `json:"unit_of_temperature"` // string
	}
	// Data struct - child of JSONData
	type Data struct {
		Car            Car            `json:"car"`
		BatteryHealth  BatteryHealth  `json:"battery_health"`
		TeslaMateUnits TeslaMateUnits `json:"units"`
	}
	// JSONData struct - main
	type JSONData struct {
		Data Data `json:"data"`
	}

	// creating required vars
	var (
		CarName                       NullString
		Efficiency                    float64
		MaxRangeRated                 float64
		MaxRangeIdeal                 float64
		CurrentRangeRated             float64
		CurrentRangeIdeal             float64
		MaxCapacity                   float64
		CurrentCapacity               float64
		PreferredRange                string
		UnitsLength, UnitsTemperature string
	)

	query := `
	WITH Aux as (
		SELECT 
			car_id,
			COALESCE(derived_efficiency, car_efficiency) AS efficiency
		FROM (
			SELECT
				ROUND((charge_energy_added / NULLIF(end_rated_range_km - start_rated_range_km, 0))::numeric, 3) * 100 AS derived_efficiency,
				COUNT(*) as count,
				cars.id as car_id,
				cars.efficiency * 100 AS car_efficiency
			FROM cars
				LEFT JOIN charging_processes ON
					cars.id = charging_processes.car_id 
					AND duration_min > 10
					AND end_battery_level <= 95
					AND start_rated_range_km IS NOT NULL
					AND end_rated_range_km IS NOT NULL
					AND charge_energy_added > 0
			WHERE cars.id = $1
			GROUP BY 1, 3, 4
			ORDER BY 2 DESC
			LIMIT 1
		) AS Efficiency
	),
	CurrentCapacity AS (
		SELECT
			AVG(Capacity) AS Capacity
		FROM (
			SELECT 
				c.rated_battery_range_km * aux.efficiency / c.usable_battery_level AS Capacity
			FROM charging_processes cp
				INNER JOIN charges c ON c.charging_process_id = cp.id 
				INNER JOIN aux ON cp.car_id = aux.car_id
			WHERE
				cp.car_id = $1
				AND cp.end_date IS NOT NULL
				AND cp.charge_energy_added >= aux.efficiency
				AND c.usable_battery_level > 0
			ORDER BY cp.end_date DESC, c.date desc
			LIMIT 100
		) AS lastCharges
	),
	MaxCapacity AS (
		SELECT 
			MAX(c.rated_battery_range_km * aux.efficiency / c.usable_battery_level) AS Capacity
		FROM charging_processes cp
			INNER JOIN (
				SELECT
					charging_process_id,
					MAX(date) as date FROM charges WHERE usable_battery_level > 0 GROUP BY charging_process_id
			) AS gcharges ON
				cp.id = gcharges.charging_process_id
			INNER JOIN charges c ON
				c.charging_process_id = cp.id
				AND c.date = gcharges.date
			INNER JOIN aux ON cp.car_id = aux.car_id
		WHERE
			cp.car_id = $1
			AND cp.end_date IS NOT NULL
			AND cp.charge_energy_added >= aux.efficiency
	),
	CurrentRangeRated AS (
		SELECT
			(range * 100.0 / usable_battery_level) AS range
		FROM (
			(
				SELECT
					date,
					rated_battery_range_km AS range,
					usable_battery_level AS usable_battery_level
				FROM positions
				WHERE
					car_id = $1
					AND ideal_battery_range_km IS NOT NULL
					AND usable_battery_level > 0 
				ORDER BY date DESC
				LIMIT 1
			)
			UNION ALL
			(
				SELECT date,
					rated_battery_range_km AS range,
					usable_battery_level as usable_battery_level
				FROM charges c
					INNER JOIN charging_processes p ON p.id = c.charging_process_id
				WHERE
					p.car_id = $1
					AND usable_battery_level > 0
				ORDER BY date DESC
				LIMIT 1
			)
		) AS data
		ORDER BY date DESC
		LIMIT 1
	),
	CurrentRangeIdeal AS (
		SELECT
			(range * 100.0 / usable_battery_level) AS range
		FROM (
			(
				SELECT
					date,
					ideal_battery_range_km AS range,
					usable_battery_level AS usable_battery_level
				FROM positions
				WHERE
					car_id = $1
					AND ideal_battery_range_km IS NOT NULL
					AND usable_battery_level > 0 
				ORDER BY date DESC
				LIMIT 1
			)
			UNION ALL
			(
				SELECT date,
					ideal_battery_range_km AS range,
					usable_battery_level as usable_battery_level
				FROM charges c
					INNER JOIN charging_processes p ON p.id = c.charging_process_id
				WHERE
					p.car_id = $1
					AND usable_battery_level > 0
				ORDER BY date DESC
				LIMIT 1
			)
		) AS data
		ORDER BY date DESC
		LIMIT 1
	),
	MaxRangeRated AS (
		SELECT
			CASE
				WHEN sum(usable_battery_level) = 0 THEN sum(rated_battery_range_km) * 100
				ELSE sum(rated_battery_range_km) / sum(usable_battery_level) * 100
			END AS range
		FROM (
			SELECT
				battery_level,
				usable_battery_level,
				date,
				rated_battery_range_km
			FROM charges c 
				INNER JOIN charging_processes p ON p.id = c.charging_process_id 
			WHERE
				p.car_id = $1
				AND usable_battery_level IS NOT NULL
		) AS data
		GROUP BY floor(extract(epoch from date)/86400)*86400
		ORDER BY range DESC
		LIMIT 1
	),
	MaxRangeIdeal AS (
		SELECT
			CASE
				WHEN sum(usable_battery_level) = 0 THEN sum(ideal_battery_range_km) * 100
				ELSE sum(ideal_battery_range_km) / sum(usable_battery_level) * 100
			END AS range
		FROM (
			SELECT
				battery_level,
				usable_battery_level,
				date,
				ideal_battery_range_km
			FROM charges c 
				INNER JOIN charging_processes p ON p.id = c.charging_process_id 
			WHERE
				p.car_id = $1
				AND usable_battery_level IS NOT NULL
		) AS data
		GROUP BY floor(extract(epoch from date)/86400)*86400
		ORDER BY range DESC
		LIMIT 1
	)
	SELECT
		COALESCE(MaxRangeRated.range, 0) as max_range_rated,
		COALESCE(MaxRangeIdeal.range, 0) as max_range_ideal,
		COALESCE(CurrentRangeRated.range, 0) as current_range_rated,
		COALESCE(CurrentRangeIdeal.range, 0) as current_range_ideal,
		COALESCE(MaxCapacity.Capacity, 0) as max_capacity,
		COALESCE(CurrentCapacity.Capacity, 0) as current_capacity,
		COALESCE(aux.efficiency, 0) as efficiency,
		(SELECT preferred_range FROM settings LIMIT 1) as preferred_range,
		(SELECT unit_of_length FROM settings LIMIT 1) as unit_of_length,
		(SELECT unit_of_temperature FROM settings LIMIT 1) as unit_of_temperature,
		cars.name
	FROM cars
		LEFT JOIN MaxRangeRated ON true
		LEFT JOIN MaxRangeIdeal ON true
		LEFT JOIN CurrentRangeRated ON true
		LEFT JOIN CurrentRangeIdeal ON true
		LEFT JOIN Aux ON cars.id = aux.car_id
		LEFT JOIN MaxCapacity ON true
		LEFT JOIN CurrentCapacity ON true
	WHERE cars.id = $1;`

	// execute query
	err := db.QueryRow(query, CarID).Scan(
		&MaxRangeRated,
		&MaxRangeIdeal,
		&CurrentRangeRated,
		&CurrentRangeIdeal,
		&MaxCapacity,
		&CurrentCapacity,
		&Efficiency,
		&PreferredRange,
		&UnitsLength,
		&UnitsTemperature,
		&CarName,
	)

	// checking for errors in query
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsBatteryHealthV1", CarsBatteryHealthError1, err.Error())
		return
	}

	// Create battery health object
	batteryHealth := BatteryHealth{
		CurrentCapacity: CurrentCapacity,
		MaxCapacity:     MaxCapacity,
		RatedEfficiency: Efficiency,
	}

	// Select the correct range based on preferred_range setting
	if PreferredRange == "ideal" {
		batteryHealth.MaxRange = MaxRangeIdeal
		batteryHealth.CurrentRange = CurrentRangeIdeal
	} else {
		batteryHealth.MaxRange = MaxRangeRated
		batteryHealth.CurrentRange = CurrentRangeRated
	}

	// Calculate battery health percentage
	if MaxCapacity > 0 {
		batteryHealth.BatteryHealthPercentage = (CurrentCapacity / MaxCapacity) * 100
	} else {
		batteryHealth.BatteryHealthPercentage = 100
	}

	// converting values based on settings UnitsLength
	if UnitsLength == "mi" {
		batteryHealth.MaxRange = kilometersToMiles(batteryHealth.MaxRange)
		batteryHealth.CurrentRange = kilometersToMiles(batteryHealth.CurrentRange)
	}

	jsonData := JSONData{
		Data{
			Car: Car{
				CarID:   CarID,
				CarName: CarName,
			},
			BatteryHealth: batteryHealth,
			TeslaMateUnits: TeslaMateUnits{
				UnitsLength:      UnitsLength,
				UnitsTemperature: UnitsTemperature,
			},
		},
	}

	TeslaMateAPIHandleSuccessResponse(c, "TeslaMateAPICarsBatteryHealthV1", jsonData)
}
