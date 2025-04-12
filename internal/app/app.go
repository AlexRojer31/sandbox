package app

import (
	"os"
	"sync"

	"github.com/AlexRojer31/sandbox/internal/environment"
	"github.com/sirupsen/logrus"
)

type App struct {
	Env    *environment.Env
	Logger *logrus.Logger

	wg sync.WaitGroup
}

var (
	instance *App
	once     sync.Once
)

func GetInstance(args ...any) *App {
	once.Do(func() {
		for _, a := range args {
			switch v := a.(type) {
			case []string:
				app := App{}
				env, err := environment.New(v)
				if err != nil {
					panic(err)
				}
				app.Env = env
				app.setLogger()
				app.wg.Add(1)

				instance = &app
				return
			}
		}
	})

	return instance
}

func (app *App) setLogger() {
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
