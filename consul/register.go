// Copyright 2023 Sun Quan
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package consul

import (
	"os"
	"strconv"

	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/log"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/sq325/kitComplement/tool"
)

type Registrar interface {
	Register(svc *Service)
	Deregister(svc *Service)
}

type registrar struct {
	client consulsd.Client
	logger log.Logger
}

func NewRegistrar(consulClient consulsd.Client, logger log.Logger) Registrar {

	return &registrar{
		client: consulClient,
		logger: logger,
	}
}

type Service struct {
	Name  string
	ID    string // default is svcName_UUID
	IP    string // svc ip, default hostAdmIp
	Port  int    // svc port
	Tags  []string
	Check struct {
		Path     string // default /health
		Interval string // "60s"
		Timeout  string // "10s"
	}
}

func (rg *registrar) Register(svc *Service) {
	if svc.Check.Path == "" {
		svc.Check.Path = "/health"
	}
	if svc.IP == "" {
		svc.IP, _ = tool.HostAdmIp(nil)
	}
	if svc.ID == "" {
		uuid := uuid.NewString()
		svc.ID = svc.Name + "_" + uuid
	}

	checkUrl := "http://" + svc.IP + ":" + strconv.Itoa(svc.Port) + svc.Check.Path

	check := api.AgentServiceCheck{
		HTTP:     checkUrl,
		Interval: svc.Check.Interval,
		Timeout:  svc.Check.Timeout,
		Notes:    "svc health checks",
	}
	asr := api.AgentServiceRegistration{
		Name:    svc.Name,
		ID:      svc.ID,
		Address: svc.IP,
		Port:    svc.Port,
		Tags:    append(svc.Tags, svc.Name),
		Check:   &check,
	}
	sdRegistrar := consulsd.NewRegistrar(rg.client, &asr, rg.logger)
	sdRegistrar.Register()
}

func (rg *registrar) Deregister(svc *Service) {
	asr := api.AgentServiceRegistration{
		ID: svc.ID,
	}
	sdRegistrar := consulsd.NewRegistrar(rg.client, &asr, rg.logger)
	sdRegistrar.Deregister()
}

func NewConsulClient(consulIP string, consulPort int) consulsd.Client {

	consulConfig := api.DefaultConfig()
	port := strconv.Itoa(consulPort)
	consulConfig.Address = "http://" + consulIP + ":" + port
	consulClient, _ := api.NewClient(consulConfig)

	return consulsd.NewClient(consulClient)
}

func NewLogger() (logger log.Logger) {

	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestamp)
	logger = log.With(logger, "caller", log.DefaultCaller)

	return
}
