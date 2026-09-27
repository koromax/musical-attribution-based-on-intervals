package repository

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"interval_attribution/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetComposerFeedItem(currentID int, next bool) (ds.Composer, error) {
	var composer ds.Composer

	if currentID == 0 {
		err := r.db.Where("status = ?", "published").Order("id asc").Limit(1).First(&composer).Error
		return composer, err
	}

	if next {
		err := r.db.Where("status = ? AND id > ?", "published", currentID).Order("id asc").Limit(1).First(&composer).Error
		if err != nil {
			err = r.db.Where("status = ?", "published").Order("id asc").Limit(1).First(&composer).Error
		}
		return composer, err
	}

	err := r.db.Where("status = ? AND id = ?", "published", currentID).First(&composer).Error
	return composer, err
}

func (r *Repository) GetComposerDraft(userID uint) (ds.Composer, error) {
	var composer ds.Composer
	err := r.db.Where("creator_id = ? AND status = ?", userID, "draft").First(&composer).Error
	if err != nil {
		return ds.Composer{}, fmt.Errorf("черновик не найден")
	}
	return composer, nil
}

func (r *Repository) GetPublishedComposers(minFreq *float64, maxFreq *float64) []ds.Composer {
	var composers []ds.Composer
	query := r.db.Where("status = ?", "published")

	if minFreq != nil {
		query = query.Where("freq1 >= ?", *minFreq)
	}
	if maxFreq != nil {
		query = query.Where("freq1 <= ?", *maxFreq)
	}

	query.Order("id asc").Find(&composers)
	return composers
}

func (r *Repository) CreateDraftComposer(userID uint, name string) (ds.Composer, error) {
	var existing ds.Composer
	err := r.db.Where("creator_id = ? AND status = ?", userID, "draft").First(&existing).Error
	if err == nil {
		return existing, nil
	}

	newComposer := ds.Composer{
		Name:       name,
		Status:     "draft",
		CreatorID:  userID,
		DateCreate: time.Now(),
		DateFormed: time.Now(),
	}

	if err := r.db.Create(&newComposer).Error; err != nil {
		return ds.Composer{}, err
	}
	return newComposer, nil
}

func (r *Repository) PublishComposer(id uint, description string, freq1, freq2 float64) error {
	return r.db.Model(&ds.Composer{}).
		Where("id = ? AND status = ?", id, "draft").
		Updates(map[string]interface{}{
			"description": description,
			"freq1":       freq1,
			"freq2":       freq2,
			"status":      "published",
			"date_formed": time.Now(),
		}).Error
}

func (r *Repository) DeleteComposerSQL(id int) error {
	query := "UPDATE composers SET status = 'deleted', date_formed = NOW() WHERE id = $1 AND status != 'deleted'"
	res := r.db.Exec(query, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("карточка не найдена или уже удалена")
	}
	return nil
}

func (r *Repository) GetLikesCount(composerID uint) int64 {
	var count int64
	r.db.Model(&ds.ComposerLike{}).Where("composer_id = ?", composerID).Count(&count)
	return count
}
