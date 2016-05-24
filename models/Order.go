package models

import (
    "time"
    "os"
    "log"
    "container/list"
    "path/filepath"

    "github.com/satori/go.uuid"
)

const (
    ORDER_TYPE_BUY = "buy_market" //市价买
    ORDER_TYPE_SELL = "sell_market" //市价卖
    ORDER_TYPE_ASK = "buy"  //限价买
    ORDER_TYPE_BID = "sell"  //限价卖
)

const (
    ORDER_STATUS_NEW = "Created"
    ORDER_STATUS_DONE = "Exchanged"
    ORDER_STATUS_DONE_PART = "PartExchanged"
    ORDER_STATUS_WITHDRAW = "Withdrawed"
    ORDER_STATUS_WITHDRAW_REST = "RestWithdrawed"
)

type Order struct {
    serial string
    symbol string
    type_ string
    price float32
    amount int
    status string
    timestamp time.Time
    next *Order // For depth, Since the same price orders is merged , if not to lose order, linked them.
}

func NewOrder(symbol string, type_ string, price float32, amount int) *Order {
    return &Order{symbol: symbol, serial: uuid.NewV4().String(), type_: type_,
        price: price, amount: amount, status: ORDER_STATUS_NEW, timestamp: time.Now()}
}

func (o *Order) GetSerial() string {
    return o.serial
}

func (o *Order) GetSymbol() string {
    return o.symbol
}

func (o *Order) GetType() string {
    return o.type_
}

func (o *Order) GetPrice() float32 {
    return o.price
}
func (o *Order) SetPrice(p float32) {
    o.price = p
}

func (o *Order) GetAmount() int {
    return o.amount
}

func (o *Order) GetNext() *Order {
    return o.next
}

func (o *Order) GetTimestamp() time.Time {
    return o.timestamp
}

func (o *Order) SetNext(next *Order) {
    o.next = next
}

func (o *Order) AmountSum() int {
    amount := o.amount
    for {
        order := o.GetNext()
        if order == nil {
            return amount
        }
        amount = amount + order.amount
    }
}

func (o *Order) Deal(amount int) *Order {
    if o.status == ORDER_STATUS_WITHDRAW || o.status == ORDER_STATUS_WITHDRAW_REST {
        return o
    }
    o.amount = o.amount - amount;
    if o.amount == 0 {
        o.status = ORDER_STATUS_DONE
        o.log(amount)
    } else {
        o.status = ORDER_STATUS_DONE_PART
    }
    return o
}

func (o *Order) End() *Order {
    if o.status == ORDER_STATUS_NEW {
        o.status = ORDER_STATUS_WITHDRAW
    } else {
        o.status = ORDER_STATUS_WITHDRAW_REST
    }
    o.log(0)
    return o
}

func (o *Order) log(amountPart int) {
    r, _ := filepath.Abs("logs")
    f, err := os.OpenFile(r + "/" + o.symbol + "_order.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Error opening file: %v", err)
    }
    defer f.Close()
    logger := log.New(f, o.symbol + " ", 0)
    amount := amountPart
    if amount == 0 {
        amount = o.amount
    }
    logger.Printf("%s, %s, %s, %f %d %s", time.Now(), o.serial, o.type_, o.price, amount, o.status)
    logger.Println("")
}

/**
 #########################################
    Store ASK and BID orders in memory
 #########################################*/
var buyOrdersMap = make (map[string]*list.List)
var sellOrdersMap = make (map[string]*list.List)

func GetBuyOrders(symbol string) *list.List{
    return buyOrdersMap[symbol]
}

func GetSellOrders(symbol string) *list.List{
    return sellOrdersMap[symbol]
}

func PushOrder(o *Order)  {
    buyOrder := true;
    orderList := buyOrdersMap[o.symbol]
    if o.type_ == ORDER_TYPE_BID || o.type_ == ORDER_TYPE_SELL {
        buyOrder = false;
        orderList = sellOrdersMap[o.symbol]
    }
    if orderList == nil {
        orderList = list.New()
        if buyOrder {
            buyOrdersMap[o.symbol] = orderList
        } else {
            sellOrdersMap[o.symbol] = orderList
        }
    }
    pushOrders(orderList, o)
}

func pushOrders(orders *list.List, o *Order) {
    if(orders.Len() == 0){
        orders.PushBack(o)
        return
    }
    baseOrderEle := orders.Front()
    for {
        baseOrder := baseOrderEle.Value.(Order)
        if baseOrder.price < o.price {
            orders.InsertBefore(o, baseOrderEle)
            return
        }
        if baseOrder.price == o.price {
            baseOrder.SetNext(o)
            return
        }
        if baseOrderEle = baseOrderEle.Next(); baseOrderEle == nil {
            orders.PushBack(o)
            break;
        }
    }
}

func PopOrder(o *Order){
    l, ele := findEleInList(o.symbol, o.serial, o.type_ == ORDER_TYPE_ASK || o.type_ == ORDER_TYPE_BUY)
    order := ele.Value.(*Order)
    for {
        var nextEle *list.Element
        if order.status == ORDER_STATUS_NEW || order.status == ORDER_STATUS_DONE_PART {
            return
        } else {
            nextEle = ele.Next()
            l.Remove(ele)
        }
        if order.next != nil {
            return
        }
        ele = l.InsertBefore(order.next, nextEle);
        order = order.next
    }
}

func findEleInList(symbol string, serial string, buy bool) (*list.List, *list.Element) {
    var orderList *list.List
    if buy {
        orderList = GetBuyOrders(symbol)
    } else {
        orderList = GetSellOrders(symbol)
    }
    orderEle := orderList.Front()
    for{
        order := orderEle.Value.(*Order)
        if order.serial == serial {
            return orderList, orderEle
            break;
        } else {
            order = order.next
        }
        if orderEle = orderEle.Next(); orderEle == nil {
            break;
        }
    }
    return orderList, nil
}

func WithDraw(symbol string, serial string) int {
    return 0
}