package helper

import (
	"errors"
	"fmt"
	"go-ecommerce-app/internal/domain"
	"log"
	"strings"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Secret string
}

func SetupAuth(s string) Auth {
	return Auth{Secret: s}
}

func (a Auth) CreateHashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("password must be at least 6 characters long")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		// Log error
		log.Printf("error while creating hash password: %v\n", err)
		return "", err
	}

	return string(hash), nil
}

func (a Auth) VerifyPassword(plainPassword string, hash string) error {
	if len(plainPassword) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainPassword))

	if err != nil {
		return errors.New("invalid password")
	}

	return nil
}

func (a Auth) GenerateToken(id uint, email string, role string) (string, error) {
	if id == 0 || email == "" || role == "" {
		return "", errors.New("invalid user data")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.Secret))

	if err != nil {
		log.Printf("error while signing token: %v\n", err)
		return "", errors.New("error signing token")
	}

	return tokenString, nil
}

func (a Auth) VerifyToken(t string) (domain.User, error) {
	tokenArr := strings.Split(t, " ")

	if len(tokenArr) != 2 {
		return domain.User{}, errors.New("invalid token format")
	}

	if tokenArr[0] != "Bearer" {
		return domain.User{}, errors.New("invalid token type")
	}

	tokenStr := tokenArr[1]

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header)
		}

		return []byte(a.Secret), nil
	})

	if err != nil {
		return domain.User{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return domain.User{}, errors.New("token expired")
		}

		user := domain.User{}
		user.ID = uint(claims["user_id"].(float64))
		user.Email = claims["email"].(string)
		user.UserType = claims["role"].(string)

		return user, nil
	}

	return domain.User{}, errors.New("invalid token")
}

func (a Auth) Authorize(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	user, err := a.VerifyToken(authHeader)

	if err != nil || user.ID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
			"reason":  err,
		})
	}

	ctx.Locals("user", user)
	return ctx.Next()
}

func (a Auth) GetCurrentUser(ctx *fiber.Ctx) (domain.User, error) {
	user := ctx.Locals("user").(domain.User)
	return user, nil
}
