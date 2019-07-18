package main

import (
	"github.com/gocolly/colly"
	"fmt"
)

func main() {
	fmt.Println("ok")
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "chrome")
	})
}
