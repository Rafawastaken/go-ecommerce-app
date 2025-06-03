package service

import (
	"errors"
	"fmt"
	"go-ecommerce-app/config"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/pkg/notifications"
	"log"
	"time"
)

type UserService struct {
	Repo   repository.UserRepository
	Auth   helper.Auth
	Config config.AppConfig
}

func (s UserService) findUserByEmail(email string) (*domain.User, error) {
	user, err := s.Repo.FindUser(email)
	return &user, err
}

func (s UserService) SignUp(input dto.UserSignup) (string, error) {
	log.Printf("create user: %v\n", input)

	hPassword, err := s.Auth.CreateHashPassword(input.Password)
	if err != nil {
		return "", err
	}

	user, err := s.Repo.CreateUser(domain.User{
		Email:    input.Email,
		Password: hPassword,
		Phone:    input.Phone,
	})

	//Generate token
	log.Printf("user created: %v\n", user)

	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s UserService) Login(email string, password string) (string, error) {
	log.Printf("Login attempt for user: %s\n", email)

	user, err := s.findUserByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := s.Auth.VerifyPassword(password, user.Password); err != nil {
		return "", errors.New("wrong password")
	}

	token, err := s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
	if err != nil {
		return "", errors.New("error while generating token")
	}

	return token, nil
}

func (s UserService) isVerifiedUser(id uint) bool {
	currentUser, err := s.Repo.FindUserById(id)
	return err == nil && currentUser.Verified
}

func (s UserService) GetVerificationCode(e domain.User) error {
	//check if user is verified
	if s.isVerifiedUser(e.ID) {
		return errors.New("user is verified")
	}

	//generate verification code
	code, err := s.Auth.GenerateCode()

	if err != nil {
		return nil
	}

	//update user
	user := domain.User{
		Expiry: time.Now().Add(30 * time.Minute),
		Code:   code,
	}

	_, err = s.Repo.UpdateUser(e.ID, user)

	if err != nil {
		return errors.New("unable to update verification code")
	}

	//get the user phone
	user, _ = s.Repo.FindUserById(e.ID)

	//send sms
	notificationClient := notifications.NewNotificationClient(s.Config)
	message := fmt.Sprintf("Verification code: %v", code)
	err = notificationClient.SendSMS(user.Phone, message)

	if err != nil {
		log.Printf("Unable to send verification code: %v", err)
		return errors.New("unable to send verification code")
	}

	return nil
}

func (s UserService) VerifyCode(id uint, code int) error {
	if s.isVerifiedUser(id) {
		return errors.New("user is already verified")
	}

	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return err
	}

	if user.Code != code {
		return errors.New("invalid code")
	}

	if !time.Now().Before(user.Expiry) {
		return errors.New("user is expired")
	}

	updateUser := domain.User{
		Verified: true,
	}

	_, err = s.Repo.UpdateUser(id, updateUser)
	if err != nil {
		return errors.New("unable to update verification code")
	}

	return nil
}

func (s UserService) CreateProfile(id uint, input any) error {
	return nil
}

func (s UserService) GetProfile(id uint) (*domain.User, error) {
	return nil, nil
}

func (s UserService) UpdateProfile(id uint, input any) (*domain.User, error) {
	return nil, nil
}

func (s UserService) BecomeSeller(id uint, input dto.SellerInput) (string, error) {
	user, _ := s.Repo.FindUserById(id)

	if user.UserType == domain.SELLER {
		return "", errors.New("you have already joined seller program")
	}

	seller, err := s.Repo.UpdateUser(id, domain.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		UserType:  domain.SELLER,
	})

	if err != nil {
		return "", err
	}

	token, err := s.Auth.GenerateToken(user.ID, user.Email, seller.UserType)

	err = s.Repo.CreateBankAccount(domain.BankAccount{
		BankAccountNumber: input.BankAccountNumber,
		SwiftCode:         input.SwiftCode,
		PaymentType:       input.PaymentType,
		UserID:            id,
	})

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s UserService) FindCart(id uint) ([]interface{}, error) {
	return nil, nil
}

func (s UserService) CreateCard(input any, u domain.User) ([]interface{}, error) {
	return nil, nil
}

func (s UserService) CreateOrder(u domain.User) (int, error) {
	return 0, nil
}

func (s UserService) GetOrders(input any, u domain.User) ([]interface{}, error) {
	return nil, nil
}

func (s UserService) GetOrderById(id uint, uId uint) ([]interface{}, error) {
	return nil, nil
}
