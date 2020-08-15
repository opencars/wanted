package model

import (
	"crypto/md5"
	"encoding/hex"
)

type WantedVehicle struct {
	OVD           string  `db:"ovd" json:"organunit"`
	Brand         string  `db:"brand" json:"brandmodel"`
	Color         *string `db:"color" json:"color"`
	CarType       *string `db:"car_type" json:"cartype"`
	Number        *string `db:"number" json:"vehiclenumber"`
	BodyNumber    *string `db:"body_number" json:"bodynumber"`
	ChassisNumber *string `db:"chassis_number" json:"chassisnumber"`
	EngineNumber  *string `db:"engine_number" json:"enginenumber"`
	TheftDate     string  `db:"theft_date" json:"illegalseizuredate"`
	InsertDate    string  `db:"insert_date" json:"insertdate"`
	CheckSum      string  `db:"check_sum" json:"check_sum"`
}

func (w *WantedVehicle) CalculateCheckSum() string {
	data := ""

	if w.Number != nil && *w.Number != "" {
		data += *w.Number
	}

	if w.ChassisNumber != nil && *w.ChassisNumber != "" {
		data += *w.ChassisNumber
	}

	if w.EngineNumber != nil && *w.EngineNumber != "" {
		data += *w.EngineNumber
	}

	checksum := md5.Sum([]byte(data))
	return hex.EncodeToString(checksum[:])
}
