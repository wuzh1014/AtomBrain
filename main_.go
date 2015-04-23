package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strconv"
	"time"
)

var outText, cUnix *walk.TextEdit
var inText, cClock *walk.LineEdit
var cDate *walk.DateEdit

const shortForm = "15:04:05"

func trans2date() {
	ctime, _ := strconv.ParseInt(inText.Text(), 10, 64)
	str := ""
	cdate := time.Unix(ctime, 0)
	str += cdate.String()
	outText.SetText(str)
}

func trans2unix() {
	tDate := cDate.Date().Unix()
	tClock, _ := time.Parse(shortForm, cClock.Text())
	tUnix := tClock.Hour()*3600 + tClock.Minute()*60 + tClock.Second()
	str := strconv.FormatFloat(float64(tUnix)+float64(tDate), 'f', 0, 64)
	cUnix.SetText(str)
}

func main_() {
	var mw *walk.MainWindow

	MainWindow{
		AssignTo: &mw,
		Title:    "Process",
		MinSize:  Size{400, 300},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{

					PushButton{
						ColumnSpan: 2,
						Text:       "->",
						OnClicked: func() {
							trans2date()
						},
					},

					HSplitter{
						ColumnSpan: 2,
						Children: []Widget{

							LineEdit{
								AssignTo: &inText, Text: "input Millisecond",
							},
							TextEdit{AssignTo: &outText, Text: "input Date"},
						},
					},

					DateEdit{
						AssignTo: &cDate,
						OnDateChanged: func() {
						},
					},

					PushButton{
						Text: "->",
						OnClicked: func() {
							trans2unix()
						},
					},

					LineEdit{
						AssignTo: &cClock,
						Text:     "00:00:00",
					},

					TextEdit{
						AssignTo: &cUnix,
						Text:     "0000000000",
					},
				},
			},
		},
	}.Run()
}
