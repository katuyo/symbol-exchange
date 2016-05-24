package res

type JSONResult struct {
    Result bool `json:"result"`
    Order_Id string `json:"order_id"`
    Msg string `json:"message"`
}
