package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	"log"
	"time"

	"avito_chat/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) (*PostgresRepository, error) {
	err := db.Ping()
	if err != nil {
		log.Println(err)
		return &PostgresRepository{}, err
	}

	q := `
	CREATE TABLE IF NOT EXISTS member (
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS chat (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS member_chat(
		chat_id INTEGER REFERENCES chat,
		member_id INTEGER REFERENCES member,
		UNIQUE(chat_id,member_id)
	);

	CREATE TABLE IF NOT EXISTS message (
		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL,
		member_id INTEGER NOT NULL REFERENCES member,
		chat_id INTEGER NOT NULL REFERENCES chat,
		created_at TIMESTAMP NOT NULL DEFAULT now()
	);
	`

	_, err = db.Exec(q)
	if err != nil {
		log.Println(err)
		return &PostgresRepository{}, err
	}
	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) AddUser(user models.User) (int, error) {
	userInsertQ := "INSERT INTO member (username) VALUES ($1) RETURNING member.id;"
	stmt, err := r.db.Prepare(userInsertQ)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer stmt.Close()

	var userID int
	err = stmt.QueryRow(user.Username).Scan(&userID)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return userID, nil
}

func (r *PostgresRepository) AddChat(chat models.Chat) (int, error) {
	//insert chat
	chatInsertQ := "INSERT INTO chat (name) VALUES ($1) RETURNING chat.id;"
	stmt, err := r.db.Prepare(chatInsertQ)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer stmt.Close()

	var chatID int
	err = stmt.QueryRow(chat.Name).Scan(&chatID)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	//insert users and chat relations
	usersInsertQ := "INSERT INTO member_chat (chat_id,member_id) VALUES ($1,$2);"
	stmt2, err := r.db.Prepare(usersInsertQ)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer stmt2.Close()

	for _, userID := range chat.Users {
		_, err := stmt2.Exec(chatID, userID)
		if err != nil {
			log.Println(err)
			return 0, err
		}
	}
	return chatID, nil
}

func (r *PostgresRepository) AddMessage(message models.Message) (int, error) {
	stmtCheck, err := r.db.Prepare(`
SELECT
	EXISTS(
		SELECT
			*
		FROM
			member_chat
		WHERE
			member_id = $1 AND chat_id = $2
	);
`)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer stmtCheck.Close()

	var exist bool
	err = stmtCheck.QueryRow(message.AuthorID, message.ChatID).Scan(&exist)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if !exist {
		err = fmt.Errorf("user (id:%d) not in chat (id:%d)", message.AuthorID, message.ChatID)
		log.Println(err)
		return 0, err
	}

	messageInsertQ := "INSERT INTO message (content, member_id, chat_id) VALUES ($1, $2, $3) RETURNING message.id;"
	stmt, err := r.db.Prepare(messageInsertQ)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer stmt.Close()

	var messageID int
	err = stmt.QueryRow(message.Text, message.AuthorID, message.ChatID).Scan(&messageID)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return messageID, nil
}

func (r *PostgresRepository) GetChats(user models.User) ([]models.Chat, error) {
	query := `
	SELECT
		id,
		name,
		created_at
	FROM
		chat JOIN member_chat ON chat.id = member_chat.chat_id
	WHERE
		member_id = $1;
	`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	var chatID int
	var chatName string
	var chatTime time.Time
	rows, err := stmt.Query(user.ID)
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		rows.Scan(&chatID, &chatName, &chatTime)
		chats = append(chats, models.Chat{ID: chatID, Name: chatName, CreatedAt: chatTime})
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return chats, nil
}

func (r *PostgresRepository) GetMessages(chat models.Chat) ([]models.Message, error) {
	query := `
SELECT * FROM message WHERE chat_id = $1;
	`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	var messageID int
	var messageText string
	var chatID int
	var authorID int
	var messageTime time.Time
	rows, err := stmt.Query(chat.ID)
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		rows.Scan(&messageID, &messageText, &chatID, &authorID, &messageTime)
		messages = append(messages, models.Message{
			ID:        messageID,
			Text:      messageText,
			CreatedAt: messageTime,
			ChatID:    chatID,
			AuthorID:  authorID,
		})
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return messages, nil
}
