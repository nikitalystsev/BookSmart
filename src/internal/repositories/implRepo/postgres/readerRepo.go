package postgres

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errsRepo"
	"BookSmart/internal/repositories/intfRepo"
	"context"
	"database/sql"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

type ReaderRepo struct {
	db     *sqlx.DB
	client *redis.Client
	logger *logrus.Entry
}

func NewReaderRepo(db *sqlx.DB, client *redis.Client, logger *logrus.Entry) intfRepo.IReaderRepo {
	return &ReaderRepo{db: db, client: client, logger: logger}
}

func (rr *ReaderRepo) Create(ctx context.Context, reader *models.ReaderModel) error {
	rr.logger.Infof("inserting reader with ID: %s", reader.ID)

	query := `INSERT INTO bs.reader VALUES ($1, $2, $3, $4, $5)`

	rr.logger.Infof("executing query: %s", query)

	_, err := rr.db.ExecContext(ctx, query, reader.ID, reader.Fio, reader.PhoneNumber, reader.Age, reader.Password)
	if err != nil {
		rr.logger.Errorf("error inserting reader: %v", err)
		return err
	}

	rr.logger.Infof("reader with ID: %s inserted successfully", reader.ID)

	return nil
}

func (rr *ReaderRepo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.ReaderModel, error) {
	rr.logger.Infof("selected reader by phoneNumber: %s", phoneNumber)

	query := `SELECT id, fio, phone_number, age, password FROM bs.reader WHERE phone_number = $1`

	rr.logger.Infof("executing query: %s", query)

	var reader models.ReaderModel
	err := rr.db.GetContext(ctx, &reader, query, phoneNumber)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selected reader by phoneNumber: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("no reader found by phoneNumber: %s", phoneNumber)
		return nil, errsRepo.ErrNotFound
	}

	rr.logger.Infof("successfully selected reader: %v", reader)

	return &reader, nil
}

func (rr *ReaderRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.ReaderModel, error) {
	rr.logger.Infof("select reader with ID: %s", id)

	query := `SELECT id, fio, phone_number, age, password FROM bs.reader WHERE id = $1`

	rr.logger.Infof("executing query: %s", query)

	var reader models.ReaderModel
	err := rr.db.GetContext(ctx, &reader, query, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selected reader by id: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("no reader found by id: %v", id)
		return nil, errsRepo.ErrNotFound
	}

	rr.logger.Infof("successfully selected reader: %v", reader)

	return &reader, nil
}

func (rr *ReaderRepo) IsFavorite(ctx context.Context, readerID, bookID uuid.UUID) (bool, error) {
	rr.logger.Infof("book already is favorite?")

	query := `SELECT COUNT(*) FROM bs.favorite_books WHERE reader_id = $1 AND book_id = $2`

	rr.logger.Infof("executing query: %s", query)

	var count int
	err := rr.db.GetContext(ctx, &count, query, readerID, bookID)
	if err != nil {
		rr.logger.Errorf("error checking favorite book: %v", err)
		return false, err
	}

	rr.logger.Infof("successfully checked favorite book")

	return count > 0, nil
}

func (rr *ReaderRepo) AddToFavorites(ctx context.Context, readerID, bookID uuid.UUID) error {
	rr.logger.Infof("add book to favorites: %s", bookID)

	query := `INSERT INTO bs.favorite_books (reader_id, book_id) VALUES ($1, $2)`

	rr.logger.Infof("executing query: %s", query)

	_, err := rr.db.ExecContext(ctx, query, readerID, bookID)
	if err != nil {
		rr.logger.Errorf("error adding book to favorites: %v", err)
		return err
	}

	rr.logger.Infof("book with ID: %s successfully added in favorites", bookID)

	return nil
}

func (rr *ReaderRepo) SaveRefreshToken(ctx context.Context, id uuid.UUID, token string, ttl time.Duration) error {
	rr.logger.Infof("saving refresh token in redis")

	err := rr.client.Set(ctx, token, id.String(), ttl).Err()
	if err != nil {
		rr.logger.Errorf("error saving refresh token: %v", err)
		return err
	}

	rr.logger.Infof("successfully saving refresh token in redis")

	return nil
}

func (rr *ReaderRepo) GetByRefreshToken(ctx context.Context, token string) (*models.ReaderModel, error) {
	rr.logger.Infof("getting reader by refresh token: %s", token)

	var readerID uuid.UUID

	readerIDStr, err := rr.client.Get(ctx, token).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		rr.logger.Errorf("error getting reader by refresh token: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, redis.Nil) {
		rr.logger.Errorf("no reader found by refresh token: %s", token)
		return nil, errsRepo.ErrNotFound
	}

	readerID, err = uuid.Parse(readerIDStr)
	if err != nil {
		rr.logger.Errorf("error parsing reader by refresh token: %v", err)
		return nil, err
	}

	var reader models.ReaderModel

	query := `SELECT id, fio, phone_number, age, password FROM bs.reader WHERE id = $1`

	rr.logger.Infof("executing query: %s", query)

	err = rr.db.GetContext(ctx, &reader, query, readerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selected reader by id: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("no reader found by id: %v", readerID)
		return nil, errsRepo.ErrNotFound
	}

	rr.logger.Infof("successfully getting reader by refresh token: %v", reader)

	return &reader, nil
}
