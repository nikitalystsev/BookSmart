package impl

import (
	"BookSmart-services/core/models"
	"BookSmart-services/errs"
	"BookSmart-services/impl"
	"BookSmart-services/intfRepo"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ReservationRepo struct {
	db     *mongo.Collection
	logger *logrus.Entry
}

func NewReservationRepo(db *mongo.Database, logger *logrus.Entry) intfRepo.IReservationRepo {
	return &ReservationRepo{db: db.Collection("reservation"), logger: logger}
}

func (rr *ReservationRepo) Create(ctx context.Context, reservation *models.ReservationModel) error {
	rr.logger.Infof("inserting reservation with ID: %s", reservation.ID)

	_, err := rr.db.InsertOne(ctx, reservation)
	if err != nil {
		rr.logger.Errorf("error inserting reservation: %v", err)
		return err
	}

	rr.logger.Infof("inserted reservation with ID: %s", reservation.ID)

	return nil
}

func (rr *ReservationRepo) GetByReaderAndBook(ctx context.Context, readerID, bookID uuid.UUID) (*models.ReservationModel, error) {
	rr.logger.Infof("find reservation with readerID и bookID: %s и %s", readerID, bookID)

	one := rr.db.FindOne(ctx, bson.M{"reader_id": readerID, "book_id": bookID})

	if one.Err() != nil && !errors.Is(one.Err(), mongo.ErrNoDocuments) {
		rr.logger.Errorf("error selecting reservation: %v", one.Err())
		return nil, one.Err()
	}
	if one.Err() != nil && errors.Is(one.Err(), mongo.ErrNoDocuments) {
		rr.logger.Warnf("reservation with this readerID и bookID not found: %s и %s", readerID, bookID)
		return nil, errs.ErrReservationDoesNotExists
	}

	var reservation models.ReservationModel
	if err := one.Decode(&reservation); err != nil {
		rr.logger.Errorf("error decoding reservation: %v", err)
		return nil, err
	}

	rr.logger.Infof("found reservation with readerID и bookID: %s и %s", readerID, bookID)

	return &reservation, nil
}

func (rr *ReservationRepo) GetByID(ctx context.Context, ID uuid.UUID) (*models.ReservationModel, error) {
	rr.logger.Infof("find reservation with ID: %s", ID)

	one := rr.db.FindOne(ctx, bson.M{"_id": ID})

	if one.Err() != nil && !errors.Is(one.Err(), mongo.ErrNoDocuments) {
		rr.logger.Errorf("error find reservation: %v", one.Err())
		return nil, one.Err()
	}
	if one.Err() != nil && errors.Is(one.Err(), mongo.ErrNoDocuments) {
		rr.logger.Warnf("reservation with this ID not found: %s", ID)
		return nil, errs.ErrReservationDoesNotExists
	}

	var reservation models.ReservationModel
	if err := one.Decode(&reservation); err != nil {
		rr.logger.Errorf("error decoding reservation: %v", err)
		return nil, err
	}

	rr.logger.Infof("found reservation with ID: %s", ID)

	return &reservation, nil
}

func (rr *ReservationRepo) Update(ctx context.Context, reservation *models.ReservationModel) error {
	rr.logger.Infof("updating reservation with ID: %s", reservation.ID)

	one, err := rr.db.UpdateOne(
		ctx, bson.M{"_id": reservation.ID},
		bson.M{"$set": bson.M{"issue_date": reservation.IssueDate, "return_date": reservation.ReturnDate, "state": reservation.State}},
	)
	if err != nil {
		rr.logger.Errorf("error updating reservation with ID: %v", err)
		return err
	}

	if one.MatchedCount == 0 {
		rr.logger.Warnf("reservation with this ID not found: %v", reservation.ID)
		return errs.ErrReservationDoesNotExists
	}

	rr.logger.Infof("updated reservation with ID: %s", reservation.ID)

	return nil
}

func (rr *ReservationRepo) GetExpiredByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	rr.logger.Infof("find expired reservations with readerID: %s", readerID)

	filter := bson.M{
		"reader_id":   readerID,
		"return_date": bson.M{"$lt": time.Now()},
	}

	cursor, err := rr.db.Find(ctx, filter)
	if err != nil {
		rr.logger.Errorf("error find expired reservations: %v", err)
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			fmt.Println("error close cursor")
		}
	}(cursor, ctx)

	var reservations []*models.ReservationModel
	if err = cursor.All(ctx, &reservations); err != nil {
		rr.logger.Printf("error decoding reservations: %v", err)
		return nil, err
	}

	if len(reservations) == 0 {
		rr.logger.Warnf("expired reservations with this readerID not found: %s", readerID)
		return nil, errs.ErrReservationDoesNotExists
	}

	rr.logger.Infof("found %d expired reservations with readerID %s", len(reservations), readerID)

	return reservations, nil
}

func (rr *ReservationRepo) GetActiveByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	rr.logger.Infof("find active reservations with readerID: %s", readerID)

	filter := bson.M{
		"reader_id": readerID,
		"state": bson.M{
			"$nin": []string{impl.ReservationExpired, impl.ReservationClosed},
		},
	}

	cursor, err := rr.db.Find(ctx, filter)
	if err != nil {
		rr.logger.Errorf("error find active reservations: %v", err)
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			fmt.Println("error close cursor")
		}
	}(cursor, ctx)

	var reservations []*models.ReservationModel
	if err = cursor.All(ctx, &reservations); err != nil {
		rr.logger.Printf("error decoding reservations: %v", err)
		return nil, err
	}

	if len(reservations) == 0 {
		rr.logger.Warnf("active reservations with this readerID not found: %s", readerID)
		return nil, errs.ErrReservationDoesNotExists
	}

	rr.logger.Infof("found %d active reservations with readerID %s", len(reservations), readerID)

	return reservations, nil
}
