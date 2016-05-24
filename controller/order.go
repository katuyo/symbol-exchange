package controller

import (
    "fmt"
    "gopkg.in/macaron.v1"

    "github.com/katuyo/symbol-exchange/models"
    "github.com/katuyo/symbol-exchange/models/req"
    "github.com/katuyo/symbol-exchange/models/res"
)

type OrderController struct{}

func (oc *OrderController) Exchange(ctx *macaron.Context, o req.Order) {
    r := oc.validateOrder(o)
    if !r.Result {
	ctx.Render.JSON(200, res.JSONResult {Result: false, Msg: r.Msg})
    } else {
	newO := models.NewOrder(o.Symbol, o.Type, o.Price, o.Amount)
	models.PushInMarket(newO)
	ctx.Render.JSON(200, res.JSONResult {Result: true, Order_Id: newO.GetSerial()});
    }
}

func (oc *OrderController) Cancel(ctx *macaron.Context, w req.Withdraw) {
    amount := models.WithDraw(w.Symbol, w.Serial, true)
    if amount == 0 {
        amount = models.WithDraw(w.Symbol, w.Serial, false)
    }
    if amount == 0 {
	ctx.Render.JSON(200, res.JSONResult {Result: false, Msg: "Exchanged order."})
    } else {
	ctx.Render.JSON(200, res.JSONResult {Result: true, Msg: fmt.Sprintf("Withdrawed amount: %d", amount)})
    }
}

func (oc *OrderController) validateOrder(o req.Order) res.Result {
    var s *models.Stock
    if s := models.GetStock(o.Symbol); s == nil {
	return res.Result { Result: false, Code:2, Msg: "Stock symbol not exists."}
    }
    max := s.Open * 1.1
    min := s.Open * 0.9
    if o.Price < min || o.Price > max {
        return res.Result{ Result: false, Code:1, Msg: "Order price overflow."}
    }
    return res.Result { Result: true, Code: 0, Msg: ""}
}