package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/Masterminds/squirrel"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// UserRepository defines the interface that the PostgresUserRepository must implement
type UserRepository interface {
	CreateUser(user *User) error
	UpdateUser(user *User) error
	GetUserByID(userID int) (*User, error)
	PatchUser(userID int, updates map[string]interface{}) error
	GetAllUsers() ([]User, error)
}

// Ensure PostgresUserRepository implements UserRepository
var _ UserRepository = &PostgresUserRepository{}

// ErrUserNotFound is returned when a user is not found in the database.
var ErrUserNotFound = errors.New("user not found")

type PostgresUserRepository struct {
	db   *sql.DB
	psql squirrel.StatementBuilderType
}

// NewPostgresUserRepository initializes a new PostgresUserRepository with an optional *sql.DB parameter.
func NewPostgresUserRepository(db *sql.DB) (*PostgresUserRepository, error) {
	if db == nil {
		// Load .env file if present
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file, using system environment variables")
		}

		// Construct the connection string
		connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_SSLMODE"))

		// Open the database connection
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %v", err)
		}
	}

	// Return the initialized repository
	return &PostgresUserRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

// Close closes the database connection when done.
func (r *PostgresUserRepository) Close() {
	r.db.Close()
}

// CreateUser inserts a new user into the database.
func (r *PostgresUserRepository) CreateUser(user *User) error {
	query, args, err := r.psql.Insert("public.users").
		Columns("user_name", "first_name", "last_name", "email", "user_status", "department").
		Values(user.User_name, user.First_name, user.Last_name, user.Email, user.User_status, user.Department).
		Suffix("RETURNING user_id").ToSql()

	if err != nil {
		return err
	}
	ProduceKafkaMessage(user, "prism-user-create")
	return r.db.QueryRow(query, args...).Scan(&user.User_id)
}

// UpdateUser updates the entire user record in the database.
func (r *PostgresUserRepository) UpdateUser(user *User) error {
	query, args, err := r.psql.Update("public.users").
		Set("user_name", user.User_name).
		Set("first_name", user.First_name).
		Set("last_name", user.Last_name).
		Set("email", user.Email).
		Set("user_status", user.User_status).
		Set("department", user.Department).
		Where(squirrel.Eq{"user_id": user.User_id}).ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, args...)
	ProduceKafkaMessage(user, "prism-user-update")
	return err
}

// PatchUser updates specific fields of a user in the database.
func (r *PostgresUserRepository) PatchUser(userID int, updates map[string]interface{}) error {
	queryBuilder := r.psql.Update("public.users").Where(squirrel.Eq{"user_id": userID})

	// Dynamically add fields to update
	for key, value := range updates {
		queryBuilder = queryBuilder.Set(key, value)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, args...)
	return err
}

// GetUserByID fetches a user by their ID.
func (r *PostgresUserRepository) GetUserByID(userID int) (*User, error) {
	var user User
	query, args, err := r.psql.Select("user_id", "user_name", "first_name", "last_name", "email", "user_status", "department").
		From("public.users").
		Where(squirrel.Eq{"user_id": userID}).ToSql()

	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(query, args...).Scan(&user.User_id, &user.User_name, &user.First_name, &user.Last_name, &user.Email, &user.User_status, &user.Department)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetAllUsers fetches all users from the database.
func (r *PostgresUserRepository) GetAllUsers() ([]User, error) {
	query, args, err := r.psql.Select("user_id", "user_name", "first_name", "last_name", "email", "user_status", "department").
		From("public.users").ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.User_id, &user.User_name, &user.First_name, &user.Last_name, &user.Email, &user.User_status, &user.Department)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// DeleteUserByID deletes a user by their ID from the database.
func (r *PostgresUserRepository) DeleteUserByID(userID int) error {
	// get the user info from database so it can be used in Kafka
	user, err2 := r.GetUserByID(userID)

	if err2 != nil {
		return err2
	}

	query, args, err := r.psql.Delete("public.users").
		Where(squirrel.Eq{"user_id": userID}).ToSql()

	if err != nil {
		return err
	}

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	ProduceKafkaMessage(user, "prism-user-delete")
	return nil
}

// prism-user-create
// prism-user-delete
// prism-user-update
func ProduceKafkaMessage(user *User, topic string) error {
	broker := "localhost:9092"

	// Serialize user struct to JSON
	userData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error serializing user data: %w", err)
	}

	// Create Kafka producer config with idempotent option enabled
	config := sarama.NewConfig()
	config.Producer.Idempotent = true // Enable idempotent producer
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	// Create Kafka producer
	producer, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		return fmt.Errorf("error creating Kafka producer: %w", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("error closing Kafka producer: %v", err)
		}
	}()

	// Generate a unique key using user_id and datetime stamp
	timestamp := time.Now().Format(time.RFC3339) // Using RFC3339 for a standard format
	key := fmt.Sprintf("%d_%s", user.User_id, timestamp)

	// Create Kafka message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(userData),
	}

	// Produce the message
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("error producing to Kafka: %w", err)
	}

	fmt.Printf("Produced message to topic %s, partition %d, offset %d, key: %s\n", topic, partition, offset, key)
	return nil
}
