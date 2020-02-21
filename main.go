package main

import (
	"context"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"image/color"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	startBtn, endBtn *widget.Button
	checkTimes       *canvas.Text
	foundText        *canvas.Text
	entry            *widget.Entry
	count            = 0
	cancelFunc       context.CancelFunc
)

func start() {
	if !strings.HasPrefix(entry.Text, "https://www.amazon.co.jp/") {
		foundText.Text = "url is not amazon japan!"
		refresh()
		return
	}
	var ctx context.Context
	ctx, cancelFunc = context.WithCancel(context.Background())
	reset()
	refresh()
	go run(ctx)
	startBtn.Disable()
	entry.Disable()
	endBtn.Enable()
}

func run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			count++
			s, has, err := crawl(entry.Text)
			if err != nil {
				foundText.Text = err.Error()
				foundText.Refresh()
				time.Sleep(5 * time.Second)
				continue
			}
			foundText.Text = s
			if has {
				startBtn.Enable()
				entry.Enable()
				endBtn.Disable()
				foundText.Refresh()
				refresh()
				return
			}
			refresh()
			time.Sleep(5 * time.Second)
		}
	}
}

func end() {
	cancelFunc()
	reset()
	refresh()
	startBtn.Enable()
	entry.Enable()
	endBtn.Disable()
}

func reset() {
	foundText.Text = "No found"
	count = 0
}

func crawl(url string) (string, bool, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", false, err
	}
	c, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", false, err
	}
	content := string(c)
	if strings.Contains(content, "buy-now-button") {
		return "You can buy!!", true, nil
	}
	return "Out of stock", false, nil
}

func refresh() {
	checkTimes.Text = fmt.Sprintf("Checking: %d", count)
	checkTimes.Refresh()
	foundText.Refresh()
}

func Render() *fyne.Container {
	checkTimes = canvas.NewText(fmt.Sprintf("Checking: %d", count), color.White)
	foundText = canvas.NewText("No found", color.RGBA{0, 255, 0, 0})
	foundText.TextSize = 30
	foundText.Alignment = fyne.TextAlignCenter
	checkTimes.TextSize = 30
	checkTimes.Alignment = fyne.TextAlignCenter
	startBtn = widget.NewButton("start", start)
	endBtn = widget.NewButton("end", end)
	endBtn.Disable()
	entry = widget.NewEntry()
	entry.MultiLine = false
	entry.PlaceHolder = "Url"
	return fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(1),
		entry,
		startBtn,
		endBtn,
		foundText,
		checkTimes,
	)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Amazon Checker")
	content := Render()
	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.SetFixedSize(true)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
