package action

import (
	"github.com/nilpntr/mekong/pkg/cli"
	"log"
	"net/http"
)

func NewRun(settings *cli.EnvSettings) error {
	cfg, err := newConfig(settings.ConfigFile)
	if err != nil {
		return err
	}

	if err := validateConfig(cfg); err != nil {
		return err
	}

	proxy := newProxyServer(cfg)

	finish := make(chan bool)

	for _, listenPort := range proxy.Config.ListenPorts {
		go func(port string) {
			log.Fatalln(http.ListenAndServe(port, proxy))
		}(listenPort)
	}

	<-finish

	return nil
}