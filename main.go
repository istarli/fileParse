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
		// "GBT 26768-2011.txt",
		"GBT 29110-2012.txt",
	}
	// ps := parser.NewParser()

	// ps := parser.NewParserWithCheckFunc(func(str string) (bool, bool) {
	// 	if strings.Index(str, "中文名称：") == 0 {
	// 		return true, false
	// 	}
	// 	return false, false
	// })

	ps := parser.NewParserWithCheckFunc(func(str *string) (bool, bool) {
		s := *str
		if strings.Index(s, "7.") == 0 {
			items := strings.SplitN(s, " ", 2)
			if len(items) < 2 {
				fmt.Println("begin parse error!")
				return false, false
			}
			newStr := "中文名称：" + strings.TrimSpace(items[1])
			*str = newStr
			return true, false
		}
		return false, false
	})

	for _, input := range files {
		output := strings.Split(input, ".")[0] + ".csv"
		err := ps.Parse(input, output)
		if err != nil {
			fmt.Printf("Parse file err, %+v \n", err)
		}
	}
}
