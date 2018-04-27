package html_parser

import (
	"html"
)

func EscapeWithEmojis(str string) string {
	parser := NewEmojiParser()
	return parser.ToHtmlEntities(html.EscapeString(str))
}
