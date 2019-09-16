//
// Copyright (c) 2019 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package appsdk

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/util"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

const (
	ValueDescriptors = "valuedescriptors"
	DeviceNames      = "devicenames"
	Key              = "key"
	InitVector       = "initvector"
	Url              = "url"
	MimeType         = "mimetype"
	PersistOnError   = "persistOnError"
	Cert             = "cert"
	Qos              = "qos"
	Retain           = "retain"
	AutoReconnect    = "autoreconnect"
	DeviceName       = "devicename"
	ReadingName      = "readingname"
)

// AppFunctionsSDKConfigurable contains the helper functions that return the function pointers for building the configurable function pipeline.
// They transform the parameters map from the Pipeline configuration in to the actual actual parameters required by the function.
type AppFunctionsSDKConfigurable struct {
	Sdk *AppFunctionsSDK
}

// FilterByDeviceName - Specify the devices of interest to  filter for data coming from certain sensors.
// The Filter by Device transform looks at the Event in the message and looks at the devices of interest list,
// provided by this function, and filters out those messages whose Event is for devices not on the
// devices of interest.
// This function will return an error and stop the pipeline if a non-edgex
// event is received or if no data is recieved.
// For example, data generated by a motor does not get passed to functions only interested in data from a thermostat.
// This function is a configuration function and returns a function pointer.
func (dynamic AppFunctionsSDKConfigurable) FilterByDeviceName(parameters map[string]string) appcontext.AppFunction {
	deviceNames, ok := parameters[DeviceNames]
	if !ok {
		dynamic.Sdk.LoggingClient.Error("Could not find " + DeviceNames)
		return nil
	}
	deviceNamesCleaned := util.DeleteEmptyAndTrim(strings.FieldsFunc(deviceNames, util.SplitComma))
	transform := transforms.Filter{
		FilterValues: deviceNamesCleaned,
	}
	dynamic.Sdk.LoggingClient.Debug("Device Name Filters", DeviceNames, strings.Join(deviceNamesCleaned, ","))

	return transform.FilterByDeviceName
}

// FilterByValueDescriptor - Specify the value descriptors of interest to filter for data from certain types of IoT objects,
// such as temperatures, motion, and so forth, that may come from an array of sensors or devices. The Filter by Value Descriptor assesses
// the data in each Event and Reading, and removes readings that have a value descriptor that is not in the list of
// value descriptors of interest for the application.
// This function will return an error and stop the pipeline if a non-edgex
// event is received or if no data is recieved.
// For example, pressure reading data does not go to functions only interested in motion data.
// This function is a configuration function and returns a function pointer.
func (dynamic AppFunctionsSDKConfigurable) FilterByValueDescriptor(parameters map[string]string) appcontext.AppFunction {
	valueDescriptors, ok := parameters[ValueDescriptors]
	if !ok {
		dynamic.Sdk.LoggingClient.Error("Could not find " + ValueDescriptors)
		return nil
	}
	valueDescriptorsCleaned := util.DeleteEmptyAndTrim(strings.FieldsFunc(valueDescriptors, util.SplitComma))
	transform := transforms.Filter{
		FilterValues: valueDescriptorsCleaned,
	}
	dynamic.Sdk.LoggingClient.Debug("Value Descriptors Filter", ValueDescriptors, strings.Join(valueDescriptorsCleaned, ","))
	return transform.FilterByValueDescriptor
}

// TransformToXML transforms an EdgeX event to XML.
// It will return an error and stop the pipeline if a non-edgex
// event is received or if no data is recieved.
// This function is a configuration function and returns a function pointer.
func (dynamic AppFunctionsSDKConfigurable) TransformToXML() appcontext.AppFunction {
	transform := transforms.Conversion{}
	return transform.TransformToXML
}

// TransformToJSON transforms an EdgeX event to JSON.
// It will return an error and stop the pipeline if a non-edgex
// event is received or if no data is recieved.
// This function is a configuration function and returns a function pointer.
func (dynamic AppFunctionsSDKConfigurable) TransformToJSON() appcontext.AppFunction {
	transform := transforms.Conversion{}
	return transform.TransformToJSON
}

// MarkAsPushed will make a request to CoreData to mark the event that triggered the pipeline as pushed.
// This function is a configuration function and returns a function pointer.
func (dynamic AppFunctionsSDKConfigurable) MarkAsPushed() appcontext.AppFunction {
	transform := transforms.CoreData{}
	return transform.MarkAsPushed
}

// PushToCore pushes the provided value as an event to CoreData using the device name and reading name that have been set. If validation is turned on in
// CoreServices then your deviceName and readingName must exist in the CoreMetadata and be properly registered in EdgeX.
// This function is a configuration function and returns a function pointer.
func (dynamic AppFunctionsSDKConfigurable) PushToCore(parameters map[string]string) appcontext.AppFunction {
	deviceName, ok := parameters[DeviceName]
	if !ok {
		dynamic.Sdk.LoggingClient.Error("Could not find " + DeviceName)
		return nil
	}
	readingName, ok := parameters[ReadingName]
	if !ok {
		dynamic.Sdk.LoggingClient.Error("Could not find " + readingName)
		return nil
	}
	deviceName = strings.TrimSpace(deviceName)
	readingName = strings.TrimSpace(readingName)
	dynamic.Sdk.LoggingClient.Debug("PushToCore Parameters", DeviceName, deviceName, ReadingName, readingName)
	transform := transforms.CoreData{
		DeviceName:  deviceName,
		ReadingName: readingName,
	}
	return transform.PushToCoreData
}

// CompressWithGZIP compresses data received as either a string,[]byte, or json.Marshaler using gzip algorithm and returns a base64 encoded string as a []byte.
// This function is a configuration function and returns a function pointer.
func (dynamic AppFunctionsSDKConfigurable) CompressWithGZIP() appcontext.AppFunction {
	transform := transforms.Compression{}
	return transform.CompressWithGZIP
}

// CompressWithZLIB compresses data received as either a string,[]byte, or json.Marshaler using zlib algorithm and returns a base64 encoded string as a []byte.
// This function is a configuration function and returns a function pointer.
func (dynamic AppFunctionsSDKConfigurable) CompressWithZLIB() appcontext.AppFunction {
	transform := transforms.Compression{}
	return transform.CompressWithZLIB
}

// EncryptWithAES encrypts either a string, []byte, or json.Marshaller type using AES encryption.
// It will return a byte[] of the encrypted data.
// This function is a configuration function and returns a function pointer.
func (dynamic AppFunctionsSDKConfigurable) EncryptWithAES(parameters map[string]string) appcontext.AppFunction {
	key, ok := parameters[Key]
	if !ok {
		dynamic.Sdk.LoggingClient.Error("Could not find " + Key)
		return nil
	}
	initVector, ok := parameters[InitVector]
	if !ok {
		dynamic.Sdk.LoggingClient.Error("Could not find " + InitVector)
		return nil
	}
	transforms := transforms.Encryption{
		Key:                  key,
		InitializationVector: initVector,
	}
	return transforms.EncryptWithAES
}

// HTTPPost will send data from the previous function to the specified Endpoint via http POST. If no previous function exists,
// then the event that triggered the pipeline will be used. Passing an empty string to the mimetype
// method will default to application/json.
// This function is a configuration function and returns a function pointer.
func (dynamic AppFunctionsSDKConfigurable) HTTPPost(parameters map[string]string) appcontext.AppFunction {
	var err error

	url, ok := parameters[Url]
	if !ok {
		dynamic.Sdk.LoggingClient.Error("Could not find " + Url)
		return nil
	}
	mimeType, ok := parameters[MimeType]
	if !ok {
		dynamic.Sdk.LoggingClient.Error("Could not find " + MimeType)
		return nil
	}

	// PersistOnError is optional and is false by default.
	persistOnError := false
	value, ok := parameters[PersistOnError]
	if ok {
		persistOnError, err = strconv.ParseBool(value)
		if err != nil {
			dynamic.Sdk.LoggingClient.Error(fmt.Sprintf("Could not parse '%s' to a bool for '%s' parameter", value, PersistOnError), "error", err)
			return nil
		}
	}

	url = strings.TrimSpace(url)
	mimeType = strings.TrimSpace(mimeType)

	transform := transforms.HTTPSender{
		URL:            url,
		MimeType:       mimeType,
		PersistOnError: persistOnError,
	}
	dynamic.Sdk.LoggingClient.Debug("HTTP Post Parameters", Url, transform.URL, MimeType, transform.MimeType)
	return transform.HTTPPost
}

// HTTPPostJSON sends data from the previous function to the specified Endpoint via http POST with a mime type of application/json.
// If no previous function exists, then the event that triggered the pipeline will be used.
// This function is a configuration function and returns a function pointer.
func (dynamic AppFunctionsSDKConfigurable) HTTPPostJSON(parameters map[string]string) appcontext.AppFunction {
	url, ok := parameters[Url]
	if !ok {
		dynamic.Sdk.LoggingClient.Error("Could not find " + Url)
		return nil
	}

	// PersistOnError is optional and is false by default.
	persistOnError := false
	value, ok := parameters[PersistOnError]
	if ok {
		var err error
		persistOnError, err = strconv.ParseBool(value)
		if err != nil {
			dynamic.Sdk.LoggingClient.Error(fmt.Sprintf("Could not parse '%s' to a bool for '%s' parameter", value, PersistOnError), "error", err)
			return nil
		}
	}

	url = strings.TrimSpace(url)
	dynamic.Sdk.LoggingClient.Debug("HTTP Post JSON Parameters", Url, url)
	return transforms.NewHTTPSender(url, "application/json", persistOnError).HTTPPost
}

// HTTPPostXML sends data from the previous function to the specified Endpoint via http POST with a mime type of application/xml.
// If no previous function exists, then the event that triggered the pipeline will be used.
// This function is a configuration function and returns a function pointer.
func (dynamic AppFunctionsSDKConfigurable) HTTPPostXML(parameters map[string]string) appcontext.AppFunction {
	url, ok := parameters[Url]
	if !ok {
		dynamic.Sdk.LoggingClient.Error("Could not find " + Url)
		return nil
	}

	// PersistOnError is optional and is false by default.
	persistOnError := false
	value, ok := parameters[PersistOnError]
	if ok {
		var err error
		persistOnError, err = strconv.ParseBool(value)
		if err != nil {
			dynamic.Sdk.LoggingClient.Error(fmt.Sprintf("Could not parse '%s' to a bool for '%s' parameter", value, PersistOnError), "error", err)
			return nil
		}
	}

	url = strings.TrimSpace(url)
	dynamic.Sdk.LoggingClient.Debug("HTTP Post XML Parameters", Url, url)
	return transforms.NewHTTPSender(url, "application/xml", persistOnError).HTTPPost
}

// MQTTSend sends data from the previous function to the specified MQTT broker.
// If no previous function exists, then the event that triggered the pipeline will be used.
// This function is a configuration function and returns a function pointer.
func (dynamic AppFunctionsSDKConfigurable) MQTTSend(parameters map[string]string, addr models.Addressable) appcontext.AppFunction {
	var err error
	qos := 0
	retain := false
	autoreconnect := false
	// optional string params
	cert := parameters[Cert]
	key := parameters[Key]

	qosVal, ok := parameters[Qos]
	if ok {
		qos, err = strconv.Atoi(qosVal)
		if err != nil {
			dynamic.Sdk.LoggingClient.Error("Unable to parse " + Qos + " value")
			return nil
		}
	}
	retainVal, ok := parameters[Retain]
	if ok {
		retain, err = strconv.ParseBool(retainVal)
		if err != nil {
			dynamic.Sdk.LoggingClient.Error("Unable to parse " + Retain + " value")
			return nil
		}
	}
	autoreconnectVal, ok := parameters[AutoReconnect]
	if ok {
		autoreconnect, err = strconv.ParseBool(autoreconnectVal)
		if err != nil {
			dynamic.Sdk.LoggingClient.Error("Unable to parse " + AutoReconnect + " value")
			return nil
		}
	}
	dynamic.Sdk.LoggingClient.Debug("MQTT Send Parameters", "Address", addr, Qos, qosVal, Retain, retainVal, AutoReconnect, autoreconnectVal, Cert, cert, Key, key)

	var pair *transforms.KeyCertPair

	if len(cert) > 0 && len(key) > 0 {
		pair = &transforms.KeyCertPair{
			CertFile: cert,
			KeyFile:  key,
		}
	}

	// PersistOnError os optional and is false by default.
	persistOnError := false
	value, ok := parameters[PersistOnError]
	if ok {
		persistOnError, err = strconv.ParseBool(value)
		if err != nil {
			dynamic.Sdk.LoggingClient.Error(fmt.Sprintf("Could not parse '%s' to a bool for '%s' parameter", value, PersistOnError), "error", err)
			return nil
		}
	}

	mqttconfig := transforms.NewMqttConfig()
	mqttconfig.SetQos(byte(qos))
	mqttconfig.SetRetain(retain)
	mqttconfig.SetAutoreconnect(autoreconnect)
	sender := transforms.NewMQTTSender(dynamic.Sdk.LoggingClient, addr, pair, mqttconfig, persistOnError)
	return sender.MQTTSend
}

// SetOutputData sets the output data to that passed in from the previous function.
// It will return an error and stop the pipeline if data passed in is not of type []byte, string or json.Mashaler
// This function is a configuration function and returns a function pointer.
func (dynamic AppFunctionsSDKConfigurable) SetOutputData() appcontext.AppFunction {
	transform := transforms.OutputData{}
	return transform.SetOutputData
}
