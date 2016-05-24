package models

import (
    "time"
    "os"
    "log"
    "path/filepath"
)

const (
    TRADE_TYPE_ALL = "All"
    TRADE_TYPE_PART = "Part"
    TRADE_TYPE_REST = "Rest"
)

type Trade struct {
    buyOrder *Order
    sellOrder *Order
    type_ string
    price float32
    amount int
    timestamp time.Time
}

func (t *Trade) Log() {
    r, _ := filepath.Abs("logs")
    symbol := t.sellOrder.GetSymbol()
    f, err := os.OpenFile(r + "/" + symbol + "_trade.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Error opening file: %v", err)
    }
    defer f.Close()
    logger := log.New(f, symbol + " ", log.Ldate|log.Ltime)
    logger.Printf(" %f %d ", t.price, t.amount)
    logger.Println("")
}

func PushInMarket (o *Order){
    if o.GetType() == ORDER_TYPE_BUY {
        sellEle := GetSellOrders(o.GetSymbol()).Back()
        sellOrder := sellEle.Value.(Order)
        amount := orderMarket(o, &sellOrder)
        trade := Trade {buyOrder: o, sellOrder: &sellOrder, type_: TRADE_TYPE_ALL,
            price: sellOrder.GetPrice(), amount: amount, timestamp: time.Now()}
        trade.Log()
    } else if o.GetType() == ORDER_TYPE_SELL {
        buyEle := GetBuyOrders(o.GetSymbol()).Front()
        buyOrder := buyEle.Value.(Order)
        amount := orderMarket(o, &buyOrder)
        trade := Trade {buyOrder: &buyOrder, sellOrder: o, type_: TRADE_TYPE_ALL,
            price: buyOrder.GetPrice(), amount: amount, timestamp: time.Now()}
        trade.Log()
    } else {
        PushOrder(o)
    }
    hedgeOrders(o.GetSymbol())
}

// Exchange immediately, withdraw rest
func orderMarket(o *Order, marketOrder *Order) int {
    o.SetPrice(marketOrder.GetPrice())
    amount := 0
    for{
        if o.GetAmount() < marketOrder.GetAmount() {
            o.Deal(o.GetAmount())
            marketOrder.Deal(o.GetAmount())
            amount = amount + o.GetAmount()
            break;
        } else if o.GetAmount() == marketOrder.GetAmount() {
            o.Deal(o.GetAmount())
            marketOrder.Deal(o.GetAmount())
            amount = amount + o.GetAmount()
            break;
        } else {
            marketOrder.Deal(marketOrder.GetAmount())
            o.Deal(marketOrder.GetAmount())
            amount = amount + marketOrder.GetAmount()
            marketOrder = marketOrder.GetNext()
        }
        //WithDraw Rest
        if marketOrder == nil {
            o.End()
            break;
        }
    }
    PopOrder(marketOrder)
    return amount;
}


func hedgeOrders(symbol string) {
    buyMaxEle := GetBuyOrders(symbol).Front()
    sellMinEle := GetSellOrders(symbol).Back()
    buyMax := buyMaxEle.Value.(Order)
    sellMin := sellMinEle.Value.(Order)
    if buyMax.GetPrice() < sellMin.GetPrice() {
        return
    }
}

