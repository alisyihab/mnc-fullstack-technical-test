package handlers

import (
	"mnc-fullstack-technical-test/tahap-2/internal/delivery/http/response"
	"mnc-fullstack-technical-test/tahap-2/internal/domain/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(u usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: u}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with first name, last name, phone number, address, and pin
// @Tags users
// @Accept json
// @Produce json
// @Param request body usecase.RegisterRequest true "Register request"
// @Success 200 {object} response.JSONResponse
// @Failure 400 {object} response.JSONResponse
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req usecase.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userUsecase.Register(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, registerResponse{
		UserID:      user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
		CreatedDate: user.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

// Login godoc
// @Summary Login user
// @Description Login user with phone number and pin
// @Tags users
// @Accept json
// @Produce json
// @Param request body usecase.LoginRequest true "Login request"
// @Success 200 {object} response.JSONResponse
// @Failure 400 {object} response.JSONResponse
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req usecase.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.userUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, res)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user's first name, last name, or address
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecase.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} response.JSONResponse
// @Failure 400 {object} response.JSONResponse
// @Router /profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req usecase.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userUsecase.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, updateProfileResponse{
		UserID:      user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Address:     user.Address,
		UpdatedDate: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}
