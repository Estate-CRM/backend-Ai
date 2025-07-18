package model

import (
	"time"
)

type Contact struct {
	ID                        int       `gorm:"primaryKey" json:"id"`
	ClientID                  int       `gorm:"not null" validate:"required" json:"client_id"`
	Latitude                  float64   `gorm:"not null" validate:"required" json:"latitude"`
	Longitude                 float64   `gorm:"not null" validate:"required" json:"longitude"`
	MinBudget                 int       `gorm:"not null" validate:"required,gte=0" json:"min_budget"`
	MaxBudget                 int       `gorm:"not null" validate:"required,gtefield=MinBudget" json:"max_budget"`
	DesiredAreaMin            float64   `validate:"required,gte=0" json:"desired_area_min"`
	DesiredAreaMax            float64   `validate:"required,gtefield=DesiredAreaMin" json:"desired_area_max"`
	PropertyType              string    `gorm:"not null" validate:"required" json:"property_type"`
	Floors                    int       `gorm:"not null" validate:"gte=0" json:"floors"`
	Rooms                     int       `gorm:"not null" validate:"gte=0" json:"rooms"`
	HasParking                bool      `json:"has_parking"`
	DistanceToCityCenter      float64   `validate:"gte=0" json:"distance_to_city_center"`
	HospitalNearby            bool      `json:"hospital_nearby"`
	PoliceStationNearby       bool      `json:"police_station_nearby"`
	FireStationNearby         bool      `json:"fire_station_nearby"`
	PublicTransportAccessible bool      `json:"public_transport_accessible"`
	Is_active                 bool      `gorm:"default:true" json:"is_active"`
	CreatedAt                 time.Time `gorm:"autoCreateTime" json:"created_at"`
}
