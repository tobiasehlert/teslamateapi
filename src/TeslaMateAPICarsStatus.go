package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
)

// defining a lot of vars that are used my the MQTT MessageHandler
var (
	MQTTDataDisplayName                string
	MQTTDataState                      string
	MQTTDataStateSince                 string
	MQTTDataHealthy                    bool
	MQTTDataVersion                    string
	MQTTDataUpdateAvailable            bool
	MQTTDataUpdateVersion              string
	MQTTDataModel                      string
	MQTTDataTrimBadging                string
	MQTTDataExteriorColor              string
	MQTTDataWheelType                  string
	MQTTDataSpoilerType                string
	MQTTDataGeofence                   string
	MQTTDataLatitude                   float64
	MQTTDataLongitude                  float64
	MQTTDataShiftState                 string
	MQTTDataSpeed                      int
	MQTTDataHeading                    int
	MQTTDataElevation                  int
	MQTTDataLocked                     bool
	MQTTDataSentryMode                 bool
	MQTTDataWindowsOpen                bool
	MQTTDataDoorsOpen                  bool
	MQTTDataTrunkOpen                  bool
	MQTTDataFrunkOpen                  bool
	MQTTDataIsUserPresent              bool
	MQTTDataIsClimateOn                bool
	MQTTDataInsideTemp                 float64
	MQTTDataOutsideTemp                float64
	MQTTDataIsPreconditioning          bool
	MQTTDataOdometer                   float64
	MQTTDataEstBatteryRange            float64
	MQTTDataRatedBatteryRange          float64
	MQTTDataIdealBatteryRange          float64
	MQTTDataBatteryLevel               int
	MQTTDataUsableBatteryLevel         int
	MQTTDataPluggedIn                  bool
	MQTTDataChargeEnergyAdded          float64
	MQTTDataChargeLimitSoc             int
	MQTTDataChargePortDoorOpen         bool
	MQTTDataChargerActualCurrent       float64
	MQTTDataChargerPhases              int
	MQTTDataChargerPower               float64
	MQTTDataChargerVoltage             int
	MQTTDataScheduledChargingStartTime string
	MQTTDataTimeToFullCharge           float64
)

//define a function for the default message handler
var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

	// extracting the last part of topic
	_, MqttTopic := path.Split(msg.Topic())

	// running if-else statements to collect data and put into overall vars..
	if MqttTopic == "display_name" {
		MQTTDataDisplayName = string(msg.Payload())
	} else if MqttTopic == "state" {
		MQTTDataState = string(msg.Payload())
	} else if MqttTopic == "since" {
		MQTTDataStateSince = string(msg.Payload())
	} else if MqttTopic == "healthy" {
		MQTTDataHealthy = convertStringToBool(string(msg.Payload()))
	} else if MqttTopic == "version" {
		MQTTDataVersion = string(msg.Payload())
	} else if MqttTopic == "update_available" {
		MQTTDataUpdateAvailable = convertStringToBool(string(msg.Payload()))
	} else if MqttTopic == "update_version" {
		MQTTDataUpdateVersion = string(msg.Payload())
	} else if MqttTopic == "model" {
		MQTTDataModel = string(msg.Payload())
	} else if MqttTopic == "trim_badging" {
		MQTTDataTrimBadging = string(msg.Payload())
	} else if MqttTopic == "exterior_color" {
		MQTTDataExteriorColor = string(msg.Payload())
	} else if MqttTopic == "wheel_type" {
		MQTTDataWheelType = string(msg.Payload())
	} else if MqttTopic == "spoiler_type" {
		MQTTDataSpoilerType = string(msg.Payload())
	} else if MqttTopic == "geofence" {
		MQTTDataGeofence = string(msg.Payload())
	} else if MqttTopic == "latitude" {
		MQTTDataLatitude = convertStringToFloat(string(msg.Payload()))
	} else if MqttTopic == "longitude" {
		MQTTDataLongitude = convertStringToFloat(string(msg.Payload()))
	} else if MqttTopic == "shift_state" {
		MQTTDataShiftState = string(msg.Payload())
	} else if MqttTopic == "speed" {
		MQTTDataSpeed = convertStringToInteger(string(msg.Payload()))
	} else if MqttTopic == "heading" {
		MQTTDataHeading = convertStringToInteger(string(msg.Payload()))
	} else if MqttTopic == "elevation" {
		MQTTDataElevation = convertStringToInteger(string(msg.Payload()))
	} else if MqttTopic == "locked" {
		MQTTDataLocked = convertStringToBool(string(msg.Payload()))
	} else if MqttTopic == "sentry_mode" {
		MQTTDataSentryMode = convertStringToBool(string(msg.Payload()))
	} else if MqttTopic == "windows_open" {
		MQTTDataWindowsOpen = convertStringToBool(string(msg.Payload()))
	} else if MqttTopic == "doors_open" {
		MQTTDataDoorsOpen = convertStringToBool(string(msg.Payload()))
	} else if MqttTopic == "trunk_open" {
		MQTTDataTrunkOpen = convertStringToBool(string(msg.Payload()))
	} else if MqttTopic == "frunk_open" {
		MQTTDataFrunkOpen = convertStringToBool(string(msg.Payload()))
	} else if MqttTopic == "is_user_present" {
		MQTTDataIsUserPresent = convertStringToBool(string(msg.Payload()))
	} else if MqttTopic == "is_climate_on" {
		MQTTDataIsClimateOn = convertStringToBool(string(msg.Payload()))
	} else if MqttTopic == "inside_temp" {
		MQTTDataInsideTemp = convertStringToFloat(string(msg.Payload()))
	} else if MqttTopic == "outside_temp" {
		MQTTDataOutsideTemp = convertStringToFloat(string(msg.Payload()))
	} else if MqttTopic == "is_preconditioning" {
		MQTTDataIsPreconditioning = convertStringToBool(string(msg.Payload()))
	} else if MqttTopic == "odometer" {
		MQTTDataOdometer = convertStringToFloat(string(msg.Payload()))
	} else if MqttTopic == "est_battery_range_km" {
		MQTTDataEstBatteryRange = convertStringToFloat(string(msg.Payload()))
	} else if MqttTopic == "rated_battery_range_km" {
		MQTTDataRatedBatteryRange = convertStringToFloat(string(msg.Payload()))
	} else if MqttTopic == "ideal_battery_range_km" {
		MQTTDataIdealBatteryRange = convertStringToFloat(string(msg.Payload()))
	} else if MqttTopic == "battery_level" {
		MQTTDataBatteryLevel = convertStringToInteger(string(msg.Payload()))
	} else if MqttTopic == "usable_battery_level" {
		MQTTDataUsableBatteryLevel = convertStringToInteger(string(msg.Payload()))
	} else if MqttTopic == "plugged_in" {
		MQTTDataPluggedIn = convertStringToBool(string(msg.Payload()))
	} else if MqttTopic == "charge_energy_added" {
		MQTTDataChargeEnergyAdded = convertStringToFloat(string(msg.Payload()))
	} else if MqttTopic == "charge_limit_soc" {
		MQTTDataChargeLimitSoc = convertStringToInteger(string(msg.Payload()))
	} else if MqttTopic == "charge_port_door_open" {
		MQTTDataChargePortDoorOpen = convertStringToBool(string(msg.Payload()))
	} else if MqttTopic == "charger_actual_current" {
		MQTTDataChargerActualCurrent = convertStringToFloat(string(msg.Payload()))
	} else if MqttTopic == "charger_phases" {
		MQTTDataChargerPhases = convertStringToInteger(string(msg.Payload()))
	} else if MqttTopic == "charger_power" {
		MQTTDataChargerPower = convertStringToFloat(string(msg.Payload()))
	} else if MqttTopic == "charger_voltage" {
		MQTTDataChargerVoltage = convertStringToInteger(string(msg.Payload()))
	} else if MqttTopic == "scheduled_charging_start_time" {
		MQTTDataScheduledChargingStartTime = string(msg.Payload())
	} else if MqttTopic == "time_to_full_charge" {
		MQTTDataTimeToFullCharge = convertStringToFloat(string(msg.Payload()))
	} else {
		log.Printf("Error: extraction of data for %s not implemented in mqtt.MessageHandler yet.", MqttTopic)
	}

}

// TeslaMateAPICarsStatus func
func TeslaMateAPICarsStatus(CarID int) (string, bool) {

	// default values that get might get overwritten..
	MQTTPort := 0
	MQTTUserstring := ""
	MQTTProtocol := "tcp"

	// creating connection string towards mqtt
	MQTTTLS := getEnvAsBool("MQTT_TLS", false)
	if MQTTTLS {
		MQTTPort = getEnvAsInt("MQTT_PORT", 8883)
		MQTTProtocol = "tls"
	} else {
		MQTTPort = getEnvAsInt("MQTT_PORT", 1883)
	}
	MQTTHost := getEnv("MQTT_HOST", "mosquitto")
	MQTTUser := getEnv("MQTT_USERNAME", "")
	MQTTPass := getEnv("MQTT_PASSWORD", "")
	// MQTTInvCert := getEnvAsBool("MQTT_TLS_ACCEPT_INVALID_CERTS", false)
	MQTTNameSpace := getEnv("MQTT_NAMESPACE", "")

	// creating mqttURL to connect with
	// mqtt[s]://[username][:password]@host.domain[:port]
	if len(MQTTUser) > 0 {
		MQTTUserstring = MQTTUser
	}
	if len(MQTTPass) > 0 {
		MQTTUserstring = (MQTTUserstring + ":" + MQTTPass)
	}
	if len(MQTTUserstring) > 0 {
		MQTTUserstring = (MQTTUserstring + "@")
	}

	mqttURL := fmt.Sprintf("%s://%s%s:%d", MQTTProtocol, MQTTUserstring, MQTTHost, MQTTPort)

	// adding MQTTNameSpace info
	if len(MQTTNameSpace) > 0 {
		MQTTNameSpace = ("/" + MQTTNameSpace)
	}

	/*
		// if some more logging is needed.. which we skip for now
		mqtt.DEBUG = log.New(os.Stdout, "", 0)
		mqtt.ERROR = log.New(os.Stdout, "", 0)
	*/

	// create options for the MQTT client connection
	opts := mqtt.NewClientOptions().AddBroker(mqttURL).SetClientID("teslamateapi")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	// creating MQTT connection with options
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Panic(token.Error())
	}

	// showing mqtt successfully connected
	log.Println("successfully connected to mqtt.")

	// creating structs for /cars
	// BatteryDetails struct - child of MQTTInformation
	type BatteryDetails struct {
		EstBatteryRange    float64 `json:"est_battery_range"`    // 372.5 - Estimated Range in km
		RatedBatteryRange  float64 `json:"rated_battery_range"`  // 401.63 - Rated Range in km
		IdealBatteryRange  float64 `json:"ideal_battery_range"`  // 335.79 - Ideal Range in km
		BatteryLevel       int     `json:"battery_level"`        // 88 - Battery Level Percentage
		UsableBatteryLevel int     `json:"usable_battery_level"` // 85 - Usable battery level percentage
	}
	// CarDetails struct - child of MQTTInformation
	type CarDetails struct {
		Model       string `json:"model"`        // character varying(255)
		TrimBadging string `json:"trim_badging"` // P100D - Trim badging
	}
	// CarExterior struct - child of MQTTInformation
	type CarExterior struct {
		ExteriorColor string `json:"exterior_color"` // DeepBlue - The exterior color
		SpoilerType   string `json:"spoiler_type"`   // None - The spoiler type
		WheelType     string `json:"wheel_type"`     // Pinwheel18 - The wheel type
	}
	// CarGeodata struct - child of MQTTInformation
	type CarGeodata struct {
		Geofence  string  `json:"geofence"`  // Home - The name of the Geo-fence, if one exists at the current position
		Latitude  float64 `json:"latitude"`  // 35.278131 - Last reported car latitude
		Longitude float64 `json:"longitude"` // 29.744801 - Last reported car longitude
	}
	// CarStatus struct - child of MQTTInformation
	type CarStatus struct {
		Healthy       bool `json:"healthy"`         // true - Health status of the logger for that vehicle
		Locked        bool `json:"locked"`          // true - Indicates if the car is locked
		SentryMode    bool `json:"sentry_mode"`     // false - Indicates if Sentry Mode is active
		WindowsOpen   bool `json:"windows_open"`    // false - Indicates if any of the windows are open
		DoorsOpen     bool `json:"doors_open"`      // false - Indicates if any of the doors are open
		TrunkOpen     bool `json:"trunk_open"`      // false - Indicates if the trunk is open
		FrunkOpen     bool `json:"frunk_open"`      // false - Indicates if the frunk is open
		IsUserPresent bool `json:"is_user_present"` // false - Indicates if a user is present in the vehicle
	}
	// CarVersions struct - child of MQTTInformation
	type CarVersions struct {
		Version         string `json:"version"`          // 2019.32.12.2 - Software Version
		UpdateAvailable bool   `json:"update_available"` // false - Indicates if a software update is available
		UpdateVersion   string `json:"update_version"`   // 2019.32.12.3 - Software version of the available update
	}
	// ChargingDetails struct - child of MQTTInformation
	type ChargingDetails struct {
		PluggedIn                  bool    `json:"plugged_in"`                    // true - If car is currently plugged into a charger
		ChargeEnergyAdded          float64 `json:"charge_energy_added"`           // 5.06 - Last added energy in kWh
		ChargeLimitSoc             int     `json:"charge_limit_soc"`              // 90 - Charge Limit Configured in Percentage
		ChargePortDoorOpen         bool    `json:"charge_port_door_open"`         // true - Indicates if the charger door is open
		ChargerActualCurrent       float64 `json:"charger_actual_current"`        // 2.05 - Current amperage supplied by charger
		ChargerPhases              int     `json:"charger_phases"`                // 3 - Number of charger power phases (1-3)
		ChargerPower               float64 `json:"charger_power"`                 // 48.9 - Charger Power
		ChargerVoltage             int     `json:"charger_voltage"`               // 240 - Charger Voltage
		ScheduledChargingStartTime string  `json:"scheduled_charging_start_time"` // 2019-02-29T23:00:07Z - Start time of the scheduled charge
		TimeToFullCharge           float64 `json:"time_to_full_charge"`           // 1.83 - Hours remaining to full charge
	}
	// ClimateDetails struct - child of MQTTInformation
	type ClimateDetails struct {
		IsClimateOn       bool    `json:"is_climate_on"`      // true - Indicates if the climate control is on
		InsideTemp        float64 `json:"inside_temp"`        // 20.8 - Inside Temperature in °C
		OutsideTemp       float64 `json:"outside_temp"`       // 18.4 - Temperature in °C
		IsPreconditioning bool    `json:"is_preconditioning"` // false - Indicates if the vehicle is being preconditioned
	}
	// DrivingDetails struct - child of MQTTInformation
	type DrivingDetails struct {
		ShiftState string `json:"shift_state"` // D - Current/Last Shift State (D/N/R/P)
		Speed      int    `json:"speed"`       // 12 - Current Speed in km/h
		Heading    int    `json:"heading"`     // 340 - Last reported car direction
		Elevation  int    `json:"elevation"`   // 70 - Current elevation above sea level in meters
	}
	// MQTTInformation struct - child of Cars
	type MQTTInformation struct {
		DisplayName     string          `json:"display_name"`     // Blue Thunder - Vehicle Name
		State           string          `json:"state"`            // asleep - Status of the vehicle (e.g. online, asleep, charging)
		StateSince      string          `json:"state_since"`      // 2019-02-29T23:00:07Z - Date of the last status change
		Odometer        float64         `json:"odometer"`         // 1653 - Car odometer in km
		CarStatus       CarStatus       `json:"car_status"`       // struct
		CarDetails      CarDetails      `json:"car_details"`      // struct
		CarExterior     CarExterior     `json:"car_exterior"`     // struct
		CarGeodata      CarGeodata      `json:"car_geodata"`      // struct
		CarVersions     CarVersions     `json:"car_versions"`     // struct
		DrivingDetails  DrivingDetails  `json:"driving_details"`  // struct
		ClimateDetails  ClimateDetails  `json:"climate_details"`  // struct
		BatteryDetails  BatteryDetails  `json:"battery_details"`  // struct
		ChargingDetails ChargingDetails `json:"charging_details"` // struct
	}
	// Cars struct - child of Data
	type Car struct {
		ID   int    `json:"id"`   // smallint
		Name string `json:"name"` // text
	}
	// TeslaMateUnits struct - child of Data
	type TeslaMateUnits struct {
		UnitsLength      string `json:"unit_of_length"`      // string
		UnitsTemperature string `json:"unit_of_temperature"` // string
	}
	// Data struct - child of JSONData
	type Data struct {
		Car             Car             `json:"car"`
		MQTTInformation MQTTInformation `json:"status"`
		TeslaMateUnits  TeslaMateUnits  `json:"units"`
	}
	// JSONData struct - main
	type JSONData struct {
		Data Data `json:"data"`
	}

	// creating required vars
	var CarData Car
	var MQTTInformationData MQTTInformation
	var UnitsLength, UnitsTemperature string
	var ValidResponse bool // default is false

	// getting data from database
	query := `
		SELECT
			id,
			name,
			(SELECT unit_of_length FROM settings LIMIT 1) as unit_of_length,
			(SELECT unit_of_temperature FROM settings LIMIT 1) as unit_of_temperature
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

		// scanning row and putting values into the car
		err = rows.Scan(
			&CarData.ID,
			&CarData.Name,
			&UnitsLength,
			&UnitsTemperature,
		)

		// checking for errors after scanning
		if err != nil {
			log.Fatal(err)
		}

		if CarID != 0 && CarID == CarData.ID || CarID == 0 {

			// creating lots of subscribe to get values from every topic..
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/display_name", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/state", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/since", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/healthy", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/version", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/update_available", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/update_version", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/model", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/trim_badging", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/exterior_color", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/wheel_type", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/spoiler_type", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/geofence", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/latitude", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/longitude", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/shift_state", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/speed", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/heading", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/elevation", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/locked", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/sentry_mode", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/windows_open", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/doors_open", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/trunk_open", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/frunk_open", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/is_user_present", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/is_climate_on", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/inside_temp", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/outside_temp", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/is_preconditioning", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/odometer", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/est_battery_range_km", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/rated_battery_range_km", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/ideal_battery_range_km", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/battery_level", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/usable_battery_level", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/plugged_in", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/charge_energy_added", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/charge_limit_soc", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/charge_port_door_open", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/charger_actual_current", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/charger_phases", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/charger_power", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/charger_voltage", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/scheduled_charging_start_time", MQTTNameSpace, CarID), 0, nil)
			c.Subscribe(fmt.Sprintf("teslamate%s/cars/%d/time_to_full_charge", MQTTNameSpace, CarID), 0, nil)

			// adding some short sleep before unsubscribe
			time.Sleep(100 * time.Millisecond)

			// doing unsubscribe to not receive info anymore..
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/display_name", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/state", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/since", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/healthy", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/version", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/update_available", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/update_version", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/model", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/trim_badging", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/exterior_color", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/wheel_type", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/spoiler_type", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/geofence", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/latitude", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/longitude", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/shift_state", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/speed", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/heading", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/elevation", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/locked", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/sentry_mode", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/windows_open", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/doors_open", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/trunk_open", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/frunk_open", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/is_user_present", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/is_climate_on", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/inside_temp", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/outside_temp", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/is_preconditioning", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/odometer", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/est_battery_range_km", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/rated_battery_range_km", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/ideal_battery_range_km", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/battery_level", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/usable_battery_level", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/plugged_in", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/charge_energy_added", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/charge_limit_soc", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/charge_port_door_open", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/charger_actual_current", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/charger_phases", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/charger_power", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/charger_voltage", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/scheduled_charging_start_time", MQTTNameSpace, CarID))
			c.Unsubscribe(fmt.Sprintf("teslamate%s/cars/%d/time_to_full_charge", MQTTNameSpace, CarID))

			// disconnecting from MQTT
			c.Disconnect(250)

			// setting data from MQTT into data fields to return
			MQTTInformationData.DisplayName = MQTTDataDisplayName
			MQTTInformationData.State = MQTTDataState
			MQTTInformationData.StateSince = MQTTDataStateSince
			MQTTInformationData.CarStatus.Healthy = MQTTDataHealthy
			MQTTInformationData.CarVersions.Version = MQTTDataVersion
			MQTTInformationData.CarVersions.UpdateAvailable = MQTTDataUpdateAvailable
			MQTTInformationData.CarVersions.UpdateVersion = MQTTDataUpdateVersion
			MQTTInformationData.CarDetails.Model = MQTTDataModel
			MQTTInformationData.CarDetails.TrimBadging = MQTTDataTrimBadging
			MQTTInformationData.CarExterior.ExteriorColor = MQTTDataExteriorColor
			MQTTInformationData.CarExterior.WheelType = MQTTDataWheelType
			MQTTInformationData.CarExterior.SpoilerType = MQTTDataSpoilerType
			MQTTInformationData.CarGeodata.Geofence = MQTTDataGeofence
			MQTTInformationData.CarGeodata.Latitude = MQTTDataLatitude
			MQTTInformationData.CarGeodata.Longitude = MQTTDataLongitude
			MQTTInformationData.DrivingDetails.ShiftState = MQTTDataShiftState
			MQTTInformationData.DrivingDetails.Speed = MQTTDataSpeed
			MQTTInformationData.DrivingDetails.Heading = MQTTDataHeading
			MQTTInformationData.DrivingDetails.Elevation = MQTTDataElevation
			MQTTInformationData.CarStatus.Locked = MQTTDataLocked
			MQTTInformationData.CarStatus.SentryMode = MQTTDataSentryMode
			MQTTInformationData.CarStatus.WindowsOpen = MQTTDataWindowsOpen
			MQTTInformationData.CarStatus.DoorsOpen = MQTTDataDoorsOpen
			MQTTInformationData.CarStatus.TrunkOpen = MQTTDataTrunkOpen
			MQTTInformationData.CarStatus.FrunkOpen = MQTTDataFrunkOpen
			MQTTInformationData.CarStatus.IsUserPresent = MQTTDataIsUserPresent
			MQTTInformationData.ClimateDetails.IsClimateOn = MQTTDataIsClimateOn
			MQTTInformationData.ClimateDetails.InsideTemp = MQTTDataInsideTemp
			MQTTInformationData.ClimateDetails.OutsideTemp = MQTTDataOutsideTemp
			MQTTInformationData.ClimateDetails.IsPreconditioning = MQTTDataIsPreconditioning
			MQTTInformationData.Odometer = MQTTDataOdometer
			MQTTInformationData.BatteryDetails.EstBatteryRange = MQTTDataEstBatteryRange
			MQTTInformationData.BatteryDetails.RatedBatteryRange = MQTTDataRatedBatteryRange
			MQTTInformationData.BatteryDetails.IdealBatteryRange = MQTTDataIdealBatteryRange
			MQTTInformationData.BatteryDetails.BatteryLevel = MQTTDataBatteryLevel
			MQTTInformationData.BatteryDetails.UsableBatteryLevel = MQTTDataUsableBatteryLevel
			MQTTInformationData.ChargingDetails.PluggedIn = MQTTDataPluggedIn
			MQTTInformationData.ChargingDetails.ChargeEnergyAdded = MQTTDataChargeEnergyAdded
			MQTTInformationData.ChargingDetails.ChargeLimitSoc = MQTTDataChargeLimitSoc
			MQTTInformationData.ChargingDetails.ChargePortDoorOpen = MQTTDataChargePortDoorOpen
			MQTTInformationData.ChargingDetails.ChargerActualCurrent = MQTTDataChargerActualCurrent
			MQTTInformationData.ChargingDetails.ChargerPhases = MQTTDataChargerPhases
			MQTTInformationData.ChargingDetails.ChargerPower = MQTTDataChargerPower
			MQTTInformationData.ChargingDetails.ChargerVoltage = MQTTDataChargerVoltage
			MQTTInformationData.ChargingDetails.ScheduledChargingStartTime = MQTTDataScheduledChargingStartTime
			MQTTInformationData.ChargingDetails.TimeToFullCharge = MQTTDataTimeToFullCharge

			// converting values based of settings UnitsLength
			if UnitsLength == "mi" {
				// drive.OdometerDetails.OdometerStart = kilometersToMiles(drive.OdometerDetails.OdometerStart)
				MQTTInformationData.Odometer = kilometersToMiles(MQTTInformationData.Odometer)
				MQTTInformationData.BatteryDetails.EstBatteryRange = kilometersToMiles(MQTTInformationData.BatteryDetails.EstBatteryRange)
				MQTTInformationData.BatteryDetails.RatedBatteryRange = kilometersToMiles(MQTTInformationData.BatteryDetails.RatedBatteryRange)
				MQTTInformationData.BatteryDetails.IdealBatteryRange = kilometersToMiles(MQTTInformationData.BatteryDetails.IdealBatteryRange)
			}
			// converting values based of settings UnitsTemperature
			if UnitsTemperature == "F" {
				MQTTInformationData.ClimateDetails.InsideTemp = celsiusToFahrenheit(MQTTInformationData.ClimateDetails.InsideTemp)
				MQTTInformationData.ClimateDetails.OutsideTemp = celsiusToFahrenheit(MQTTInformationData.ClimateDetails.OutsideTemp)
			}

			// setting response as valid
			ValidResponse = true
		}
	}

	// checking for errors in the rows result
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	//
	// build the data-blob
	jsonData := JSONData{
		Data{
			Car:             CarData,
			MQTTInformation: MQTTInformationData,
			TeslaMateUnits: TeslaMateUnits{
				UnitsLength:      UnitsLength,
				UnitsTemperature: UnitsTemperature,
			},
		},
	}

	// print readable output to log
	log.Printf("data for /cars/%d/status created:", CarID)

	js, _ := json.Marshal(jsonData)
	fmt.Printf("%s\n", js)
	return string(js), ValidResponse
}
