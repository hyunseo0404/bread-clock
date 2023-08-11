package models

type Bakery struct {
	ID           int     `gorm:"primaryKey" json:"id"`
	Name         string  `gorm:"size:128;not null" json:"name"`
	Address      string  `gorm:"size:255" json:"address"`
	OpeningHours string  `gorm:"size:70" json:"openingHours"`
	Latitude     float64 `gorm:"type:decimal(8,6)" json:"latitude"`
	Longitude    float64 `gorm:"type:decimal(9,6)" json:"longitude"`
}

type BakeryDetail struct {
	Bakery
	Distance     *float64      `json:"distance,omitempty"`
	Favorite     bool          `json:"favorite,omitempty"`
	BreadDetails []BreadDetail `json:"breads"`
	PhotoURLs    []string      `json:"photoUrls"`
}
