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
	"time"
)

var (
	startBtn, endBtn *widget.Button
	checkTimes       *canvas.Text
	foundText        *canvas.Text
	entry            *widget.Entry
	count            = 0
	found            = 0
	cancelFunc       context.CancelFunc
)

func start() {
	var ctx context.Context
	ctx, cancelFunc = context.WithCancel(context.Background())
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
	count = 0
}

func refresh() {
	checkTimes.Text = fmt.Sprintf("Checking: %d", count)
	foundText.Text = fmt.Sprintf("Found: %d", found)
	checkTimes.Refresh()
	foundText.Refresh()
}

func Render() *fyne.Container {
	checkTimes = canvas.NewText(fmt.Sprintf("Checking: %d", count), color.White)
	foundText = canvas.NewText(fmt.Sprintf("Found: %d", found), color.RGBA{0, 255, 0, 0})
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
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
