package scan

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
)

func New(config *Config) *App {
	return &App{
		config: config,
	}
}

func (app *App) Entry() error {
	log.SetFormatter(&log.JSONFormatter{})
	dirName := app.config.BugReportDirector
	if dirName == "" {
		log.Errorln("Directory not specified.")
		return errors.New("directory not specified")
	} else {
		dir, err := os.Stat(dirName)
		if err != nil {
			log.Errorf("failed to open directory, error: %w", err)
			return err
		}
		if !dir.IsDir() {
			log.Errorf("%q is not a directory", dir.Name())
			return errors.New(dir.Name() + " is not a directory")
		}
	}

	err := ScanForNsAndDeployments(dirName, app.config.GenerateFakeService)
	if err != nil {
		log.Errorln("Scan of %s directory failed.", dirName)
		return err
	}

	return nil
}
