package function

import (
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
)

func chartData() map[string][]map[string]interface{} {
	return map[string][]map[string]interface{}{
		"Locomotion": {
			{"name": "a", "creator": "深圳市越疆科技股份有限公司", "year": 2021},
			{"name": "b", "creator": "トヨタ自動車株式会社", "year": 2021},
			{"name": "c", "creator": "Tesla, Inc.", "year": 2024},
		},
		"VLA Model": {
			{"name": "d", "creator": "Guangzhou Moju Technology Co Ltd", "year": 2025},
		},
	}
}

func TestPrintYearChart(t *testing.T) {
	out := string(PrintYearChart(chartData(), "Patents by year"))
	assert.Contains(t, out, "```mermaid\nxychart-beta\n")
	assert.Contains(t, out, `title "Patents by year"`)
	// gaps filled with zero: 2021..2025 -> 5 points
	assert.Contains(t, out, "x-axis [2021, 2022, 2023, 2024, 2025]")
	assert.Contains(t, out, "line [2, 0, 0, 1, 1]")
	assert.Empty(t, PrintYearChart("bad", "x"))
}

func TestPrintRegionChart(t *testing.T) {
	out := string(PrintRegionChart(chartData(), "Domestic vs International"))
	assert.Contains(t, out, "```mermaid\npie showData\n")
	assert.Contains(t, out, `"国内公司" : 2`)
	assert.Contains(t, out, `"国际公司" : 2`) // Toyota (kana) and Tesla
	assert.Empty(t, PrintRegionChart(nil, "x"))
}

func TestPrintChartsReturnTemplateHTML(t *testing.T) {
	assert.IsType(t, template.HTML(""), PrintYearChart(chartData(), "t"))
	assert.IsType(t, template.HTML(""), PrintRegionChart(chartData(), "t"))
}

func TestIsChineseAssignee(t *testing.T) {
	assert.True(t, isChineseAssignee("深圳市越疆科技股份有限公司"))
	assert.True(t, isChineseAssignee("Guangzhou Moju Technology Co Ltd"))
	assert.False(t, isChineseAssignee("トヨタ自動車株式会社"))
	assert.False(t, isChineseAssignee("(주)뉴빌리티"))
	assert.False(t, isChineseAssignee("Boston Dynamics, Inc."))
}
