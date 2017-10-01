package main

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"

	"github.com/jessevdk/go-flags"
)

func main() {
	args := parseArgs()
	url := "http://gatherer.wizards.com/Pages/Search/Default.aspx?action=advanced"

	full_result := fmt.Sprintf("%v%v", url, args)

	fmt.Println(full_result)

	if args.InBrowser {
		openUrl(full_result)
	}
}

func parseArgs() QueryStruct {
	var pargs QueryStruct

	_, err := flags.ParseArgs(&pargs, os.Args)

	if err != nil {
		panic(err)
	}

	return pargs
}

func openUrl(url string) {
	var args []string

	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windwos":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"sensible-browser"}
	}

	cmd := exec.Command(args[0], append(args[1:], url)...)
	cmd.Start()
}

type Color int

const (
	White Color = iota
	Blue
	Black
	Red
	Green
)

func (c Color) String() string {
	switch c {
	case White:
		return "W"
	case Blue:
		return "U"
	case Black:
		return "B"
	case Red:
		return "R"
	case Green:
		return "G"
	}
	return ""
}

func ToColor(s string) Color {
	switch lower := strings.ToLower(s); lower {
	case "w":
		return White
	case "u":
		return Blue
	case "b":
		return Black
	case "r":
		return Red
	case "g":
		return Green
	default:
		panic("Invalid color!") /// TODO ?
	}

}

func getJoin(query string, default_join string) (string, string) {
	ret := default_join
	q := query
	switch c := query[0]; c {
	case '!':
		ret = "+!"
		q = query[1:]
		break
	case '+':
		ret = "+"
		q = query[1:]
		break
	case '|':
		ret = "|"
		q = query[1:]
		break
	}
	return q, ret
}

func getComparison(query string, default_comparison string) (string, string) {
	ret := default_comparison
	q := query
	switch c := query[0]; c {
	case '=':
		fallthrough
	case '<':
		fallthrough
	case '>':
		ret = query[:1]
		q = query[1:]
	}

	return q, ret
}

type QueryStruct struct {
	Type              string `short:"t" description:"The card type" join:"+" split:" " query:"type"`
	Suptype           string `long:"st" description:":The card's subtype" join:"+" split:" " query:"subtype"`
	Name              string `short:"n" description:"The card name" join:"+" split:" " query:"name"`
	ConvertedManaCost string `long:"cmc" description:"The card's converted mana cast" join:"+=" query:"cmc" comparison:"="`
	Color             string `short:"c" description:"The card's color" join:"+" split:"" query:"color" converter:"color"`
	ColorIdentiy      string `short:"i" description:"The card's color identity" join:"+" split:"" query:"coloridentiy" converter:"color"`
	Rules             string `short:"r" description:"The card's rules test" join:"+" split:" " query:"text"`
	InBrowser         bool   `short:"b" description:"Open the result in the default web browser" skip:"yes"`
	Power             string `long:"pow" description:"The creature's power" query:"power" comparison:"="`
	Toughness         string `long:"tough" description:"The creature's toughness" query:"power" comparison:"="`
}

func (query QueryStruct) String() string {
	ptr := reflect.ValueOf(&query)
	val := ptr.Elem()

	result := ""

	/// loop over the arguments
	for i := 0; i < val.NumField(); i++ {
		key := val.Type().Field(i)
		val := val.Field(i)

		/// if it wasn't present, move on
		if val.String() == "" {
			continue
		}

		/// if it's supposed to be skipped, move on
		if key.Tag.Get("skip") != "" {
			continue
		}

		/// set up the query string
		current := fmt.Sprintf("&%v=", key.Tag.Get("query"))

		/// if we're looking at a string field
		if key.Tag.Get("split") != "" {
			/// perform the split and get the join
			parts := strings.Split(val.String(), key.Tag.Get("split"))

			join := key.Tag.Get("join")

			parts[0], join = getJoin(parts[0], join)

			/// clear out the join if it was given
			if parts[0] == "" {
				parts = parts[1:]
			}

			/// join the parts and format
			for _, c := range parts {
				current = fmt.Sprintf("%v%v[%v]", current, join, c)
			}
		} else {
			/// this is a numeric
			number, comparison := getComparison(val.String(), key.Tag.Get("comparison"))

			current = fmt.Sprintf("%v+%v[%v]", current, comparison, number)
		}

		/// append it to the result
		result = fmt.Sprintf("%v%v", result, current)
	}

	return result
}
