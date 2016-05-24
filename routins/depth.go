package routins

import (
    "os"
    "log"
    "path/filepath"
    "container/list"

    "../models"
)

func PrintDepth(symbol string){
    r, _ := filepath.Abs("logs")
    f, err := os.OpenFile(r + "/" + symbol + "_depth.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
	log.Fatalf("Error opening file: %v", err)
    }
    defer f.Close()
    sellList := models.GetSellOrders(symbol)
    buyList := models.GetBuyOrders(symbol)
    if( sellList == nil && buyList == nil) {
        return
    }
    logger := log.New(f, symbol + " ", 0)
    logger.Println("**********************************************")
    logger.Println("Type, Price, Amount")
    printDepth(sellList, true)
    logger.Println("##########  I'm The Cool Cut-off Line ######## ")
    printDepth(buyList, false)
    logger.Println("**********************************************")
}

func printDepth(l *list.List, sell bool){
    if l == nil {
        return
    }
    count := 20
    orderEle := l.Front()
    if sell && l.Len() > count {
        index := 0
        for index < l.Len() - count{
            orderEle = orderEle.Next()
            index = index + 1
        }
    }
    o := orderEle.Value.(models.Order);
    for count > 0 {
        log.Printf("%s, %s, %d", o.GetType(), o.GetPrice(), o.AmountSum());
	if orderEle = orderEle.Next(); orderEle == nil {
	    break;
	}
        count = count -1
    }
}
