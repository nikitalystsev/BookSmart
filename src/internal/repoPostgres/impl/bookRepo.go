package impl

import (
	"BookSmart-repositories/errs"
	"BookSmart-repositories/intf"
	"BookSmart-services/dto"
	"BookSmart-services/models"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type BookRepo struct {
	db     *sqlx.DB
	logger *logrus.Entry
}

func NewBookRepo(db *sqlx.DB, logger *logrus.Entry) intf.IBookRepo {
	return &BookRepo{db: db, logger: logger}
}

func (br *BookRepo) Create(ctx context.Context, book *models.BookModel) error {
	br.logger.Infof("inserting book with ID: %s", book.ID)

	query := `INSERT INTO bs.book VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := br.db.ExecContext(ctx, query, book.ID, book.Title, book.Author, book.Publisher,
		book.CopiesNumber, book.Rarity, book.Genre, book.PublishingYear, book.Language, book.AgeLimit)
	if err != nil {
		br.logger.Errorf("error inserting book: %v", err)
		return err
	}

	br.logger.Infof("inserted book with ID: %s", book.ID)

	return nil
}

func (br *BookRepo) GetByID(ctx context.Context, ID uuid.UUID) (*models.BookModel, error) {
	br.logger.Infof("selecting book with ID: %s", ID)

	query := `SELECT * FROM bs.book WHERE id = $1`

	var book models.BookModel
	err := br.db.GetContext(ctx, &book, query, ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		br.logger.Errorf("error selecting book with ID: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		br.logger.Warnf("book with this ID not found %s", ID)
		return nil, errs.ErrNotFound
	}

	br.logger.Infof("selected book with ID: %s", ID)

	return &book, nil
}

func (br *BookRepo) GetByTitle(ctx context.Context, title string) (*models.BookModel, error) {
	br.logger.Infof("selecting book by title: %s", title)

	query := `SELECT * FROM bs.book WHERE title = $1`

	var book models.BookModel
	err := br.db.GetContext(ctx, &book, query, title)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		br.logger.Errorf("error selecting book by title: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		br.logger.Warnf("book with this title not found: %s", title)
		return nil, errs.ErrNotFound
	}

	br.logger.Infof("selected book with title: %s", title)

	return &book, nil
}

func (br *BookRepo) Delete(ctx context.Context, ID uuid.UUID) error {
	br.logger.Infof("deleting book with ID: %s", ID)

	query := `DELETE FROM bs.book WHERE id = $1`

	_, err := br.db.ExecContext(ctx, query, ID)
	if err != nil {
		br.logger.Errorf("error deleting book: %v", err)
		return err
	}

	br.logger.Infof("deleted book with ID: %s", ID)

	return nil
}

func (br *BookRepo) Update(ctx context.Context, book *models.BookModel) error {
	br.logger.Infof("updating book with ID: %s", book.ID)

	query := `UPDATE bs.book SET copies_number = $1 WHERE id = $2`

	_, err := br.db.ExecContext(ctx, query, book.CopiesNumber, book.ID)
	if err != nil {
		br.logger.Errorf("error updating book: %v", err)
		return err
	}

	br.logger.Infof("updated book with ID: %s", book.ID)

	return nil
}

func (br *BookRepo) GetByParams(ctx context.Context, params *dto.BookParamsDTO) ([]*models.BookModel, error) {
	br.logger.Infof("selecting books with params")

	query := `SELECT * 
	          FROM bs.book 
	          WHERE ($1 = '' OR title ILIKE '%' || $1 || '%') AND 
	                ($2 = '' OR author ILIKE '%' || $2 || '%') AND 
	                ($3 = '' OR publisher ILIKE '%' || $3 || '%') AND 
	                ($4 = 0 OR copies_number = $4) AND 
	                ($5 = '' OR rarity::text = $5) AND 
	                ($6 = '' OR genre ILIKE '%' || $6 || '%') AND 
	                ($7 = 0 OR publishing_year = $7) AND 
	                ($8 = '' OR language ILIKE '%' || $8 || '%') AND 
	                ($9 = 0 OR age_limit = $9)
	          LIMIT $10 OFFSET $11`

	var books []*models.BookModel

	err := br.db.SelectContext(ctx, &books, query,
		params.Title,
		params.Author,
		params.Publisher,
		params.CopiesNumber,
		params.Rarity,
		params.Genre,
		params.PublishingYear,
		params.Language,
		params.AgeLimit,
		params.Limit,
		params.Offset,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		br.logger.Errorf("error selecting books with params")
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		br.logger.Warnf("books not found with this params")
		return nil, errs.ErrNotFound
	}

	br.logger.Infof("found %d books", len(books))

	return books, nil
}
