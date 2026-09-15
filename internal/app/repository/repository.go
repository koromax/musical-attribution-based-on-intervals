package repository

import "fmt"

type Composition struct {
	ID                int     `json:"id"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	ImageURL          string  `json:"image_url"`
	VideoURL          string  `json:"video_url"`
	Status            string  `json:"status"`
	LikesCount        int     `json:"likes_count"`
	IntervalFrequency float64 `json:"interval_frequency"`
	AttributedAuthor  string  `json:"attributed_author"`
}

type Repository struct {
	compositions []Composition
}

func NewRepository() *Repository {
	return &Repository{
		compositions: []Composition{
			{
				ID:                1,
				Title:             "Токката и фуга ре минор",
				Description:       "Высокая частотность квинты характерна для полифонии барокко.",
				ImageURL:          "http://localhost:9000/music/cover1.jpg",
				VideoURL:          "http://localhost:9000/music/1.mp4",
				Status:            "published",
				LikesCount:        142,
				IntervalFrequency: 38.5,
				AttributedAuthor:  "И. С. Бах",
			},
			{
				ID:                2,
				Title:             "Ноктюрн Op. 9 No. 2",
				Description:       "Частотный анализ интервалов гармонии периода романтизма.",
				ImageURL:          "http://localhost:9000/music/cover2.jpg",
				VideoURL:          "http://localhost:9000/music/2.mp4",
				Status:            "published",
				LikesCount:        89,
				IntervalFrequency: 12.0,
				AttributedAuthor:  "Ф. Шопен",
			},
			{
				ID:                3,
				Title:             "Черновик партитуры",
				Description:       "Неоконченный анализ частотности гармонических интервалов.",
				ImageURL:          "http://localhost:9000/music/draft_cover.jpg",
				VideoURL:          "http://localhost:9000/music/3.mp4",
				Status:            "draft",
				LikesCount:        0,
				IntervalFrequency: 25.4,
				AttributedAuthor:  "Неизвестный автор",
			},
		},
	}
}

func (r *Repository) GetDraft() (Composition, error) {
	for _, item := range r.compositions {
		if item.Status == "draft" {
			return item, nil
		}
	}
	return Composition{}, fmt.Errorf("черновик не найден")
}

func (r *Repository) GetFeedItem(id int, next bool) (Composition, error) {
	var published []Composition
	for _, item := range r.compositions {
		if item.Status == "published" {
			published = append(published, item)
		}
	}
	if len(published) == 0 {
		return Composition{}, fmt.Errorf("нет публикаций")
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

func (r *Repository) GetPublishedGrid(minFreq *float64, maxFreq *float64) []Composition {
	var result []Composition
	for _, item := range r.compositions {
		if item.Status != "published" {
			continue
		}
		if minFreq != nil && item.IntervalFrequency < *minFreq {
			continue
		}
		if maxFreq != nil && item.IntervalFrequency > *maxFreq {
			continue
		}
		result = append(result, item)
	}
	return result
}
