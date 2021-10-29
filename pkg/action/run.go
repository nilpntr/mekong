package action

import (
	"github.com/getsentry/sentry-go"
	"github.com/nilpntr/mekong/pkg/cli"
	"log"
	"net/http"
	"time"
)

func NewRun(settings *cli.EnvSettings) error {
	cfg, err := newConfig(settings.ConfigFile)
	if err != nil {
		return err
	}

	if err := validateConfig(cfg); err != nil {
		return err
	}

	if cfg.Sentry != nil {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:         cfg.Sentry.Dsn,
			Environment: pointerString(cfg.Sentry.Environment),
			Release:     pointerString(cfg.Sentry.Release),
			Debug:       pointerBool(cfg.Sentry.Debug),
		}); err != nil {
			return err
		}
		defer sentry.Flush(2 * time.Second)
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
