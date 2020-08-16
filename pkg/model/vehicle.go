package model

import (
	"strings"
	"time"
	"unicode"

	"github.com/opencars/translit"

	"github.com/opencars/wanted/pkg/cleansing"
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
	CheckSum      string    `db:"id" json:"id"`
	Brand         string    `db:"brand" json:"brand"`
	Maker         *string   `db:"maker" json:"maker,omitempty"`
	Model         *string   `db:"model" json:"model,omitempty"`
	Color         *string   `db:"color" json:"color,omitempty"`
	Number        *string   `db:"number" json:"number,omitempty"`
	BodyNumber    *string   `db:"body_number" json:"body_number,omitempty"`
	ChassisNumber *string   `db:"chassis_number" json:"chassis_number,omitempty"`
	EngineNumber  *string   `db:"engine_number" json:"engine_number,omitempty"`
	OVD           string    `db:"ovd" json:"ovd"`
	Kind          *string   `db:"kind" json:"kind,omitempty"`
	Status        Status    `db:"status" json:"status"`
	RevisionID    string    `db:"revision_id" json:"revision_id"`
	TheftDate     string    `db:"theft_date" json:"theft_date"`
	InsertDate    time.Time `db:"insert_date" json:"insert_date"`
}

func (v *Vehicle) BeforeCreate(c *cleansing.Cleansing) {
	if v.Brand == "" {
		return
	}

	maker, model, err := c.Brand(v.Brand)
	if err != nil {
		return
	}

	v.Maker = &maker
	v.Model = &model
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
	kind := vehicle.CarType
	if kind != nil {
		*kind = strings.TrimFunc(*kind, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		// Remove unnecessary lexemes from vehicle kind.
		for _, lexeme := range []string{"АВТОБУС ", "АВТОТРАНСПОРТ"} {
			*kind = strings.ReplaceAll(strings.ToUpper(*kind), lexeme, "")
		}

		*kind = strings.TrimFunc(*kind, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})
	}

	brand := ParseBrand(vehicle.Brand)

	inserted, err := time.Parse(TimeLayout, vehicle.InsertDate)
	if err != nil {
		return nil, err
	}

	return &Vehicle{
		CheckSum:      vehicle.CheckSum,
		Brand:         strings.ToUpper(strings.TrimSpace(brand)),
		Color:         fixedColor(vehicle.Color),
		Number:        fixedNumber(vehicle.Number),
		BodyNumber:    utils.Trim(vehicle.BodyNumber),
		ChassisNumber: utils.Trim(vehicle.ChassisNumber),
		EngineNumber:  utils.Trim(vehicle.EngineNumber),
		OVD:           strings.TrimSpace(vehicle.OVD),
		Kind:          vehicle.CarType,
		Status:        StatusStolen,
		RevisionID:    revision,
		TheftDate:     vehicle.TheftDate[0:10],
		InsertDate:    inserted,
	}, nil
}
