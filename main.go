package main

import (
	"context"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	models "project1/model"
	"time"
)

var bot *linebot.Client

func main() {

	// Get the mongoClient
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the MongoDB server
	err = mongoClient.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Create a new mongoRepository object
	mongoRepository := models.NewMongoRepository(mongoClient)

	// Create a new Gin router
	r := gin.Default()

	// Define your API routes here

	// Handle the LINE webhook
	r.POST("/webhook", func(c *gin.Context) {

		var err error
		secret := "9a157e33616ad20e32a9dfa66f69cd79"
		token := "cgqL7YQYGjIf1YuBLiMJCJNarWaylBPVmZC7YbO3qAtI+qkVNQa/2fViOGZV6JQWeSvNGloXJEA9+JOWbs+SwW63NuoRkT60dON32BKKoi+QkrN9q1Ut3PK95fLAGyPVPsCtAapu32u4XflPH0GVnAdB04t89/1O/w1cDnyilFU="
		bot, err = linebot.New(secret, token)
		if err != nil {
			log.Println(err)
		}

		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				log.Print(err)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do()

					mongoRepository := models.NewMongoRepository(mongoClient)
					messageObj := models.Message{
						UserID:    event.Source.UserID,
						Message:   message.Text,
						CreatedAt: time.Now(),
					}
					err = mongoRepository.SaveMessage(context.Background(), &messageObj)
					if err != nil {
						return
					}
				}
			}
		}
	})

	// Send a message to a user
	r.POST("/send-message", sendMessageToLine(mongoRepository))

	// Get all of the user messages for a given user ID
	r.GET("/messages/:user_id", getMessagesFromMongoDB(mongoRepository))

	// Start the Gin server
	log.Fatal(r.Run(":8080"))
}

// Send a message to a user
func sendMessageToLine(mongoRepository *models.MongoRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		var err error

		secret := "9a157e33616ad20e32a9dfa66f69cd79"
		token := "cgqL7YQYGjIf1YuBLiMJCJNarWaylBPVmZC7YbO3qAtI+qkVNQa/2fViOGZV6JQWeSvNGloXJEA9+JOWbs+SwW63NuoRkT60dON32BKKoi+QkrN9q1Ut3PK95fLAGyPVPsCtAapu32u4XflPH0GVnAdB04t89/1O/w1cDnyilFU="
		bot, err = linebot.New(secret, token)
		if err != nil {
			log.Println(err)
		}

		http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
			events, err := bot.ParseRequest(req)
			if err != nil {
				if err == linebot.ErrInvalidSignature {
					w.WriteHeader(400)
				} else {
					w.WriteHeader(500)
				}
				return
			}
			for _, event := range events {
				if event.Type == linebot.EventTypeMessage {
					switch message := event.Message.(type) {
					case *linebot.TextMessage:

						// 遇到用戶輸入文字，傳一樣的訊息給用戶
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
							log.Print(err)
						}
					}
				}
			}
		})
	}
}

// Get all of the user messages for a given user ID
func getMessagesFromMongoDB(mongoRepository *models.MongoRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user ID
		userID := c.Query("user_id")

		var err error

		// Get the message list from MongoDB
		messageList, err := mongoRepository.GetMessagesByUserID(context.Background(), userID)
		if err != nil {
			log.Println(err)
		}

		// Respond with the message list
		c.JSON(http.StatusOK, messageList)
	}
}
