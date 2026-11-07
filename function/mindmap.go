package function

import (
	"fmt"
	"html/template"
	"sort"
	"strings"
)

// PrintMindmap renders grouped items into a Mermaid mindmap code block which
// can be rendered natively by GitHub.
//
// data      - the grouped data map when the '--group-by' flag is used, i.e. map[group][]item
// root      - the title of the root node
// nameField - the item field rendered as the child leaf, e.g. "name"
// subField  - an optional second item field appended to the leaf, e.g. "creator".
// Pass an empty string to skip it.
//
// The returned value is template.HTML so that group names containing characters
// like '&' are not HTML-escaped by the template engine.
func PrintMindmap(data interface{}, root, nameField, subField string) template.HTML {
	grouped, ok := data.(map[string][]map[string]interface{})
	if !ok {
		return ""
	}

	groups := make([]string, 0, len(grouped))
	for key := range grouped {
		groups = append(groups, key)
	}
	sort.Strings(groups)

	var buf strings.Builder
	buf.WriteString("```mermaid\nmindmap\n")
	buf.WriteString("  root((" + sanitizeMindmapText(root) + "))\n")

	for _, group := range groups {
		buf.WriteString("    " + sanitizeMindmapText(group) + "\n")
		for _, item := range grouped[group] {
			leaf := mindmapFieldString(item, nameField)
			if subField != "" {
				if sub := mindmapFieldString(item, subField); sub != "" {
					leaf += " · " + sub
				}
			}
			if leaf != "" {
				buf.WriteString("      " + leaf + "\n")
			}
		}
	}
	buf.WriteString("```\n")
	return template.HTML(buf.String())
}

func mindmapFieldString(item map[string]interface{}, field string) string {
	val, ok := item[field]
	if !ok || val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return sanitizeMindmapText(v)
	case int:
		return sanitizeMindmapText(fmt.Sprintf("%d", v))
	}
	return ""
}

// sanitizeMindmapText replaces characters that Mermaid treats as node shape
// delimiters or that break the indentation-based mindmap syntax.
func sanitizeMindmapText(text string) string {
	replacer := strings.NewReplacer(
		"(", "（",
		")", "）",
		"[", "【",
		"]", "】",
		"{", "｛",
		"}", "｝",
		"\n", " ",
		"\r", " ",
		"\t", " ",
	)
	return strings.TrimSpace(replacer.Replace(text))
}
