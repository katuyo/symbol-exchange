package models

import "time"

type Stock struct {
    Symbol string
    Open  float64
    InitAmount int
    date  time.Time
}

func (s *Stock) GetDate() time.Time {
    return s.date;
}

func (s *Stock) Issue() {
    PushOrder(NewOrder(s.Symbol, ORDER_TYPE_BID, s.Open, s.InitAmount))
}

/** #########################################
    Store Stocks list in memory
 ############################################*/

var stockMap = make (map[string]*Stock)

func PushStock(s *Stock){
    stockMap[s.Symbol] = s;
}

func GetStock(symbol string) *Stock {
    return stockMap[symbol];
}

func GetStockMap() map[string]*Stock {
    return stockMap
}

