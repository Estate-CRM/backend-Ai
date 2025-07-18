package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Estate-CRM/backend-go/internal/model"
)

func ExportContactsToCSV(contacts []model.Contact, filename string) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	headers := []string{
		"id", "client_id", "latitude", "longitude", "min_budget", "max_budget",
		"desired_area_min", "desired_area_max", "property_type", "floors", "rooms",
		"has_parking", "distance_to_city_center", "hospital_nearby", "police_station_nearby",
		"fire_station_nearby", "public_transport_accessible", "is_active", "created_at",
	}
	writer.Write(headers)

	for _, c := range contacts {
		record := []string{
			strconv.Itoa(c.ID),
			strconv.Itoa(c.ClientID),
			fmt.Sprintf("%f", c.Latitude),
			fmt.Sprintf("%f", c.Longitude),
			strconv.Itoa(c.MinBudget),
			strconv.Itoa(c.MaxBudget),
			fmt.Sprintf("%f", c.DesiredAreaMin),
			fmt.Sprintf("%f", c.DesiredAreaMax),
			c.PropertyType,
			strconv.Itoa(c.Floors),
			strconv.Itoa(c.Rooms),
			strconv.FormatBool(c.HasParking),
			fmt.Sprintf("%f", c.DistanceToCityCenter),
			strconv.FormatBool(c.HospitalNearby),
			strconv.FormatBool(c.PoliceStationNearby),
			strconv.FormatBool(c.FireStationNearby),
			strconv.FormatBool(c.PublicTransportAccessible),
			strconv.FormatBool(c.Is_active),
			c.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(record)
	}

	return nil
}
