package usecase

import (
	"context"
	"errors"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/models"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/repository"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/usecase"
	"mnc-fullstack-technical-test/tahap-2/internal/infrastructure/auth"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo   repository.UserRepository
	jwtService *auth.JWTService
}

func NewUserUsecase(userRepo repository.UserRepository, jwtService *auth.JWTService) usecase.UserUsecase {
	return &userUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (u *userUsecase) Register(ctx context.Context, req *usecase.RegisterRequest) (*models.User, error) {
	existing, _ := u.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if existing != nil {
		return nil, errors.New("Phone Number already registered")
	}

	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:          uuid.New(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		PIN:         string(hashedPIN),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) Login(ctx context.Context, req *usecase.LoginRequest) (*usecase.LoginResponse, error) {
	user, err := u.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return nil, errors.New("Phone Number and PIN doesn't match")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PIN), []byte(req.PIN))
	if err != nil {
		return nil, errors.New("Phone Number and PIN doesn't match")
	}

	accessToken, err := u.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &usecase.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *userUsecase) UpdateProfile(ctx context.Context, userID uuid.UUID, req *usecase.UpdateProfileRequest) (*models.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Address != "" {
		user.Address = req.Address
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return u.userRepo.FindByID(ctx, userID)
}
