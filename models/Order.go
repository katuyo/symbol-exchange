package models

import (
    "time"
    "os"
    "log"
    "container/list"

    "github.com/satori/go.uuid"
    "gopkg.in/macaron.v1"
    "github.com/go-macaron/binding"

    "./result"
)

const (
    ORDER_TYPE_BUY = "buy_market" //市价买
    ORDER_TYPE_SELL = "sell_market" //市价卖
    ORDER_TYPE_ASK = "buy"  //限价买
    ORDER_TYPE_BID = "sell"  //限价卖
)

const (
    ORDER_STATUS_DONE = ""
    ORDER_STATUS_DONE_PART = ""
    ORDER_STATUS_WITHDRAW = ""
    ORDER_STATUS_WITHDRAW_PART = ""
)

type Order struct {
    serial string
    Symbol string `json:"symbol" binding:"Required;MaxSize(10)"`
    Type string   `json:"type" binding:"In(buy, sell, buy_market, sell_market)"`
    Price float32 `json:"price"`
    Amount int    `json:"amount" binding:"Required;Range(1, 999)"`
    timestamp time.Time
    next *Order // For depth, Since the same price orders is merged , if not to lose order, linked them.
}

func (o *Order) GetSerial() string {
    return o.serial
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

func (o *Order) Refactor() *Order {
    if o.serial != "" {
        return o
    }
    o.serial = uuid.NewV4().String()
    o.timestamp = time.Now()
    return o
}

func (o *Order) AmountSum() int {
    amount := o.Amount
    for {
        order := o.GetNext()
        if order == nil {
            return amount
        }
        amount = amount + order.Amount
    }
}

func (o *Order) Log(status string) {
    f, err := os.OpenFile("../logs/" + o.Symbol + "_depth.log", os.O_APPEND | os.O_CREATE, 0666)
    if err != nil {
        log.Fatalf("Error opening file: %v", err)
    }
    defer f.Close()
    log.SetOutput(f)
    log.Printf("%s, %s, %s, %f %d %s", time.Now(), o.serial, o.Type, o.Price, o.Amount, status)
    log.Println("")
}

/**
 #########################################
  Macron's Validation Framework implements
 #########################################*/

func (o Order) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
    if (o.Type == ORDER_TYPE_ASK || o.Type == ORDER_TYPE_BID) && o.Price == 0 {
        errs = append(errs, binding.Error{
            FieldNames:     []string{"Price"},
            Classification: "UnexpectedZeroError",
            Message:        "While ask and bid, price should be provided",
        })
    }
    return errs
}

func (o Order) Error(ctx *macaron.Context, errs binding.Errors) {
    if errs.Len() > 0 {
        ctx.Render.JSON(200, result.JSONResult{Result:false, Msg: "Parameter format error."})
    } else {
        ctx.Next()
    };
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

func PushInOrders(o *Order)  {
    buyOrder := true;
    orderList := buyOrdersMap[o.Symbol]
    if o.Type == ORDER_TYPE_BID || o.Type == ORDER_TYPE_SELL {
        buyOrder = false;
        orderList = sellOrdersMap[o.Symbol]
    }
    if orderList == nil {
        orderList = list.New()
        if buyOrder {
            buyOrdersMap[o.Symbol] = orderList
        } else {
            sellOrdersMap[o.Symbol] = orderList
        }
    }
    pushInOrders(orderList, o)
}

func pushInOrders(orders *list.List, o *Order) {
    if(orders.Len() == 0){
        orders.PushBack(o)
        return
    }
    baseOrderEle := orders.Front()
    for {
        baseOrder := baseOrderEle.Value.(Order)
        if baseOrder.Price < o.Price {
            orders.InsertBefore(o, baseOrderEle)
            return
        }
        if baseOrder.Price == o.Price {
            baseOrder.SetNext(o)
            return
        }
        if baseOrderEle = baseOrderEle.Next(); baseOrderEle == nil {
            orders.PushBack(o)
            break;
        }
    }
}
