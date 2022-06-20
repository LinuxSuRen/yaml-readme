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
		output = fmt.Sprintf("[![twitter](%s)](https://twitter.com/%s)", GStatic("twitter"), user)
	}
	return
}

// YouTubeLink returns a Markdown style link of YouTube
func YouTubeLink(id string) (output string) {
	if id != "" {
		output = fmt.Sprintf("[![youtube](%s)](https://www.youtube.com/%s)", GStatic("youtube"), id)
	}
	return
}

// GStatic returns the known gstatic image URL
func GStatic(id string) (output string) {
	var tbn string
	switch id {
	case "youtube":
		tbn = "ANd9GcRY4no9kYJtEAHXBEY2GDprV__HH1zc94olyS6G6fT5isS71bPyqvIi7-9VE1MMy3_3vsNOQLAerwcSQqGNyADWfxKpd2hLc8HuacZdgEjgZc_WLN8"
	case "twitter":
		tbn = "ANd9GcTA3XDrUCnqJvmP3gfZKpXtV8ZO23EalnKszft6-V73d8G2Lt54v9TEnnkeO_MXseXmT5ERutOo0yPqoODJkFPtvxCeQbg_PYDJjXDAFfIMzM2p4bI"
	}
	output = fmt.Sprintf("https://encrypted-tbn3.gstatic.com/favicon-tbn?q=tbn:%s", tbn)
	return
}
