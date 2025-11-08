package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TeslaMateAPICarsDrivesSummaryV1 func
func TeslaMateAPICarsDrivesSummaryV1(c *gin.Context) {
	// define error messages
	const errMsg1 = "Unable to load drives summary."
	const errMsg2 = "Car not found."

	// getting CarID param from URL
	var CarID = convertStringToInteger(c.Param("CarID"))

	// check if CarID is not zero
	if CarID == 0 {
		TeslaMateAPIHandleErrorResponse(c, "CarID", errMsg1, "CarID is not valid.")
		return
	}

	// getting optional startDate and endDate query parameters
	startDateParam := c.Query("startDate")
	endDateParam := c.Query("endDate")

	// parse dates if provided
	startDateParsed := ""
	endDateParsed := ""
	if startDateParam != "" {
		parsed, err := parseDateParam(startDateParam)
		if err != nil {
			TeslaMateAPIHandleErrorResponse(c, "startDate", errMsg1, err.Error())
			return
		}
		startDateParsed = parsed
	}
	if endDateParam != "" {
		parsed, err := parseDateParam(endDateParam)
		if err != nil {
			TeslaMateAPIHandleErrorResponse(c, "endDate", errMsg1, err.Error())
			return
		}
		endDateParsed = parsed
	}

	// creating structs for /cars/:CarID/drives/summary
	// Car struct
	type Car struct {
		CarID   int        `json:"car_id"`
		CarName NullString `json:"car_name"`
	}
	// DrivesSummary struct
	type DrivesSummary struct {
		TotalDrives     int      `json:"total_drives"`
		TotalDistance   float64  `json:"total_distance"`
		TotalDurationHr float64  `json:"total_duration_hr"`
		StartDate       *string  `json:"start_date"`
		EndDate         *string  `json:"end_date"`
	}

	// creating required vars
	var (
		CarDetails            Car
		UnitsLength           string
		DrivesSummaryData     DrivesSummary
		totalDistanceKm       float64
		totalDurationMin      int
		earliestStartDate     sql.NullTime
		latestEndDate         sql.NullTime
	)

	// getting car details
	query := `
		SELECT
			cars.id,
			cars.name,
			settings.unit_of_length
		FROM cars
		LEFT JOIN settings ON settings.id = 1
		WHERE cars.id = $1
		LIMIT 1;`

	err := db.QueryRow(query, CarID).Scan(
		&CarDetails.CarID,
		&CarDetails.CarName,
		&UnitsLength,
	)

	// checking if car exists
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "CarID", errMsg2, err.Error())
		return
	}

	// build summary query with dynamic filtering
	summaryQuery := `
		SELECT
			COUNT(drives.id) as total_drives,
			COALESCE(SUM(drives.distance), 0) as total_distance_km,
			COALESCE(SUM(EXTRACT(epoch FROM (drives.end_date - drives.start_date)) / 60), 0) as total_duration_min,
			MIN(drives.start_date) as earliest_start_date,
			MAX(drives.end_date) as latest_end_date
		FROM drives
		WHERE drives.car_id = $1
			AND drives.end_date IS NOT NULL`

	// add date filtering if provided
	queryParams := []interface{}{CarID}
	paramIndex := 2

	if startDateParsed != "" {
		summaryQuery += fmt.Sprintf(" AND drives.start_date >= $%d", paramIndex)
		queryParams = append(queryParams, startDateParsed)
		paramIndex++
	}
	if endDateParsed != "" {
		summaryQuery += fmt.Sprintf(" AND drives.end_date <= $%d", paramIndex)
		queryParams = append(queryParams, endDateParsed)
		paramIndex++
	}

	summaryQuery += ";"

	// execute query
	err = db.QueryRow(summaryQuery, queryParams...).Scan(
		&DrivesSummaryData.TotalDrives,
		&totalDistanceKm,
		&totalDurationMin,
		&earliestStartDate,
		&latestEndDate,
	)

	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "Database", errMsg1, err.Error())
		return
	}

	// convert distance based on unit settings
	if UnitsLength == "mi" {
		DrivesSummaryData.TotalDistance = kilometersToMiles(totalDistanceKm)
	} else {
		DrivesSummaryData.TotalDistance = totalDistanceKm
	}

	// convert duration from minutes to hours
	DrivesSummaryData.TotalDurationHr = float64(totalDurationMin) / 60.0

	// handle nullable dates and convert to timezone
	if earliestStartDate.Valid {
		startDateStr := getTimeInTimeZone(earliestStartDate.Time.Format(time.RFC3339))
		DrivesSummaryData.StartDate = &startDateStr
	}
	if latestEndDate.Valid {
		endDateStr := getTimeInTimeZone(latestEndDate.Time.Format(time.RFC3339))
		DrivesSummaryData.EndDate = &endDateStr
	}

	// create the data response
	TeslaMateAPIHandleSuccessResponse(c, "ok", gin.H{
		"car": gin.H{
			"car_id":   CarDetails.CarID,
			"car_name": CarDetails.CarName,
		},
		"drives_summary": DrivesSummaryData,
		"units": gin.H{
			"unit_of_length": UnitsLength,
		},
	})

	// return http statuscode 200
	c.Status(http.StatusOK)
}
