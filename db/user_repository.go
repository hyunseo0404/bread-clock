package db

import (
	e "bread-clock/error"
	"bread-clock/models"
	"context"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, uid string, provider string, emailAddress string, avatarURL string) (*models.User, error)
	FindOrCreate(ctx context.Context, uid string, provider string, emailAddress string, avatarURL string) (*models.User, error)
	Find(ctx context.Context, emailAddress string, provider string) (*models.User, error)
	FindByUserID(ctx context.Context, id int) (*models.User, error)
	Delete(ctx context.Context, uid string, provider string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindOrCreate(ctx context.Context, uid string, provider string, emailAddress string, avatarURL string) (*models.User, error) {
	tx := r.db.WithContext(ctx)

	user := models.User{
		UID:       uid,
		Provider:  provider,
		Email:     emailAddress,
		AvatarURL: avatarURL,
	}

	if err := tx.Where("uid = ? AND provider = ?", uid, provider).FirstOrCreate(&user, &user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, uid string, provider string, emailAddress string, avatarURL string) (*models.User, error) {
	tx := r.db.WithContext(ctx)

	user := models.User{
		UID:       uid,
		Provider:  provider,
		Email:     emailAddress,
		AvatarURL: avatarURL,
	}

	if err := tx.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Find(ctx context.Context, emailAddress string, provider string) (*models.User, error) {
	tx := r.db.WithContext(ctx)

	var user models.User
	if err := tx.Model(&user).Where("email = ? AND provider = ?", emailAddress, provider).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUserID(ctx context.Context, id int) (*models.User, error) {
	tx := r.db.WithContext(ctx)

	user := models.User{ID: id}
	if tx = tx.Find(&user); tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, e.ErrDBNotFound
	}

	return &user, nil
}

func (r *userRepository) Delete(ctx context.Context, uid string, provider string) error {
	tx := r.db.WithContext(ctx)

	user := models.User{
		UID:      uid,
		Provider: provider,
	}

	if err := tx.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}
