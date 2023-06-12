package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// TeslaMateAPICarsV1 func
func TeslaMateAPICarsV1(c *gin.Context) {

	// define error messages
	var CarsError1 = "Unable to load cars."

	// getting CarID param from URL
	ParamCarID := c.Param("CarID")
	var CarID int
	if ParamCarID != "" {
		CarID = convertStringToInteger(ParamCarID)
	}

	// creating structs for /cars
	// CarDetails struct - child of Cars
	type CarDetails struct {
		EID         int64       `json:"eid"`          // bigint
		VID         int64       `json:"vid"`          // bigint
		Vin         string      `json:"vin"`          // text
		Model       string      `json:"model"`        // character varying(255)
		TrimBadging NullString  `json:"trim_badging"` // text
		Efficiency  NullFloat64 `json:"efficiency"`   // double precision
	}
	// CarExterior struct - child of Cars
	type CarExterior struct {
		ExteriorColor string `json:"exterior_color"` // text
		SpoilerType   string `json:"spoiler_type"`   // text
		WheelType     string `json:"wheel_type"`     // text
	}
	// CarSettings struct - child of Cars
	type CarSettings struct {
		SuspendMin          int  `json:"suspend_min"`            // int
		SuspendAfterIdleMin int  `json:"suspend_after_idle_min"` // int
		ReqNotUnlocked      bool `json:"req_not_unlocked"`       // bool
		FreeSupercharging   bool `json:"free_supercharging"`     // bool
		UseStreamingAPI     bool `json:"use_streaming_api"`      // bool
	}
	// TeslaMateDetails struct - child of Cars
	type TeslaMateDetails struct {
		InsertedAt string `json:"inserted_at"` // timestamp(0) without time zone
		UpdatedAt  string `json:"updated_at"`  // timestamp(0) without time zone
	}
	// TeslaMateStats struct - child of Cars
	type TeslaMateStats struct {
		TotalCharges int `json:"total_charges"` // int
		TotalDrives  int `json:"total_drives"`  // int
		TotalUpdates int `json:"total_updates"` // int
	}
	// Cars struct - child of Data
	type Cars struct {
		CarID            int              `json:"car_id"`            // smallint
		Name             NullString       `json:"name"`              // text
		CarDetails       CarDetails       `json:"car_details"`       // struct
		CarExterior      CarExterior      `json:"car_exterior"`      // struct
		CarSettings      CarSettings      `json:"car_settings"`      // struct
		TeslaMateDetails TeslaMateDetails `json:"teslamate_details"` // struct
		TeslaMateStats   TeslaMateStats   `json:"teslamate_stats"`   // struct
	}
	// Information struct - child of JSONData
	type Data struct {
		Cars []Cars `json:"cars"`
	}
	// JSONData struct - main
	type JSONData struct {
		Data Data `json:"data"`
	}

	// creating required vars
	var CarsData []Cars

	// getting data from database
	query := `
		SELECT
			cars.id,
			eid,
			vid,
			model,
			efficiency,
			inserted_at,
			updated_at,
			vin,
			name,
			trim_badging,
			exterior_color,
			spoiler_type,
			wheel_type,
			suspend_min,
			suspend_after_idle_min,
			req_not_unlocked,
			free_supercharging,
			use_streaming_api,
			(SELECT COUNT(*) FROM charging_processes WHERE car_id=cars.id) as total_charges,
			(SELECT COUNT(*) FROM drives WHERE car_id=cars.id) as total_drives,
			(SELECT COUNT(*) FROM updates WHERE car_id=cars.id) as total_charges
		FROM cars
		LEFT JOIN car_settings ON cars.id = car_settings.id
		ORDER BY id;`
	rows, err := db.Query(query)

	// checking for errors in query
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsV1", CarsError1, err.Error())
		return
	}

	// defer closing rows
	defer rows.Close()

	// looping through all results
	for rows.Next() {

		// creating car object based on struct
		car := Cars{}

		// scanning row and putting values into the car
		err = rows.Scan(
			&car.CarID,
			&car.CarDetails.EID,
			&car.CarDetails.VID,
			&car.CarDetails.Model,
			&car.CarDetails.Efficiency,
			&car.TeslaMateDetails.InsertedAt,
			&car.TeslaMateDetails.UpdatedAt,
			&car.CarDetails.Vin,
			&car.Name,
			&car.CarDetails.TrimBadging,
			&car.CarExterior.ExteriorColor,
			&car.CarExterior.SpoilerType,
			&car.CarExterior.WheelType,
			&car.CarSettings.SuspendMin,
			&car.CarSettings.SuspendAfterIdleMin,
			&car.CarSettings.ReqNotUnlocked,
			&car.CarSettings.FreeSupercharging,
			&car.CarSettings.UseStreamingAPI,
			&car.TeslaMateStats.TotalCharges,
			&car.TeslaMateStats.TotalDrives,
			&car.TeslaMateStats.TotalUpdates,
		)

		// checking for errors after scanning
		if err != nil {
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsV1", CarsError1, err.Error())
			return
		}

		// appending car to CarsData if CarID is 0 or is CarID matches car.CarID
		if CarID == 0 && len(ParamCarID) == 0 || CarID != 0 && CarID == car.CarID {

			// adjusting to timezone differences from UTC to be userspecific
			car.TeslaMateDetails.InsertedAt = getTimeInTimeZone(car.TeslaMateDetails.InsertedAt)
			car.TeslaMateDetails.UpdatedAt = getTimeInTimeZone(car.TeslaMateDetails.UpdatedAt)

			CarsData = append(CarsData, car)
		}
	}

	// checking for errors in the rows result
	err = rows.Err()
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsV1", CarsError1, err.Error())
		return
	}

	//
	// build the data-blob
	jsonData := JSONData{
		Data{
			Cars: CarsData,
		},
	}

	// return jsonData
	TeslaMateAPIHandleSuccessResponse(c, "TeslaMateAPICarsV1", jsonData)

}
