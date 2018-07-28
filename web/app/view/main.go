// Template generated by reactGen

package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"thirdopinion/internal/pkg/config"

	"myitcv.io/react"

	"honnef.co/go/js/dom"
)

//go:generate reactGen

var document = dom.GetWindow().Document()

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	argumentID := generateFilter()
	ch, err := fetchPost(argumentID)
	if err != nil {
		renderError(&wg)
	}
	go readChan(ch, &wg)
	wg.Wait()
}

func generateFilter() (argumentID string) {
	baseURI := document.BaseURI()
	search := strings.Split(baseURI, "/view/")
	if len(search) == 2 {
		argumentID = search[1]
	}
	return
}

func readChan(ch chan *messageEvent, wg *sync.WaitGroup) {
	domTarget := document.GetElementByID("app")

	args := []*config.Argument{}
	defer wg.Done()
	select {
	case msg := <-ch:
		err := json.Unmarshal([]byte(msg.Data.String()), &args)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	lp := LayoutProps{
		arguments: args,
	}
	react.Render(Layout(lp), domTarget)
}

func renderError(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Broken. Fix this Alex!")
}
