package controller

import (
    "gopkg.in/macaron.v1"

    "../models"
    "../models/result"
)

type OrderController struct{}

func (oc *OrderController) Exchange(ctx *macaron.Context, o models.Order) {
    r := oc.validateOrder(o)
    if !r.Result {
	ctx.Render.JSON(200, result.JSONResult {Result: false, Msg: r.Msg})
    } else {
	newO := o.Refactor()
	models.PushInMarket(newO)
	ctx.Render.JSON(200, result.JSONResult {Result: true, Order_Id: newO.GetSerial()});
    }
}

func (oc *OrderController) Cancel(ctx *macaron.Context) {
    ctx.Render.JSON(200, result.JSONResult {Result: true});
}

func (oc *OrderController) validateOrder(o models.Order) result.Result {
    var s *models.Stock
    if s := models.GetStock(o.Symbol); s == nil {
	return result.Result { Result: false, Code:2, Msg: "Stock symbol not exists."}
    }
    max := s.Open * 1.1
    min := s.Open * 0.9
    if o.Price < min || o.Price > max {
        return result.Result{ Result: false, Code:1, Msg: "Order price overflow."}
    }
    return result.Result { Result: true, Code: 0, Msg: ""}
}