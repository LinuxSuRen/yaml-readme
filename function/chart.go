package function

import (
	"fmt"
	"html/template"
	"regexp"
	"sort"
	"strings"
)

var (
	kanaRe   = regexp.MustCompile(`[\x{3040}-\x{30ff}]`)
	hangulRe = regexp.MustCompile(`[\x{ac00}-\x{d7af}]`)
	hanRe    = regexp.MustCompile(`[\x{4e00}-\x{9fff}]`)
	// Latin transliterations of Chinese companies that do not contain Han characters
	cnLatinRe = regexp.MustCompile(`(?i)(guangzhou|beijing|shanghai|shenzhen|hangzhou|suzhou|shenyang|china)`)
)

// PrintYearChart renders the items grouped by year as a Mermaid xychart-beta
// line chart, which GitHub renders natively.
//
// data - the grouped data map when the '--group-by' flag is used; all items
// are flattened inside, so the grouping field does not matter here.
func PrintYearChart(data interface{}, title string) template.HTML {
	items := flattenGrouped(data)
	counts := map[int]int{}
	for _, item := range items {
		if y, ok := item["year"].(int); ok && y > 0 {
			counts[y]++
		}
	}
	if len(counts) == 0 {
		return ""
	}
	years := make([]int, 0, len(counts))
	for y := range counts {
		years = append(years, y)
	}
	sort.Ints(years)

	// fill gaps so the line reflects the timeline instead of skipping years
	first, last := years[0], years[len(years)-1]
	var xLabels []string
	var values []string
	for y := first; y <= last; y++ {
		xLabels = append(xLabels, fmt.Sprintf("%d", y))
		values = append(values, fmt.Sprintf("%d", counts[y]))
	}

	var buf strings.Builder
	buf.WriteString("```mermaid\nxychart-beta\n")
	buf.WriteString("    title \"" + sanitizeMindmapText(title) + "\"\n")
	buf.WriteString("    x-axis [" + strings.Join(xLabels, ", ") + "]\n")
	buf.WriteString("    y-axis \"Patents\" 0 --> " + fmt.Sprintf("%d", maxCount(counts)+1) + "\n")
	buf.WriteString("    line [" + strings.Join(values, ", ") + "]\n")
	buf.WriteString("```\n")
	return template.HTML(buf.String())
}

// PrintRegionChart renders a Mermaid pie chart splitting the items into
// domestic (China) and international assignees. The assignee is classified
// by its characters: Han without kana/hangul reads as Chinese; a small set
// of Latin-transliterated Chinese city prefixes is also recognized.
func PrintRegionChart(data interface{}, title string) template.HTML {
	items := flattenGrouped(data)
	domestic, international := 0, 0
	for _, item := range items {
		creator, _ := item["creator"].(string)
		if isChineseAssignee(creator) {
			domestic++
		} else {
			international++
		}
	}
	if domestic+international == 0 {
		return ""
	}
	var buf strings.Builder
	buf.WriteString("```mermaid\npie showData\n")
	buf.WriteString("    title \"" + sanitizeMindmapText(title) + "\"\n")
	fmt.Fprintf(&buf, "    \"国内公司\" : %d\n", domestic)
	fmt.Fprintf(&buf, "    \"国际公司\" : %d\n", international)
	buf.WriteString("```\n")
	return template.HTML(buf.String())
}

func isChineseAssignee(creator string) bool {
	if cnLatinRe.MatchString(creator) {
		return true
	}
	return hanRe.MatchString(creator) && !kanaRe.MatchString(creator) && !hangulRe.MatchString(creator)
}

func flattenGrouped(data interface{}) []map[string]interface{} {
	grouped, ok := data.(map[string][]map[string]interface{})
	if !ok {
		return nil
	}
	var items []map[string]interface{}
	for _, group := range grouped {
		items = append(items, group...)
	}
	return items
}

func maxCount(counts map[int]int) int {
	m := 0
	for _, c := range counts {
		if c > m {
			m = c
		}
	}
	return m
}
