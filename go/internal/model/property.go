package model

import "time"

type Property struct {
    ID           int     `gorm:"primaryKey" json:"id"`
    AgentID      int     `gorm:"not null" validate:"required" json:"agent_id"`
    Latitude     float64 `gorm:"not null" validate:"required" json:"latitude"`
    Longitude    float64 `gorm:"not null" validate:"required" json:"longitude"`
    Price        float64 `gorm:"not null" validate:"required,gt=0" json:"price"`
    AreaSurface  float64 `gorm:"not null" validate:"required,gt=0" json:"area_surface"`
    PropertyType string  `gorm:"not null" validate:"required" json:"property_type"`
    Floors       int     `gorm:"not null" validate:"gte=0" json:"floors"`
    Rooms        int     `gorm:"not null" validate:"gte=0" json:"rooms"`
    Description  string  `gorm:"type:text" validate:"required" json:"description"`
    HasParking                bool      `json:"has_parking"`
    DistanceToCityCenter      float64   `validate:"gte=0" json:"distance_to_city_center"`
    HospitalNearby            bool      `json:"hospital_nearby"`
    PoliceStationNearby       bool      `json:"police_station_nearby"`
    FireStationNearby         bool      `json:"fire_station_nearby"`
    PublicTransportAccessible bool      `json:"public_transport_accessible"`
    CreatedAt                 time.Time `gorm:"autoCreateTime" json:"created_at"`
}
