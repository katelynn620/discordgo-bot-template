package repo

import (
	"discordbot/pkg/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepo struct {
	BaseRepo
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) FindOrCreate(user *model.User) *model.User {
	u.db.FirstOrCreate(user, model.User{ID: user.ID})
	return user
}

func (u *UserRepo) UpdateById(id string, updateUser *model.User) error {
	logger := zap.L().Sugar()
	defer logger.Sync()
	logger.Debugf("update user by id: %v", id)

	return u.db.Model(&model.User{}).Where("id = ?", id).Updates(updateUser).Error
}
