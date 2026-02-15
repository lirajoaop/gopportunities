package schemas

import (
	"time"

	"gorm.io/gorm"
)

type Opening struct {
	gorm.Model
	Role     string
	Company  string
	Location string
	Remote   bool
	Link     string
	Salary   int64
}

func (o *Opening) ToResponse() OpeningResponse {
	var deletedAt *time.Time
	if o.DeletedAt.Valid {
		deletedAt = &o.DeletedAt.Time
	}

	return OpeningResponse{
		ID:        o.ID,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		DeletedAt: deletedAt,
		Role:      o.Role,
		Company:   o.Company,
		Location:  o.Location,
		Remote:    o.Remote,
		Link:      o.Link,
		Salary:    o.Salary,
	}
}

type OpeningResponse struct {
	ID        uint       `json:"id" example:"1"`
	CreatedAt time.Time  `json:"createdAt" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"updatedAt" example:"2024-01-01T00:00:00Z"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	Role      string     `json:"role" example:"Backend Developer"`
	Company   string     `json:"company" example:"Google"`
	Location  string     `json:"location" example:"Remote"`
	Remote    bool       `json:"remote" example:"true"`
	Link      string     `json:"link" example:"https://google.com/careers/123"`
	Salary    int64      `json:"salary" example:"120000"`
}
