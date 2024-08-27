package integrationTests

import (
	"BookSmart-services/errs"
	"context"
	"github.com/google/uuid"
)

func (s *IntegrationTestSuite) TestLibCard_Create_Success() {

	readerID, _ := uuid.Parse("8d9b001f-5760-4c40-bc60-988e0ca54d18")

	err := s.libCardService.Create(context.Background(), readerID)
	s.NoError(err)
}

func (s *IntegrationTestSuite) TestLibCard_Create_Error() {

	readerID, _ := uuid.Parse("3885b2d3-ef6e-4f62-8f86-d1454d108207")

	err := s.libCardService.Create(context.Background(), readerID)
	s.Error(err)
	s.Equal(errs.ErrLibCardAlreadyExist, err)
}

func (s *IntegrationTestSuite) TestLibCard_Update_Success() {
	readerID, err := uuid.Parse("6800b3ee-9810-450e-9ca5-776aa1c6191d")

	libCard, err := s.libCardRepo.GetByReaderID(context.Background(), readerID)
	s.NoError(err)
	s.NotNil(libCard)

	err = s.libCardService.Update(context.Background(), libCard)
	s.NoError(err)
}

func (s *IntegrationTestSuite) TestLibCard_Update_Error() {
	readerID, err := uuid.Parse("75919792-c2d9-4685-92b2-e2a80b2ed5be")

	libCard, err := s.libCardService.GetByReaderID(context.Background(), readerID)
	s.NoError(err)
	s.NotNil(libCard)

	err = s.libCardService.Update(context.Background(), libCard)
	s.Error(err)
	s.Equal(errs.ErrLibCardIsValid, err)
}
