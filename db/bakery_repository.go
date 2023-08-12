package db

import (
	"bread-clock/models"
	"context"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

type SortOption string

const (
	SortByName     SortOption = "name"
	SortByDistance            = "distance"
)

const distanceFormula = "6371 * ACOS(COS(RADIANS(%f)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS(%f)) + SIN(RADIANS(%f)) * SIN(RADIANS(latitude)))"

type BakeryRepository interface {
	List(ctx context.Context, sortOption SortOption, size int, offset int, latitude float64, longitude float64, needsDistance bool, userID int) ([]models.BakeryDetail, error)
	ListForBreads(ctx context.Context, q string, sortOption SortOption, size int, offset int, latitude float64, longitude float64, needsDistance bool, userID int) ([]models.BakeryDetail, error)
	Get(ctx context.Context, bakeryID int, latitude float64, longitude float64, needsDistance bool, userID int) (*models.BakeryDetail, error)
	MarkAsFavorite(ctx context.Context, bakeryID int, userID int) error
	UnmarkAsFavorite(ctx context.Context, bakeryID int, userID int) error
	UpdateBreadAvailabilities(ctx context.Context, breadAvailabilities []models.BreadAvailability) error
}

type bakeryRepository struct {
	db *gorm.DB
}

type BakeryDAO struct {
	ID           int
	Name         string
	Address      string
	Latitude     float64
	Longitude    float64
	Distance     float64
	OpeningHours string
	UserID       *int
	PhotoURLs    string
}

type BreadDAO struct {
	BakeryID       int
	ID             int
	Name           string
	Available      bool
	AvailableHours string
	PhotoURL       string
}

func NewBakeryRepository(db *gorm.DB) BakeryRepository {
	return &bakeryRepository{db: db}
}

func (r *bakeryRepository) List(ctx context.Context, sortOption SortOption, size int, offset int, latitude float64, longitude float64, needsDistance bool, userID int) ([]models.BakeryDetail, error) {
	tx := r.db.WithContext(ctx)

	var bakeryDAOs []BakeryDAO
	distance := fmt.Sprintf(distanceFormula, latitude, longitude, latitude)
	subQuery := tx.
		Table("bakeries AS b").
		Select(fmt.Sprintf("id, name, address, opening_hours, latitude, longitude, (%s) AS distance, user_id", distance)).
		Joins(fmt.Sprintf("LEFT JOIN favorite_bakeries AS fb ON b.id = fb.bakery_id AND fb.user_id = %d", userID)).
		Offset(offset).
		Limit(size)
	query := tx.
		Select("b.id, b.name, b.address, b.opening_hours, b.latitude, b.longitude, b.distance, user_id, GROUP_CONCAT(url) as photo_urls").
		Table("(?) AS b", subQuery).
		Joins("LEFT JOIN bakery_photos AS bp ON bp.bakery_id = b.id").
		Order(string(sortOption)).
		Group("b.id")
	err := query.Find(&bakeryDAOs).Error
	if err != nil {
		return nil, err
	}

	var bakeryIDs []int
	for _, bakeryDAO := range bakeryDAOs {
		bakeryIDs = append(bakeryIDs, bakeryDAO.ID)
	}

	var breadDAOs []BreadDAO
	err = tx.Table("bread_availabilities AS ba").
		Select("ba.bakery_id, b.id, b.name, ba.available, ba.available_hours, url AS photo_url").
		Joins("LEFT JOIN bread_photos AS bp ON bp.bakery_id = ba.bakery_id AND bp.bread_id = ba.bread_id").
		Joins("INNER JOIN breads AS b ON ba.bread_id = b.id").
		Where("ba.bakery_id IN (?)", bakeryIDs).
		Find(&breadDAOs).Error
	if err != nil {
		return nil, err
	}

	breadDetailMap := make(map[int][]models.BreadDetail)
	for _, breadDAO := range breadDAOs {
		bread := models.BreadDetail{
			Bread: models.Bread{
				ID:   breadDAO.ID,
				Name: breadDAO.Name,
			},
			Available:      breadDAO.Available,
			AvailableHours: strings.Split(breadDAO.AvailableHours, ","),
			PhotoURL:       breadDAO.PhotoURL,
		}
		breadDetailMap[breadDAO.BakeryID] = append(breadDetailMap[breadDAO.BakeryID], bread)
	}

	var bakeries []models.BakeryDetail
	for _, bakeryDAO := range bakeryDAOs {
		var bakeryDistance *float64
		if needsDistance {
			bakeryDistance = &bakeryDAO.Distance
		}

		bakeries = append(bakeries, models.BakeryDetail{
			Bakery: models.Bakery{
				ID:        bakeryDAO.ID,
				Name:      bakeryDAO.Name,
				Address:   bakeryDAO.Address,
				Latitude:  bakeryDAO.Latitude,
				Longitude: bakeryDAO.Longitude,
			},
			OpeningHours: convertOpeningHours(bakeryDAO.OpeningHours),
			Distance:     bakeryDistance,
			Favorite:     bakeryDAO.UserID != nil,
			BreadDetails: breadDetailMap[bakeryDAO.ID],
			PhotoURLs:    strings.Split(bakeryDAO.PhotoURLs, ","),
		})
	}

	return bakeries, nil
}

func (r *bakeryRepository) ListForBreads(ctx context.Context, q string, sortOption SortOption, size int, offset int, latitude float64, longitude float64, needsDistance bool, userID int) ([]models.BakeryDetail, error) {
	tx := r.db.WithContext(ctx)

	if regex, err := regexp.Compile(`\s+`); err == nil {
		q = fmt.Sprintf("%%%s%%", regex.ReplaceAllString(q, "%"))
	} else {
		return nil, err
	}

	var breadDAOs []BreadDAO
	subQuery := tx.Table("breads").Where("name LIKE (?)", q)
	query := tx.
		Select("b.id, b.name, ba.bakery_id, ba.available, ba.available_hours, bp.url AS photo_url").
		Table("(?) AS b", subQuery).
		Joins("JOIN bread_availabilities AS ba ON ba.bread_id = b.id").
		Joins("LEFT JOIN bread_photos AS bp ON bp.bread_id = b.id AND bp.bakery_id = ba.bakery_id")
	if err := query.Find(&breadDAOs).Error; err != nil {
		return nil, err
	}

	var bakeryIDs []int
	breadDetailMap := make(map[int][]models.BreadDetail)
	for _, breadDAO := range breadDAOs {
		bread := models.BreadDetail{
			Bread: models.Bread{
				ID:   breadDAO.ID,
				Name: breadDAO.Name,
			},
			Available:      breadDAO.Available,
			AvailableHours: strings.Split(breadDAO.AvailableHours, ","),
			PhotoURL:       breadDAO.PhotoURL,
		}
		breadDetailMap[breadDAO.BakeryID] = append(breadDetailMap[breadDAO.BakeryID], bread)

		bakeryIDs = append(bakeryIDs, breadDAO.BakeryID)
	}

	var bakeryDAOs []BakeryDAO
	distance := fmt.Sprintf(distanceFormula, latitude, longitude, latitude)
	subQuery = tx.
		Table("bakeries").
		Select(fmt.Sprintf("*, (%s) AS distance", distance)).
		Where("id IN (?)", bakeryIDs).
		Offset(offset).
		Limit(size)
	query = tx.
		Select("b.id, b.name, b.address, b.latitude, b.longitude, b.distance, b.opening_hours, user_id, GROUP_CONCAT(url) AS photo_urls").
		Table("(?) AS b", subQuery).
		Joins(fmt.Sprintf("LEFT JOIN favorite_bakeries AS fb ON fb.bakery_id = b.id AND fb.user_id = %d", userID)).
		Joins("LEFT JOIN bakery_photos AS bp ON bp.bakery_id = b.id").
		Order(string(sortOption)).
		Group("b.id")
	if err := query.Find(&bakeryDAOs).Error; err != nil {
		return nil, err
	}

	var bakeries []models.BakeryDetail
	for _, bakeryDAO := range bakeryDAOs {
		var bakeryDistance *float64
		if needsDistance {
			bakeryDistance = &bakeryDAO.Distance
		}

		bakeries = append(bakeries, models.BakeryDetail{
			Bakery: models.Bakery{
				ID:        bakeryDAO.ID,
				Name:      bakeryDAO.Name,
				Address:   bakeryDAO.Address,
				Latitude:  bakeryDAO.Latitude,
				Longitude: bakeryDAO.Longitude,
			},
			OpeningHours: convertOpeningHours(bakeryDAO.OpeningHours),
			Distance:     bakeryDistance,
			Favorite:     bakeryDAO.UserID != nil,
			BreadDetails: breadDetailMap[bakeryDAO.ID],
			PhotoURLs:    strings.Split(bakeryDAO.PhotoURLs, ","),
		})
	}

	return bakeries, nil
}

func (r *bakeryRepository) Get(ctx context.Context, bakeryID int, latitude float64, longitude float64, needsDistance bool, userID int) (*models.BakeryDetail, error) {
	tx := r.db.WithContext(ctx)

	var bakeryDAO BakeryDAO
	distance := fmt.Sprintf(distanceFormula, latitude, longitude, latitude)
	err := tx.Table("bakeries AS b").
		Select(fmt.Sprintf("b.id, b.name, b.address, b.latitude, b.longitude, (%s) AS distance, b.opening_hours, user_id, GROUP_CONCAT(url) AS photo_urls", distance)).
		Joins(fmt.Sprintf("LEFT JOIN favorite_bakeries AS fb ON b.id = fb.bakery_id AND fb.user_id = %d", userID)).
		Joins("LEFT JOIN bakery_photos AS bp ON bp.bakery_id = b.id").
		Group("id").
		Where("b.id = (?)", bakeryID).Find(&bakeryDAO).Error
	if err != nil {
		return nil, err
	}

	var breadDAOs []BreadDAO
	err = tx.Table("bread_availabilities AS ba").
		Select("b.id, b.name, ba.available, ba.available_hours, url AS photo_url").
		Joins("LEFT JOIN bread_photos AS bp ON bp.bakery_id = ba.bakery_id AND bp.bread_id = ba.bread_id").
		Joins("INNER JOIN breads AS b ON ba.bread_id = b.id").
		Where("ba.bakery_id = (?)", bakeryID).
		Find(&breadDAOs).Error
	if err != nil {
		return nil, err
	}

	var breadDetails []models.BreadDetail
	for _, breadDAO := range breadDAOs {
		breadDetails = append(breadDetails, models.BreadDetail{
			Bread: models.Bread{
				ID:   breadDAO.ID,
				Name: breadDAO.Name,
			},
			Available:      breadDAO.Available,
			AvailableHours: strings.Split(breadDAO.AvailableHours, ","),
			PhotoURL:       breadDAO.PhotoURL,
		})
	}

	var bakeryDistance *float64
	if needsDistance {
		bakeryDistance = &bakeryDAO.Distance
	}

	bakery := models.BakeryDetail{
		Bakery: models.Bakery{
			ID:        bakeryDAO.ID,
			Name:      bakeryDAO.Name,
			Address:   bakeryDAO.Address,
			Latitude:  bakeryDAO.Latitude,
			Longitude: bakeryDAO.Longitude,
		},
		OpeningHours: convertOpeningHours(bakeryDAO.OpeningHours),
		Distance:     bakeryDistance,
		Favorite:     bakeryDAO.UserID != nil,
		BreadDetails: breadDetails,
		PhotoURLs:    strings.Split(bakeryDAO.PhotoURLs, ","),
	}

	return &bakery, nil
}

func (r *bakeryRepository) MarkAsFavorite(ctx context.Context, bakeryID int, userID int) error {
	tx := r.db.WithContext(ctx)

	favoriteBakery := models.FavoriteBakery{
		UserID:   userID,
		BakeryID: bakeryID,
	}

	if err := tx.Create(&favoriteBakery).Error; err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) && mySQLError.Number == 1062 {
			return nil
		}
		return err
	}

	return nil
}

func (r *bakeryRepository) UnmarkAsFavorite(ctx context.Context, bakeryID int, userID int) error {
	tx := r.db.WithContext(ctx)

	favoriteBakery := models.FavoriteBakery{
		UserID:   userID,
		BakeryID: bakeryID,
	}

	if err := tx.Delete(&favoriteBakery).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}

func (r *bakeryRepository) UpdateBreadAvailabilities(ctx context.Context, breadAvailabilities []models.BreadAvailability) error {
	tx := r.db.WithContext(ctx)

	if err := tx.Updates(breadAvailabilities).Error; err != nil {
		return err
	}

	return nil
}

func convertOpeningHours(val string) []models.OpeningHours {
	var openingHours []models.OpeningHours
	openingHoursForDays := strings.Split(val, ",")
	for _, openingHoursForDay := range openingHoursForDays {
		hours := strings.Split(openingHoursForDay, "-")
		if len(hours) != 2 {
			return nil
		}
		openingHours = append(openingHours, models.OpeningHours{
			Open:  hours[0],
			Close: hours[1],
		})
	}
	return openingHours
}
