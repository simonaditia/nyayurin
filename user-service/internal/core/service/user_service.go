package service

import (
	"context"
	"errors"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/simonaditia/nyayurin/user-service/config"
	"github.com/simonaditia/nyayurin/user-service/internal/adapter/repository"
	"github.com/simonaditia/nyayurin/user-service/internal/core/domain/entity"
	"github.com/simonaditia/nyayurin/user-service/utils/conv"
)

type UserServiceInterface interface {
	SigIn(ctx context.Context, req entity.UserEntity) (*entity.UserEntity, string, error)
}

type userService struct {
	repo       repository.UserRepositoryInterface
	cfg        *config.Config
	jwtService JwtServiceInterface
}

// SigIn implements UserServiceInterface.
func (u *userService) SigIn(ctx context.Context, req entity.UserEntity) (*entity.UserEntity, string, error) {
	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Errorf("[UserService-1] SigIn failed for email %s: %v", req.Email, err)
		return nil, "", err
	}

	if checkPass := conv.CheckPasswordHash(req.Password, user.Password); !checkPass {
		err = errors.New("password is incorrect")
		log.Errorf("[UserService-2] Invalid password for email %s", req.Email)
		return nil, "", err
	}

	token, err := u.jwtService.GenerateToken(user.ID)
	if err != nil {
		log.Errorf("[UserService-3] Failed to generate token for email %s: %v", req.Email, err)
		return nil, "", err
	}

	sessionData := map[string]interface{}{
		"user_id":    user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"logged_in":  true,
		"created_at": time.Now().Unix(),
		"token":      token,
	}

	redisConn := config.NewRedisClient()
	err = redisConn.HSet(ctx, token, sessionData).Err()
	if err != nil {
		log.Errorf("[UserService-4] Failed to store session in Redis for email %s: %v", req.Email, err)
		return nil, "", err
	}

	return user, token, nil
}

func NewUserService(repo repository.UserRepositoryInterface, cfg *config.Config, jwtService JwtServiceInterface) UserServiceInterface {
	return &userService{repo: repo, cfg: cfg, jwtService: jwtService}
}
