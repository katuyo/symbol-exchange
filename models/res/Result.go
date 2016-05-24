package res

type Result struct {
    Result bool `json:"result"`
    Code int32 `json:"code"`
    Msg string `json:"msg"`
}
