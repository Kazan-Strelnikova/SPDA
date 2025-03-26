package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/models/user"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetUser(ctx context.Context, email string) (user.User, error)
	InsertUser(ctx context.Context, usr user.User) error
}

func (s *Service) Login(ctx context.Context, email, password string) (user.User, string, error) {
	const op = "service.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	usr, err := s.usrRepo.GetUser(ctx, email)
	if err != nil {
		log.Error("error getting user", slog.String("error", err.Error()))
		return user.User{}, "", fmt.Errorf("invalid email or password")
	}

	if !checkPasswordHash(password, usr.Password) {
		log.Error("incorrect password")
		return user.User{}, "", fmt.Errorf("invalid email or password")
	}

	usr.Password = ""

	token, err := generateTokens(s.tknScrt, usr)
	if err != nil {
		log.Error("error creating token", slog.String("error", err.Error()))
		return user.User{}, "", fmt.Errorf("error creating token")
	}

	return usr, token, nil
}

func (s *Service) LoginByToken(ctx context.Context, token string) (user.User, error) {
	const op = "service.LoginByToken"

	log := s.log.With(
		slog.String("op", op),
	)

	secret := []byte(s.tknScrt)

	data, err := jwt.ParseWithClaims(token, &user.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok && token.Method.Alg() == jwt.SigningMethodHS256.Alg() {
			return secret, nil
		}
		return nil, ErrInvalidToken
	})

	if err != nil {
		log.Error("invalid jwt token", slog.String("err", err.Error()))
		return user.User{}, ErrInvalidToken
	}

	if claims, ok := data.Claims.(*user.UserClaims); ok && data.Valid {
		log.Info("token validated successfully", slog.Any("email", claims.Payload.Email))
		return claims.Payload, nil
	}

	log.Error("invalid jwt token")
	return user.User{}, ErrInvalidToken
}

func (s *Service) Register(ctx context.Context, usr user.User) (string, error) {
	const op = "service.Register"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", usr.Email),
	)

	var err error

	usr.Password, err = hashPassword(usr.Password)
	if err != nil {
		log.Error("error hashing password", slog.String("error", err.Error()))
		return "", fmt.Errorf("invalid password")
	}

	err = s.usrRepo.InsertUser(ctx, usr)
	if err != nil {
		log.Error("error creating user", slog.String("error", err.Error()))
		return "", fmt.Errorf("error creating user")
	}

	usr.Password = ""

	token, err := generateTokens(s.tknScrt, usr)
	if err != nil {
		log.Error("error creating user", slog.String("error", err.Error()))
		return "", fmt.Errorf("error creating user")
	}

	return token, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateTokens(accessSecret string, usr user.User) (string, error) {

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(time.Hour * 24 * 365).Unix(),
		"payload": usr,
	}).SignedString([]byte(accessSecret))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
