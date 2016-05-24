package routins

import (
    "os"
    "log"
    "container/list"

    "../models"
)

func PrintDepth(symbol string){
    f, err := os.OpenFile("../logs/" + symbol + "_depth.log", os.O_APPEND | os.O_CREATE, 0666)
    if err != nil {
	log.Fatalf("Error opening file: %v", err)
    }
    defer f.Close()
    log.SetOutput(f)
    log.Println("**********************************************")
    log.Println("Type, Price, Amount")
    printDepth(models.GetSellOrders(symbol))
    log.Println("##########  I'm The Cool Cut-off Line ######## ")
    printDepth(models.GetBuyOrders(symbol))
    log.Println("**********************************************")
}

func printDepth(l *list.List){
    orderEle := l.Front()
    o := orderEle.Value.(models.Order);
    for{
        log.Printf("%s, %s, %d", o.Type, o.Price, o.AmountSum());
	if orderEle = orderEle.Next(); orderEle == nil {
	    break;
	}
    }
}
