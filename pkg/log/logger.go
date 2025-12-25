package log

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	// Load config.env file
	err := godotenv.Load("config.env")
	if err != nil {
		log.Println("No config.env file found, using system environment variables")
	}

	// Read environment variables
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	// Initialize logger
	Log = logrus.New()
	Log.SetOutput(os.Stdout)

	if debug {
		Log.SetLevel(logrus.DebugLevel)
		Log.SetFormatter(&logrus.TextFormatter{
			ForceColors:               true,
			EnvironmentOverrideColors: true,
			FullTimestamp:             true,
		})
		Log.Debug("Verbose logging enabled")
	} else {
		Log.SetLevel(logrus.InfoLevel)
		Log.SetFormatter(&logrus.TextFormatter{
			ForceColors:               true,
			EnvironmentOverrideColors: true,
			DisableTimestamp:          true,
		})
	}
}
