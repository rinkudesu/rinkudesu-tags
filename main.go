package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rinkudesu/go-kafka/configuration"
	"github.com/rinkudesu/go-kafka/subscriber"
	log "github.com/sirupsen/logrus"
	"os"
	"rinkudesu-tags/authorisation"
	"rinkudesu-tags/controllers"
	"rinkudesu-tags/data"
	"rinkudesu-tags/data/migrations"
	"rinkudesu-tags/message_handlers"
	"rinkudesu-tags/models"
	"rinkudesu-tags/repositories"
	"rinkudesu-tags/services"
)

var (
	routables   []controllers.Routable
	router      *gin.Engine
	config      *models.Configuration
	state       *services.GlobalState
	jwtHandler  *authorisation.JWTHandler
	subscribers []subscriber.Subscriber
)

func init() {
	config = models.NewConfiguration()

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
	var connection = data.DbConnection{}
	err := connection.Initialise(config.DbConnection)
	if err != nil {
		log.Panicf("Failed to initialise database connection: %s", err.Error())
	}

	if !config.IgnoreAuthorisation {
		jwtHandler, err = authorisation.NewJWTHandler(config)
		if err != nil {
			log.Panicf("Failed to initialise jwt handler: %s", err.Error())
		}
	}

	state = services.NewGlobalState(&connection)
}

func migrate(connection data.DbConnector) {
	migrator := migrations.NewExecutor(connection)
	migrator.Migrate()
}

func createControllers() {
	routables = make([]controllers.Routable, 3)
	routables[0] = controllers.CreateLinksController(state)
	routables[1] = controllers.CreateTagsController(state)
	routables[2] = controllers.CreateLinkTagsController(state)
}

func setupRouter() {
	router = gin.New()
	router.Use(gin.Recovery())
	router.Use(services.GetGinLogger())
	router.Use(services.GetHealthcheck(services.CreateHealthcheck(state)))
	router.Use(authorisation.GetGinAuthorisationFilter(jwtHandler, config))
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
	_ = linkDeleteSubscriber.Subscribe(message_handlers.NewLinkDeletedHandler(repositories.CreateLinksRepository(state)))
	_ = linkDeleteSubscriber.BeginHandle()

	userDeleteSubscriber, _ := subscriber.NewKafkaSubscriber(kafkaConfig)
	_ = userDeleteSubscriber.Subscribe(message_handlers.CreateUserDeletedHandler(state))
	_ = userDeleteSubscriber.BeginHandle()

	subscribers = []subscriber.Subscriber{linkDeleteSubscriber, userDeleteSubscriber}
}
