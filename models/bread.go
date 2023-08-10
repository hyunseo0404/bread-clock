package models

type Bread struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:20;not null;uniqueIndex" json:"name"`
}

type BreadDetail struct {
	Bread
	Available      bool     `json:"available"`
	AvailableHours []string `json:"availableHours"`
	PhotoURL       string   `json:"photoURL"`
}
