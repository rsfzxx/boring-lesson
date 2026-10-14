package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"

	"go-starter/internal/apperr"
	"go-starter/internal/auth"
	"go-starter/internal/dto"
	"go-starter/internal/model"
	"go-starter/internal/repository"
)

var (
	phonePattern = regexp.MustCompile(`^\+?[0-9]{8,15}$`)
	emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

const minPasswordLength = 8

type UserService struct {
	users  *repository.UserRepository
	roles  *repository.RoleRepository
	groups *repository.GroupRepository
}

func NewUserService(users *repository.UserRepository, roles *repository.RoleRepository, groups *repository.GroupRepository) *UserService {
	return &UserService{users: users, roles: roles, groups: groups}
}

func canSeeUser(actor auth.User, target *model.User) bool {
	return actor.RoleName == auth.RoleAdmin || target.GroupID.Hex() == actor.GroupID
}

func (s *UserService) Create(ctx context.Context, actor auth.User, req dto.CreateUserRequest) (*dto.UserData, error) {
	req.Name = strings.TrimSpace(req.Name)
	switch {
	case req.Name == "":
		return nil, validation("name is required")
	case !phonePattern.MatchString(req.PhoneNumber):
		return nil, validation("phone_number must be 8-15 digits, optionally starting with +")
	case len(req.Password) < minPasswordLength:
		return nil, validation(fmt.Sprintf("password must be at least %d characters", minPasswordLength))
	case req.Email != nil && !emailPattern.MatchString(*req.Email):
		return nil, validation("email is not valid")
	}

	role, err := s.findRole(ctx, req.RoleID)
	if err != nil {
		return nil, err
	}
	groupID, err := s.findGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	switch actor.RoleName {
	case auth.RoleAdmin:
	case auth.RoleLeader:
		if role.Name != auth.RoleMember {
			return nil, fmt.Errorf("%w: leaders can only create members", apperr.ErrForbidden)
		}
		if groupID.Hex() != actor.GroupID {
			return nil, fmt.Errorf("%w: leaders can only create users in their own group", apperr.ErrForbidden)
		}
	default:
		return nil, apperr.ErrForbidden
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &model.User{
		ID:          bson.NewObjectID(),
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    string(hash),
		RoleID:      role.ID,
		GroupID:     groupID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, fmt.Errorf("%w: phone_number or email is already registered", apperr.ErrConflict)
		}
		return nil, err
	}
	return toUserData(user), nil
}

func (s *UserService) List(ctx context.Context, actor auth.User) ([]dto.UserData, error) {
	var groupFilter *bson.ObjectID
	if actor.RoleName != auth.RoleAdmin {
		groupID, err := bson.ObjectIDFromHex(actor.GroupID)
		if err != nil {
			return []dto.UserData{}, nil
		}
		groupFilter = &groupID
	}

	users, err := s.users.FindAll(ctx, groupFilter)
	if err != nil {
		return nil, err
	}
	res := make([]dto.UserData, 0, len(users))
	for i := range users {
		res = append(res, *toUserData(&users[i]))
	}
	return res, nil
}

func (s *UserService) GetByID(ctx context.Context, actor auth.User, id string) (*dto.UserData, error) {
	user, err := s.findUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canSeeUser(actor, user) {
		return nil, apperr.ErrForbidden
	}
	return toUserData(user), nil
}

func (s *UserService) CheckAccess(ctx context.Context, actor auth.User, id string) error {
	_, err := s.GetByID(ctx, actor, id)
	return err
}

func (s *UserService) Update(ctx context.Context, actor auth.User, id string, req dto.UpdateUserRequest) (*dto.UserData, error) {
	if actor.RoleName != auth.RoleAdmin {
		return nil, apperr.ErrForbidden
	}
	userID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperr.ErrInvalidID
	}

	set := bson.M{}
	unset := bson.M{}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, validation("name cannot be empty")
		}
		set["name"] = name
	}
	if req.Email != nil {
		switch {
		case *req.Email == "":
			unset["email"] = ""
		case !emailPattern.MatchString(*req.Email):
			return nil, validation("email is not valid")
		default:
			set["email"] = *req.Email
		}
	}
	if req.PhoneNumber != nil {
		if !phonePattern.MatchString(*req.PhoneNumber) {
			return nil, validation("phone_number must be 8-15 digits, optionally starting with +")
		}
		set["phone_number"] = *req.PhoneNumber
	}
	if req.Password != nil {
		if len(*req.Password) < minPasswordLength {
			return nil, validation(fmt.Sprintf("password must be at least %d characters", minPasswordLength))
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		set["password"] = string(hash)
	}
	if req.RoleID != nil {
		role, err := s.findRole(ctx, *req.RoleID)
		if err != nil {
			return nil, err
		}
		set["role_id"] = role.ID
	}
	if req.GroupID != nil {
		groupID, err := s.findGroup(ctx, *req.GroupID)
		if err != nil {
			return nil, err
		}
		set["group_id"] = groupID
	}

	if len(set) == 0 && len(unset) == 0 {
		return nil, validation("no fields to update")
	}
	set["updated_at"] = time.Now().UTC()

	user, err := s.users.Update(ctx, userID, set, unset)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return nil, apperr.ErrNotFound
	case errors.Is(err, repository.ErrDuplicate):
		return nil, fmt.Errorf("%w: phone_number or email is already registered", apperr.ErrConflict)
	case err != nil:
		return nil, err
	}
	return toUserData(user), nil
}

func (s *UserService) Delete(ctx context.Context, actor auth.User, id string) error {
	if actor.RoleName != auth.RoleAdmin {
		return apperr.ErrForbidden
	}
	if id == actor.ID {
		return fmt.Errorf("%w: you cannot delete your own account", apperr.ErrForbidden)
	}
	userID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return apperr.ErrInvalidID
	}
	err = s.users.Delete(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return apperr.ErrNotFound
	}
	return err
}

func validation(msg string) error {
	return fmt.Errorf("%w: %s", apperr.ErrValidation, msg)
}

func (s *UserService) findUser(ctx context.Context, id string) (*model.User, error) {
	userID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperr.ErrInvalidID
	}
	user, err := s.users.FindByID(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, apperr.ErrNotFound
	}
	return user, err
}

func (s *UserService) findRole(ctx context.Context, id string) (*model.Role, error) {
	roleID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, validation("role_id is not valid")
	}
	role, err := s.roles.FindByID(ctx, roleID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, validation("role_id does not exist")
	}
	return role, err
}

func (s *UserService) findGroup(ctx context.Context, id string) (bson.ObjectID, error) {
	groupID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.ObjectID{}, validation("group_id is not valid")
	}
	exists, err := s.groups.ExistsByID(ctx, groupID)
	if err != nil {
		return bson.ObjectID{}, err
	}
	if !exists {
		return bson.ObjectID{}, validation("group_id does not exist")
	}
	return groupID, nil
}

func toUserData(u *model.User) *dto.UserData {
	return &dto.UserData{
		ID:          u.ID.Hex(),
		Name:        u.Name,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		RoleID:      u.RoleID.Hex(),
		GroupID:     u.GroupID.Hex(),
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}