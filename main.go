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
	args.ToQueryArgs()
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
	//Type              string `short:"t" description:"The card type"`
	//Suptype           string `long:"st" description:":The card's subtype"`
	Name string `short:"n" description:"The card name" join:"+" split:" "`
	//ConvertedManaCost string `long:"cmc" description:"The card's converted mana cast"`
	//Color             string `short:"c" description:"The card's color"`
	//ColorIdentiy      string `short:"i" description:"The card's color identity"`
}

func (query QueryStruct) ToQueryArgs() string {
	ptr := reflect.ValueOf(&query)
	val := ptr.Elem()

	for i := 0; i < val.NumField(); i++ {
		key := val.Type().Field(i)
		val := val.Field(i)

		parts := strings.Split(val.String(), key.Tag.Get("split"))

		for _, c := range parts {
			fmt.Println(key.Tag.Get("join"), c)
		}
	}

	return ""
}
