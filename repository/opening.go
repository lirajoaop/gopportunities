package repository

import (
	"github.com/lirajoaop/gopportunities/schemas"
	"gorm.io/gorm"
)

type OpeningRepository struct {
	db *gorm.DB
}

func NewOpeningRepository(db *gorm.DB) *OpeningRepository {
	return &OpeningRepository{db: db}
}

func (r *OpeningRepository) Create(opening *schemas.Opening) error {
	return r.db.Create(opening).Error
}

func (r *OpeningRepository) FindByID(id string) (*schemas.Opening, error) {
	opening := schemas.Opening{}
	if err := r.db.First(&opening, id).Error; err != nil {
		return nil, err
	}
	return &opening, nil
}

func (r *OpeningRepository) FindAll() ([]schemas.Opening, error) {
	var openings []schemas.Opening
	if err := r.db.Find(&openings).Error; err != nil {
		return nil, err
	}
	return openings, nil
}

func (r *OpeningRepository) Save(opening *schemas.Opening) error {
	return r.db.Save(opening).Error
}

func (r *OpeningRepository) Delete(opening *schemas.Opening) error {
	return r.db.Delete(opening).Error
}
