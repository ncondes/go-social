package handlers

import (
	"net/http"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/dtos"
	"github.com/ncondes/go/social/internal/logging"
)

type AuthHandler struct {
	userService domain.UserServiceInterface
	validator   *Validator
	logger      logging.Logger
}

func NewAuthHandler(userService domain.UserServiceInterface, validator *Validator, logger logging.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		validator:   validator,
		logger:      logger,
	}
}

// RegisterUser godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user with an invitation sent via email or SMS
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.RegisterUserInput		true	"Register user input"
//	@Success		201		{object}	domain.UserWithInvitationToken	"User registered successfully"
//	@Failure		400		{object}	dtos.ErrorsResponseDTO			"Validation errors"
//	@Failure		409		{object}	dtos.ErrorResponseDTO			"Conflict error"
//	@Failure		500		{object}	dtos.ErrorResponseDTO			"Internal server error"
//	@Router			/auth/register [post]
func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Validate request body
	var registerUserInput *domain.RegisterUserInput

	if err := jsonDecode(w, r, &registerUserInput); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := h.validator.validateStruct(registerUserInput); err != nil {
		respondWithErrors(w, http.StatusBadRequest, err, h.logger)
		return
	}

	registerUserInput = &domain.RegisterUserInput{
		FirstName:        registerUserInput.FirstName,
		LastName:         registerUserInput.LastName,
		Username:         registerUserInput.Username,
		Email:            registerUserInput.Email,
		Password:         registerUserInput.Password,
		InvitationMethod: registerUserInput.InvitationMethod,
	}

	userWithInvitationToken, err := h.userService.RegisterUserWithInvitation(r.Context(), registerUserInput)
	if err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	respondWithData(w, http.StatusCreated, userWithInvitationToken, h.logger)
}

// ActivateUser godoc
//
//	@Summary		Activate a user account
//	@Description	Activate a user account using the token sent via invitation
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dtos.ActivateUserDTO	true	"Activation token"
//	@Success		204		"No content"
//	@Failure		400		{object}	dtos.ErrorsResponseDTO	"Validation errors"
//	@Failure		404		{object}	dtos.ErrorResponseDTO	"Not found error"
//	@Failure		500		{object}	dtos.ErrorResponseDTO	"Internal server error"
//	@Router			/auth/activate [put]
func (h *AuthHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	var activateUserDTO *dtos.ActivateUserDTO

	if err := jsonDecode(w, r, &activateUserDTO); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := h.validator.validateStruct(activateUserDTO); err != nil {
		respondWithErrors(w, http.StatusBadRequest, err, h.logger)
		return
	}

	if err := h.userService.ActivateUser(r.Context(), activateUserDTO.Token); err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
