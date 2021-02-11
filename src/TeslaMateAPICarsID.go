package main

import (
	"log"

	_ "github.com/lib/pq"
)

// TeslaMateAPICarsID func
func TeslaMateAPICarsID(CarID int) bool {

	// creating structs for storing data
	type Car struct {
		CarID int `json:"car_id"` // smallint
	}

	// creating required vars
	var ValidCarID bool

	// getting data from database
	query := `
		SELECT
			id,
		FROM cars
		WHERE id=$1
		LIMIT 1;`
	rows, err := db.Query(query, CarID)

	// checking for errors in query
	if err != nil {
		log.Fatal(err)
	}

	// defer closing rows
	defer rows.Close()

	// looping through all results
	for rows.Next() {

		// creating car object based on struct
		car := Car{}

		// scanning row and putting values into the car
		err = rows.Scan(
			&car.CarID,
		)

		// checking for errors after scanning
		if err != nil {
			log.Fatal(err)
		}

		if CarID != 0 && CarID == car.CarID {
			// return true that CarID exists.. to do correct redirect to status page of particular car
			ValidCarID = true
		}
	}

	// checking for errors in the rows result
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// returning true/false if car exists in database or not..
	return ValidCarID
}
