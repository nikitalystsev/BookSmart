package integrationTests

import (
	"BookSmart/internal/services/errsService"
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
	s.Equal(errsService.ErrLibCardAlreadyExist, err)
}

func (s *IntegrationTestSuite) TestLibCard_Update_Success() {

	// чит (так нельзя)
	libCard, err := s.libCardRepo.GetByNum(context.Background(), "7945544456734")
	s.NoError(err)
	s.NotNil(libCard)

	err = s.libCardService.Update(context.Background(), libCard)
	s.NoError(err)
}

func (s *IntegrationTestSuite) TestLibCard_Update_Error() {

	// чит (так нельзя)
	libCard, err := s.libCardRepo.GetByNum(context.Background(), "4654645456328")
	s.NoError(err)
	s.NotNil(libCard)

	err = s.libCardService.Update(context.Background(), libCard)
	s.Error(err)
	s.Equal(errsService.ErrLibCardIsValid, err)
}
