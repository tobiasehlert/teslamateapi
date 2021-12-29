package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// TeslaMateAPICarsUpdatesV1 func
func TeslaMateAPICarsUpdatesV1(c *gin.Context) {

	// getting CarID param from URL
	CarID := convertStringToInteger(c.Param("CarID"))
	// query options to modify query when collecting data
	ResultPage := convertStringToInteger(c.DefaultQuery("page", "1"))
	ResultShow := convertStringToInteger(c.DefaultQuery("show", "100"))

	// creating structs for /cars/<CarID>/updates
	// Car struct - child of Data
	type Car struct {
		CarID   int    `json:"car_id"`   // smallint
		CarName string `json:"car_name"` // text
	}
	// Updates struct - child of Data
	type Updates struct {
		UpdateID  int    `json:"update_id"`  // smallint
		StartDate string `json:"start_date"` // string
		EndDate   string `json:"end_date"`   // string
		Version   string `json:"version"`    // string
	}
	// Data struct - child of JSONData
	type Data struct {
		Car     Car       `json:"car"`
		Updates []Updates `json:"updates"`
	}
	// JSONData struct - main
	type JSONData struct {
		Data Data `json:"data"`
	}

	// creating required vars
	var UpdatesData []Updates
	var CarData Car

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
			updates.id,
			cars.name,
			start_date,
			end_date,
			version
		FROM updates
		LEFT JOIN cars ON car_id = cars.id
		WHERE car_id = $1
		ORDER BY start_date DESC
		LIMIT $2 OFFSET $3;`
	rows, err := db.Query(query, CarID, ResultShow, ResultPage)

	// defer closing rows
	defer rows.Close()

	// checking for errors in query
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsUpdatesV1", "Unable to load updates.", err.Error())
		return

	}

	// looping through all results
	for rows.Next() {

		// creating update object based on struct
		update := Updates{}

		// scanning row and putting values into the update
		err = rows.Scan(
			&update.UpdateID,
			&CarData.CarName,
			&update.StartDate,
			&update.EndDate,
			&update.Version,
		)

		// checking for errors after scanning
		if err != nil {
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsUpdatesV1", "Unable to load updates.", err.Error())
			return
		}

		// adjusting to timezone differences from UTC to be userspecific
		update.StartDate = getTimeInTimeZone(update.StartDate)
		update.EndDate = getTimeInTimeZone(update.EndDate)

		// appending update to UpdatesData
		UpdatesData = append(UpdatesData, update)
		CarData.CarID = CarID
	}

	// checking for errors in the rows result
	err = rows.Err()
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsUpdatesV1", "Unable to load updates.", err.Error())
		return
	}

	//
	// build the data-blob
	jsonData := JSONData{
		Data{
			Car:     CarData,
			Updates: UpdatesData,
		},
	}

	// return jsonData
	TeslaMateAPIHandleSuccessResponse(c, http.StatusOK, "TeslaMateAPICarsUpdatesV1", jsonData)
}
