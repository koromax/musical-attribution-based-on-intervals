package repository

import "fmt"

type Composer struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	VideoURL    string  `json:"video_url"`
	Status      string  `json:"status"`
	Likes       []int   `json:"likes"`
	LikesCount  int     `json:"likes_count"`
	Freq1       float64 `json:"freq1"`
	Freq2       float64 `json:"freq2"`
}

type Repository struct {
	composers []Composer
}

func NewRepository() *Repository {
	return &Repository{
		composers: []Composer{
			{
				ID:          1,
				Name:        "И. С. Бах",
				Description: "Композитор эпохи барокко.",
				ImageURL:    "http://localhost:9000/music/cover1.jpg",
				VideoURL:    "http://localhost:9000/music/1.mp4",
				Status:      "published",
				Likes:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				LikesCount:  142,
				Freq1:       38.5,
				Freq2:       12.4,
			},
			{
				ID:          2,
				Name:        "Ф. Шопен",
				Description: "Композитор эпохи романтизма.",
				ImageURL:    "http://localhost:9000/music/cover2.jpg",
				VideoURL:    "http://localhost:9000/music/2.mp4",
				Status:      "published",
				Likes:       []int{1, 2, 3},
				LikesCount:  89,
				Freq1:       12.0,
				Freq2:       7.8,
			},
			{
				ID:          3,
				Name:        "Новый композитор",
				Description: "Черновик.",
				ImageURL:    "",
				VideoURL:    "http://localhost:9000/music/3.mp4",
				Status:      "draft",
				Likes:       []int{},
				LikesCount:  0,
				Freq1:       25.4,
				Freq2:       10.2,
			},
		},
	}
}

func (r *Repository) GetDraft() (Composer, error) {
	for _, item := range r.composers {
		if item.Status == "draft" {
			return item, nil
		}
	}
	return Composer{}, fmt.Errorf("черновик не найден")
}

func (r *Repository) GetFeedItem(id int, next bool) (Composer, error) {
	var published []Composer
	for _, item := range r.composers {
		if item.Status == "published" {
			published = append(published, item)
		}
	}
	if len(published) == 0 {
		return Composer{}, fmt.Errorf("нет публикаций")
	}
	if id == 0 {
		return published[0], nil
	}
	for i, item := range published {
		if item.ID == id {
			if next {
				nextIdx := (i + 1) % len(published)
				return published[nextIdx], nil
			}
			return item, nil
		}
	}
	return published[0], nil
}

func (r *Repository) GetPublishedGrid(minFreq *float64, maxFreq *float64) []Composer {
	var result []Composer
	for _, item := range r.composers {
		if item.Status != "published" {
			continue
		}
		if minFreq != nil && item.Freq1 < *minFreq {
			continue
		}
		if maxFreq != nil && item.Freq1 > *maxFreq {
			continue
		}
		result = append(result, item)
	}
	return result
}
