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
    callPrice float32
    dealPrice float32
    amount int
    status string
    timestamp time.Time
    next *Order // For depth, Since the same price orders is merged , if not to lose order, linked them.
}

func NewOrder(symbol string, type_ string, price float32, amount int) *Order {
    return &Order{symbol: symbol, serial: uuid.NewV4().String(), type_: type_,
        callPrice: price, amount: amount, status: ORDER_STATUS_NEW, timestamp: time.Now()}
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

func (o *Order) CallPrice() float32 {
    return o.callPrice
}

func (o *Order) GetPrice() float32 {
    return o.dealPrice
}

func (o *Order) DealPrice(p float32) {
    o.dealPrice = p
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
        o.log(amount)
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
    logger.Printf("%s, %s, %s, %f %d %s", time.Now(), o.serial, o.type_, o.dealPrice, amount, o.status)
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
        if baseOrder.callPrice < o.callPrice {
            orders.InsertBefore(o, baseOrderEle)
            return
        }
        if baseOrder.callPrice == o.callPrice {
            baseOrder.SetNext(o)
            return
        }
        if baseOrderEle = baseOrderEle.Next(); baseOrderEle == nil {
            orders.PushBack(o)
            break;
        }
    }
}


func OrderOne (symbol string, buy bool) (*list.List, *list.Element) {
    var orderList *list.List
    if buy {
        orderList = GetBuyOrders(symbol)
    } else {
        orderList = GetSellOrders(symbol)
    }
    if (orderList == nil) {
        return nil, nil
    }
    if(buy) {
        return orderList, orderList.Front()
    } else {
        return orderList, orderList.Back();
    }
}
func PopOrderOne(symbol string){
    PopOrder(symbol, true)
    PopOrder(symbol, false)
}
func PopOrder(symbol string, buy bool){
    list, ele := OrderOne(symbol, buy)
    if (list == nil) {
        return
    }
    order := ele.Value.(*Order)
    var preOrder *Order
    for {
        if order.status == ORDER_STATUS_NEW || order.status == ORDER_STATUS_DONE_PART {
            return
        }
        putOrderOut(order, preOrder, list, ele)
        preOrder = order
        order = order.next
        if order == nil {
            return
        }
    }
}

func putOrderOut(order *Order, preOrder *Order, orderList *list.List, orderEle *list.Element){
    order.End()
    if preOrder == nil && order.next == nil {
        orderList.Remove(orderEle)
    } else if preOrder != nil {
        preOrder.SetNext(order.next)
    } else if order.next != nil {
        orderList.InsertAfter(order.next, orderEle)
        orderList.Remove(orderEle)
    }
}

func WithDraw(symbol string, serial string, buy bool) int {
    var orderList *list.List
    if buy {
        orderList = GetBuyOrders(symbol)
    } else {
        orderList = GetSellOrders(symbol)
    }
    if orderList == nil {
        return 0
    }
    orderEle := orderList.Front()
    reset := true;
    var order *Order
    var preOrder *Order
    for{/**Find next order in same price, if not exist, reset, go next in order list*/
        if orderEle == nil {
            return 0
        }
        if reset {
            order = orderEle.Value.(*Order)
            preOrder = nil
            reset = false
        }
        if order.serial == serial {
            putOrderOut(order, preOrder, orderList, orderEle)
            return order.amount
        }
        preOrder = order
        order = order.next
        if order == nil {
            orderEle = orderEle.Next()
            reset = true
        }
    }
}
