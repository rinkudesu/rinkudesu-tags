package Models

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Configuration struct {
	BasePath            string
	LogLevel            log.Level
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

	logLevel := log.InfoLevel
	if logLevelString, isPresent := os.LookupEnv("TAGS_LOG-LEVEL"); isPresent {
		loadedLogLevel, err := log.ParseLevel(logLevelString)
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

	ssoAuthority := os.Getenv("TAGS_AUTHORITY")

	_, ignoreAuthorisation := os.LookupEnv("TAGS_IGNORE_AUTHORISATION_UNSAFE")

	if ignoreAuthorisation {
		log.Warning("Authorisation is being ignored, THIS IS UNSAFE, DON'T USE IN PRODUCTION")
	}

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
