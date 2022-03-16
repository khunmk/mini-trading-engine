package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Trade struct {
	ID           uint64  `gorm:"column:id;primaryKey" json:"id"`
	TakerOrderID uint64  `json:"taker_order_id"`
	MakerOrderID uint64  `json:"maker_order_id"`
	Amount       float64 `json:"amount"`
	Price        float64 `json:"price"`
}

func (trade *Trade) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, trade)
}

func (trade *Trade) ToJSON() []byte {
	str, _ := json.Marshal(trade)
	return str
}

func (trade *Trade) Create(tx *gorm.DB) error {
	if err := tx.Create(&trade).Error; err != nil {
		return err
	}
	return nil
}
