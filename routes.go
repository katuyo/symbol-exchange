package main

import (
    "gopkg.in/macaron.v1"
    "github.com/go-macaron/binding"

    "./controller"
    "./models"
)

func configRoutes(m *macaron.Macaron){

    m.Get("/", func() string {
	return "Hello world!"
    });

    orderController := new (controller.OrderController);
    m.Post("/trade.do", binding.Json(models.Order{}), orderController.Exchange);
    m.Post("/order.exc", binding.Json(models.Order{}), orderController.Exchange);

    m.Post("/order.cancel", orderController.Cancel);
    m.Post("/cancel_order.do", orderController.Cancel);
}
