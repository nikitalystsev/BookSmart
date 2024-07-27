package postgres

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/intfRepo"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BookRepo struct {
	db *sqlx.DB
}

func NewBookRepo(db *sqlx.DB) intfRepo.IBookRepo {
	return &BookRepo{db: db}
}

func (br BookRepo) Create(ctx context.Context, book *models.BookModel) error {
	query := `INSERT INTO book VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := br.db.ExecContext(ctx, query, book.ID, book.Title, book.Author, book.Publisher,
		book.CopiesNumber, book.Rarity, book.Genre, book.PublishingYear, book.Language, book.AgeLimit)
	if err != nil {
		return fmt.Errorf("error inserting book: %w", err)
	}

	return nil
}

func (br BookRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.BookModel, error) {
	var book models.BookModel

	query := `SELECT id, title, author, publisher, copiesnumber, rarity, genre, publishingyear, language, agelimit FROM book WHERE id = $1`

	err := br.db.GetContext(ctx, &book, query, id)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (br BookRepo) GetByTitle(ctx context.Context, title string) (*models.BookModel, error) {
	var book models.BookModel

	query := `SELECT id, title, author, publisher, copiesnumber, rarity, genre, publishingyear, language, agelimit FROM book WHERE title = $1`

	err := br.db.GetContext(ctx, &book, query, title)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (br BookRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM book WHERE id = $1`

	_, err := br.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error deleting book: %v", err)
	}

	return nil
}

func (br BookRepo) Update(ctx context.Context, book *models.BookModel) error {
	query := `UPDATE book SET copiesnumber = $1 WHERE id = $2`

	_, err := br.db.ExecContext(ctx, query, book.CopiesNumber, book.ID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error updating book copies: %v", err)
	}

	return nil
}

// GetByParams будет уточняться
func (br BookRepo) GetByParams(ctx context.Context, params *dto.BookParamsDTO) ([]*models.BookModel, error) {
	var books []*models.BookModel
	query := `SELECT id, title, author, publisher, copiesnumber, rarity, genre, publishingyear, language, agelimit 
	          FROM book 
	          WHERE ($1 = '' OR title ILIKE '%' || $1 || '%') AND 
	                ($2 = '' OR author ILIKE '%' || $2 || '%') AND 
	                ($3 = '' OR publisher ILIKE '%' || $3 || '%') AND 
	                ($4 = 0 OR copies_number = $4) AND 
	                ($5 = '' OR rarity ILIKE '%' || $5 || '%') AND 
	                ($6 = '' OR genre ILIKE '%' || $6 || '%') AND 
	                ($7 = 0 OR publishing_year = $7) AND 
	                ($8 = '' OR language ILIKE '%' || $8 || '%') AND 
	                ($9 = 0 OR age_limit = $9)
	          LIMIT $10 OFFSET $11`

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

	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var book models.BookModel
		if err = rows.StructScan(&book); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return books, nil
}
