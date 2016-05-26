package routins

import (
    "os"
    "log"
    "path/filepath"
    "container/list"

    "github.com/katuyo/symbol-exchange/models"
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
    printDepth(logger, sellList, true)
    logger.Println("##########  I'm The Cool Cut-off Line ######## ")
    printDepth(logger, buyList, false)
    logger.Println("**********************************************")
    logger.Println("")
}

func printDepth(logger *log.Logger, l *list.List, sell bool){
    if l == nil {
        logger.Println("No Depth in this order list.");
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
    if (orderEle == nil) {
        logger.Println("No Depth in this order list.");
        return
    }
    o := orderEle.Value.(*models.Order);
    for count > 0 {
        logger.Printf("%s, %.2f, %d", o.GetType(), o.CallPrice(), o.AmountSum());
	if orderEle = orderEle.Next(); orderEle == nil {
	    break;
	}
        count = count - 1
    }
}
