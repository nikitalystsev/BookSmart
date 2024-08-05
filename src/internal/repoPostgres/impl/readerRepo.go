package impl

import (
	"BookSmart-repositories/errs"
	"BookSmart-repositories/intf"
	"BookSmart-services/models"
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

func NewReaderRepo(db *sqlx.DB, client *redis.Client, logger *logrus.Entry) intf.IReaderRepo {
	return &ReaderRepo{db: db, client: client, logger: logger}
}

func (rr *ReaderRepo) Create(ctx context.Context, reader *models.ReaderModel) error {
	rr.logger.Infof("inserting reader with ID: %s", reader.ID)

	query := `INSERT INTO bs.reader VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := rr.db.ExecContext(ctx, query, reader.ID, reader.Fio, reader.PhoneNumber,
		reader.Age, reader.Password, reader.Role)
	if err != nil {
		rr.logger.Errorf("error inserting reader: %v", err)
		return err
	}

	rr.logger.Infof("inserted reader with ID: %s", reader.ID)

	return nil
}

func (rr *ReaderRepo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.ReaderModel, error) {
	rr.logger.Infof("selecting reader with phoneNumber: %s", phoneNumber)

	query := `SELECT id, fio, phone_number, age, password, role FROM bs.reader WHERE phone_number = $1`

	var reader models.ReaderModel
	err := rr.db.GetContext(ctx, &reader, query, phoneNumber)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selecting reader by phoneNumber: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("reader with this phoneNumber not found: %s", phoneNumber)
		return nil, errs.ErrNotFound
	}

	rr.logger.Infof("selected reader with phoneNumber: %s", phoneNumber)

	return &reader, nil
}

func (rr *ReaderRepo) GetByID(ctx context.Context, ID uuid.UUID) (*models.ReaderModel, error) {
	rr.logger.Infof("selecting reader with ID: %s", ID)

	query := `SELECT * FROM bs.reader WHERE id = $1`

	var reader models.ReaderModel
	err := rr.db.GetContext(ctx, &reader, query, ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selecting reader with ID: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("reader with this ID not found: %v", ID)
		return nil, errs.ErrNotFound
	}

	rr.logger.Infof("selected reader with ID: %s", ID)

	return &reader, nil
}

func (rr *ReaderRepo) IsFavorite(ctx context.Context, readerID, bookID uuid.UUID) (bool, error) {
	rr.logger.Infof("book with ID = %s already is favorite?", bookID)

	query := `SELECT COUNT(*) FROM bs.favorite_books WHERE reader_id = $1 AND book_id = $2`

	var count int
	err := rr.db.GetContext(ctx, &count, query, readerID, bookID)
	if err != nil {
		rr.logger.Errorf("error checking favorite book: %v", err)
		return false, err
	}

	rr.logger.Infof("checked favorite book")

	return count > 0, nil
}

func (rr *ReaderRepo) AddToFavorites(ctx context.Context, readerID, bookID uuid.UUID) error {
	rr.logger.Infof("reader (ID = %s) adding book (ID = %s) to favorites", readerID, bookID)

	query := `INSERT INTO bs.favorite_books (reader_id, book_id) VALUES ($1, $2)`

	_, err := rr.db.ExecContext(ctx, query, readerID, bookID)
	if err != nil {
		rr.logger.Errorf("error adding book to favorites: %v", err)
		return err
	}

	rr.logger.Infof("reader (ID = %s) added book (ID = %s) to favorites", readerID, bookID)

	return nil
}

func (rr *ReaderRepo) SaveRefreshToken(ctx context.Context, id uuid.UUID, token string, ttl time.Duration) error {
	rr.logger.Infof("saving refresh token in redis")

	err := rr.client.Set(ctx, token, id.String(), ttl).Err()
	if err != nil {
		rr.logger.Errorf("error saving refresh token: %v", err)
		return err
	}

	rr.logger.Infof("refresh token saved in redis")

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
		rr.logger.Errorf("reader with this refresh token not found: %s", token)
		return nil, errs.ErrNotFound
	}

	readerID, err = uuid.Parse(readerIDStr)
	if err != nil {
		rr.logger.Errorf("error parsing readerID by refresh token: %v", err)
		return nil, err
	}

	var reader models.ReaderModel

	query := `SELECT * FROM bs.reader WHERE id = $1`

	err = rr.db.GetContext(ctx, &reader, query, readerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selecting reader by id: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("reader with this ID not found: %v", readerID)
		return nil, errs.ErrNotFound
	}

	rr.logger.Infof("getting reader by refresh token: %v", token)

	return &reader, nil
}
