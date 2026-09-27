package ds

type ComposerLike struct {
	ID         uint `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint `gorm:"not null;uniqueIndex:idx_user_composer" json:"user_id"`
	ComposerID uint `gorm:"not null;uniqueIndex:idx_user_composer" json:"composer_id"`

	// Каскадное удаление отключено
	User     User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user"`
	Composer Composer `gorm:"foreignKey:ComposerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"composer"`
}
