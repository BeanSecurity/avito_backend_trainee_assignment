package repository

import (
	"database/sql"
	_ "github.com/bmizerany/pq"

	"avito_chat/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) (*PostgresRepository, error) {
	// 	err := db.Ping()
	// 	if err != nil {
	// 		log.Println(err)
	// 		return &PostgresRepository{}, err
	// 	}

	// 	q := `
	// CREATE TABLE IF NOT EXISTS user (
	// 	id SERIAL PRIMARY KEY,
	// 	username TEXT NOT NULL,
	// 	created_at TIMESTAMP NOT NULL
	// );

	// CREATE TABLE IF NOT EXISTS chat (
	// 	id SERIAL PRIMARY KEY,
	// 	username TEXT NOT NULL,
	// 	created_at TIMESTAMP NOT NULL
	// );

	// CREATE TABLE IF NOT EXISTS chat_user (
	// 	chat_id INTEGER REFERENCES chat,
	// 	user_id INTEGER REFERENCES user,
	// );

	// CREATE TABLE IF NOT EXISTS message (
	// 	id SERIAL PRIMARY KEY,
	// 	content TEXT NOT NULL,
	// 	user_id INTEGER NOT NULL REFERENCES user,
	// 	chat_id INTEGER NOT NULL REFERENCES chat,
	// 	UNIQUE(name,user_id)
	// );
	// `

	// 	_, err = db.Exec(q)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return &PostgresRepository{}, err
	// 	}
	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) AddUser(user models.User) (int, error) {
	return 0, nil
}

func (r *PostgresRepository) AddChat(chat models.Chat) (int, error) {
	return 0, nil
}

func (r *PostgresRepository) AddMessage(chat models.Message) (int, error) {
	return 0, nil
}

func (r *PostgresRepository) GetChats(user models.User) ([]models.Chat, error) {
	return []models.Chat{
			models.Chat{ID: 2, Name: "a"},
			models.Chat{ID: 1, Name: "ab"}},
		nil
}

func (r *PostgresRepository) GetMessages(chat models.Chat) ([]models.Message, error) {
	return []models.Message{
			models.Message{ID: 2, Text: "a"},
			models.Message{ID: 1, Text: "ab"}},
		nil
}
