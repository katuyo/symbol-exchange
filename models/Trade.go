package models

import (
    "time"
    "os"
    "log"
    "path/filepath"
)

type Trade struct {
    buyOrder *Order
    sellOrder *Order
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
    logger.Printf("%s %f %d ", t.timestamp.String(), t.price, t.amount)
    logger.Println("")
}

func PushInMarket (o *Order) {
    if(o.GetType() == ORDER_TYPE_ASK || o.GetType() == ORDER_TYPE_BID) {
        PushOrder(o)
    } else {
        list, ele := OrderOne(o.GetSymbol(), o.GetType() == ORDER_TYPE_BUY)
        if  list == nil {
            o.End()
            return
        }
        exchange(o, ele.Value.(*Order))
    }
    hedgeOrders(o.GetSymbol())
}

func exchange(o *Order, marketOrder *Order) int {
    o.DealPrice(marketOrder.CallPrice())
    amount := 0
    currentMarketOrder := marketOrder
    for{
        if o.GetAmount() < currentMarketOrder.GetAmount() {
            o.Deal(o.GetAmount())
            currentMarketOrder.Deal(o.GetAmount())
            amount = amount + o.GetAmount()
            break;
        } else if o.GetAmount() == currentMarketOrder.GetAmount() {
            o.Deal(o.GetAmount())
            currentMarketOrder.Deal(o.GetAmount())
            amount = amount + o.GetAmount()
            break;
        }
        currentMarketOrder.Deal(currentMarketOrder.GetAmount())
        o.Deal(currentMarketOrder.GetAmount())
        amount = amount + currentMarketOrder.GetAmount()
        if (currentMarketOrder.GetNext() != nil) { //Continue;
            currentMarketOrder = currentMarketOrder.GetNext();
            continue
        } else if (o.GetType() == ORDER_TYPE_BUY || o.GetType() == ORDER_TYPE_SELL){
            o.End();// Withdraw Rest Amount.
        }
        break
    }
    PopOrderOne(o.symbol)
    b, s := buySell(o, marketOrder)
    trade := &Trade {buyOrder: b, sellOrder: s, timestamp: time.Now(),
        price: marketOrder.CallPrice(), amount: amount}
    trade.Log()
    return amount
}

func hedgeOrders(symbol string) {
    list, buyMaxEle := OrderOne(symbol, true)
    if (list == nil) {
        return
    }
    list, sellMinEle := OrderOne(symbol, false)
    if (list == nil) {
        return
    }

    buyMax := buyMaxEle.Value.(*Order)
    sellMin := sellMinEle.Value.(*Order)
    if buyMax.GetPrice() < sellMin.GetPrice() {
        return
    }
    exchange(chooseBaseOrder(buyMax, sellMin))
    hedgeOrders(symbol)
}

//TODO While buyMaxPrice > sellMinPrice , which would be the deal price ?
func chooseBaseOrder (b *Order, s *Order) (*Order, *Order) {
    if (b.AmountSum() > s.AmountSum()){
        return s, b;
    } else {
        return b, s;
    }
}

func buySell (l *Order, r *Order) (*Order, *Order) {
    if l.GetType() == ORDER_TYPE_BUY || l.GetType() == ORDER_TYPE_ASK {
        return l, r
    } else {
        return r, l
    }
}