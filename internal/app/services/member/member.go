package member

import (
	"booking/internal/app/types"
	"booking/configs"
	"booking/internal/pkg/glog"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/pkg/errors"
)

// Repository is an interface of a member repository
type Repository interface {
	FindByID(ctx context.Context, id string) (*types.Member, error)
	Insert(ctx context.Context, Member types.Member) error
	UpdateMemberByID(ctx context.Context, Member types.Member) error
}

// Service is an member service
type Service struct {
	conf   *configs.Configs
	em     *configs.ErrorMessage
	repo   Repository
	logger glog.Logger
}

// NewService return a new member service
func NewService(c *configs.Configs,e *configs.ErrorMessage,r Repository, l glog.Logger) *Service {
	return &Service{
		conf:   c,
		em:     e,
		repo:   r,
		logger: l,
	}
}

// Get return given member by his/her id
func (s *Service) Get(ctx context.Context, id string) (*types.Member, error) {
	return s.repo.FindByID(ctx, id)
}

// Post basic 
func (s *Service) InsertMember(ctx context.Context, memreq types.MemberRequest) (*types.Member,error){

	user := types.Member{
		ID : 		primitive.NewObjectID(),
		Name: 		memreq.Name,
		Password:	memreq.Password,
		Email:		memreq.Email,
	}

	err := s.repo.Insert(ctx,user) 
	if err != nil {
		s.logger.Errorf("Can't create member", err)
		return nil, errors.Wrap(err, "Can't create member")
	}

	s.logger.Infof("Create succesfully!!!", memreq)
	return &user, nil
}

// Put service update info for member by ID
func (s *Service) UpdateMemberByID(ctx context.Context, mem types.Member ) error {

	err := s.repo.UpdateMemberByID(ctx, mem)

	if err != nil {
		s.logger.Errorf("Failed when update member by id", err)
		return err
	}

	s.logger.Infof("Updated member is completed !!!")
	return err
}
