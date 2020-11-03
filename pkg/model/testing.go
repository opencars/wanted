package model

import (
	"testing"
	"time"
)

func TestRevision(t *testing.T) *Revision {
	t.Helper()

	layout, str := "2006-01-02 15:04:05", "2018-10-01 22:38:14"
	createdAt, err := time.Parse(layout, str)
	if err != nil {
		t.FailNow()
	}

	return &Revision{
		ID:          "01102018_2",
		URL:         "https://data.gov.ua/dataset/9b0e87e0-eaa3-4f14-9547-03d61b70abb6/resource/06e65b06-3120-4713-8003-7905a83f95f5/revision/01102018_2",
		FileHashSum: nil,
		Removed:     11,
		Added:       12,
		CreatedAt:   createdAt,
	}
}

func TestVehicle(t *testing.T) *Vehicle {
	t.Helper()

	color := "СІРИЙ"
	number := "СВ5501ВХ"
	bodyNumber := "5YJSA1E28HF176944"
	kind := "ЛЕГКОВИЙ"
	maker := "TESLA"
	model := "MODEL S"
	brand := "TESLA - MODEL S"

	insertDate, err := time.Parse(TimeLayout, "2019-08-16T15:37:54")
	if err != nil {
		t.FailNow()
	}

	return &Vehicle{
		CheckSum:      "e8d9268a0d5b0f8235a3401013d72d9f",
		Brand:         &brand,
		Maker:         &maker,
		Model:         &model,
		Color:         &color,
		Number:        &number,
		BodyNumber:    &bodyNumber,
		ChassisNumber: nil,
		EngineNumber:  nil,
		OVD:           "СОЛОМ’ЯНСЬКЕ УПРАВЛІННЯ ПОЛІЦІЇ ГУНП В М. КИЄВІ",
		Kind:          &kind,
		Status:        "stolen",
		RevisionID:    "17082019_1",
		TheftDate:     "2019-08-16",
		InsertDate:    insertDate,
	}
}
