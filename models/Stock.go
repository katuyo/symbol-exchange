package models

import "time"

type Stock struct {
    Symbol string
    Open  float32
    date  time.Time
}

func (s *Stock) GetDate() time.Time {
    return s.date;
}

/** #########################################
    Store Stocks list in memory
 ############################################*/

var stockMap = make (map[string]*Stock)

func PushStock(s Stock){
    stockMap[s.Symbol] = &s;
}

func GetStock(symbol string) *Stock {
    return stockMap[symbol];
}

func GetStockMap() map[string]*Stock {
    return stockMap
}