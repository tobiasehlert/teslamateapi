package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TeslaMateAPIChargesSummaryV1 func
func TeslaMateAPIChargesSummaryV1(c *gin.Context) {
	// define error messages
	const errMsg1 = "Unable to load charges summary."
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

	// creating structs for /cars/:CarID/charges/summary
	// Car struct
	type Car struct {
		CarID   int        `json:"car_id"`
		CarName NullString `json:"car_name"`
	}
	// ChargesSummary struct
	type ChargesSummary struct {
		TotalCharges       int      `json:"total_charges"`
		TotalEnergyAdded   float64  `json:"total_energy_added"`
		TotalEnergyUsed    float64  `json:"total_energy_used"`
		TotalCost          float64  `json:"total_cost"`
		StartDate          *string  `json:"start_date"`
		EndDate            *string  `json:"end_date"`
	}

	// creating required vars
	var (
		CarDetails            Car
		ChargesSummaryData    ChargesSummary
		earliestStartDate     sql.NullTime
		latestEndDate         sql.NullTime
	)

	// getting car details
	query := `
		SELECT
			cars.id,
			cars.name
		FROM cars
		WHERE cars.id = $1
		LIMIT 1;`

	err := db.QueryRow(query, CarID).Scan(
		&CarDetails.CarID,
		&CarDetails.CarName,
	)

	// checking if car exists
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "CarID", errMsg2, err.Error())
		return
	}

	// build summary query with dynamic filtering
	summaryQuery := `
		SELECT
			COUNT(charging_processes.id) as total_charges,
			COALESCE(SUM(charging_processes.charge_energy_added), 0) as total_energy_added,
			COALESCE(SUM(charging_processes.charge_energy_used), 0) as total_energy_used,
			COALESCE(SUM(charging_processes.cost), 0) as total_cost,
			MIN(charging_processes.start_date) as earliest_start_date,
			MAX(charging_processes.end_date) as latest_end_date
		FROM charging_processes
		WHERE charging_processes.car_id = $1
			AND charging_processes.end_date IS NOT NULL`

	// add date filtering if provided
	queryParams := []interface{}{CarID}
	paramIndex := 2

	if startDateParsed != "" {
		summaryQuery += fmt.Sprintf(" AND charging_processes.start_date >= $%d", paramIndex)
		queryParams = append(queryParams, startDateParsed)
		paramIndex++
	}
	if endDateParsed != "" {
		summaryQuery += fmt.Sprintf(" AND charging_processes.end_date <= $%d", paramIndex)
		queryParams = append(queryParams, endDateParsed)
		paramIndex++
	}

	summaryQuery += ";"

	// execute query
	err = db.QueryRow(summaryQuery, queryParams...).Scan(
		&ChargesSummaryData.TotalCharges,
		&ChargesSummaryData.TotalEnergyAdded,
		&ChargesSummaryData.TotalEnergyUsed,
		&ChargesSummaryData.TotalCost,
		&earliestStartDate,
		&latestEndDate,
	)

	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "Database", errMsg1, err.Error())
		return
	}

	// handle nullable dates and convert to timezone
	if earliestStartDate.Valid {
		startDateStr := getTimeInTimeZone(earliestStartDate.Time.Format(time.RFC3339))
		ChargesSummaryData.StartDate = &startDateStr
	}
	if latestEndDate.Valid {
		endDateStr := getTimeInTimeZone(latestEndDate.Time.Format(time.RFC3339))
		ChargesSummaryData.EndDate = &endDateStr
	}

	// create the data response
	TeslaMateAPIHandleSuccessResponse(c, "ok", gin.H{
		"car": gin.H{
			"car_id":   CarDetails.CarID,
			"car_name": CarDetails.CarName,
		},
		"charges_summary": ChargesSummaryData,
	})

	// return http statuscode 200
	c.Status(http.StatusOK)
}
