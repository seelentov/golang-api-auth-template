package services

import (
	"fmt"
	"go.uber.org/zap"
	"golang-api-auth-template/models"
	"gorm.io/gorm"
)

type RoleServiceImpl struct {
	db          *gorm.DB
	logger      *zap.Logger
	defaultRole string
}

func (r RoleServiceImpl) GetDefaultRole() (*models.Role, error) {
	role, err := r.GetRoleByName(r.defaultRole)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r RoleServiceImpl) GetRoleByName(name string) (*models.Role, error) {
	roles, err := r.getBy("name", name)
	if err != nil {
		r.logger.Error("Error retrieving role by name", zap.Error(err))
		return nil, fmt.Errorf("%w:%w", ErrRetrievingRole, err)
	}

	if len(roles) == 0 {
		r.logger.Error("Role by name not found", zap.Error(gorm.ErrRecordNotFound))
		return nil, fmt.Errorf("%w:%w", gorm.ErrRecordNotFound, ErrRetrievingRole)
	}

	return roles[0], nil
}

func (r RoleServiceImpl) GetUsersByName(name string) ([]*models.User, error) {
	role, err := r.GetRoleByName(name)
	if err != nil {
		r.logger.Error("Error retrieving role by name", zap.Error(err))
		return nil, fmt.Errorf("%w:%w", ErrRetrievingRole, err)
	}

	var users []*models.User

	result := r.db.Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", role.Name).
		Find(&users)

	if result.Error != nil {
		r.logger.Error("Error retrieving roles by user ID", zap.Error(result.Error))
		return nil, fmt.Errorf("%w: %w", ErrRetrievingRole, result.Error)
	}

	return users, nil
}

func (r RoleServiceImpl) GetRolesByUserId(userId uint) ([]*models.Role, error) {
	var roles []*models.Role

	result := r.db.Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userId).
		Find(&roles)

	if result.Error != nil {
		r.logger.Error("Error retrieving roles by user ID", zap.Error(result.Error))
		return nil, fmt.Errorf("%w: %w", ErrRetrievingRole, result.Error)
	}

	return roles, nil
}

func (r RoleServiceImpl) getBy(key string, value interface{}) ([]*models.Role, error) {
	return r.getByWhere(fmt.Sprintf("%s = ?", key), value)
}

func (r RoleServiceImpl) getByWhere(query interface{}, args ...interface{}) ([]*models.Role, error) {
	var roles []*models.Role

	result := r.db.Where(query, args...).Find(&roles)

	if result.Error != nil {
		r.logger.Error("Error retrieving roles", zap.Error(result.Error))
		return nil, fmt.Errorf("%w: %w", ErrRetrievingRole, result.Error)
	}

	return roles, nil
}

func NewRoleServiceImpl(db *gorm.DB, logger *zap.Logger, defaultRole string) *RoleServiceImpl {
	return &RoleServiceImpl{
		db:          db,
		logger:      logger,
		defaultRole: defaultRole,
	}
}
