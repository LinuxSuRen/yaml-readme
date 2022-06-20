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

// TwitterLink returns a Markdown style link of Twitter
func TwitterLink(user string) (output string) {
	if user != "" {
		output = fmt.Sprintf("[![twitter](https://encrypted-tbn3.gstatic.com/favicon-tbn?q=tbn:ANd9GcTA3XDrUCnqJvmP3gfZKpXtV8ZO23EalnKszft6-V73d8G2Lt54v9TEnnkeO_MXseXmT5ERutOo0yPqoODJkFPtvxCeQbg_PYDJjXDAFfIMzM2p4bI)](https://twitter.com/%s)", user)
	}
	return
}
