package consul

import (
	"os"
	"strconv"

	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/log"
	"github.com/hashicorp/consul/api"
)

type Registrar struct {
	client consulsd.Client
	logger log.Logger
}

func NewRegistrar(consulIP string, consulPort int) *Registrar {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Service discovery domain. In this example we use Consul.
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()
		port := strconv.Itoa(consulPort)
		consulConfig.Address = consulIP + ":" + port
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	return &Registrar{
		client: client,
		logger: logger,
	}
}

type Service struct {
	Name  string
	IP    string
	Port  int
	Check struct {
		Path     string // default /health
		Interval string
		Timeout  string
	}
}

func (rg *Registrar) Register(svc Service) {
	if svc.Check.Path == "" {
		svc.Check.Path = "/health"
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
		Address: svc.IP,
		Port:    svc.Port,
		Tags:    []string{svc.Name},
		Check:   &check,
	}
	sdRegistrar := consulsd.NewRegistrar(rg.client, &asr, rg.logger)
	sdRegistrar.Register()
}

func (rg *Registrar) Deregister(svc Service) {
	asr := api.AgentServiceRegistration{
		ID: svc.Name,
	}
	sdRegistrar := consulsd.NewRegistrar(rg.client, &asr, rg.logger)
	sdRegistrar.Deregister()
}
