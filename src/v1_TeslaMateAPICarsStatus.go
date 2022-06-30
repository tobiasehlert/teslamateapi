package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/thanhpk/randstr"
)

// statusInfo holds the status info for a car
type statusInfo struct {
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
	MQTTDataPower                      int
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
	MQTTDataChargeCurrentRequest       int
	MQTTDataChargeCurrentRequestMax    int
	MQTTDataScheduledChargingStartTime string
	MQTTDataTimeToFullCharge           float64
	MQTTDataTpmsPressureFL             float64
	MQTTDataTpmsPressureFR             float64
	MQTTDataTpmsPressureRL             float64
	MQTTDataTpmsPressureRR             float64
}

type statusCache struct {
	mqttDisabled bool

	topicScan string // scan parameter (expect it to generate car ID then relevant parameter)

	cache map[int]*statusInfo
	mu    sync.Mutex
}

func startMQTT() (*statusCache, error) {
	s := statusCache{
		cache: make(map[int]*statusInfo),
	}
	// getting mqtt flag
	s.mqttDisabled = getEnvAsBool("DISABLE_MQTT", false)
	if s.mqttDisabled {
		return nil, errors.New("[notice] TeslaMateAPICarsStatusV1 DISABLE_MQTT is set to true.. can not return status for car without mqtt")
	}

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

	// create options for the MQTT client connection
	opts := mqtt.NewClientOptions().AddBroker(mqttURL)
	// setting generic MQTT settings in opts
	opts.SetKeepAlive(2 * time.Second)                    // setting keepalive for client
	opts.SetDefaultPublishHandler(s.newMessage)           // using f mqtt.MessageHandler function
	opts.SetPingTimeout(1 * time.Second)                  // setting pingtimeout for client
	opts.SetClientID("teslamateapi-" + randstr.String(4)) // setting mqtt client id for TeslaMateApi
	opts.SetCleanSession(true)                            // removal of all subscriptions on disconnect
	opts.SetOrderMatters(false)                           // don't care about order (removes need for callbacks to return immediately)
	opts.SetAutoReconnect(true)                           // if connection drops automatically re-establish it
	opts.AutoReconnect = true

	// creating MQTT connection with options
	m := mqtt.NewClient(opts)
	if token := m.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("[error] TeslaMateAPICarsStatusV1 failed to connect to MQTT: %w", token.Error())
		// Note : May want to use opts.ConnectRetry which will keep trying the connection
	}

	// showing mqtt successfully connected
	if gin.IsDebugging() {
		log.Println("[debug] TeslaMateAPICarsStatusV1 successfully connected to mqtt.")
	}

	// adding MQTTNameSpace info
	MQTTNameSpace := getEnv("MQTT_NAMESPACE", "")
	if len(MQTTNameSpace) > 0 {
		MQTTNameSpace = ("/" + MQTTNameSpace)
	}

	// Subscribe - we will accept info on any car...
	topic := fmt.Sprintf("teslamate%s/cars/#", MQTTNameSpace)
	if token := m.Subscribe(topic, 0, s.newMessage); token.Wait() && token.Error() != nil {
		log.Panic(token.Error()) // Note : May want to use opts.ConnectRetry which will keep trying the connection
	}
	s.topicScan = fmt.Sprintf("teslamate%s/cars/%%d/%%s", MQTTNameSpace)

	// Thats all - newMessage will be called when something new arrives
	return &s, nil
}

// newMessage - called by mqtt package when new message received
func (s *statusCache) newMessage(c mqtt.Client, msg mqtt.Message) {
	// topic is in the format teslamateMQTT_NAMESPACE/cars/carID/display_name
	var (
		carID     int
		MqttTopic string
	)
	_, err := fmt.Sscanf(msg.Topic(), s.topicScan, &carID, &MqttTopic)
	if err != nil {
		log.Printf("[warning] TeslaMateAPICarsStatusV1 unexpected topic format (%s) - ignoring message: %v", msg.Topic(), err)
		return
	}

	// extracting the last part of topic
	s.mu.Lock()
	defer s.mu.Unlock()
	stat := s.cache[carID]
	if stat == nil {
		stat = &statusInfo{}
		s.cache[carID] = stat
	}

	// running if-else statements to collect data and put into overall vars..
	switch MqttTopic {
	case "display_name":
		stat.MQTTDataDisplayName = string(msg.Payload())
	case "state":
		stat.MQTTDataState = string(msg.Payload())
	case "since":
		stat.MQTTDataStateSince = string(msg.Payload())
	case "healthy":
		stat.MQTTDataHealthy = convertStringToBool(string(msg.Payload()))
	case "version":
		stat.MQTTDataVersion = string(msg.Payload())
	case "update_available":
		stat.MQTTDataUpdateAvailable = convertStringToBool(string(msg.Payload()))
	case "update_version":
		stat.MQTTDataUpdateVersion = string(msg.Payload())
	case "model":
		stat.MQTTDataModel = string(msg.Payload())
	case "trim_badging":
		stat.MQTTDataTrimBadging = string(msg.Payload())
	case "exterior_color":
		stat.MQTTDataExteriorColor = string(msg.Payload())
	case "wheel_type":
		stat.MQTTDataWheelType = string(msg.Payload())
	case "spoiler_type":
		stat.MQTTDataSpoilerType = string(msg.Payload())
	case "geofence":
		stat.MQTTDataGeofence = string(msg.Payload())
	case "latitude":
		stat.MQTTDataLatitude = convertStringToFloat(string(msg.Payload()))
	case "longitude":
		stat.MQTTDataLongitude = convertStringToFloat(string(msg.Payload()))
	case "shift_state":
		stat.MQTTDataShiftState = string(msg.Payload())
	case "power":
		stat.MQTTDataPower = convertStringToInteger(string(msg.Payload()))
	case "speed":
		stat.MQTTDataSpeed = convertStringToInteger(string(msg.Payload()))
	case "heading":
		stat.MQTTDataHeading = convertStringToInteger(string(msg.Payload()))
	case "elevation":
		stat.MQTTDataElevation = convertStringToInteger(string(msg.Payload()))
	case "locked":
		stat.MQTTDataLocked = convertStringToBool(string(msg.Payload()))
	case "sentry_mode":
		stat.MQTTDataSentryMode = convertStringToBool(string(msg.Payload()))
	case "windows_open":
		stat.MQTTDataWindowsOpen = convertStringToBool(string(msg.Payload()))
	case "doors_open":
		stat.MQTTDataDoorsOpen = convertStringToBool(string(msg.Payload()))
	case "trunk_open":
		stat.MQTTDataTrunkOpen = convertStringToBool(string(msg.Payload()))
	case "frunk_open":
		stat.MQTTDataFrunkOpen = convertStringToBool(string(msg.Payload()))
	case "is_user_present":
		stat.MQTTDataIsUserPresent = convertStringToBool(string(msg.Payload()))
	case "is_climate_on":
		stat.MQTTDataIsClimateOn = convertStringToBool(string(msg.Payload()))
	case "inside_temp":
		stat.MQTTDataInsideTemp = convertStringToFloat(string(msg.Payload()))
	case "outside_temp":
		stat.MQTTDataOutsideTemp = convertStringToFloat(string(msg.Payload()))
	case "is_preconditioning":
		stat.MQTTDataIsPreconditioning = convertStringToBool(string(msg.Payload()))
	case "odometer":
		stat.MQTTDataOdometer = convertStringToFloat(string(msg.Payload()))
	case "est_battery_range_km":
		stat.MQTTDataEstBatteryRange = convertStringToFloat(string(msg.Payload()))
	case "rated_battery_range_km":
		stat.MQTTDataRatedBatteryRange = convertStringToFloat(string(msg.Payload()))
	case "ideal_battery_range_km":
		stat.MQTTDataIdealBatteryRange = convertStringToFloat(string(msg.Payload()))
	case "battery_level":
		stat.MQTTDataBatteryLevel = convertStringToInteger(string(msg.Payload()))
	case "usable_battery_level":
		stat.MQTTDataUsableBatteryLevel = convertStringToInteger(string(msg.Payload()))
	case "plugged_in":
		stat.MQTTDataPluggedIn = convertStringToBool(string(msg.Payload()))
	case "charge_energy_added":
		stat.MQTTDataChargeEnergyAdded = convertStringToFloat(string(msg.Payload()))
	case "charge_limit_soc":
		stat.MQTTDataChargeLimitSoc = convertStringToInteger(string(msg.Payload()))
	case "charge_port_door_open":
		stat.MQTTDataChargePortDoorOpen = convertStringToBool(string(msg.Payload()))
	case "charger_actual_current":
		stat.MQTTDataChargerActualCurrent = convertStringToFloat(string(msg.Payload()))
	case "charger_phases":
		stat.MQTTDataChargerPhases = convertStringToInteger(string(msg.Payload()))
	case "charger_power":
		stat.MQTTDataChargerPower = convertStringToFloat(string(msg.Payload()))
	case "charger_voltage":
		stat.MQTTDataChargerVoltage = convertStringToInteger(string(msg.Payload()))
	case "charge_current_request":
		stat.MQTTDataChargeCurrentRequest = convertStringToInteger(string(msg.Payload()))
	case "charge_current_request_max":
		stat.MQTTDataChargeCurrentRequestMax = convertStringToInteger(string(msg.Payload()))
	case "scheduled_charging_start_time":
		stat.MQTTDataScheduledChargingStartTime = string(msg.Payload())
	case "time_to_full_charge":
		stat.MQTTDataTimeToFullCharge = convertStringToFloat(string(msg.Payload()))
	case "tpms_pressure_fl":
		stat.MQTTDataTpmsPressureFL = convertStringToFloat(string(msg.Payload()))
	case "tpms_pressure_fr":
		stat.MQTTDataTpmsPressureFR = convertStringToFloat(string(msg.Payload()))
	case "tpms_pressure_rl":
		stat.MQTTDataTpmsPressureRL = convertStringToFloat(string(msg.Payload()))
	case "tpms_pressure_rr":
		stat.MQTTDataTpmsPressureRR = convertStringToFloat(string(msg.Payload()))
	default:
		log.Printf("[warning] TeslaMateAPICarsStatusV1 mqtt.MessageHandler issue.. extraction of data for %s not implemented!", MqttTopic)
	}
}

// TeslaMateAPICarsStatusV1 func
func (s *statusCache) TeslaMateAPICarsStatusV1(c *gin.Context) {
	if s.mqttDisabled {
		log.Println("[notice] TeslaMateAPICarsStatusV1 DISABLE_MQTT is set to true.. can not return status for car without mqtt!")
		TeslaMateAPIHandleOtherResponse(c, http.StatusNotImplemented, "TeslaMateAPICarsStatusV1", gin.H{"error": "mqtt disabled.. status not accessible!"})
		return
	}

	// getting CarID param from URL
	carID := convertStringToInteger(c.Param("CarID"))

	// Now see what data we have on the car
	s.mu.Lock()
	stat := s.cache[carID]
	s.mu.Unlock()

	if stat == nil {
		// or should it be http.StatusNoContent instead?
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsStatusV1", "no info on this car ID", "-")
		return
	}

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
		ChargeCurrentRequest       int     `json:"charge_current_request"`        // 40 - How many amps the car wants
		ChargeCurrentRequestMax    int     `json:"charge_current_request_max"`    // 40 - How many amps the car can have
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
		Power      int    `json:"power"`       // -9 Current battery power in watts. Positive value on discharge, negative value on charge
		Speed      int    `json:"speed"`       // 12 - Current Speed in km/h
		Heading    int    `json:"heading"`     // 340 - Last reported car direction
		Elevation  int    `json:"elevation"`   // 70 - Current elevation above sea level in meters
	}
	// TpmsDetails struct - child of MQTTInformatiojn
	type TpmsDetails struct {
		TpmsPressureFL float64 `json:"tpms_pressure_fl"` // 2.9 - Tire pressure measure in BAR, front left tire
		TpmsPressureFR float64 `json:"tpms_pressure_fr"` // 2.8 - Tire pressure measure in BAR, front right tire
		TpmsPressureRL float64 `json:"tpms_pressure_rl"` // 2.9 - Tire pressure measure in BAR, rear left tire
		TpmsPressureRR float64 `json:"tpms_pressure_rr"` // 2.8 - Tire pressure measure in BAR, rear right tire
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
		TpmsDetails     TpmsDetails     `json:"tpms_details"`     // struct
	}
	// Cars struct - child of Data
	type Car struct {
		CarID   int    `json:"car_id"`   // smallint
		CarName string `json:"car_name"` // text
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
	var (
		CarData                       Car
		MQTTInformationData           MQTTInformation
		UnitsLength, UnitsTemperature string
	)

	// getting data from database (assume that carID is unique!)
	query := `
		SELECT
			id,
			name,
			(SELECT unit_of_length FROM settings LIMIT 1) as unit_of_length,
			(SELECT unit_of_temperature FROM settings LIMIT 1) as unit_of_temperature
		FROM cars
		WHERE id=$1
		LIMIT 1;`
	err := db.QueryRow(query, carID).Scan(&CarData.CarID,
		&CarData.CarName,
		&UnitsLength,
		&UnitsTemperature)

	// checking for errors in query (this will include no rows found)
	if err != nil {
		TeslaMateAPIHandleErrorResponse(c, "TeslaMateAPICarsStatusV1", "Unable to load cars.", err.Error())
		return
	}

	// setting data from MQTT into data fields to return
	MQTTInformationData.DisplayName = stat.MQTTDataDisplayName
	MQTTInformationData.State = stat.MQTTDataState
	MQTTInformationData.StateSince = stat.MQTTDataStateSince
	MQTTInformationData.CarStatus.Healthy = stat.MQTTDataHealthy
	MQTTInformationData.CarVersions.Version = stat.MQTTDataVersion
	MQTTInformationData.CarVersions.UpdateAvailable = stat.MQTTDataUpdateAvailable
	MQTTInformationData.CarVersions.UpdateVersion = stat.MQTTDataUpdateVersion
	MQTTInformationData.CarDetails.Model = stat.MQTTDataModel
	MQTTInformationData.CarDetails.TrimBadging = stat.MQTTDataTrimBadging
	MQTTInformationData.CarExterior.ExteriorColor = stat.MQTTDataExteriorColor
	MQTTInformationData.CarExterior.WheelType = stat.MQTTDataWheelType
	MQTTInformationData.CarExterior.SpoilerType = stat.MQTTDataSpoilerType
	MQTTInformationData.CarGeodata.Geofence = stat.MQTTDataGeofence
	MQTTInformationData.CarGeodata.Latitude = stat.MQTTDataLatitude
	MQTTInformationData.CarGeodata.Longitude = stat.MQTTDataLongitude
	MQTTInformationData.DrivingDetails.ShiftState = stat.MQTTDataShiftState
	MQTTInformationData.DrivingDetails.Power = stat.MQTTDataPower
	MQTTInformationData.DrivingDetails.Speed = stat.MQTTDataSpeed
	MQTTInformationData.DrivingDetails.Heading = stat.MQTTDataHeading
	MQTTInformationData.DrivingDetails.Elevation = stat.MQTTDataElevation
	MQTTInformationData.CarStatus.Locked = stat.MQTTDataLocked
	MQTTInformationData.CarStatus.SentryMode = stat.MQTTDataSentryMode
	MQTTInformationData.CarStatus.WindowsOpen = stat.MQTTDataWindowsOpen
	MQTTInformationData.CarStatus.DoorsOpen = stat.MQTTDataDoorsOpen
	MQTTInformationData.CarStatus.TrunkOpen = stat.MQTTDataTrunkOpen
	MQTTInformationData.CarStatus.FrunkOpen = stat.MQTTDataFrunkOpen
	MQTTInformationData.CarStatus.IsUserPresent = stat.MQTTDataIsUserPresent
	MQTTInformationData.ClimateDetails.IsClimateOn = stat.MQTTDataIsClimateOn
	MQTTInformationData.ClimateDetails.InsideTemp = stat.MQTTDataInsideTemp
	MQTTInformationData.ClimateDetails.OutsideTemp = stat.MQTTDataOutsideTemp
	MQTTInformationData.ClimateDetails.IsPreconditioning = stat.MQTTDataIsPreconditioning
	MQTTInformationData.Odometer = stat.MQTTDataOdometer
	MQTTInformationData.BatteryDetails.EstBatteryRange = stat.MQTTDataEstBatteryRange
	MQTTInformationData.BatteryDetails.RatedBatteryRange = stat.MQTTDataRatedBatteryRange
	MQTTInformationData.BatteryDetails.IdealBatteryRange = stat.MQTTDataIdealBatteryRange
	MQTTInformationData.BatteryDetails.BatteryLevel = stat.MQTTDataBatteryLevel
	MQTTInformationData.BatteryDetails.UsableBatteryLevel = stat.MQTTDataUsableBatteryLevel
	MQTTInformationData.ChargingDetails.PluggedIn = stat.MQTTDataPluggedIn
	MQTTInformationData.ChargingDetails.ChargeEnergyAdded = stat.MQTTDataChargeEnergyAdded
	MQTTInformationData.ChargingDetails.ChargeLimitSoc = stat.MQTTDataChargeLimitSoc
	MQTTInformationData.ChargingDetails.ChargePortDoorOpen = stat.MQTTDataChargePortDoorOpen
	MQTTInformationData.ChargingDetails.ChargerActualCurrent = stat.MQTTDataChargerActualCurrent
	MQTTInformationData.ChargingDetails.ChargerPhases = stat.MQTTDataChargerPhases
	MQTTInformationData.ChargingDetails.ChargerPower = stat.MQTTDataChargerPower
	MQTTInformationData.ChargingDetails.ChargerVoltage = stat.MQTTDataChargerVoltage
	MQTTInformationData.ChargingDetails.ChargeCurrentRequest = stat.MQTTDataChargeCurrentRequest
	MQTTInformationData.ChargingDetails.ChargeCurrentRequestMax = stat.MQTTDataChargeCurrentRequestMax
	MQTTInformationData.ChargingDetails.ScheduledChargingStartTime = stat.MQTTDataScheduledChargingStartTime
	MQTTInformationData.ChargingDetails.TimeToFullCharge = stat.MQTTDataTimeToFullCharge
	MQTTInformationData.TpmsDetails.TpmsPressureFL = stat.MQTTDataTpmsPressureFL
	MQTTInformationData.TpmsDetails.TpmsPressureFR = stat.MQTTDataTpmsPressureFR
	MQTTInformationData.TpmsDetails.TpmsPressureRL = stat.MQTTDataTpmsPressureRL
	MQTTInformationData.TpmsDetails.TpmsPressureRR = stat.MQTTDataTpmsPressureRR

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

	// adjusting to timezone differences from UTC to be userspecific
	MQTTInformationData.StateSince = getTimeInTimeZone(MQTTInformationData.StateSince)
	MQTTInformationData.ChargingDetails.ScheduledChargingStartTime = getTimeInTimeZone(MQTTInformationData.ChargingDetails.ScheduledChargingStartTime)

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

	// return jsonData
	TeslaMateAPIHandleSuccessResponse(c, "TeslaMateAPICarsStatusV1", jsonData)
}
