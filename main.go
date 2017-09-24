package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/jessevdk/go-flags"
)

func main() {
	args := parse_args()
	url := "http://gatherer.wizards.com/Pages/Search/Default.aspx?action=advanced"

	full_result := fmt.Sprintf("%v%v", url, args)

	fmt.Println(full_result)
}

func parse_args() QueryStruct {
	var pargs QueryStruct

	_, err := flags.ParseArgs(&pargs, os.Args)

	if err != nil {
		panic(err)
	}

	return pargs
}

type Color int

const (
	White Color = iota
	Blue
	Black
	Red
	Green
)

type QueryStruct struct {
	Type              string `short:"t" description:"The card type" join:"+" split:" " query:"type"`
	Suptype           string `long:"st" description:":The card's subtype" join:"+" split:" " query:"subtype"`
	Name              string `short:"n" description:"The card name" join:"+" split:" " query:"name"`
	ConvertedManaCost string `long:"cmc" description:"The card's converted mana cast" join:"+=" query:"cmc"`
	Color             string `short:"c" description:"The card's color" join:"+" split:"" query:"color"`
	ColorIdentiy      string `short:"i" description:"The card's color identity" join:"+" split:"" query:"coloridentiy"`
	Rules             string `short:"r" description:"The card's rules test" join:"+" split:" " query:"text"`
}

func (query QueryStruct) String() string {
	ptr := reflect.ValueOf(&query)
	val := ptr.Elem()

	result := ""

	for i := 0; i < val.NumField(); i++ {
		key := val.Type().Field(i)
		val := val.Field(i)

		if val.String() == "" {
			continue
		}

		current := fmt.Sprintf("&%v=", key.Tag.Get("query"))

		parts := strings.Split(val.String(), key.Tag.Get("split"))

		join := key.Tag.Get("join")

		for _, c := range parts {
			current = fmt.Sprintf("%v%v[%v]", current, join, c)
		}

		result = fmt.Sprintf("%v%v", result, current)
	}

	return result
}
