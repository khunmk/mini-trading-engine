package engine

import (
	"fmt"
	"strconv"

	"github.com/khunmk/mini-trading-engine/model"
)

func (book *OrderBook) Process(order *model.Order) []*model.Trade {
	if order.Side == 1 {
		return book.processLimitBuy(order)
	}
	return book.processLimitSell(order)
}

func (book *OrderBook) processLimitBuy(order *model.Order) []*model.Trade {
	trades := make([]*model.Trade, 0, 1)
	n := len(book.SellOrders)

	if n != 0 && book.SellOrders[n-1].Price <= order.Price {
		for i := n - 1; i >= 0; i-- {
			sellOrder := book.SellOrders[i]
			if sellOrder.Price > order.Price {
				break
			}
			if sellOrder.Amount >= order.Amount {
				trades = append(trades, &model.Trade{
					TakerOrderID: order.ID,
					MakerOrderID: sellOrder.ID,
					Amount:       order.Amount,
					Price:        sellOrder.Price,
				})
				amtStr := fmt.Sprintf("%.2f", sellOrder.Amount-order.Amount)
				amt, _ := strconv.ParseFloat(amtStr, 64)
				sellOrder.Amount = amt
				if sellOrder.Amount == 0 {
					book.removeSellOrder(i)
				}
				return trades
			}

			if sellOrder.Amount < order.Amount {
				trades = append(trades, &model.Trade{
					TakerOrderID: order.ID,
					MakerOrderID: sellOrder.ID,
					Amount:       sellOrder.Amount,
					Price:        sellOrder.Price,
				})
				amtStr := fmt.Sprintf("%.2f", order.Amount-sellOrder.Amount)
				amt, _ := strconv.ParseFloat(amtStr, 64)
				order.Amount = amt
				book.removeSellOrder(i)
				continue
			}
		}
	}

	book.addBuyOrder(order)

	return trades
}

func (book *OrderBook) processLimitSell(order *model.Order) []*model.Trade {
	trades := make([]*model.Trade, 0, 1)
	n := len(book.BuyOrders)
	if n != 0 && book.BuyOrders[n-1].Price >= order.Price {
		for i := n - 1; i >= 0; i-- {
			buyOrder := book.BuyOrders[i]
			if buyOrder.Price < order.Price {
				break
			}

			if buyOrder.Amount >= order.Amount {
				trades = append(trades, &model.Trade{
					TakerOrderID: order.ID,
					MakerOrderID: buyOrder.ID,
					Amount:       order.Amount,
					Price:        buyOrder.Price,
				})
				amtStr := fmt.Sprintf("%.2f", buyOrder.Amount-order.Amount)
				amt, _ := strconv.ParseFloat(amtStr, 64)
				buyOrder.Amount = amt
				if buyOrder.Amount == 0 {
					book.removeBuyOrder(i)
				}
				return trades
			}

			if buyOrder.Amount < order.Amount {
				trades = append(trades, &model.Trade{
					TakerOrderID: order.ID,
					MakerOrderID: buyOrder.ID,
					Amount:       order.Amount,
					Price:        buyOrder.Price,
				})
				amtStr := fmt.Sprintf("%.2f", order.Amount-buyOrder.Amount)
				amt, _ := strconv.ParseFloat(amtStr, 64)
				order.Amount = amt
				book.removeBuyOrder(i)
				continue
			}
		}
	}

	book.addSellOrder(order)

	return trades
}
