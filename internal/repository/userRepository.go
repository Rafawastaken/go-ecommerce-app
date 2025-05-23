package repository

import (
	"errors"
	"go-ecommerce-app/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type UserRepository interface {
	CreateUser(usr domain.User) (domain.User, error)
	FindUser(email string) (domain.User, error)
	FindUserById(id uint) (domain.User, error)
	UpdateUser(id uint, usr domain.User) (domain.User, error)
	GetVerificationCode(email string) (int, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r userRepository) CreateUser(usr domain.User) (domain.User, error) {
	err := r.db.Create(&usr).Error

	if err != nil {
		log.Printf("database error while creating user: %v\n", err)
		return domain.User{}, errors.New("failed to create user")
	}

	return usr, nil
}

func (r userRepository) FindUser(email string) (domain.User, error) {
	var user domain.User

	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		log.Printf("database error while finding user: %v\n", err)
		return domain.User{}, errors.New("user does not exist")
	}

	return user, nil
}

func (r userRepository) FindUserById(id uint) (domain.User, error) {
	var user domain.User

	err := r.db.First(&user, "id = ?", id).Error

	if err != nil {
		log.Printf("database error while finding user: %v\n", err)
		return domain.User{}, errors.New("user does not exist")
	}

	return user, nil
}

func (r userRepository) UpdateUser(id uint, usr domain.User) (domain.User, error) {
	var user domain.User

	err := r.db.Model(&user).Clauses(clause.Returning{}).Where("id = ?", id).Updates(&usr).Error

	if err != nil {
		log.Printf("database error while updating user: %v\n", err)
		return domain.User{}, errors.New("cannot update user")
	}

	return user, nil
}

func (r userRepository) GetVerificationCode(email string) (int, error) {
	var user domain.User

	err := r.db.First(&user, "email = ?", email).Error

	if err != nil {
		log.Printf("database error while finding user: %v\n", err)
		return 0, errors.New("user does not exist")
	}

	if user.Verified {
		log.Printf("user is verified\n")
		return user.Code, nil
	}

	return 1245, nil
}
