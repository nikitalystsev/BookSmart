package integrationTests

import (
	"context"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
)

func (s *IntegrationTestSuite) TestBook_Create_Success() {
	book := &models.BookModel{
		ID:             uuid.New(),
		Title:          "Test Book",
		Author:         "Test Author",
		Publisher:      "Test Publisher",
		CopiesNumber:   10,
		Rarity:         "Common",
		Genre:          "Fiction",
		PublishingYear: 2021,
		Language:       "English",
		AgeLimit:       12,
	}

	err := s.bookService.Create(context.Background(), book)
	s.NoError(err)

	expectedBook, err := s.bookService.GetByID(context.Background(), book.ID)
	s.NoError(err)
	s.Equal(book, expectedBook)
}

func (s *IntegrationTestSuite) TestBook_Create_Error() {
	book := &models.BookModel{
		ID:             uuid.New(),
		Title:          "Test Book",
		Author:         "Test Author",
		Publisher:      "Test Publisher",
		CopiesNumber:   10,
		Genre:          "Fiction",
		PublishingYear: 2021,
		Language:       "English",
		AgeLimit:       12,
	}

	err := s.bookService.Create(context.Background(), book)
	s.Error(err)
	s.Equal(errs.ErrEmptyBookRarity, err)
}

func (s *IntegrationTestSuite) TestBook_Delete_Success() {

	id, err := uuid.Parse("55cbe7ec-66f1-4383-b52d-9bc8c7d33480")

	book, err := s.bookService.GetByID(context.Background(), id)
	s.NoError(err)

	err = s.bookService.Delete(context.Background(), book.ID)
	s.NoError(err)

	expectedBook, err := s.bookService.GetByID(context.Background(), id)
	s.Error(errs.ErrBookDoesNotExists)
	s.Nil(expectedBook)
}

func (s *IntegrationTestSuite) TestBook_Delete_Error() {

	id, err := uuid.Parse("305c0d87-6599-4589-8337-d55ba937898a")

	book, err := s.bookService.GetByID(context.Background(), id)
	s.Error(err)
	s.Equal(errs.ErrBookDoesNotExists, err)
	s.Nil(book)

	err = s.bookService.Delete(context.Background(), uuid.Nil)
	s.Error(err)
	s.Equal(errs.ErrBookObjectIsNil, err)
}
