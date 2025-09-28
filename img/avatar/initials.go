package avatar

import (
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"html"
	"strings"
)

var colors = []string{"#E53935", "#D81B60", "#8E24AA", "#5E35B1", "#3949AB", "#1E88E5", "#039BE5", "#00ACC1", "#00897B", "#43A047", "#F57F17", "#6D4C41", "#546E7A"}

func pickColor(name string) string {
	h := fnv.New32a()
	h.Write([]byte(name))
	return colors[int(h.Sum32())%len(colors)]
}

func SVGWithInitials(name string) string {
	words := strings.Fields(name)
	if len(words) == 0 {
		return ""
	}

	first := []rune(strings.ToUpper(words[0]))[0]
	var last rune
	if len(words) > 1 {
		last = []rune(strings.ToUpper(words[len(words)-1]))[0]
	}

	initials := fmt.Sprintf("%c", first)

	color := pickColor(name)

	if last != 0 {
		initials += fmt.Sprintf("%c", last)
	}

	initials = strings.ToUpper(initials)

	escapedText := html.EscapeString(initials)

	svg := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64">
  <rect width="100%%" height="100%%" fill="%s"/>
  <text x="32" y="36" dominant-baseline="middle" text-anchor="middle" font-size="32" fill="#ffffff" font-family="Roboto, sans-serif">%s</text>
</svg>`, color, escapedText)

	base64SVG := base64.StdEncoding.EncodeToString([]byte(svg))
	return "data:image/svg+xml;base64," + base64SVG
}
