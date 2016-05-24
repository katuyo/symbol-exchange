package controller

import (
    "fmt"
    "github.com/go-macaron/renders"

    "github.com/katuyo/symbol-exchange/models"
    "github.com/katuyo/symbol-exchange/models/req"
    "github.com/katuyo/symbol-exchange/models/res"
)

type OrderController struct{}

func (oc *OrderController) Exchange(render renders.Render, o req.Order) {
    r := oc.validateOrder(o)
    if !r.Result {
	render.JSON(200, res.JSONResult {Result: false, Msg: r.Msg})
    } else {
	newO := models.NewOrder(o.Symbol, o.Type, o.Price, o.Amount)
	models.PushInMarket(newO)
	render.JSON(200, res.JSONResult {Result: true, Order_Id: newO.GetSerial()});
    }
}

func (oc *OrderController) Cancel(ren renders.Render, w req.Withdraw) {
    amount := models.WithDraw(w.Symbol, w.Serial, true)
    if amount == 0 {
        amount = models.WithDraw(w.Symbol, w.Serial, false)
    }
    if amount == 0 {
	ren.JSON(200, res.JSONResult {Result: false, Msg: "Exchanged order."})
    } else {
	ren.JSON(200, res.JSONResult {Result: true, Msg: fmt.Sprintf("Withdrawed amount: %d", amount)})
    }
}

func (oc *OrderController) validateOrder(o req.Order) res.Result {
    s := models.GetStock(o.Symbol);
    if s == nil {
	return res.Result { Result: false, Code:2, Msg: "Stock symbol not exists."}
    }
    maxPrice := s.Open * 1.1
    minPrice := s.Open * 0.9
    if o.Price < minPrice || o.Price > maxPrice {
        return res.Result{ Result: false, Code:1, Msg: "Order price overflow."}
    }
    return res.Result { Result: true, Code: 0, Msg: ""}
}