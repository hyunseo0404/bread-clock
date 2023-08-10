package models

type Bakery struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"size:128;not null" json:"name"`
	Address      string `gorm:"size:255" json:"address"`
	Coordinates  string `gorm:"size:128" json:"coordinates"`
	OpeningHours string `gorm:"size:70" json:"openingHours"`
}

type BakeryDetail struct {
	Bakery
	Favorite     bool          `json:"favorite" json:"favorite,omitempty"`
	BreadDetails []BreadDetail `json:"breads"`
	PhotoURLs    []string      `json:"photoUrls"`
}
