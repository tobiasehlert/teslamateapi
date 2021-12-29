package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// TeslaMateAPIGlobalsettingsV1 func
func TeslaMateAPIGlobalsettingsV1(c *gin.Context) {

	// define error messages
	var CarsGlobalsettingsError1 = "Unable to load settings."

	// creating structs for /globalsettings
	// AccountInfo struct - child of GlobalSettings
	type AccountInfo struct {
		InsertedAt string `json:"inserted_at"` // string
		UpdatedAt  string `json:"updated_at"`  // string
	}
	// TeslaMateUnits struct - child of GlobalSettings
	type TeslaMateUnits struct {
		UnitsLength      string `json:"unit_of_length"`      // string
		UnitsTemperature string `json:"unit_of_temperature"` // string
	}
	// TeslaMateGUI struct - child of GlobalSettings
	type TeslaMateGUI struct {
		PreferredRange string `json:"preferred_range"` // string
		Language       string `json:"language"`        // string
	}
	// TeslaMateURLs struct - child of GlobalSettings
	type TeslaMateURLs struct {
		BaseURL    string `json:"base_url"`    // string
		GrafanaURL string `json:"grafana_url"` // string
	}
	// GlobalSettings struct - child of Data
	type GlobalSettings struct {
		SettingID      int            `json:"setting_id"`       // smallint
		AccountInfo    AccountInfo    `json:"account_info"`     // struct
		TeslaMateUnits TeslaMateUnits `json:"teslamate_units"`  // struct
		TeslaMateGUI   TeslaMateGUI   `json:"teslamate_webgui"` // struct
		TeslaMateURLs  TeslaMateURLs  `json:"teslamate_urls"`   // struct
	}
	// Data struct - child of JSONData
	type Data struct {
		GlobalSettings GlobalSettings `json:"settings"`
	}
	// JSONData struct - main
	type JSONData struct {
		Data Data `json:"data"`
	}

	// creating required vars
	var GlobalSettingData GlobalSettings

	// getting data from database
	query := `
		SELECT
			id,
			inserted_at,
			updated_at,
			unit_of_length,
			unit_of_temperature,
			preferred_range,
			language,
			base_url,
			grafana_url
		FROM settings
		LIMIT 1;`
	rows, err := db.Query(query)

	// checking for errors in query
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPIGlobalsettingsV1", CarsGlobalsettingsError1, err.Error())
		return
	}

	// defer closing rows
	defer rows.Close()

	// looping through all results
	for rows.Next() {

		// creating GlobalSetting object based on struct
		GlobalSetting := GlobalSettings{}

		// scanning row and putting values into the GlobalSetting
		err = rows.Scan(
			&GlobalSetting.SettingID,
			&GlobalSetting.AccountInfo.InsertedAt,
			&GlobalSetting.AccountInfo.UpdatedAt,
			&GlobalSetting.TeslaMateUnits.UnitsLength,
			&GlobalSetting.TeslaMateUnits.UnitsTemperature,
			&GlobalSetting.TeslaMateGUI.PreferredRange,
			&GlobalSetting.TeslaMateGUI.Language,
			&GlobalSetting.TeslaMateURLs.BaseURL,
			&GlobalSetting.TeslaMateURLs.GrafanaURL,
		)

		// checking for errors after scanning
		if err != nil {
			TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPIGlobalsettingsV1", CarsGlobalsettingsError1, err.Error())
			return
		}

		// adjusting to timezone differences from UTC to be userspecific
		GlobalSetting.AccountInfo.InsertedAt = getTimeInTimeZone(GlobalSetting.AccountInfo.InsertedAt)
		GlobalSetting.AccountInfo.UpdatedAt = getTimeInTimeZone(GlobalSetting.AccountInfo.UpdatedAt)

		// setting response to valid
		GlobalSettingData = GlobalSetting
	}

	// checking for errors in the rows result
	err = rows.Err()
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPIGlobalsettingsV1", CarsGlobalsettingsError1, err.Error())
		return
	}

	//
	// build the data-blob
	jsonData := JSONData{
		Data{
			GlobalSettings: GlobalSettingData,
		},
	}

	// return jsonData
	TeslaMateAPIHandleSuccessResponse(c, "TeslaMateAPIGlobalsettingsV1", jsonData)
}
