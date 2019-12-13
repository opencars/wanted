package model

type WantedVehicle struct {
	ID            string  `db:"id" json:"ID"`
	OVD           string  `db:"ovd" json:"OVD"`
	Brand         string  `db:"brand" json:"BRAND"`
	Color         *string `db:"color" json:"COLOR"`
	Number        *string `db:"number" json:"VEHICLENUMBER"`
	BodyNumber    *string `db:"body_number" json:"BODYNUMBER"`
	ChassisNumber *string `db:"chassis_number" json:"CHASSISNUMBER"`
	EngineNumber  *string `db:"engine_number" json:"ENGINENUMBER"`
	TheftDate     string  `db:"theft_date" json:"THEFT_DATA"`
	InsertDate    string  `db:"insert_date" json:"INSERT_DATE"`
}
