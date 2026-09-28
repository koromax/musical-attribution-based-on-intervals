package repository

import (
	"fmt"

	"interval_attribution/internal/app/ds"
)

func (r *Repository) GetPublishedComposers(minFreq, maxFreq *float64, currentUserID uint) ([]ds.Composer, error) {
	var composers []ds.Composer
	query := r.db.Where("status = ?", "published")

	if minFreq != nil {
		query = query.Where("freq1 >= ?", *minFreq)
	}
	if maxFreq != nil {
		query = query.Where("freq1 <= ?", *maxFreq)
	}

	err := query.Order("id asc").Find(&composers).Error
	if err != nil {
		return nil, err
	}

	for i := range composers {
		composers[i].IsCreator = (composers[i].CreatorID == currentUserID)
	}

	return composers, nil
}

func (r *Repository) GetComposerFeedItem(currentID int, next bool, currentUserID uint) (ds.Composer, error) {
	var composer ds.Composer

	if currentID == 0 {
		err := r.db.Where("status = ?", "published").Order("id asc").First(&composer).Error
		if err == nil {
			composer.IsCreator = (composer.CreatorID == currentUserID)
		}
		return composer, err
	}

	if next {
		err := r.db.Where("status = ? AND id > ?", "published", currentID).Order("id asc").First(&composer).Error
		if err != nil {
			err = r.db.Where("status = ?", "published").Order("id asc").First(&composer).Error
		}
		if err == nil {
			composer.IsCreator = (composer.CreatorID == currentUserID)
		}
		return composer, err
	}

	err := r.db.Where("status = ? AND id = ?", "published", currentID).First(&composer).Error
	if err == nil {
		composer.IsCreator = (composer.CreatorID == currentUserID)
	}
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

func (r *Repository) CreateDraftComposer(userID uint, name string) (ds.Composer, error) {
	var existing ds.Composer
	err := r.db.Where("creator_id = ? AND status = ?", userID, "draft").First(&existing).Error
	if err == nil {
		return existing, nil
	}

	newComposer := ds.Composer{
		Name:      name,
		Status:    "draft",
		CreatorID: userID,
	}

	if err := r.db.Create(&newComposer).Error; err != nil {
		return ds.Composer{}, err
	}
	return newComposer, nil
}

func (r *Repository) UpdateComposerMedia(composerID uint, imageURL, videoURL string) error {
	updates := map[string]interface{}{}
	if imageURL != "" {
		updates["image_url"] = imageURL
	}
	if videoURL != "" {
		updates["video_url"] = videoURL
	}
	return r.db.Model(&ds.Composer{}).Where("id = ?", composerID).Updates(updates).Error
}

func (r *Repository) PublishComposer(id uint, userID uint, description string, freq1, freq2 float64) error {
	res := r.db.Model(&ds.Composer{}).
		Where("id = ? AND creator_id = ? AND status = ?", id, userID, "draft").
		Updates(map[string]interface{}{
			"description": description,
			"freq1":       freq1,
			"freq2":       freq2,
			"status":      "published",
		})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("черновик не найден или уже опубликован")
	}
	return nil
}

func (r *Repository) DeleteComposerSoft(id uint, userID uint) error {
	res := r.db.Model(&ds.Composer{}).
		Where("id = ? AND creator_id = ? AND status != ?", id, userID, "deleted").
		Update("status", "deleted")

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("запись не найдена, уже удалена или принадлежит другому пользователю")
	}
	return nil
}
