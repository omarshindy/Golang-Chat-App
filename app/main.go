package main

import (
	"app/controllers"
	"app/database"
	"app/entities"
	"app/redis"
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

var DB *gorm.DB
var ctx = context.TODO()

func main() {

	// Load Configurations from config.json using Viper
	LoadAppConfig()

	// Initialize Database
	database.Connect(AppConfig.ConnectionString)

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)

	// Caching to redis
	c := cron.New()
	c.AddFunc("*/5 * * * *", func() { cronCacheSave() })
	c.Start()

	// Register Routes
	RegisterProductRoutes(router)

	// Start the server
	log.Println(fmt.Sprintf("Starting Server on port %s", AppConfig.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", AppConfig.Port), router))
}

// Application and Chat Routers
func RegisterProductRoutes(router *mux.Router) {
	// Creating application
	router.HandleFunc("/api/application", controllers.CreateApplications).Methods("POST")

	// Updating application
	router.HandleFunc("/api/application/{application_token}", controllers.UpdateApplication).Methods("POST")

	// listing all applications
	router.HandleFunc("/api/applications", controllers.GetApplications).Methods("GET")

	// Creating chat for specific app
	router.HandleFunc("/api/application/{application_token}/chat", controllers.CreateChats).Methods("POST")

	// creating message for specific chat
	router.HandleFunc("/api/application/{application_token}/{chat_number}/message", controllers.CreateMessages).Methods("POST")

	// getting all messages for specific chat
	router.HandleFunc("/api/application/{application_token}/{chat_number}/messages", controllers.GetMessages).Methods("GET")

	// Serching for Messages via ElasticSearch
	router.HandleFunc("/api/search/{message}", controllers.SearchMessage).Methods("GET")

}

func cronCacheSave() {
	iter := redis.Ring.Scan(ctx, 0, "", 0).Iterator()

	var chatsCount int64
	var messagesCount int64

	for iter.Next(ctx) {

		key := iter.Val()
		if key == "ping" {
			continue
		}

		app_match, _ := regexp.MatchString("^app", key)
		if app_match == true {
			exactKey := strings.Split(key, ".")
			chatsCount = redis.GetFromRedis(key)
			var application entities.Application
			database.Instance.Model(&application).Where("token = ?", string(exactKey[1])).Update("chats_count", chatsCount)
		} else {
			exactKey := strings.Split(key, ".")
			chatID, err := strconv.Atoi(exactKey[1])
			if err != nil {
				fmt.Println("error:", err)
			}
			messagesCount = redis.GetFromRedis(key)
			var chat entities.Chat

			database.Instance.Model(&chat).Where("id = ?", chatID).Update("messages_count", messagesCount)
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	time.Sleep(2000)
}
