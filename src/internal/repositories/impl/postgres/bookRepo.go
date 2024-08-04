package postgres

import (
	"BookSmart-repositories/errs"
	"BookSmart-repositories/intf"
	"BookSmart-services/dto"
	"BookSmart-services/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	br.logger.Infof("executing query: %s", query)

	_, err := br.db.ExecContext(ctx, query, book.ID, book.Title, book.Author, book.Publisher,
		book.CopiesNumber, book.Rarity, book.Genre, book.PublishingYear, book.Language, book.AgeLimit)
	if err != nil {
		br.logger.Errorf("error inserting book: %v", err)
		return err
	}

	br.logger.Infof("book with ID: %s inserted successfully", book.ID)

	return nil
}

func (br *BookRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.BookModel, error) {
	br.logger.Infof("select book with ID: %s", id)

	query := `SELECT id, title, author, publisher, copies_number, rarity, genre, publishing_year, language, age_limit FROM bs.book WHERE id = $1`

	br.logger.Infof("executing query: %s", query)

	var book models.BookModel
	err := br.db.GetContext(ctx, &book, query, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		br.logger.Errorf("error selected book by id: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		br.logger.Warnf("no book found by id: %s", id.String())
		return nil, errs.ErrNotFound
	}

	br.logger.Infof("successfully selected book: %v", book)

	return &book, nil
}

func (br *BookRepo) GetByTitle(ctx context.Context, title string) (*models.BookModel, error) {
	br.logger.Infof("selected book by title: %s", title)

	query := `SELECT id, title, author, publisher, copies_number, rarity, genre, publishing_year, language, age_limit FROM bs.book WHERE title = $1`

	br.logger.Infof("executing query: %s", query)

	var book models.BookModel
	err := br.db.GetContext(ctx, &book, query, title)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		br.logger.Errorf("error selected book by title: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		br.logger.Warnf("no book found by title: %s", title)
		return nil, errs.ErrNotFound
	}

	br.logger.Infof("successfully selected book: %v", book)

	return &book, nil
}

func (br *BookRepo) Delete(ctx context.Context, id uuid.UUID) error {
	br.logger.Infof("deleting book with ID: %s", id)

	query := `DELETE FROM bs.book WHERE id = $1`

	br.logger.Infof("executing query: %s", query)

	_, err := br.db.ExecContext(ctx, query, id)
	if err != nil {
		br.logger.Errorf("error deleting book with ID %s: %v", id, err)
		return err
	}

	br.logger.Infof("successfully deleted book with ID: %s", id)

	return nil
}

func (br *BookRepo) Update(ctx context.Context, book *models.BookModel) error {
	br.logger.Infof("updating book with ID: %s", book.ID)

	query := `UPDATE bs.book SET copies_number = $1 WHERE id = $2`

	br.logger.Infof("executing query: %s", query)

	_, err := br.db.ExecContext(ctx, query, book.CopiesNumber, book.ID)
	if err != nil {
		br.logger.Errorf("error updating book copies: %v", err)
		return err
	}

	br.logger.Infof("successfully updated book copies for ID: %s", book.ID)

	return nil
}

// GetByParams будет уточняться
func (br *BookRepo) GetByParams(ctx context.Context, params *dto.BookParamsDTO) ([]*models.BookModel, error) {
	br.logger.Infof("selecting books with params: %+v", params)

	query := `SELECT id, title, author, publisher, copies_number, rarity, genre, publishing_year, language, age_limit 
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

	br.logger.Infof("executing query")

	rows, err := br.db.QueryxContext(ctx, query,
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
		br.logger.Errorf("error selecting books with params: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		br.logger.Warnf("no books found by params")
		return nil, errs.ErrNotFound
	}

	defer func(rows *sqlx.Rows) {
		err = rows.Close()
		if err != nil {
			br.logger.Errorf("error closing rows: %v", err)
			fmt.Printf("error closing rows: %v", err)
		}
	}(rows)

	var books []*models.BookModel
	for rows.Next() {
		var book models.BookModel
		if err = rows.StructScan(&book); err != nil {
			br.logger.Errorf("error scanning row: %v", err)
			return nil, err
		}

		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		br.logger.Errorf("rows iteration error: %v", err)
		return nil, err
	}

	br.logger.Infof("successfully found %d books with params: %+v", len(books), params)

	return books, nil
}
