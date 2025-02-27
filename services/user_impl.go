package services

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang-api-auth-template/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	db          *gorm.DB
	logger      *zap.Logger
	defaultRole string
}

func (u UserServiceImpl) Create(user *models.User) error {
	err := u.checkUnique(user)
	if err != nil {
		u.logger.Error("User duplicate", zap.Error(err))
		return fmt.Errorf("%w: %w", ErrDuplicate, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		u.logger.Error("Error hashing password", zap.Error(err))
		return fmt.Errorf("%w: %w", ErrHashingPassword, err)
	}

	user.Password = string(hashedPassword)

	tx := u.db.Begin()

	result := tx.Create(user)
	if result.Error != nil {
		tx.Rollback()
		u.logger.Error("Error creating user", zap.Error(result.Error))
		return fmt.Errorf("%w: %w", ErrCreatingUser, result.Error)
	}

	if err := tx.First(&user, user.ID).Error; err != nil {
		tx.Rollback()
		u.logger.Error("Error reloading user", zap.Error(err))
		return fmt.Errorf("failed to reload user: %w", err)
	}

	role := &models.Role{}
	result = u.db.Where("name = ?", u.defaultRole).First(&role)

	if result.Error != nil {
		tx.Rollback()
		u.logger.Error("Error retrieving role", zap.Error(result.Error))
		return fmt.Errorf("%w: %w", ErrRetrievingRole, result.Error)
	}

	err = tx.Model(user).Association("Roles").Append(role)
	if err != nil {
		tx.Rollback()
		u.logger.Error("Error assigning role to user", zap.Error(err))
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return tx.Commit().Error
}

func (u UserServiceImpl) Update(id uint, data map[string]interface{}) error {
	existingUser := &models.User{}
	result := u.db.First(existingUser, id).Updates(data)
	if result.Error != nil {
		u.logger.Error("Error update user", zap.Error(result.Error))
		return fmt.Errorf("%w: %w", ErrUpdateUser, result.Error)
	}

	return nil
}

func (u UserServiceImpl) GetById(id uint) (*models.User, error) {
	return u.getBy("id", id)
}

func (u UserServiceImpl) GetByName(s string) (*models.User, error) {
	return u.getBy("name", s)
}

func (u UserServiceImpl) GetByEmail(s string) (*models.User, error) {
	return u.getBy("email", s)
}

func (u UserServiceImpl) GetByNumber(s string) (*models.User, error) {
	return u.getBy("number", s)
}

func (u UserServiceImpl) GetByCredential(s string) (*models.User, error) {
	query := "number = ? OR email = ? OR name = ?"

	user, err := u.getByWhere(query, s, s, s)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRetrievingUser, err)
	}
	return user, nil
}

func (u UserServiceImpl) Verify(input string, credential string) error {
	user, err := u.GetByCredential(credential)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidCredentials, err)
	}

	err = u.VerifyPassword(input, user.Password)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidCredentials, err)
	}

	return nil
}

func (u UserServiceImpl) VerifyPassword(input string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(input))
	if err != nil {
		u.logger.Error("Error comparing password hashes", zap.Error(err))
		return fmt.Errorf("%w: %w", ErrPasswordComparison, err)
	}

	return nil
}

func (u UserServiceImpl) checkUnique(user *models.User) error {
	var existingUser models.User
	result := u.db.Where("name = ?", user.Name).First(&existingUser)
	if result.Error == nil {
		u.logger.Error("Duplicate name", zap.String("name", user.Name))
		return ErrDuplicateName
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // Some other database error occurred
		u.logger.Error("Error checking for duplicate name", zap.Error(result.Error))
		return fmt.Errorf("%w: %w", ErrRetrievingUser, result.Error)
	}

	result = u.db.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		u.logger.Error("Duplicate email", zap.String("email", user.Email))
		return ErrDuplicateEmail
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		u.logger.Error("Error checking for duplicate email", zap.Error(result.Error))
		return fmt.Errorf("error checking email uniqueness: %w", result.Error)
	}

	result = u.db.Where("number = ?", user.Number).First(&existingUser)
	if result.Error == nil {
		u.logger.Error("Duplicate number", zap.String("number", user.Number))
		return ErrDuplicateNumber
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		u.logger.Error("Error checking for duplicate number", zap.Error(result.Error))
		return fmt.Errorf("error checking number uniqueness: %w", result.Error)
	}

	return nil
}

func (u UserServiceImpl) getBy(key string, value interface{}) (*models.User, error) {
	return u.getByWhere(fmt.Sprintf("%s = ?", key), value)
}

func (u UserServiceImpl) getByWhere(query interface{}, args ...interface{}) (*models.User, error) {
	user := &models.User{}
	result := u.db.Where(query, args...).First(user)
	if result.Error != nil {
		u.logger.Error("Error retrieving user", zap.Error(result.Error))
		return nil, fmt.Errorf("%w: %w", ErrRetrievingUser, result.Error)
	}
	return user, nil
}

func NewUserServiceImpl(db *gorm.DB, logger *zap.Logger, defaultRole string) *UserServiceImpl {
	return &UserServiceImpl{db, logger, defaultRole}
}
