package service

import (
	"errors"
	"strings"
	"time"

	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	City     string `json:"city"`
	State    string `json:"state"`
	Postcode string `json:"postcode"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (s *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || len(req.Password) < 6 {
		return nil, errors.New("emel sah dan kata laluan sekurang-kurangnya 6 aksara diperlukan")
	}

	// Semak jika emel telah digunakan
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("akaun dengan emel ini telah pun wujud")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        req.Email,
		PasswordHash: string(hashed),
		Phone:        strings.TrimSpace(req.Phone),
		Role:         models.RoleCustomer,
		Address:      strings.TrimSpace(req.Address),
		City:         strings.TrimSpace(req.City),
		State:        strings.TrimSpace(req.State),
		Postcode:     strings.TrimSpace(req.Postcode),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("emel atau kata laluan tidak sah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("emel atau kata laluan tidak sah")
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := middleware.UserClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 hari
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) GetProfile(userID uint) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}
