package controllers

import (
	"app/database"
	"app/entities"
	"app/redis"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic/v7"
)

type ChatNumber struct {
	ChatNumber int32 `json:"chatNumber"`
}

type ErrorMessage struct {
	ErrorMessage string `json:"ErrorMessage"`
}

type SuccessMessage struct {
	SuccessMessage string `json:"SuccessMessage"`
}

type ElasticSearchRes struct {
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
	Number       int       `json:"number"`
	Message_body string    `json:"message_body" `
	Chat_id      int       `json:"chat_id"`
}

func CreateApplications(w http.ResponseWriter, r *http.Request) {
	var application entities.Application
	body := map[string]interface{}{}
	json.NewDecoder(r.Body).Decode(&body)

	application.Name = body["Name"].(string)
	database.Instance.Create(&application)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(application)
}

func SearchMessage(w http.ResponseWriter, r *http.Request) {
	searchString := mux.Vars(r)["message"]
	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://instabug_elasticsearch:9200"))
	if err != nil {
		// Handle error
		panic(err)
	}

	// Use the IndexExists service to check if a specified0 index exists.
	exists, err := client.IndexExists("messages").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Index exists? %v\n", exists)

	var results []map[string]interface{}

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewRegexpQuery("message_body", fmt.Sprint(".*", searchString, ".*")))

	queryStr, err := searchSource.Source()
	queryJs, err2 := json.Marshal(queryStr)

	if err != nil || err2 != nil {
		fmt.Println("[esclient][GetResponse]err during query marshal=", err, err2)
	}
	fmt.Println("[esclient]Final ESQuery=\n", string(queryJs))

	searchService := client.Search().Index("messages").SearchSource(searchSource)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		fmt.Println("[ProductsES][GetPIds]Error=", err)
		return
	}

	for _, hit := range searchResult.Hits.Hits {
		searchResults := map[string]interface{}{}
		err := json.Unmarshal(hit.Source, &searchResults)
		if err != nil {
			fmt.Println("[Getting Students][Unmarshal] Err=", err)
		}

		delete(searchResults, "created_at")
		delete(searchResults, "updated_at")
		delete(searchResults, "number")
		results = append(results, searchResults)
	}

	json.NewEncoder(w).Encode(results)

}

func GetApplications(w http.ResponseWriter, r *http.Request) {
	var applications []entities.Application
	database.Instance.Find(&applications)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(applications)
}

func UpdateApplication(w http.ResponseWriter, r *http.Request) {
	appToken := mux.Vars(r)["application_token"]
	if checkIfApplicationExists(appToken) == false {
		json.NewEncoder(w).Encode("Application Not Found!")
		return
	}
	var application entities.Application
	body := map[string]interface{}{}
	json.NewDecoder(r.Body).Decode(&body)
	database.Instance.Where("token = ?", appToken).First(&application)
	json.NewDecoder(r.Body).Decode(&application)
	database.Instance.Model(&application).Update("name", body["Name"].(string))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessMessage{
		SuccessMessage: "App Updated Successfully",
	})
}

func CreateChats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var chat entities.Chat
	var numberToInserted int32
	appToken := mux.Vars(r)["application_token"]

	if checkIfApplicationExists(appToken) == true {
		query := fmt.Sprintf("SELECT * FROM `chats` WHERE `application_token` = '%s' ORDER BY number DESC LIMIT 1;", appToken)
		queryRes := map[string]interface{}{}
		database.Instance.Raw(query).Scan(&queryRes)

		if len(queryRes) > 0 {
			numberToInserted = queryRes["number"].(int32) + 1
		} else {
			numberToInserted = 1
		}

		chat.ApplicationToken = appToken
		chat.Number = numberToInserted

		database.Instance.Create(&chat)

		query = fmt.Sprintf("select count(*) as count from `chats` where `application_token` = '%s';", appToken)

		database.Instance.Raw(query).Scan(&queryRes)

		redis.SaveInRedis(fmt.Sprintf("app.%s", appToken), queryRes["count"].(int64))

		chatNumRes := struct {
			ChatNumber int32 `json:"chatNumber"`
		}{
			ChatNumber: numberToInserted,
		}

		json.NewEncoder(w).Encode(chatNumRes)
	} else {
		json.NewEncoder(w).Encode(ErrorMessage{
			ErrorMessage: "Please check your App Token",
		})
	}
}

func CreateMessages(w http.ResponseWriter, r *http.Request) {
	var message entities.Message
	var chatNumberStr string
	var chatID int

	body := map[string]interface{}{}
	json.NewDecoder(r.Body).Decode(&body)
	appToken := mux.Vars(r)["application_token"]

	chatNumberStr = mux.Vars(r)["chat_number"]
	chatNumber, err := strconv.Atoi(chatNumberStr)
	if err != nil {
		panic(err)
	}

	queryChat := fmt.Sprintf("SELECT * FROM `chats` WHERE `application_token` = '%s' and number = '%d' ORDER BY number DESC LIMIT 1;", appToken, chatNumber)
	queryChatRes := map[string]interface{}{}
	database.Instance.Raw(queryChat).Scan(&queryChatRes)

	if len(queryChatRes) == 0 {
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(ErrorMessage{
			ErrorMessage: "Please check your App Token and Chat Number",
		})

	} else {
		chatID = int(queryChatRes["id"].(int32))
		queryMessage := fmt.Sprintf("SELECT * FROM `messages` WHERE chat_id = '%d' ORDER BY number DESC LIMIT 1;", chatID)
		queryMessageRes := map[string]interface{}{}

		database.Instance.Raw(queryMessage).Scan(&queryMessageRes)

		if queryMessageRes == nil || len(queryMessageRes) == 0 {
			message.Number = 1
		} else {
			message.Number = int(queryMessageRes["number"].(int32)) + 1
		}

		message.ChatID = int32(queryChatRes["id"].(int32))
		message.MessageBody = body["MessageBody"].(string)

		database.Instance.Create(&message)

		queryCache := fmt.Sprintf("SELECT count(*) as count FROM `messages` WHERE chat_id = '%d';", chatID)
		queryCacheRes := map[string]interface{}{}

		database.Instance.Raw(queryCache).Scan(&queryCacheRes)

		redis.SaveInRedis(fmt.Sprintf("chat.%d", chatID), queryCacheRes["count"].(int64))

		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(SuccessMessage{
			SuccessMessage: "Message Created Successfully",
		})
	}

}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	var messages []entities.Message
	var chat entities.Chat

	appToken := mux.Vars(r)["application_token"]

	chatNumberStr := mux.Vars(r)["chat_number"]

	chatNumber, err := strconv.Atoi(chatNumberStr)
	if err == nil {
		fmt.Println("chatNumber:", chatNumber)
	}

	database.Instance.Where("application_token = ? AND number = ?", appToken, chatNumber).First(&chat)
	database.Instance.Where("chat_id = ?", chat.ID).Find(&messages)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)

}

func checkIfApplicationExists(application_token string) bool {
	var application entities.Application
	query := fmt.Sprintf("SELECT * FROM `applications` WHERE `token` = '%s' LIMIT 1;", application_token)
	database.Instance.Raw(query).Scan(&application)
	if application.Token == "" {
		return false
	}
	return true
}
