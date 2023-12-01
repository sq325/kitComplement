package consul

import (
	"os"
	"strconv"

	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/log"
	"github.com/hashicorp/consul/api"
	"github.com/sq325/kitComplement/pkg/tool"
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
		Tags:    append(svc.Tags, svc.Name),
		Check:   &check,
	}
	sdRegistrar := consulsd.NewRegistrar(rg.client, &asr, rg.logger)
	sdRegistrar.Register()
}

func (rg *registrar) Deregister(svc *Service) {
	asr := api.AgentServiceRegistration{
		ID: svc.Name,
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
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	return
}
