package function

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintMindmap(t *testing.T) {
	data := map[string][]map[string]interface{}{
		"Locomotion": {
			{"name": "Robot balance control (CN103634206A)", "creator": "UBTech"},
			{"name": "Gait planning", "creator": ""},
		},
		"VLA Model": {
			{"name": "Vision-language-action policy", "creator": "Galbot", "year": 2024},
		},
	}

	out := PrintMindmap(data, "Awesome Patent", "name", "creator")
	expect := "```mermaid\n" +
		"mindmap\n" +
		"  root((Awesome Patent))\n" +
		"  Locomotion\n" +
		"    Robot balance control （CN103634206A） · UBTech\n" +
		"    Gait planning\n" +
		"  VLA Model\n" +
		"    Vision-language-action policy · Galbot\n" +
		"```\n"
	assert.Equal(t, expect, out)
}

func TestPrintMindmapEmpty(t *testing.T) {
	out := PrintMindmap(map[string][]map[string]interface{}{}, "Root", "name", "creator")
	assert.Contains(t, out, "root((Root))")
	assert.NotContains(t, out, "·")
	assert.Empty(t, PrintMindmap("not grouped", "Root", "name", ""))
}

func TestSanitizeMindmapText(t *testing.T) {
	assert.Equal(t, "a（b）", sanitizeMindmapText("a(b)"))
	assert.Equal(t, "x【y】", sanitizeMindmapText("x[y]"))
	assert.Equal(t, "a b", sanitizeMindmapText("a\tb"))
	assert.Equal(t, "", sanitizeMindmapText("  \n "))
}
