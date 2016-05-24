package models

import (
    "time"
    "os"
    "log"
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
    f, err := os.OpenFile("../logs/depth.log", os.O_APPEND | os.O_CREATE, 0666)
    if err != nil {
        log.Fatalf("Error opening file: %v", err)
    }
    defer f.Close()
    log.SetOutput(f)
    log.Println("**********************************************")
    log.Println("Type, Price, Amount")

    log.Println("##########  I'm The Cool Cut-off Line ######## ")

    log.Println("**********************************************")
}

func PushInMarket (o *Order){
    if o.Type == ORDER_TYPE_BUY {
        sellEle := GetSellOrders(o.Symbol).Back()
        sellOrder := sellEle.Value.(Order)

        amount := sellOrder.Amount
        if o.Amount <= amount {
            amount = o.Amount
        } else {
            o.Amount = o.Amount - amount
            sellOrder.Log(ORDER_STATUS_DONE)
        }
        o.Log(ORDER_STATUS_DONE)
        o.Log(ORDER_STATUS_WITHDRAW_PART)
        trade := Trade {buyOrder: o, sellOrder: &sellOrder, type_: TRADE_TYPE_ALL,
            price: sellOrder.Price, amount: amount, timestamp: time.Now()}
        trade.Log()
    } else if o.Type == ORDER_TYPE_SELL {
        buyEle := GetBuyOrders(o.Symbol).Front()
        buyOrder := buyEle.Value.(Order)

        amount := buyOrder.AmountSum()
        if o.Amount <= amount {
            amount = o.Amount
        }

        o.Log(ORDER_STATUS_DONE)
        o.Log(ORDER_STATUS_WITHDRAW_PART)
        trade := Trade {buyOrder: o, sellOrder: &buyOrder, type_: TRADE_TYPE_ALL,
            price: buyOrder.Price, amount: amount, timestamp: time.Now()}
        trade.Log()
    } else {
        PushInOrders(o)
    }
    rollTrade(o.Symbol)
}

func rollTrade(symbol string) {
    buyMaxEle := GetBuyOrders(symbol).Front()
    sellMinEle := GetSellOrders(symbol).Back()
    buyMax := buyMaxEle.Value.(Order)
    sellMin := sellMinEle.Value.(Order)
    if buyMax.Price < sellMin.Price {
        return
    }
}

