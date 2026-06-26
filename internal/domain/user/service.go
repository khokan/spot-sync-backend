package user

import (
	"fmt"
	"spotsync/intetnal/auth"
	"spotsync/intetnal/domain/user/dto"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")

type service struct {
	repo       Repository
	jwtService auth.JWTService
	bcryptCost int
}

func NewService(repo Repository, jwtService auth.JWTService, bcryptCost int) *service {
	return &service{repo: repo, jwtService: jwtService, bcryptCost: bcryptCost}
}

func (s *service) CreateUser(req dto.CreateRequest) (*dto.UserResponse, error) {

	user := User{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	if user.Role == "" {
		user.Role = "driver"
	}

	// hash password and set to user.Password
	err := user.HashPassword(req.Password, s.bcryptCost)
	if err != nil {
		return nil, err
	}

	err = s.repo.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.String(),
	}

	return &response, nil

}

func (s *service) LoginUser(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials // User not found
	}

	// check password
	err = user.CheckPassword(req.Password)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// generate token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token: token,
		User:  s.toUserResponse(user),
	}, nil

}

func (s *service) toUserResponse(user *User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.String(),
	}
}
