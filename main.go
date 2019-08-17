package main

import (
	"fmt"
	"strings"

	"github.com/istarli/fileParse/parser"
)

func main() {
	files := []string{
		// "JTT1022-2016.txt",
		// "JTT1055-2016.txt",
		// "JTT1057-2016.txt",
		// "JTT1075-2016.txt",
		// "JTT7353-2009.txt",
		// "JTT9792-2015.txt",
		// "GBT 1948-2-2008.txt",
	}
	ps := parser.NewParser()
	// ps := parserv1.NewParserWithCheckFunc(func(str string) (bool, bool) {
	// 	if strings.Index(str, "中文名称：") == 0 {
	// 		return true, false
	// 	}
	// 	return false, false
	// })

	for _, input := range files {
		output := strings.Split(input, ".")[0] + ".csv"
		err := ps.Parse(input, output)
		if err != nil {
			fmt.Printf("Parse file err, %+v \n", err)
		}
	}
}
