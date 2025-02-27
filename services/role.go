package services

import (
	"errors"
	"golang-api-auth-template/models"
)

var (
	ErrRetrievingRole   = errors.New("error retrieving role")
	ErrFailedAssignRole = errors.New("error assigning role to user")
)

type RoleService interface {
	GetRoleByName(string) (*models.Role, error)
	GetUsersByName(string) ([]*models.User, error)
	GetRolesByUserId(uint) ([]*models.Role, error)
	GetDefaultRole() (*models.Role, error)
}
