package jsonobj

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
	SubDirectory []*Directory
}

type Chapter0 struct {
	ImageUrl string
	BookTitle string
	Directories []*Directory
}
