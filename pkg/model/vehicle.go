package model

import (
	"strings"
	"time"
	"unicode"

	"github.com/opencars/translit"
	"github.com/opencars/wanted/pkg/utils"
)

const (
	// TimeLayout represent time layout for parsing time.
	TimeLayout = "2006-01-02T15:04:05"
)

// Status represent status of the vehicle search.
// Can be either "stolen" or "removed".
type Status string

const (
	// StatusStolen means, that vehicle was added to the registry.
	StatusStolen Status = "stolen"
	// StatusRemoved means, that vehicle was added, and then removed from the registry.
	StatusRemoved Status = "removed"
)

// Vehicle represents storage model for vehicles entity.
type Vehicle struct {
	ID            string    `db:"id" json:"id"`
	Brand         string    `db:"brand" json:"brand"`
	Color         *string   `db:"color" json:"color,omitempty"`
	Number        *string   `db:"number" json:"number,omitempty"`
	BodyNumber    *string   `db:"body_number" json:"body_number,omitempty"`
	ChassisNumber *string   `db:"chassis_number" json:"chassis_number,omitempty"`
	EngineNumber  *string   `db:"engine_number" json:"engine_number,omitempty"`
	OVD           string    `db:"ovd" json:"ovd"`
	Kind          string    `db:"kind" json:"kind"`
	Status        Status    `db:"status" json:"status"`
	RevisionID    string    `db:"revision_id" json:"revision_id"`
	TheftDate     string    `db:"theft_date" json:"theft_date"`
	InsertDate    time.Time `db:"insert_date" json:"insert_date"`
}

func fixedColor(color *string) *string {
	color = utils.Trim(color)

	if color != nil {
		for _, lexeme := range []string{"НЕВИЗНАЧЕНИЙ", "НЕОПРЕДЕЛЕН"} {
			*color = strings.ReplaceAll(strings.ToUpper(*color), lexeme, "")
		}
	}

	return utils.Trim(color)
}

func fixedNumber(number *string) *string {
	number = utils.Trim(number)

	if number != nil {
		*number = translit.ToUA(*number)
	}

	return number
}

func VehicleFromGov(revision string, vehicle *WantedVehicle) (*Vehicle, error) {
	brand, kind := ParseKind(vehicle.Brand)

	// Remove unnecessary lexemes from vehicle kind.
	for _, lexeme := range []string{"АВТОБУС ", "АВТОТРАНСПОРТ"} {
		kind = strings.ReplaceAll(strings.ToUpper(kind), lexeme, "")
	}

	kind = strings.TrimFunc(kind, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	inserted, err := time.Parse(TimeLayout, vehicle.InsertDate)
	if err != nil {
		return nil, err
	}

	return &Vehicle{
		ID:            vehicle.ID,
		Brand:         strings.ToUpper(strings.TrimSpace(brand)),
		Color:         fixedColor(vehicle.Color),
		Number:        fixedNumber(vehicle.Number),
		BodyNumber:    utils.Trim(vehicle.BodyNumber),
		ChassisNumber: utils.Trim(vehicle.ChassisNumber),
		EngineNumber:  utils.Trim(vehicle.EngineNumber),
		OVD:           strings.TrimSpace(vehicle.OVD),
		Kind:          kind,
		Status:        StatusStolen,
		RevisionID:    revision,
		TheftDate:     vehicle.TheftDate[0:10],
		InsertDate:    inserted,
	}, nil
}
