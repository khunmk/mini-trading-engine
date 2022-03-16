package engine

import (
	"github.com/khunmk/mini-trading-engine/model"
)

type OrderBook struct {
	BuyOrders  []*model.Order
	SellOrders []*model.Order
}

type buyTrade struct {
	Id          uint64  `json:"id"`
	Price       float64 `json:"price"`
	OrderAmount float64 `json:"order_amount"`
	TradeAmount float64 `json:"trade_amount"`
	Side        int8    `json:"side"`
}

type sellTrade struct {
	Id          uint64  `json:"id"`
	Price       float64 `json:"price"`
	OrderAmount float64 `json:"order_amount"`
	TradeAmount float64 `json:"trade_amount"`
	Side        int8    `json:"side"`
}

func NewOrderBook() *OrderBook {

	//buy trade
	bt := make([]*buyTrade, 0)
	db := model.DB.Model(&model.Order{})
	db = db.Select(`
		orders.id, 
		orders.price, 
		orders.side,
		orders.amount as order_amount, 
		t.amount as trade_amount
	`)
	db = db.Joins(`
		LEFT JOIN (
			SELECT 
			taker_order_id,
			sum(amount) as amount
		FROM trades
		GROUP BY taker_order_id
		) t ON t.taker_order_id = orders.id
	`)
	db = db.Where("side = ?", 1)
	db = db.Order("orders.price ASC")
	if err := db.Find(&bt).Error; err != nil {
		return nil
	}

	buyOrders := make([]*model.Order, 0)
	for _, item := range bt {
		if item.TradeAmount != 0 {
			item.OrderAmount = item.OrderAmount - item.TradeAmount
		}
		if item.OrderAmount > 0 {
			buyOrders = append(buyOrders, &model.Order{
				ID:     item.Id,
				Price:  item.Price,
				Amount: item.OrderAmount,
				Side:   item.Side,
			})
		}
	}

	//sell trade
	st := make([]*sellTrade, 0)
	db = model.DB.Model(&model.Order{})
	db = db.Select(`
		orders.id, 
		orders.price, 
		orders.side,
		orders.amount as order_amount, 
		t.amount as trade_amount
	`)
	db = db.Joins(`
		LEFT JOIN (
			SELECT 
			maker_order_id,
			sum(amount) as amount
		FROM trades
		GROUP BY maker_order_id
		) t ON t.maker_order_id = orders.id
	`)
	db = db.Where("side = ?", 2)
	db = db.Order("orders.price DESC")
	if err := db.Find(&st).Error; err != nil {
		return nil
	}

	sellOrders := make([]*model.Order, 0)
	for _, item := range st {
		if item.TradeAmount != 0 {
			item.OrderAmount = item.OrderAmount - item.TradeAmount
		}
		if item.OrderAmount > 0 {
			sellOrders = append(sellOrders, &model.Order{
				ID:     item.Id,
				Price:  item.Price,
				Amount: item.OrderAmount,
				Side:   item.Side,
			})
		}
	}

	return &OrderBook{
		BuyOrders:  buyOrders,
		SellOrders: sellOrders,
	}
}

func (book *OrderBook) addBuyOrder(order *model.Order) {
	n := len(book.BuyOrders)
	var i int
	for i = n - 1; i >= 0; i-- {
		buyOrder := book.BuyOrders[i]
		if buyOrder.Price < order.Price {
			break
		}
	}
	if i == n-1 {
		book.BuyOrders = append(book.BuyOrders, order)
	} else {
		l1 := len(book.BuyOrders[0 : i+1])
		l2 := len(book.BuyOrders[i+1:])
		arr1 := make([]*model.Order, l1)
		arr2 := make([]*model.Order, l2)
		copy(arr1, book.BuyOrders[0:i+1])
		copy(arr2, book.BuyOrders[i+1:])
		arr1 = append(arr1, order)
		arr1 = append(arr1, arr2...)
		book.BuyOrders = arr1
	}
}

func (book *OrderBook) addSellOrder(order *model.Order) {
	n := len(book.SellOrders)
	var i int
	for i = n - 1; i >= 0; i-- {
		sellOrder := book.SellOrders[i]
		if sellOrder.Price > order.Price {
			break
		}
	}
	if i == n-1 {
		book.SellOrders = append(book.SellOrders, order)
	} else {
		l1 := len(book.SellOrders[0 : i+1])
		l2 := len(book.SellOrders[i+1:])
		arr1 := make([]*model.Order, l1)
		arr2 := make([]*model.Order, l2)
		copy(arr1, book.SellOrders[0:i+1])
		copy(arr2, book.SellOrders[i+1:])
		arr1 = append(arr1, order)
		arr1 = append(arr1, arr2...)
		book.SellOrders = arr1
	}
}

func (book *OrderBook) removeBuyOrder(index int) {
	book.BuyOrders = append(book.BuyOrders[:index], book.BuyOrders[index+1:]...)
}

func (book *OrderBook) removeSellOrder(index int) {
	book.SellOrders = append(book.SellOrders[:index], book.SellOrders[index+1:]...)
}
