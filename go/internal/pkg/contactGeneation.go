package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Estate-CRM/backend-go/internal/model"
	"github.com/jung-kurt/gofpdf"
)

func GenerateContractPDF(contactID, clientID, propertyID int, property model.Property) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(40, 10, "Real Estate Contract")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)

	pdf.Cell(0, 10, fmt.Sprintf("Property ID: %d", property.ID))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Agent ID: %d", property.AgentID))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Latitude: %.6f", property.Latitude))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Longitude: %.6f", property.Longitude))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Price: %.2f DZD", property.Price))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Area Surface: %.2f m²", property.AreaSurface))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Property Type: %s", property.PropertyType))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Floors: %d", property.Floors))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Rooms: %d", property.Rooms))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Description: %s", property.Description))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Has Parking: %t", property.HasParking))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Distance to City Center: %.2f km", property.DistanceToCityCenter))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Hospital Nearby: %t", property.HospitalNearby))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Police Station Nearby: %t", property.PoliceStationNearby))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Fire Station Nearby: %t", property.FireStationNearby))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Public Transport Accessible: %t", property.PublicTransportAccessible))
	pdf.Ln(8)

	pdf.Cell(0, 10, fmt.Sprintf("Created At: %s", property.CreatedAt.Format("2006-01-02 15:04")))
	pdf.Ln(10)

	filename := fmt.Sprintf("%d_%d_contract.pdf", contactID, propertyID)
	path := filepath.Join("static", "contracts", filename)

	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return "", err
	}

	err = pdf.OutputFileAndClose(path)
	if err != nil {
		return "", err
	}

	// in future it wwill stored n a blob in supabase instead
	return fmt.Sprintf("https://example.com/static/contracts/%s", filename), nil
}
