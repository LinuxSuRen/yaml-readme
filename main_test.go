package main

import "testing"
import "github.com/stretchr/testify/assert"

func Test_sortBy(t *testing.T) {
	type args struct {
		items  []map[string]interface{}
		sortBy string
	}
	tests := []struct {
		name   string
		args   args
		verify func([]map[string]interface{}, *testing.T)
	}{{
		name: "normal",
		args: args{
			items: []map[string]interface{}{{
				"name": "b",
			}, {
				"name": "c",
			}, {
				"name": "a",
			}},
			sortBy: "name",
		},
		verify: func(data []map[string]interface{}, t *testing.T) {
			assert.Equal(t, map[string]interface{}{
				"name": "a",
			}, data[0])
		},
	}, {
		name: "number values",
		args: args{
			items: []map[string]interface{}{{
				"name": "12",
			}, {
				"name": "13",
			}, {
				"name": "11",
			}, {
				"name": "1",
			}},
			sortBy: "name",
		},
		verify: func(data []map[string]interface{}, t *testing.T) {
			assert.Equal(t, map[string]interface{}{
				"name": "1",
			}, data[0])
		},
	}, {
		name: "slice values",
		args: args{
			items: []map[string]interface{}{{
				"name": []string{},
			}, {
				"name": []string{},
			}},
			sortBy: "name",
		},
		verify: func(data []map[string]interface{}, t *testing.T) {
			assert.Equal(t, map[string]interface{}{
				"name": []string{},
			}, data[0])
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortBy(tt.args.items, tt.args.sortBy, true)
			tt.verify(tt.args.items, t)
		})
	}
}
