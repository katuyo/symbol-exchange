package main

import (
    "time"

    "gopkg.in/macaron.v1"
    "github.com/go-macaron/renders"

    "github.com/katuyo/symbol-exchange/models"
    "github.com/katuyo/symbol-exchange/routins"
)

func prepareStocks(){
    stock := models.Stock{Symbol: "WSCN", Open: 100}
    models.PushStock(stock)
}

func main() {
    m := macaron.Classic()

    prepareStocks()
    configRoutes(m)

    go func() {
        for {
            for k, _ := range models.GetStockMap() {
                routins.PrintDepth(k)
            }
            time.Sleep(1 * time.Second)
        }
    }()

    m.Use(macaron.Static("public"))
    /**  Since no html UI, just for Render JSON*/
    m.Use(renders.Renderer(
        renders.Options{
            Directory:  "views",                // Specify what path to load the templates from.
            Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
            //Funcs:           FuncMap,                    // Specify helper function maps for templates to access.
            Charset:         "UTF-8",     // Sets encoding for json and html content-types. Default is "UTF-8".
            IndentJSON:      true,        // Output human readable JSON
            IndentXML:       true,        // Output human readable XML
            HTMLContentType: "text/html", // Output XHTML content type instead of default "text/html"
        }))

    m.Run()
}