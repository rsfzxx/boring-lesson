package service

import (
	"context"
	"errors"
	"slices"

	"go.mongodb.org/mongo-driver/v2/bson"

	"go-starter/internal/apperr"
	"go-starter/internal/auth"
	"go-starter/internal/dto"
	"go-starter/internal/repository"
)

type RoleService struct {
	roles *repository.RoleRepository
}

func NewRoleService(roles *repository.RoleRepository) *RoleService {
	return &RoleService{roles: roles}
}

func visibleRoleNames(roleName string) []string {
	switch roleName {
	case auth.RoleAdmin:
		return []string{auth.RoleAdmin, auth.RoleLeader, auth.RoleMember}
	case auth.RoleLeader:
		return []string{auth.RoleLeader, auth.RoleMember}
	case auth.RoleMember:
		return []string{auth.RoleMember}
	default:
		return nil
	}
}

func (s *RoleService) List(ctx context.Context, user auth.User) ([]dto.RoleResponse, error) {
	roles, err := s.roles.FindByNames(ctx, visibleRoleNames(user.RoleName))
	if err != nil {
		return nil, err
	}
	res := make([]dto.RoleResponse, 0, len(roles))
	for _, role := range roles {
		res = append(res, dto.RoleResponse{ID: role.ID.Hex(), Name: role.Name})
	}
	return res, nil
}

func (s *RoleService) GetByID(ctx context.Context, user auth.User, id string) (*dto.RoleResponse, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperr.ErrInvalidID
	}
	role, err := s.roles.FindByID(ctx, objectID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, apperr.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if !slices.Contains(visibleRoleNames(user.RoleName), role.Name) {
		return nil, apperr.ErrForbidden
	}
	return &dto.RoleResponse{ID: role.ID.Hex(), Name: role.Name}, nil
}

func (s *RoleService) CheckAccess(ctx context.Context, user auth.User, id string) error {
	_, err := s.GetByID(ctx, user, id)
	return err
}