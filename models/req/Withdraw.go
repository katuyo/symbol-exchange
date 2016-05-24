package req

type Withdraw struct {
    Symbol string `json:"symbol" binding:"Required;MaxSize(10)"`
    Serial string `json:"order_id" binding:"Required"`
}