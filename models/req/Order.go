package req

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/binding"
	"../res"
	"../../models"
)

type Order struct {
    Symbol string `json:"symbol" binding:"Required;MaxSize(10)"`
    Type string   `json:"type" binding:"In(buy, sell, buy_market, sell_market)"`
    Price float32 `json:"price"`
    Amount int    `json:"amount" binding:"Required;Range(1, 999)"`
}

func (o Order) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
    if (o.Type == models.ORDER_TYPE_ASK || o.Type == models.ORDER_TYPE_BID) && o.Price == 0 {
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
        ctx.Render.JSON(200, res.JSONResult {Result: false, Msg: "Parameter format error."})
    } else {
        ctx.Next()
    };
}