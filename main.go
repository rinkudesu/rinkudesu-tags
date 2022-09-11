package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rinkudesu/go-kafka/configuration"
	"github.com/rinkudesu/go-kafka/subscriber"
	log "github.com/sirupsen/logrus"
	"os"
	"rinkudesu-tags/Authorisation"
	"rinkudesu-tags/Controllers"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Data/Migrations"
	"rinkudesu-tags/MessageHandlers"
	"rinkudesu-tags/Models"
	"rinkudesu-tags/Repositories"
	"rinkudesu-tags/Services"
)

var (
	routables   []Controllers.Routable
	router      *gin.Engine
	config      *Models.Configuration
	state       *Services.GlobalState
	jwtHandler  *Authorisation.JWTHandler
	subscribers []subscriber.Subscriber
)

func init() {
	config = Models.NewConfiguration()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(config.LogLevel)
}

func main() {
	makeGlobalState()

	defer state.DbConnection.Close()
	migrate(state.DbConnection)

	setupMessageHandlers()
	createControllers()
	//this blocks until the server is turned off
	setupRouter()

	for _, runningSubscriber := range subscribers {
		runningSubscriber.StopHandle()
		_ = runningSubscriber.Close()
	}
}

func makeGlobalState() {
	var connection = Data.DbConnection{}
	err := connection.Initialise(config.DbConnection)
	if err != nil {
		log.Panicf("Failed to initialise database connection: %s", err.Error())
	}

	if !config.IgnoreAuthorisation {
		jwtHandler, err = Authorisation.NewJWTHandler(config)
		if err != nil {
			log.Panicf("Failed to initialise jwt handler: %s", err.Error())
		}
	}

	state = Services.NewGlobalState(&connection)
}

func migrate(connection Data.DbConnector) {
	migrator := Migrations.NewExecutor(connection)
	migrator.Migrate()
}

func createControllers() {
	routables = make([]Controllers.Routable, 3)
	routables[0] = Controllers.CreateLinksController(state)
	routables[1] = Controllers.CreateTagsController(state)
	routables[2] = Controllers.CreateLinkTagsController(state)
}

func setupRouter() {
	router = gin.New()
	router.Use(gin.Recovery())
	router.Use(Services.GetGinLogger())
	router.Use(Services.GetHealthcheck(Services.CreateHealthcheck(state)))
	router.Use(Authorisation.GetGinAuthorisationFilter(jwtHandler, config))
	err := router.SetTrustedProxies(config.TrustedProxies)
	if err != nil {
		log.Panicf("Failed to set trusted proxies: %s", err.Error())
	}

	for _, routable := range routables {
		routable.SetupRoutes(router, config.BasePath)
	}

	if err := router.Run(config.ListenAddress); err != nil {
		log.Panicf("Server failed while listening: %s", err.Error())
	}
}

func setupMessageHandlers() {
	kafkaConfig, err := configuration.NewKafkaConfigurationFromEnv()
	if err != nil {
		log.Fatalf("Failed to read kafka config from env")
	}

	linkDeleteSubscriber, _ := subscriber.NewKafkaSubscriber(kafkaConfig)
	_ = linkDeleteSubscriber.Subscribe(MessageHandlers.NewLinkDeletedHandler(Repositories.CreateLinksRepository(state)))
	_ = linkDeleteSubscriber.BeginHandle()

	userDeleteSubscriber, _ := subscriber.NewKafkaSubscriber(kafkaConfig)
	_ = userDeleteSubscriber.Subscribe(MessageHandlers.CreateUserDeletedHandler(state))
	_ = userDeleteSubscriber.BeginHandle()

	subscribers = []subscriber.Subscriber{linkDeleteSubscriber, userDeleteSubscriber}
}
