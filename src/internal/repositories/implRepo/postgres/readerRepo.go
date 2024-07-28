package postgres

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/intfRepo"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type ReaderRepo struct {
	db     *sqlx.DB
	client *redis.Client
}

func NewReaderRepo(db *sqlx.DB, client *redis.Client) intfRepo.IReaderRepo {
	return &ReaderRepo{db: db, client: client}
}

func (rr *ReaderRepo) Create(ctx context.Context, reader *models.ReaderModel) error {
	query := `INSERT INTO reader VALUES ($1, $2, $3, $4, $5)`

	_, err := rr.db.ExecContext(ctx, query, reader.ID, reader.Fio, reader.PhoneNumber, reader.Age, reader.Password)
	if err != nil {
		return fmt.Errorf("error inserting reader: %v", err)
	}

	return nil
}

func (rr *ReaderRepo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.ReaderModel, error) {
	query := `SELECT id, fio, phone_number, age, password FROM reader WHERE phone_number = $1`

	var reader models.ReaderModel
	err := rr.db.GetContext(ctx, &reader, query, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("error fetching reader by phone number: %v", err)
	}

	return &reader, nil
}

func (rr *ReaderRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.ReaderModel, error) {
	query := `SELECT id, fio, phone_number, age, password FROM reader WHERE id = $1`

	var reader models.ReaderModel
	err := rr.db.GetContext(ctx, &reader, query, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching reader by phone number: %v", err)
	}

	return &reader, nil
}

func (rr *ReaderRepo) IsFavorite(ctx context.Context, readerID, bookID uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) FROM favorite_books WHERE reader_id = $1 AND book_id = $2`

	var count int
	err := rr.db.GetContext(ctx, &count, query, readerID, bookID)
	if err != nil {
		return false, fmt.Errorf("error checking if book is favorite: %w", err)
	}

	return count > 0, nil
}

func (rr *ReaderRepo) AddToFavorites(ctx context.Context, readerID, bookID uuid.UUID) error {
	query := `INSERT INTO favorite_books (reader_id, book_id) VALUES ($1, $2)`

	_, err := rr.db.ExecContext(ctx, query, readerID, bookID)
	if err != nil {
		return fmt.Errorf("error adding book to favorites: %w", err)
	}

	return nil
}

func (rr *ReaderRepo) SaveRefreshToken(ctx context.Context, id uuid.UUID, token string, ttl time.Duration) error {
	err := rr.client.Set(ctx, token, id.String(), ttl).Err()
	if err != nil {
		return fmt.Errorf("error saving refresh token: %w", err)
	}

	return nil
}

func (rr *ReaderRepo) GetByRefreshToken(ctx context.Context, token string) (*models.ReaderModel, error) {
	var readerID uuid.UUID

	readerIDStr, err := rr.client.Get(ctx, token).Result()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("refresh token not found")
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving refresh token: %w", err)
	}

	readerID, err = uuid.Parse(readerIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid reader ID: %w", err)
	}

	var reader models.ReaderModel

	query := `SELECT id, fio, phone_number, age, password FROM reader WHERE id = $1`

	err = rr.db.GetContext(ctx, &reader, query, readerID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving reader: %w", err)
	}

	return &reader, nil
}
