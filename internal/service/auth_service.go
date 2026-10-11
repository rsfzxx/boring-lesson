package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"go-starter/internal/dto"
	"go-starter/internal/repository"
)

var ErrInvalidCredentials = errors.New("nomor telepon atau password salah")

type Claims struct {
	RoleID  string `json:"role_id"`
	GroupID string `json:"group_id"`
	jwt.RegisteredClaims
}

type AuthService struct {
	users     *repository.UserRepository
	jwtSecret []byte
	jwtTTL    time.Duration
}

func NewAuthService(users *repository.UserRepository, jwtSecret string, jwtTTL time.Duration) *AuthService {
	return &AuthService{users: users, jwtSecret: []byte(jwtSecret), jwtTTL: jwtTTL}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.users.FindByPhoneNumber(ctx, req.PhoneNumber)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	now := time.Now()
	claims := Claims{
		RoleID:  user.RoleID.Hex(),
		GroupID: user.GroupID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.Hex(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtTTL)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Message:     "Login successful",
		AccessToken: token,
		ExpiresIn:   int(s.jwtTTL.Seconds()),
	}, nil
}