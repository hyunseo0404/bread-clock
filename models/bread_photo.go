package models

type BreadPhoto struct {
	ID       int    `gorm:"primaryKey"`
	BreadID  int    `gorm:"index:ux_bread_photos_bread_id_bakery_id_url,unique"`
	BakeryID int    `gorm:"index:ux_bread_photos_bread_id_bakery_id_url,unique"`
	URL      string `gorm:"size:255;index:ux_bread_photos_bread_id_bakery_id_url,unique"`
	Bread    Bread  `gorm:"constraint:OnDelete:CASCADE"`
	Bakery   Bakery `gorm:"constraint:OnDelete:CASCADE"`
}
