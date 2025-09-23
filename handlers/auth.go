package handlers

import (
	"database/sql"
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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "не получилось разобрать json", http.StatusInternalServerError)
	}
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

func validateRegisterRequest(req *RegisterRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email не может быть пустым")
	}
	if !validateEmail(req.Email) {
		return errors.New("невалидный формат email")
	}
	if len(req.Password) < 6 {
		return errors.New("пароль должен быть не менее 6 символов")
	}
	return nil
}

func validateLoginRequest(req *LoginRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email не может быть пустым")
	}
	if !validateEmail(req.Email) {
		return errors.New("невалидный формат email")
	}
	if !validatePassword(req.Password) {
		return errors.New("невалидные учетные данные")
	}

	return nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, repo *db.Repo) {
	if r.Method != http.MethodPost {
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, response.NewErrorResp("невалидный json"))
		return
	}

	if err := validateRegisterRequest(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, response.NewErrorResp(err.Error()))
		return
	}

	user, err := repo.User().FindByEmail(req.Email)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response.NewErrorResp("ошибка бд"))
		return
	}
	if user != nil {
		writeJSON(w, http.StatusConflict, response.NewErrorResp("Пользователь с таким email уже существует"))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response.NewErrorResp("ошибка хеширования пароля"))
		return
	}

	var userID int
	err = repo.User().Create(&models.User{
		Email:    req.Email,
		Password: hashedPassword,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response.NewErrorResp("ошибка бд"))
		return
	}

	token, err := utils.GenerateJWT(strconv.Itoa(userID))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response.NewErrorResp("ошибка генерации jwt"))
		return
	}

	writeJSON(w, http.StatusCreated, AuthResponse{
		Token: token,
		User: &models.User{
			Email: req.Email,
		},
	})
}
func LoginHandler(w http.ResponseWriter, r *http.Request, repo *db.Repo) {
	if r.Method != http.MethodPost {
		http.Error(w, response.NewErrorResp("метод не поддерживается").Error, http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, response.NewErrorResp("невалидный json"))
		return
	}

	if err := validateLoginRequest(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, response.NewErrorResp(err.Error()))
		return
	}

	user, err := repo.User().FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusUnauthorized, response.NewErrorResp("невалидные учетные данные"))
			return
		}
		writeJSON(w, http.StatusInternalServerError, response.NewErrorResp("ошибка бд"))
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		writeJSON(w, http.StatusUnauthorized, response.NewErrorResp("невалидные учетные данные"))
		return
	}

	token, err := utils.GenerateJWT(strconv.Itoa(user.ID))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response.NewErrorResp("ошибка генерации jwt"))
		return
	}

	writeJSON(w, http.StatusOK, AuthResponse{
		Token: token,
		User: &models.User{
			Email: req.Email,
		},
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	expiredToken, err := utils.GenerateExpiredJWT()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, response.NewErrorResp("ошибка генерации jwt"))
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "успещный логаут",
		"Token":   expiredToken,
	})
}
