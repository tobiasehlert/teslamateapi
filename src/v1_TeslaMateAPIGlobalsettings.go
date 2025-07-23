package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// TeslaMateAPIGlobalsettingsV1 godoc
// @Summary      TeslaMate global settings endpoint
// @Description  Get the global settings from TeslaMate
// @Tags         globalsettings
// @Router       /v1/globalsettings [get]
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
	var globalSetting GlobalSettings

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
	row := db.QueryRow(query)

	// scanning row and putting values into the globalSetting
	err := row.Scan(
		&globalSetting.SettingID,
		&globalSetting.AccountInfo.InsertedAt,
		&globalSetting.AccountInfo.UpdatedAt,
		&globalSetting.TeslaMateUnits.UnitsLength,
		&globalSetting.TeslaMateUnits.UnitsTemperature,
		&globalSetting.TeslaMateGUI.PreferredRange,
		&globalSetting.TeslaMateGUI.Language,
		&globalSetting.TeslaMateURLs.BaseURL,
		&globalSetting.TeslaMateURLs.GrafanaURL,
	)

	switch err {
	case sql.ErrNoRows:
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPIGlobalsettingsV1", "No rows were returned!", err.Error())
		return
	case nil:
		// nothing wrong.. continuing
		break
	default:
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPIGlobalsettingsV1", CarsGlobalsettingsError1, err.Error())
		return
	}

	// adjusting to timezone differences from UTC to be userspecific
	globalSetting.AccountInfo.InsertedAt = getTimeInTimeZone(globalSetting.AccountInfo.InsertedAt)
	globalSetting.AccountInfo.UpdatedAt = getTimeInTimeZone(globalSetting.AccountInfo.UpdatedAt)

	//
	// build the data-blob
	jsonData := JSONData{
		Data{
			GlobalSettings: globalSetting,
		},
	}

	// return jsonData
	TeslaMateAPIHandleSuccessResponse(c, "TeslaMateAPIGlobalsettingsV1", jsonData)
}
