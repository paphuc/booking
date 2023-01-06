package member

import (
	"booking/configs"
	"booking/internal/app/types"
	"booking/internal/pkg/glog"
	"booking/internal/pkg/jwt"
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository is an interface of a member repository
type Repository interface {
	FindByID(ctx context.Context, id string) (*types.Member, error)
	Insert(ctx context.Context, Member types.Member) error
	UpdateMemberByID(ctx context.Context, Member types.Member) error
	FindByEmail(ctx context.Context, email string) (*types.Member, error)
}

// Service is an member service
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

// Get return given member by his/her id
func (s *Service) Get(ctx context.Context, id string) (*types.Member, error) {
	return s.repo.FindByID(ctx, id)
}

// Post basic
func (s *Service) InsertMember(ctx context.Context, memreq types.MemberRequest) (*types.Member, error) {

	user := types.Member{
		ID:       primitive.NewObjectID(),
		Name:     memreq.Name,
		Password: memreq.Password,
		Email:    memreq.Email,
	}

	err := s.repo.Insert(ctx, user)
	if err != nil {
		s.logger.Errorf("Can't create member", err)
		return nil, errors.Wrap(err, "Can't create member")
	}

	s.logger.Infof("Create succesfully!!!", memreq)
	return &user, nil
}

// Put service update info for member by ID
func (s *Service) UpdateMemberByID(ctx context.Context, mem types.Member) error {

	err := s.repo.UpdateMemberByID(ctx, mem)

	if err != nil {
		s.logger.Errorf("Failed when update member by id", err)
		return err
	}

	s.logger.Infof("Updated member is completed !!!")
	return err
}
func (s *Service) Login(ctx context.Context, UserLogin types.UserLogin) (*types.UserResponseSignUp, error) {

	user, err := s.repo.FindByEmail(ctx, UserLogin.Email)
	if err != nil {
		s.logger.Errorf("Not found email exits", err)
		return nil, errors.Wrap(errors.New("Not found email exits"), "Email not exists, can't find user")
	}

	if !jwt.IsCorrectPassword(UserLogin.Password, user.Password) {
		s.logger.Errorf("Password incorrect", UserLogin.Email)
		return nil, errors.Wrap(errors.New("Password isn't like password from database"), "Password incorrect")
	}

	var tokenString string
	tokenString, error := jwt.GenToken(types.UserFieldInToken{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email}, s.conf.Jwt.Duration)

	if error != nil {
		s.logger.Errorf("Can not gen token", error)
		return nil, errors.Wrap(error, "Can't gen token")
	}
	s.logger.Infof("Login completed ", user.Email)
	return &types.UserResponseSignUp{
		Name:  user.Name,
		Email: user.Email,
		Token: tokenString}, nil
}
