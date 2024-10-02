package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/google/uuid"
	"github.com/Dartmouth-OpenAV/microservice-framework/framework"
)

// HELPER FUNCTIONS
func convertAndSend(socketKey string, method string, params interface{}) bool {
	uniqueID, err := uuid.NewRandom()
	if err != nil {
		errMsg := "k3n5oe - Error generating a unique id"
		framework.AddToErrors(socketKey, errMsg)
	}
	data := map[string]interface{}{"jsonrpc": "2.0", "id": uniqueID, "method": method, "params": params}

	encodedData, err := json.Marshal(data)

	nullValue := []byte("\x00")
	encodedData = append(encodedData, nullValue...)
	if err != nil {
		errMsg := "j3knxk - Error marshaling data"
		framework.AddToErrors(socketKey, errMsg)
	}

	sent := framework.WriteLineToSocket(socketKey, string(encodedData))

	return sent
}

func readAndConvert(socketKey string) (map[string]interface{}, error) {
	response := framework.ReadLineFromSocket(socketKey)
	response = strings.Trim(response, "\x00")

	if response == "" {
		errMsg := "45h3dr - Response was blank "
		framework.AddToErrors(socketKey, errMsg)
		return nil, errors.New(errMsg)
	}

	var data map[string]interface{}
	err := json.Unmarshal([]byte(response), &data)

	if err != nil {
		errMsg := "jxneh3 - Error unmarshaling the json response: " + err.Error()
		framework.AddToErrors(socketKey, errMsg)
	}

	_, methodPresent := data["method"]
	if methodPresent == true {
		if data["method"] == "EngineStatus" {
			response = framework.ReadLineFromSocket(socketKey)
			response = strings.Trim(response, "\x00")

			if response == "" {
				errMsg := "5i3id - Response was blank "
				framework.AddToErrors(socketKey, errMsg)
				return nil, errors.New(errMsg)
			}
			data = nil
			err := json.Unmarshal([]byte(response), &data)

			if err != nil {
				errMsg := "4j2jd - Error unmarshaling the json response: " + err.Error()
				framework.AddToErrors(socketKey, errMsg)
			}
		}
	}

	_, errorPresent := data["error"]
	if errorPresent {
		errorMap := data["error"].(map[string]interface{})
		errorMessage := errorMap["message"].(string)
		errorCode := errorMap["code"].(float64)
		fullErrorMessage := "Error Code: " + fmt.Sprint(errorCode) + ", Error Message: " + errorMessage
		framework.AddToErrors(socketKey, fullErrorMessage)
		return nil, errors.New(fullErrorMessage)
	}

	return data, nil
}

// SET FUNCTIONS
func setVolume(socketKey string, parameterName string, volume string) (string, error) {
	function := "setVolume"

	value := "notok"
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = setVolumeDo(socketKey, parameterName, volume)
		if value != "ok" { // Something went wrong - perhaps try again
			framework.Log(function + " - nw3njd retrying volume operation")
			maxRetries--
			time.Sleep(1 * time.Second)
			if maxRetries == 0 {
				errMsg := fmt.Sprintf(function + "du4hnd - max retries reached")
				framework.AddToErrors(socketKey, errMsg)
			}
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func setVolumeDo(socketKey string, parameterName string, volume string) (string, error) {
	method := "Control.Set"
	function := "setVolumeDo"
	volume = strings.Trim(volume, "\"")

	floatVolume, err := strconv.ParseFloat(volume, 64)

	if err != nil {
		errMsg := function + " - 493kjd - Error converting volume to float" + err.Error()
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, err
	}

	position := floatVolume / 100
	params := map[string]interface{}{"Name": parameterName, "Position": position}

	sent := convertAndSend(socketKey, method, params)

	if sent != true {
		errMsg := fmt.Sprintf(function + " - l5m2e - error sending command")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	response, err := readAndConvert(socketKey)

	if err != nil {
		errMsg := function + " - 1h54h - error reading response: " + err.Error()
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	framework.Log(function + "- Response: " + fmt.Sprint(response))

	// If we got here, the response was good, so successful return with the state indication
	return "ok", nil
}

func setToggle(socketKey string, parameterName string, state string) (string, error) {
	function := "setToggle"

	value := "notok"
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = setToggleDo(socketKey, parameterName, state)
		if value != "ok" { // Something went wrong - perhaps try again
			framework.Log(function + " - k4ndn3 retrying toggle operation")
			maxRetries--
			time.Sleep(1 * time.Second)
			if maxRetries == 0 {
				errMsg := fmt.Sprintf(function + "hxh5b3 - max retries reached")
				framework.AddToErrors(socketKey, errMsg)
			}
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func setToggleDo(socketKey string, parameterName string, state string) (string, error) {
	method := "Control.Set"
	function := "setToggleDo"
	state = strings.Trim(state, "\"")
	state = strings.ToLower(state)
	toggleValue := 0

	if state == "true" {
		toggleValue = 1
	} else if state == "false" {
		toggleValue = 0
	}

	params := map[string]interface{}{"Name": parameterName, "Value": toggleValue}

	sent := convertAndSend(socketKey, method, params)

	if sent != true {
		errMsg := fmt.Sprintf(function + " - qu2bs4 - error sending command for " + parameterName)
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	response, err := readAndConvert(socketKey)

	if err != nil {
		errMsg := fmt.Sprintf(function + " - wen42d - error reading response: " + err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	framework.Log(function + " " + parameterName + "- Response: " + fmt.Sprint(response))

	// If we got here, the response was good, so successful return with the state indication
	return "ok", nil
}

// paramterName is something like "decoder-1_output-2"
// input (arg2) is something like "2"
func setVideoRoute(socketKey string, parameterName string, arg2 string) (string, error) {
	function := "setVideoRoute"

	value := "notok"
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = setVideoRouteDo(socketKey, parameterName, arg2)
		if value != "ok" { // Something went wrong - perhaps try again
			framework.Log(function + " - k4ndn3 retrying setVideoRoute operation")
			maxRetries--
			time.Sleep(1 * time.Second)
			if maxRetries == 0 {
				errMsg := fmt.Sprintf(function + "hxh5b3z - max retries reached")
				framework.AddToErrors(socketKey, errMsg)
			}
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

// For setVideoRouteDo, the arguments depend on the Q-Sys component control name,
// which is something like "hdmi.out.1.select.avh.3.led"
// Thus paramterName in this case is "decoder-1_hdmi.out.1" where 'decoder-1' is Device Name and 'hdmi.out.1' is output
// input (arg2) is in this case "avh.3.led"
func setVideoRouteDo(socketKey string, parameterName string, arg2 string) (string, error) {
	// the string in parameterName following '_' is the 'output'
	// arg2 is the 'input'
	var output string
	var input string
	var outputDevice string

	parameterParts := strings.Split(parameterName, "_")
	outputDevice = parameterParts[0]
	output = parameterParts[1]
	input = arg2

	method := "Component.Set"
	function := "setVideoRouteDo"
	input = strings.Trim(input, "\"")
	inputChannel := output + ".select." + input
	channelDictionary := map[string]interface{}{"Name": inputChannel, "Value": 1}
	channelArray := [1]interface{}{channelDictionary}

	params := map[string]interface{}{"Name": outputDevice, "Controls": channelArray}

	sent := convertAndSend(socketKey, method, params)

	if sent != true {
		errMsg := fmt.Sprintf(function + " - qu2bs4z - error sending command for " + parameterName)
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	response, err := readAndConvert(socketKey)

	if err != nil {
		errMsg := fmt.Sprintf(function + " - wen42dz - error reading response: " + err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	framework.Log(function + " " + parameterName + "- Response: " + fmt.Sprint(response))

	// If we got here, the response was good, so successful return with the state indication
	return "ok", nil
}

//GET FUNCTIONS

func getVolume(socketKey string, parameterName string) (string, error) {
	function := "getVolume"

	value := `"unknown"`
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = getVolumeDo(socketKey, parameterName)
		if value == `"unknown"` { // Something went wrong - perhaps try again
			framework.Log(function + " - fq3sdvc retrying volume operation")
			maxRetries--
			time.Sleep(1 * time.Second)
			if maxRetries == 0 {
				errMsg := fmt.Sprintf(function + "f839dk4 - max retries reached")
				framework.AddToErrors(socketKey, errMsg)
			}
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func getVolumeDo(socketKey string, parameterName string) (string, error) {
	method := "Control.Get"
	function := "getVolumeDo"

	params := []string{parameterName}

	sent := convertAndSend(socketKey, method, params)

	if sent != true {
		errMsg := fmt.Sprintf(function + " - h2j3nd- error sending command")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	response, err := readAndConvert(socketKey)
	if err != nil {
		errMsg := fmt.Sprintf(function + " - 1b3dj - error reading response: " + err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	value := "unknown"

	_, resultPresent := response["result"]
	if resultPresent == true {
		result := response["result"].([]interface{})
		resultSlice := result[0].(map[string]interface{})
		_, valuePresent := resultSlice["Position"]
		if valuePresent == true {
			valueData := resultSlice["Position"].(float64) * 100
			value = strconv.FormatFloat(valueData, 'f', -1, 64)
		}
	}

	// If we got here, the response was good, so successful return with the state indication
	return value, nil
}

func getToggle(socketKey string, parameterName string) (string, error) {
	function := "getToggle"

	value := `"unknown"`
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = getToggleDo(socketKey, parameterName)
		if value == `"unknown"` { // Something went wrong - perhaps try again
			framework.Log(function + " - 34rd3i retrying toggle operation")
			maxRetries--
			time.Sleep(1 * time.Second)
			if maxRetries == 0 {
				errMsg := fmt.Sprintf(function + "sh4hd3 - max retries reached")
				framework.AddToErrors(socketKey, errMsg)
			}
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func getToggleDo(socketKey string, parameterName string) (string, error) {
	method := "Control.Get"
	function := "getToggleDo"

	params := []string{parameterName}

	sent := convertAndSend(socketKey, method, params)

	if sent != true {
		errMsg := fmt.Sprintf(function + " - 94dk34 - error sending command for " + parameterName)
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	response, err := readAndConvert(socketKey)
	if err != nil {
		errMsg := fmt.Sprintf(function + " - 4udj4 - error reading response: " + err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	value := "unknown"
	var valueData float64
	_, resultPresent := response["result"]
	if resultPresent == true {
		result := response["result"].([]interface{})
		resultSlice := result[0].(map[string]interface{})
		_, valuePresent := resultSlice["Value"]
		if valuePresent == true {
			valueData = resultSlice["Value"].(float64)
		}
	}

	if valueData == 1 {
		value = "true"
	} else if valueData == 0 {
		value = "false"
	} else {
		errMsg := function + " - result value for " + parameterName + " was not 0 or 1"
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	// If we got here, the response was good, so successful return with the state indication
	return `"` + value + `"`, nil
}

// parameterName is something like "decoder-1_hdmi.out.1"
func getVideoRoute(socketKey string, parameterName string) (string, error) {
	function := "getVideoRoute"

	value := `"unknown"`
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = getVideoRouteDo(socketKey, parameterName)
		if value == `"unknown"` { // Something went wrong - perhaps try again
			framework.Log(function + " - 34rd3iz retrying getVideoRouteDo operation")
			maxRetries--
			time.Sleep(1 * time.Second)
			if maxRetries == 0 {
				errMsg := fmt.Sprintf(function + "sh4hd3z - max retries reached")
				framework.AddToErrors(socketKey, errMsg)
			}
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

// parameterName is something like "decoder-1_hdmi.out.1"
func getVideoRouteDo(socketKey string, parameterName string) (string, error) {
	method := "Component.GetControls"
	function := "getVideoRouteDo"

	parameterParts := strings.Split(parameterName, "_")
	outputDevice := parameterParts[0]
	output := parameterParts[1]

	params := map[string]string{"Name": outputDevice}

	sent := convertAndSend(socketKey, method, params)

	if sent != true {
		errMsg := fmt.Sprintf(function + " - 94dk34z - error sending command for " + parameterName)
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	response, err := readAndConvert(socketKey)
	if err != nil {
		errMsg := fmt.Sprintf(function + " - 4udj4z - error reading response: " + err.Error())
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	value := "unknown"
	_, resultPresent := response["result"]
	if resultPresent == true {
		result := response["result"].(map[string]interface{})
		controls := result["Controls"].([]interface{})
		for index := 0; index < len(controls); index++ {
			control := controls[index].(map[string]interface{})
			controlname := control["Name"].(string)
			nameparts := strings.Split(controlname, ".select.")
			if control["Value"] == true && len(nameparts) == 2 && strings.Contains(controlname, output+".select.") {
				value = nameparts[1]
				break
			}
		}
	}
	if value == "unknown" {
		errMsg := function + " - ksdl45z - " + outputDevice + " has no inputs that are true with output set to " + output
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	// If we got here, the response was good, so successful return with the state indication
	return `"` + value + `"`, nil
}
