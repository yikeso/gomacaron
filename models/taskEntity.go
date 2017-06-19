package models

import "bytes"

type ChapterEntity struct {
	Content *bytes.Buffer
	HtmlContent *bytes.Buffer
	Paragraph int
}