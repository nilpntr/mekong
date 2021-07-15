package cli

import (
	"github.com/spf13/pflag"
	"os"
)

type EnvSettings struct {
	ConfigFile string
}

func New() *EnvSettings {
	env := &EnvSettings{
		ConfigFile: os.Getenv("MEKONG_CONFIG_FILE"),
	}

	return env
}

func (e *EnvSettings) AddFlags(f *pflag.FlagSet) {
	f.StringVarP(&e.ConfigFile, "config-file", "", e.ConfigFile, "config file")
}
