package main

import (
    "gopkg.in/macaron.v1"
    "github.com/go-macaron/binding"

    "github.com/katuyo/symbol-exchange/controller"
    "github.com/katuyo/symbol-exchange/models/req"
)

func configRoutes(m *macaron.Macaron){

    m.Get("/", func() string {
	return "Hello world!"
    });

    orderController := new (controller.OrderController);
    m.Post("/trade.do", binding.Bind(req.Order{}), orderController.Exchange);
    m.Post("/cancel_order.do", binding.Bind(req.Withdraw{}), orderController.Cancel);

    m.Post("/order.exchange", binding.Bind(req.Order{}), orderController.Exchange);
    m.Post("/order.withdraw", binding.Bind(req.Withdraw{}), orderController.Cancel);
}
