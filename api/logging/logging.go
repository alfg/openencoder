package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Log exports the configured logger.
var Log *logrus.Logger

func init() {
	log := logrus.New()

	// Set JSON formatter if production.
	if os.Getenv("GIN_MODE") == "release" {
		log.SetFormatter(&logrus.JSONFormatter{})
	}
	log.SetOutput(os.Stdout)

	Log = log
}
