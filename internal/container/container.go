package container

import (
	"fmt"
	"os"
	"sync"

	"github.com/AlexRojer31/sandbox/internal/environment"
	"github.com/sirupsen/logrus"
)

const INITIAL_FAILED = "Container initial was failed: "

type Container struct {
	Env    *environment.Env
	Logger *logrus.Logger

	initialErrors []error
}

var (
	instance *Container
	once     sync.Once
)

func GetInstance(args ...any) *Container {
	once.Do(func() {
		for _, a := range args {
			switch v := a.(type) {
			case []string:
				c := Container{
					initialErrors: make([]error, 0),
				}
				env, err := environment.New(v)
				if err != nil {
					c.initialErrors = append(c.initialErrors, err)
				}

				c.Env = env
				c.setLogger()

				if len(c.initialErrors) != 0 {
					errorsStr := fmt.Sprint(c.initialErrors)

					panic(INITIAL_FAILED + errorsStr)
				}

				instance = &c
				return
			default:
				panic(INITIAL_FAILED)
			}
		}
	})

	return instance
}

func (app *Container) setLogger() {
	app.Logger = logrus.New()
	app.Logger.SetFormatter(
		&logrus.TextFormatter{
			TimestampFormat:        "2006-01-02 15:04:05",
			FullTimestamp:          true,
			PadLevelText:           true,
			DisableLevelTruncation: true,
		})
	app.Logger.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel(app.Env.Config.LogLevel)
	if err != nil {
		level = logrus.ErrorLevel
	}

	if app.Env.Debug {
		level = logrus.DebugLevel
		app.Logger.SetReportCaller(app.Env.Debug)
	}
	app.Logger.SetLevel(level)
}
