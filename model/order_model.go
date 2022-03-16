package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Order struct {
	ID     uint64  `gorm:"column:id;primaryKey" json:"id"`
	Amount float64 `gorm:"column:amount" json:"amount" form:"amount" binding:"required"`
	Price  float64 `gorm:"column:price" json:"price" form:"price" binding:"required"`
	Side   int8    `gorm:"column:side" json:"side" form:"side" binding:"required"` // 1 buy, 2 sell
}

func (order *Order) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, order)
}

func (order *Order) ToJSON() []byte {
	str, _ := json.Marshal(order)
	return str
}

func (order *Order) Create(tx *gorm.DB) error {
	if err := tx.Create(&order).Error; err != nil {
		return err
	}
	return nil
}

func (order *Order) Update(tx *gorm.DB) error {
	if err := tx.Updates(&order).Error; err != nil {
		return err
	}
	return nil
}
