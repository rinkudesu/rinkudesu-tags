package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
)

type Configuration struct {
	BasePath       string
	LogLevel       logrus.Level
	DbConnection   string
	TrustedProxies []string
	ListenAddress  string
}

func NewConfiguration() *Configuration {
	basePath := "/api"
	if loadedPath, isPresent := os.LookupEnv("TAGS_BASE-PATH"); isPresent {
		basePath = loadedPath
	}

	logLevel := logrus.InfoLevel
	if logLevelString, isPresent := os.LookupEnv("TAGS_LOG-LEVEL"); isPresent {
		loadedLogLevel, err := logrus.ParseLevel(logLevelString)
		if err != nil {
			log.Panicf("Failed to parse log level: %s", err.Error())
		}
		logLevel = loadedLogLevel
	}

	dbConnection := os.Getenv("TAGS_DB")

	var trustedProxies []string = nil
	if trustedProxiesString, isPresent := os.LookupEnv("TAGS_PROXY"); isPresent {
		trustedProxies = strings.Split(trustedProxiesString, ",")
	}

	listenAddress := "localhost:5000"
	if loadedAddress, isPresent := os.LookupEnv("TAGS_ADDRESS"); isPresent {
		listenAddress = loadedAddress
	}

	return &Configuration{
		BasePath:       basePath,
		LogLevel:       logLevel,
		DbConnection:   dbConnection,
		TrustedProxies: trustedProxies,
		ListenAddress:  listenAddress,
	}
}
