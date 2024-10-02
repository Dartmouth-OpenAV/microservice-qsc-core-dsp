package main

import (
	"errors"
	"github.com/Dartmouth-OpenAV/microservice-framework/framework"
)

func setFrameworkGlobals() {
	// globals that change modes in the microservice framework:
	framework.MicroserviceName = "OpenAV QSC Core DSP Microservice"
	framework.DefaultSocketPort = 1710 // QSC's default socket port
	framework.CheckFunctionAppendBehavior = "Remove older instance"
	framework.GlobalDelimiter = 0
	framework.KeepAlive = true

	framework.RegisterMainGetFunc(doDeviceSpecificGet)
	framework.RegisterMainSetFunc(doDeviceSpecificSet)
}

// Every microservice using this golang microservice framework needs to provide this function to invoke functions to do sets.
// socketKey is the network connection for the framework to use to communicate with the device.
// setting is the first parameter in the URI.
// arg1 are the second and third parameters in the URI.
//
//	  Example PUT URIs that will result in this function being invoked:
//		 ":address/:setting/"
//	  ":address/:setting/:arg1"
//	  ":address/:setting/:arg1/:arg2"
func doDeviceSpecificSet(socketKey string, setting string, arg1 string, arg2 string, arg3 string) (string, error) {
	function := "doDeviceSpecificSet"

	// Add a case statement for each set function your microservice implements.  These calls can use 0, 1, or 2 arguments.
	switch setting {
	case "volume":
		return setVolume(socketKey, arg1, arg2)
	case "audiomute":
		return setToggle(socketKey, arg1, arg2)
	case "voicelift":
		return setToggle(socketKey, "voice_lift", arg1)
	case "autotracking":
		return setToggle(socketKey, "camera_tracking", arg1)
	case "toggle":
		return setToggle(socketKey, arg1, arg2)
	case "videoroute":
		return setVideoRoute(socketKey, arg1, arg2)
	}

	// If we get here, we didn't recognize the setting.  Send an error back to the config writer who had a bad URL.
	errMsg := function + " - unrecognized setting in URI: " + setting
	framework.AddToErrors(socketKey, errMsg)
	err := errors.New(errMsg)
	return setting, err
}

// Every microservice using this golang microservice framework needs to provide this function to invoke functions to do gets.
// socketKey is the network connection for the framework to use to communicate with the device.
// setting is the first parameter in the URI.
// arg1 are the second and third parameters in the URI.
//
//	  Example GET URIs that will result in this function being invoked:
//		 ":address/:setting/"
//	  ":address/:setting/:arg1"
//	  ":address/:setting/:arg1/:arg2"
//
// Every microservice using this golang microservice framework needs to provide this function to invoke functions to do gets.
func doDeviceSpecificGet(socketKey string, setting string, arg1 string, arg2 string) (string, error) {
	function := "doDeviceSpecificGet"

	switch setting {
	case "volume":
		return getVolume(socketKey, arg1)
	case "audiomute":
		return getToggle(socketKey, arg1)
	case "voicelift":
		return getToggle(socketKey, "voice_lift")
	case "autotracking":
		return getToggle(socketKey, "camera_tracking")
	case "toggle":
		return getToggle(socketKey, arg1)
	case "videoroute":
		return getVideoRoute(socketKey, arg1)
	}

	// If we get here, we didn't recognize the setting.  Send an error back to the config writer who had a bad URL.
	errMsg := function + " - unrecognized setting in URI: " + setting
	framework.AddToErrors(socketKey, errMsg)
	err := errors.New(errMsg)
	return setting, err
}

func main() {
	setFrameworkGlobals()
	framework.Startup()
 } 
