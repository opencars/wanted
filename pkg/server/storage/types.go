package storage

import (
	"time"
)

// WantedVehicle represents model.
// type WantedVehicle struct {
// 	ID            string  `db:"id"`
// 	OVD           string  `db:"ovd"`
// 	Brand         *string `db:"brand"`
// 	Model         *string `db:"model"`
// 	Kind          *string `db:"kind"`
// 	Color         string  `db:"color"`
// 	Plates        *string `db:"plates"`
// 	BodyNumber    *string `db:"body_number"`
// 	ChassisNumber *string `db:"chassis_number"`
// 	EngineNumber  *string `db:"engine_number"`
// 	TheftDate     string  `db:"theft_date"`
// 	State         string  `db:"state"`
// 	InsertDate    string  `db:"insert_date"`
// }

//
type State string

// State values.
const (
	Stolen  State = "STOLEN"
	Removed State = "REMOVED"
)

type WantedVehicle struct {
	ID            string    `db:"id" json:"id"`
	Brand         string    `db:"brand" json:"brand"`
	Color         string    `db:"color" json:"color"`
	Number        string    `db:"number" json:"number"`
	BodyNumber    string    `db:"body_number" json:"body_number"`
	ChassisNumber string    `db:"chassis_number" json:"chassis_number"`
	EngineNumber  string    `db:"engine_number" json:"engine_number"`
	OVD           string    `db:"ovd" json:"ovd"`
	Kind          string    `db:"kind" json:"kind"`
	TheftDate     time.Time `db:"theft_date" json:"theft_date"`
	InsertDate    time.Time `db:"insert_date" json:"insert_date"`
	State         State     `db:"state" json:"state"`
}
