package environment

import (
	"flag"
	"fmt"

	"github.com/AlexRojer31/sandbox/internal/config"
)

const (
	APP_NAME = "sandbox"
)

type Env struct {
	Debug      bool
	configFile string
	Config     *config.Config
}

func New(args []string) (*Env, error) {
	env := Env{}

	if err := env.fromArgs(args); err != nil {
		return &env, err
	}

	if err := env.loadConfig(); err != nil {
		return &env, err
	}

	return &env, nil
}

func (env *Env) fromArgs(args []string) error {
	flagSet := flag.NewFlagSet(APP_NAME, flag.ContinueOnError)
	flagSet.BoolVar(&env.Debug, "d", false, "Print debug messages")
	flagSet.StringVar(&env.configFile, "c", "configs/dev.yaml", "Config file")

	if err := flagSet.Parse(args); err != nil {
		return fmt.Errorf("environment: can't parse flags: %w", err)
	}

	return nil
}

func (env *Env) loadConfig() error {
	config, err := config.New(env.configFile)
	if err != nil {
		return fmt.Errorf("environment: can't create config: %w", err)
	}

	env.Config = config

	return nil
}
