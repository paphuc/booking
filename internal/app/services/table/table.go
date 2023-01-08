package table

import (
	"booking/configs"
	"booking/internal/app/types"
	"booking/internal/pkg/glog"
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Repository is an interface of a table repository
type Repository interface {
	FindByID(ctx context.Context, id string) (*types.Table, error)
	Insert(ctx context.Context, Table types.Table) error
	UpdateTableByID(ctx context.Context, UpdateTableRequest types.UpdateTableRequest) error
	DeleteTable(ctx context.Context, DeleteTableRequest types.DeleteTableRequest) error
}

// Service is an table service
type Service struct {
	conf   *configs.Configs
	em     *configs.ErrorMessage
	repo   Repository
	logger glog.Logger
}

// NewService return a new member service
func NewService(c *configs.Configs, e *configs.ErrorMessage, r Repository, l glog.Logger) *Service {
	return &Service{
		conf:   c,
		em:     e,
		repo:   r,
		logger: l,
	}
}

// Post basic
func (s *Service) InsertTable(ctx context.Context, tableReq types.TableRequest) (*types.Table, error){

	Table := types.Table{
		ID:       primitive.NewObjectID(),
		Status:	  tableReq.Status,
		Type:	  tableReq.Type,
		Slots:    tableReq.Slots,
		DelFlg:	  false,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	err := s.repo.Insert(ctx, Table)
	if err != nil {
		s.logger.Errorf("Can't create table", err)
		return nil, errors.Wrap(err, "Can't create table")
	}

	s.logger.Infof("Create succesfully!!!", tableReq)
	return &Table, nil
}

// Put service update info for table by ID
func (s *Service) UpdateTableByID(ctx context.Context, table types.UpdateTableRequest) error {

	// Check table is existed or not by ID
	if _,err := s.repo.FindByID(ctx,table.ID); err != nil {
		s.logger.Errorf("Table is not existed !!!", err)
		return errors.Wrap(err, "Table existed, can't update Table")
	}

	err := s.repo.UpdateTableByID(ctx, table)

	if err != nil {
		s.logger.Errorf("Failed when update table by id", err)
		return err
	}

	s.logger.Infof("Updated table is completed !!!")
	return err
}

// Put service delete table by ID 
func (s *Service) DeleteTable(ctx context.Context, table types.DeleteTableRequest) error {

	// Check table is existed or not by ID
	if _,err := s.repo.FindByID(ctx,table.ID); err != nil {
		s.logger.Errorf("Table is not existed !!!", err)
		return errors.Wrap(err, "Table existed, can't update Table")
	}

	err := s.repo.DeleteTable(ctx, table)

	if err != nil {
		s.logger.Errorf("Failed when delete table by id", err)
		return err
	}

	s.logger.Infof("Delete table is completed !!!")
	return err
}

