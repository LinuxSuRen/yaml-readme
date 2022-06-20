package function

import "fmt"

// LinkOrEmpty returns a Markdown style link or empty if the link is none
func LinkOrEmpty(text, link string) (output string) {
	if output = Link(text, link); output == text {
		output = ""
	}
	return
}

// Link returns a Markdown style link
func Link(text, link string) (output string) {
	output = text
	if link != "" {
		output = fmt.Sprintf("[%s](%s)", text, link)
	}
	return
}
