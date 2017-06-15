package jsonobj

import "container/list"

type Directory struct {
	Id string
	Title string
	Url string
	Anchor string
	Pid string
	ChapterParagraphNum int
	ParagraphNum int
	Level int
	NewPage bool
	SubDirectory *list.List
}