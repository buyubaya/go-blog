package main


import (
	"github.com/buyubaya/go-blog/App"
	"github.com/buyubaya/go-blog/config"
)


func main() {
	// CONFIG
	config := config.GetConfig()

	// APP
	app := &App.App{}
	app.Initialize(config)
	app.Run(":3000")
}
