package function

import (
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_printContributor(t *testing.T) {
	type args struct {
		owner string
		repo  string
	}
	tests := []struct {
		name       string
		args       args
		prepare    func()
		wantOutput func() string
	}{{
		name: "normal case",
		args: args{
			owner: "linuxsuren",
			repo:  "yaml-readme",
		},
		prepare: func() {
			gock.New("https://api.github.com").
				Get("/repos/linuxsuren/yaml-readme/contributors").
				Reply(http.StatusOK).
				File("data/yaml-readme.json")
		},
		wantOutput: func() string {
			data, _ := ioutil.ReadFile("data/yaml-readme-contributors.txt")
			return string(data)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			tt.prepare()
			assert.Equalf(t, tt.wantOutput(), PrintContributors(tt.args.owner, tt.args.repo), "printContributors(%v, %v)", tt.args.owner, tt.args.repo)
		})
	}
}

func TestGithubUserLink(t *testing.T) {
	type args struct {
		id  string
		bio bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "normal case without bio",
		args: args{
			id:  "linuxsuren",
			bio: false,
		},
		want: `[Rick](https://github.com/LinuxSuRen)`,
	}, {
		name: "normal case with bio",
		args: args{
			id:  "linuxsuren",
			bio: true,
		},
		want: `[Rick](https://github.com/LinuxSuRen) (程序员，业余开源布道者)`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gock.New("https://api.github.com").
				Get("/users/linuxsuren").Reply(http.StatusOK).File("data/linuxsuren.json")

			assert.Equalf(t, tt.want, GithubUserLink(tt.args.id, tt.args.bio), "GithubUserLink(%v, %v)", tt.args.id, tt.args.bio)
		})
	}
}
