package repository

import (
	"context"
	"errors"

	"github.com/labstack/gommon/log"
	"github.com/simonaditia/nyayurin/user-service/internal/core/domain/entity"
	"github.com/simonaditia/nyayurin/user-service/internal/core/domain/model"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
}

type userRepository struct {
	db *gorm.DB
}

// GetUserByEmail implements UserRepository.
func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	modelUser := model.User{}

	if err := u.db.Where("email = ? AND is_verified = ?", email, true).Preload("Roles").First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Warnf("[UserRepository-1] User with email %s not found", email)
			return nil, err
		}

		log.Errorf("[UserRepository-1] GetUserByEmail %s: %v", email, err)
		return nil, err
	}

	return &entity.UserEntity{
		ID:         modelUser.ID,
		Name:       modelUser.Name,
		Email:      email,
		Password:   modelUser.Password,
		RoleName:   modelUser.Roles[0].Name,
		Address:    modelUser.Address,
		Lat:        modelUser.Lat,
		Lng:        modelUser.Lng,
		Phone:      modelUser.Phone,
		Photo:      modelUser.Photo,
		IsVerified: modelUser.IsVerified,
	}, nil

	// return &entityUser, nil
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{db: db}
}
