package repository

import (
	"interval_attribution/internal/app/ds"
)

func (r *Repository) ToggleLike(userID uint, composerID uint, enable int) error {
	if enable == 1 {
		like := ds.ComposerLike{
			UserID:     userID,
			ComposerID: composerID,
		}
		return r.db.Where(ds.ComposerLike{UserID: userID, ComposerID: composerID}).FirstOrCreate(&like).Error
	}

	return r.db.Where("user_id = ? AND composer_id = ?", userID, composerID).Delete(&ds.ComposerLike{}).Error
}
