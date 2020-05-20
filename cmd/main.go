package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/trinhdaiphuc/Source-code-marking/configs"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/handlers"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/routers"
)

// Golbal variablesfor db connection

func main() {
	echoServer := configs.NewEchoServer()

	err := godotenv.Load()
	if err != nil {
		echoServer.EchoContext.Logger.Error("Error loading .env file", err)
	}

	// Set max go process
	configs.ConfigureMaxProcess()

	// Customizing Echo Logger
	configs.LoggerConfig(echoServer)

	// configsure middleware
	err = configs.ConfigureMiddleware(echoServer)
	if err != nil {
		echoServer.Logger.Error("Error when configsure middleware ", err)
	}

	// Declare Context type object for managing multiple API requests timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Database connection
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DB_HOST")))
	if err != nil {
		echoServer.Logger.Fatal("Error when connect to MongoDB.", err)
	}
	err = db.Ping(ctx, nil)
	if err != nil {
		echoServer.Logger.Fatal("Could not connect to MongoDB: %v\n", err)
	} else {
		echoServer.Logger.Info("Connected to Mongodb.")
	}

	// Create new User collection
	models.NewUserCollection(db)

	// Initialize handler
	h := &handlers.Handler{DB: db, AppLog: echoServer.Logger}

	// configsure HTTP error handler
	echoServer.EchoContext.HTTPErrorHandler = h.CustomHTTPErrorHandler

	// configs routing server
	routers.LandingPage(echoServer.EchoContext, h)
	routers.Users(echoServer.EchoContext, h)

	// Customizing Echo server
	serverCustomize := &http.Server{
		Addr:         fmt.Sprint(":", os.Getenv("PORT")),
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
	}

	// Start the server in a child routine
	go func() {
		echoServer.EchoContext.Logger.Fatal(echoServer.EchoContext.StartServer(serverCustomize))
	}()

	quit := make(chan os.Signal)

	// Relay os.Interrupt to our channel (os.Interrupt = CTRL+C)
	// Ignore other incoming signals
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block main routine until a signal is received
	// As long as user doesn't press CTRL+C a message is not passed and our main routine keeps running
	<-quit

	// After receiving CTRL+C Properly stop the server
	echoServer.Logger.Info("Stopping the server...")
	defer cancel()
	echoServer.Logger.Info("Closing MongoDB connection.")
	db.Disconnect(ctx)
	echoServer.Logger.Info("Done.")
	if err := echoServer.EchoContext.Shutdown(ctx); err != nil {
		echoServer.Logger.Fatal(err)
	}
}
