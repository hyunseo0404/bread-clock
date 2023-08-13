package models

type Bakery struct {
	ID           int     `gorm:"primaryKey" json:"id"`
	PlaceID      int     `gorm:"uniqueIndex" json:"-"`
	Name         string  `gorm:"size:128;not null" json:"name"`
	Address      string  `gorm:"size:255" json:"address"`
	OpeningHours string  `gorm:"size:70" json:"-"`
	Latitude     float64 `gorm:"type:decimal(8,6)" json:"latitude"`
	Longitude    float64 `gorm:"type:decimal(9,6)" json:"longitude"`
}

type BakeryDetail struct {
	Bakery
	OpeningHours []OpeningHours `json:"openingHours"`
	Distance     *float64       `json:"distance,omitempty"`
	Favorite     bool           `json:"favorite"`
	BreadDetails []BreadDetail  `json:"breads"`
	PhotoURLs    []string       `json:"photoUrls"`
}

type OpeningHours struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}
