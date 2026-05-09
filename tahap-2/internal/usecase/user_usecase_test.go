package usecase_test

import (
	"context"
	"errors"
	"testing"

	"mnc-fullstack-technical-test/tahap-2/internal/domain/models"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/repository/mocks"
	domainUsecase "mnc-fullstack-technical-test/tahap-2/internal/domain/usecase"
	"mnc-fullstack-technical-test/tahap-2/internal/infrastructure/auth"
	"mnc-fullstack-technical-test/tahap-2/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestUserUsecase(userRepo *mocks.UserRepository) domainUsecase.UserUsecase {
	jwtService := auth.NewJWTService("test-secret")
	return usecase.NewUserUsecase(userRepo, jwtService)
}

// ---------------------------------------------------------------------------
// Register
// ---------------------------------------------------------------------------

func TestRegister_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	uc := newTestUserUsecase(userRepo)

	req := &domainUsecase.RegisterRequest{
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "081234567890",
		Address:     "Jl. Test No. 1",
		PIN:         "123456",
	}

	// Phone not yet registered
	userRepo.On("FindByPhoneNumber", mock.Anything, req.PhoneNumber).
		Return(nil, errors.New("not found"))

	// Create succeeds
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).
		Return(nil)

	user, err := uc.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.FirstName, user.FirstName)
	assert.Equal(t, req.PhoneNumber, user.PhoneNumber)
	userRepo.AssertExpectations(t)
}

func TestRegister_PhoneAlreadyExists(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	uc := newTestUserUsecase(userRepo)

	req := &domainUsecase.RegisterRequest{
		FirstName:   "Jane",
		LastName:    "Doe",
		PhoneNumber: "081234567890",
		Address:     "Jl. Test No. 2",
		PIN:         "654321",
	}

	existing := &models.User{
		ID:          uuid.New(),
		PhoneNumber: req.PhoneNumber,
	}

	userRepo.On("FindByPhoneNumber", mock.Anything, req.PhoneNumber).
		Return(existing, nil)

	user, err := uc.Register(context.Background(), req)

	assert.Nil(t, user)
	assert.EqualError(t, err, "Phone Number already registered")
	userRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Login
// ---------------------------------------------------------------------------

func TestLogin_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	uc := newTestUserUsecase(userRepo)

	// Pre-register to get a bcrypt hash
	jwtService := auth.NewJWTService("test-secret")
	registerUC := usecase.NewUserUsecase(userRepo, jwtService)

	regReq := &domainUsecase.RegisterRequest{
		FirstName:   "Alice",
		LastName:    "Smith",
		PhoneNumber: "089876543210",
		Address:     "Jl. Test No. 3",
		PIN:         "111222",
	}
	userRepo.On("FindByPhoneNumber", mock.Anything, regReq.PhoneNumber).
		Return(nil, errors.New("not found")).Once()
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).
		Return(nil).
		Run(func(args mock.Arguments) {
			// capture the created user so we can look it up on login
		})

	var createdUser *models.User
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).
		Maybe() // already set above

	// Register to get a hashed PIN
	createdUser, err := registerUC.Register(context.Background(), regReq)
	assert.NoError(t, err)

	// Now set up login lookup
	userRepo.On("FindByPhoneNumber", mock.Anything, regReq.PhoneNumber).
		Return(createdUser, nil)

	loginReq := &domainUsecase.LoginRequest{
		PhoneNumber: regReq.PhoneNumber,
		PIN:         regReq.PIN,
	}

	resp, err := uc.Login(context.Background(), loginReq)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.NotEqual(t, resp.AccessToken, resp.RefreshToken)
}

func TestLogin_WrongPIN(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	uc := newTestUserUsecase(userRepo)

	jwtService := auth.NewJWTService("test-secret")
	registerUC := usecase.NewUserUsecase(userRepo, jwtService)

	regReq := &domainUsecase.RegisterRequest{
		FirstName:   "Bob",
		LastName:    "Builder",
		PhoneNumber: "081111111111",
		Address:     "Jl. Test No. 4",
		PIN:         "999888",
	}
	userRepo.On("FindByPhoneNumber", mock.Anything, regReq.PhoneNumber).
		Return(nil, errors.New("not found")).Once()
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).
		Return(nil)

	createdUser, _ := registerUC.Register(context.Background(), regReq)

	userRepo.On("FindByPhoneNumber", mock.Anything, regReq.PhoneNumber).
		Return(createdUser, nil)

	loginReq := &domainUsecase.LoginRequest{
		PhoneNumber: regReq.PhoneNumber,
		PIN:         "000000", // wrong PIN
	}

	resp, err := uc.Login(context.Background(), loginReq)

	assert.Nil(t, resp)
	assert.EqualError(t, err, "Phone Number and PIN doesn't match")
}
