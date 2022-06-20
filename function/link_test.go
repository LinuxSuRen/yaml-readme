package function

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLink(t *testing.T) {
	type args struct {
		text string
		link string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
	}{{
		name: "link is not empty",
		args: args{
			text: "text",
			link: "link",
		},
		wantOutput: "[text](link)",
	}, {
		name: "link is empty",
		args: args{
			text: "text",
		},
		wantOutput: "text",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantOutput, Link(tt.args.text, tt.args.link), "Link(%v, %v)", tt.args.text, tt.args.link)
		})
	}
}

func TestLinkOrEmpty(t *testing.T) {
	type args struct {
		text string
		link string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
	}{{
		name: "link is empty",
		args: args{
			text: "text",
		},
	}, {
		name: "link is not empty",
		args: args{
			text: "text",
			link: "link",
		},
		wantOutput: "[text](link)",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantOutput, LinkOrEmpty(tt.args.text, tt.args.link), "LinkOrEmpty(%v, %v)", tt.args.text, tt.args.link)
		})
	}
}

func TestTwitterLink(t *testing.T) {
	type args struct {
		user string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
	}{{
		name: "user is empty",
	}, {
		name: "user is not empty",
		args: args{
			user: "linuxsuren",
		},
		wantOutput: "[![twitter](https://encrypted-tbn3.gstatic.com/favicon-tbn?q=tbn:ANd9GcTA3XDrUCnqJvmP3gfZKpXtV8ZO23EalnKszft6-V73d8G2Lt54v9TEnnkeO_MXseXmT5ERutOo0yPqoODJkFPtvxCeQbg_PYDJjXDAFfIMzM2p4bI)](https://twitter.com/linuxsuren)",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantOutput, TwitterLink(tt.args.user), "TwitterLink(%v)", tt.args.user)
		})
	}
}

func TestYouTubeLink(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
	}{{
		name: "empty",
	}, {
		name: "not empty",
		args: args{
			id: "channel/UC63xz3pq26BBgwB3cnwCoqQ",
		},
		wantOutput: "[![youtube](https://encrypted-tbn3.gstatic.com/favicon-tbn?q=tbn:ANd9GcRY4no9kYJtEAHXBEY2GDprV__HH1zc94olyS6G6fT5isS71bPyqvIi7-9VE1MMy3_3vsNOQLAerwcSQqGNyADWfxKpd2hLc8HuacZdgEjgZc_WLN8)](https://www.youtube.com/channel/UC63xz3pq26BBgwB3cnwCoqQ)",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantOutput, YouTubeLink(tt.args.id), "YouTubeLink(%v)", tt.args.id)
		})
	}
}
