package Models

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
)

type Configuration struct {
	BasePath            string
	LogLevel            logrus.Level
	DbConnection        string
	TrustedProxies      []string
	ListenAddress       string
	SsoAuthority        string
	SsoClientId         string
	IgnoreAuthorisation bool
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

	ssoClientId := "rinkudesu"
	if loadedClientId, isPresent := os.LookupEnv("TAGS_CLIENTID"); isPresent {
		ssoClientId = loadedClientId
	}

	ssoAuthority := os.Getenv("TAGS_AUTHORITY") //todo: add to docker-compose

	_, ignoreAuthorisation := os.LookupEnv("TAGS_IGNORE_AUTHORISATION_UNSAFE")

	return &Configuration{
		BasePath:            basePath,
		LogLevel:            logLevel,
		DbConnection:        dbConnection,
		TrustedProxies:      trustedProxies,
		ListenAddress:       listenAddress,
		SsoClientId:         ssoClientId,
		SsoAuthority:        ssoAuthority,
		IgnoreAuthorisation: ignoreAuthorisation,
	}
}
