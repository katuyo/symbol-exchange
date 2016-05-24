package routins

import (
    "os"
    "log"
    "container/list"

    "../models"
    "path/filepath"
)

func PrintDepth(symbol string){
    r, _ := filepath.Abs("logs")
    f, err := os.OpenFile(r + "/" + symbol + "_depth.log", os.O_APPEND | os.O_CREATE, 0666)
    if err != nil {
	log.Fatalf("Error opening file: %v", err)
    }
    defer f.Close()
    logger := log.New(f, symbol, log.Llongfile)
    logger.Println("**********************************************")
    logger.Println("Type, Price, Amount")
    printDepth(models.GetSellOrders(symbol))
    logger.Println("##########  I'm The Cool Cut-off Line ######## ")
    printDepth(models.GetBuyOrders(symbol))
    logger.Println("**********************************************")
}

func printDepth(l *list.List){
    if l == nil {
        return
    }
    orderEle := l.Front()
    o := orderEle.Value.(models.Order);
    for{
        log.Printf("%s, %s, %d", o.Type, o.Price, o.AmountSum());
	if orderEle = orderEle.Next(); orderEle == nil {
	    break;
	}
    }
}
