//
// Copyright (c) 2020 Intel Corporation
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

package handlers

import (
	"context"
	"sync"
	"testing"

	"github.com/edgexfoundry/go-mod-registry/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/logging"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/config"
	"github.com/edgexfoundry/go-mod-bootstrap/di"

	"github.com/edgexfoundry/app-functions-sdk-go/internal"
	"github.com/edgexfoundry/app-functions-sdk-go/internal/bootstrap/container"
	"github.com/edgexfoundry/app-functions-sdk-go/internal/common"
)

func TestClientsBootstrapHandler(t *testing.T) {
	configuration := &common.ConfigurationStruct{
		Service: common.ServiceInfo{
			ClientMonitor: "5s",
		},
	}

	logger := logging.FactoryToStdout("clients-test")
	var registryClient registry.Client = nil

	dic := di.NewContainer(di.ServiceConstructorMap{
		bootstrapContainer.LoggingClientInterfaceName: func(get di.Get) interface{} {
			return logger
		},
		bootstrapContainer.RegistryClientInterfaceName: func(get di.Get) interface{} {
			return registryClient
		},
	})

	coreDataClientInfo := config.ClientInfo{
		Host:     "localhost",
		Port:     48080,
		Protocol: "http",
	}

	commandClientInfo := config.ClientInfo{
		Host:     "localhost",
		Port:     48081,
		Protocol: "http",
	}

	notificationsClientInfo := config.ClientInfo{
		Host:     "localhost",
		Port:     48082,
		Protocol: "http",
	}

	startupTimer := startup.NewStartUpTimer(internal.BootRetrySecondsDefault, internal.BootTimeoutSecondsDefault)

	tests := []struct {
		Name                    string
		CoreDataClientInfo      *config.ClientInfo
		CommandClientInfo       *config.ClientInfo
		NotificationsClientInfo *config.ClientInfo
		MonitorDuration         string
		ExpectSuccess           bool
	}{
		{
			Name:                    "All Clients",
			CoreDataClientInfo:      &coreDataClientInfo,
			CommandClientInfo:       &commandClientInfo,
			NotificationsClientInfo: &notificationsClientInfo,
			MonitorDuration:         "5s",
			ExpectSuccess:           true,
		},
		{
			Name:                    "No Clients",
			CoreDataClientInfo:      nil,
			CommandClientInfo:       nil,
			NotificationsClientInfo: nil,
			MonitorDuration:         "5s",
			ExpectSuccess:           true,
		},
		{
			Name:                    "Only Core Data Clients",
			CoreDataClientInfo:      &coreDataClientInfo,
			CommandClientInfo:       nil,
			NotificationsClientInfo: nil,
			MonitorDuration:         "5s",
			ExpectSuccess:           true,
		},
		{
			Name:                    "Invalid MonitorDuration",
			CoreDataClientInfo:      &coreDataClientInfo,
			CommandClientInfo:       nil,
			NotificationsClientInfo: nil,
			MonitorDuration:         "bogus",
			ExpectSuccess:           true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			configuration.Service.ClientMonitor = test.MonitorDuration
			configuration.Clients = make(map[string]config.ClientInfo)

			if test.CoreDataClientInfo != nil {
				configuration.Clients[common.CoreDataClientName] = coreDataClientInfo
			}

			if test.CommandClientInfo != nil {
				configuration.Clients[common.CoreCommandClientName] = commandClientInfo
			}

			if test.NotificationsClientInfo != nil {
				configuration.Clients[common.NotificationsClientName] = notificationsClientInfo
			}

			dic.Update(di.ServiceConstructorMap{
				container.ConfigurationName: func(get di.Get) interface{} {
					return configuration
				},
			})

			actualSuccess := NewClients().BootstrapHandler(context.Background(), &sync.WaitGroup{}, startupTimer, dic)
			require.Equal(t, test.ExpectSuccess, actualSuccess)
			if actualSuccess == false {
				return // Test is complete
			}

			eventClient := container.EventClientFrom(dic.Get)
			valueDescriptorClient := container.ValueDescriptorClientFrom(dic.Get)
			commandClient := container.CommandClientFrom(dic.Get)
			notificationsClient := container.NotificationsClientFrom(dic.Get)

			if test.CoreDataClientInfo != nil {
				assert.NotNil(t, eventClient)
				assert.NotNil(t, valueDescriptorClient)
			} else {
				assert.Nil(t, eventClient)
				assert.Nil(t, valueDescriptorClient)
			}

			if test.CommandClientInfo != nil {
				assert.NotNil(t, commandClient)
			} else {
				assert.Nil(t, commandClient)
			}

			if test.NotificationsClientInfo != nil {
				assert.NotNil(t, notificationsClient)
			} else {
				assert.Nil(t, notificationsClient)
			}
		})
	}
}
