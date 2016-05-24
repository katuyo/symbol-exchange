package main

import (
    "gopkg.in/macaron.v1"
    "github.com/go-macaron/binding"

    "./controller"
    "./models/req"
)

func configRoutes(m *macaron.Macaron){

    m.Get("/", func() string {
	return "Hello world!"
    });

    orderController := new (controller.OrderController);
    m.Post("/trade.do", binding.Json(req.Order{}), orderController.Exchange);
    m.Post("/cancel_order.do", binding.Json(req.Withdraw{}), orderController.Cancel);

    m.Post("/order.exchange", binding.Json(req.Order{}), orderController.Exchange);
    m.Post("/order.withdraw", binding.Json(req.Withdraw{}), orderController.Cancel);
}
