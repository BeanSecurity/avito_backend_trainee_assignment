package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"

	"avito_chat/models"
	"avito_chat/repository"
)

var repo *repository.PostgresRepository

func main() {
	var err error
	// db, err := sql.Open("postgres", "")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// defer db.Close()
	// repo, err = repository.NewPostgresRepository(db)

	repo, err = repository.NewPostgresRepository(nil)
	if err != nil {
		log.Println(err)
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Route("/users", func(r chi.Router) {
		r.Post("/add", addUser)
	})

	r.Route("/chats", func(r chi.Router) {
		r.Post("/get", getChats)
		r.Post("/add", addChat)
	})

	r.Route("/messages", func(r chi.Router) {
		r.Post("/get", getMessages)
		r.Post("/add", sendMessage)
	})

	http.ListenAndServe(":9000", r)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	_, err := r.Body.Read(body)
	log.Printf("body %s\n", body)

	user := models.User{}
	if err = json.Unmarshal(body, &user); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := repo.AddUser(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userIDRes, err := json.Marshal(userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(userIDRes)
}

func getChats(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	_, err := r.Body.Read(body)
	log.Printf("body %s\n", body)

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := models.User{ID: int(data["user"].(float64))}

	log.Printf("user id: %v\n", data)
	log.Printf("user: %+v\n", user)

	chats, err := repo.GetChats(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatsRes, err := json.Marshal(chats)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(chatsRes)
}

func addChat(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	_, err := r.Body.Read(body)
	log.Printf("body %s\n", body)

	chat := models.Chat{}

	err = json.Unmarshal(body, &chat)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, err := repo.AddChat(chat)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatIDRes, err := json.Marshal(chatID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(chatIDRes)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	_, err := r.Body.Read(body)
	log.Printf("body %s\n", body)

	message := models.Message{}

	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messageID, err := repo.AddMessage(message)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messageIDRes, err := json.Marshal(messageID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(messageIDRes)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	_, err := r.Body.Read(body)
	log.Printf("body %s\n", body)

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chat := models.Chat{ID: int(data["chat"].(float64))}

	log.Printf("chat id: %v\n", data)
	log.Printf("chat: %+v\n", chat)

	messages, err := repo.GetMessages(chat)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messagesRes, err := json.Marshal(messages)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(messagesRes)
}
