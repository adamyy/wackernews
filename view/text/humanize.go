package text

import (
	"strings"

	"github.com/fatih/color"
	"golang.org/x/net/html"
)

var (
	Italic       = color.Underline // color.Italic not widely supported, use alternative
	Bold         = color.Bold
	Faint        = color.Faint // color.Faint not widely supported, use alternative
	Underline    = color.Underline
	BlinkSlow    = color.BlinkSlow    // tested, does not work properly
	BlinkRapid   = color.BlinkRapid   // tested, does not work properly
	ReverseVideo = color.ReverseVideo // swap foreground and background color
	Concealed    = color.Concealed    // tested, does not work properly
	CrossedOut   = color.CrossedOut   // color.CrossedOut not widely supported, use alternative
)

// this function assumes that the input string is well-formed
// i.e., the start tags and end tags are balanced
// therefore, it interprets inputs such as "<i>italic<b>bold</i></b>" as "<i>italic<b>bold</b></i>"
func Humanize(htmlString string) string {
	tz := html.NewTokenizer(strings.NewReader(htmlString))
	var builder strings.Builder
	var attrs []color.Attribute

	for {
		tz.Next()
		token := tz.Token()

		switch token.Type {
		case html.ErrorToken:
			return builder.String()
		case html.StartTagToken:
			switch token.Data {
			case "i":
				attrs = append(attrs, Italic)
			case "b":
				attrs = append(attrs, Bold)
			case "u":
				attrs = append(attrs, Underline)
			case "strike":
				attrs = append(attrs, CrossedOut)
			}
		case html.TextToken:
			if len(attrs) > 0 {
				c := color.New(attrs...)
				str := c.Sprintf(token.Data)
				builder.WriteString(str)
			} else {
				str := token.Data
				builder.WriteString(str)
			}
		case html.EndTagToken:
			switch token.Data {
			case "i", "b", "u", "strike":
				attrs = attrs[:len(attrs)-1]
			}
		}
	}
}
