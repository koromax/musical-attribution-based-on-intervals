package ds

import "time"

type Composer struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(64);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(9);not null;default:'draft'" json:"status"`
	ImageURL    string    `gorm:"type:varchar(255)" json:"image_url"`
	VideoURL    string    `gorm:"type:varchar(255)" json:"video_url"`
	Freq1       float64   `gorm:"type:numeric(5,2)" json:"freq1"`
	Freq2       float64   `gorm:"type:numeric(5,2)" json:"freq2"`
	DateCreate  time.Time `gorm:"type:timestamp;not null;autoCreateTime" json:"date_create"`
	CreatorID   uint      `gorm:"not null" json:"creator_id"`
	DateFormed  time.Time `gorm:"type:timestamp;autoUpdateTime" json:"date_formed"`

	IsCreator bool `gorm:"-" json:"is_creator"`

	Creator User `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
}
