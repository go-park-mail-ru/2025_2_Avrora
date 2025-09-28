package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-park-mail-ru/2025_2_Avrora/db"
	"github.com/go-park-mail-ru/2025_2_Avrora/models"
	"github.com/go-park-mail-ru/2025_2_Avrora/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/utils"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[\p{L}\p{N}._%+-]+@[\p{L}\p{N}.-]+\.[\p{L}]{2,}$`)
	return re.MatchString(email)
}

func validatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	var hasUpper, hasLower, hasDigit bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

// --- Валидация для логина ---
func validateLoginRequest(req *LoginRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email не может быть пустым")
	}
	if strings.TrimSpace(req.Password) == "" {
		return errors.New("пароль не может быть пустым")
	}
	return nil
}

// --- Валидация для регистрации ---
func validateRegisterRequest(req *RegisterRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email не может быть пустым")
	}
	if !validateEmail(req.Email) {
		return errors.New("невалидный формат email")
	}
	if !validatePassword(req.Password) {
		return errors.New("пароль должен быть не менее 6 символов, иметь разный регистр, иметь хотя бы одну цифру")
	}
	return nil
}

// --- Регистрация ---
func RegisterHandler(w http.ResponseWriter, r *http.Request, repo *db.Repo, jwtGen *utils.JwtGenerator, passwordHasher *utils.PasswordHasher) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.HandleError(w, err, http.StatusBadRequest, "невалидный json")
		return
	}

	if err := validateRegisterRequest(&req); err != nil {
		response.HandleError(w, err, http.StatusBadRequest, err.Error())
		return
	}

	user, err := repo.User().FindByEmail(req.Email)
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка бд при поиске пользователя")
		return
	}
	if user != nil {
		response.HandleError(w, nil, http.StatusConflict, "Пользователь с таким email уже существует")
		return
	}

	hashedPassword, err := passwordHasher.Hash(req.Password)
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка хеширования пароля")
		return
	}

	newUser := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}
	err = repo.User().Create(newUser)
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка бд при создании пользователя")
		return
	}

	token, err := jwtGen.GenerateJWT(strconv.Itoa(newUser.ID))
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка генерации jwt")
		return
	}

	response.WriteJSON(w, http.StatusCreated, AuthResponse{
		Token: token,
		User: &models.User{
			Email: req.Email,
		},
	})
}

// --- Логин ---
func LoginHandler(w http.ResponseWriter, r *http.Request, repo *db.Repo, jwtGen *utils.JwtGenerator, passwordHasher *utils.PasswordHasher) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.HandleError(w, err, http.StatusBadRequest, "невалидный json")
		return
	}

	if err := validateLoginRequest(&req); err != nil {
		response.HandleError(w, err, http.StatusBadRequest, err.Error())
		return
	}

	user, err := repo.User().FindByEmail(req.Email)
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка бд при поиске пользователя")
		return
	}
	if user == nil {
		response.HandleError(w, nil, http.StatusUnauthorized, "невалидные учетные данные")
		return
	}

	if !passwordHasher.Compare(req.Password, user.Password) {
		response.HandleError(w, nil, http.StatusUnauthorized, "невалидные учетные данные")
		return
	}

	token, err := jwtGen.GenerateJWT(strconv.Itoa(user.ID))
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка генерации jwt")
		return
	}

	response.WriteJSON(w, http.StatusOK, AuthResponse{
		Token: token,
		User: &models.User{
			Email: req.Email,
		},
	})
}

// --- Логаут ---
func LogoutHandler(w http.ResponseWriter, r *http.Request, jwtGen *utils.JwtGenerator) {
	expiredToken, err := jwtGen.GenerateExpiredJWT()
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка генерации jwt")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "успешный логаут",
		"Token":   expiredToken,
	})
}
