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
