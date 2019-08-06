package main

import (
	// "bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	strange = "０１２３４５６７８９Ｊ"
)

var (
	table map[string][]string
)

func init() {
	table = make(map[string][]string)
}

func parse(input, output string) error {
	buf, err := ioutil.ReadFile(input)
	if err != nil {
		return fmt.Errorf("Open file err, %+v", err)
	}
	body := string(buf)

	blockNum, err := parseBlock(body)
	if err != nil {
		return fmt.Errorf("ParseBlock err, %+v", err)
	}
	lists := toLists(blockNum)
	writeBuf := toCsv(lists)
	err = ioutil.WriteFile(output, writeBuf, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("Write to file err, %+v", err)
	}
	// clear table
	table = make(map[string][]string)
	return nil
}

func parseBlock(body string) (int, error) {
	blockNum := 0
	lines := strings.Split(body, "\r\n")
	for _, line := range lines {
		line := strings.Trim(line, " ")
		if len(line) == 0 {
			continue
		}
		if isBegin(line) {
			blockNum++
		} else {
			items := strings.Split(line, "：")
			if len(items) != 2 {
				// empty line
				continue
			}
			key, value := items[0], strings.Trim(items[1], "。")
			for len(table[key]) < blockNum {
				table[key] = append(table[key], "")
			}
			if len(table[key][blockNum-1]) == 0 {
				table[key][blockNum-1] = value
			} else {
				table[key][blockNum-1] += (":" + value)
			}
		}
	}
	for key := range table {
		missNum := blockNum - len(table[key])
		for i := 0; i < missNum; i++ {
			table[key] = append(table[key], "")
		}
	}
	return blockNum, nil
}

func toLists(num int) [][]string {
	lists := make([][]string, num+1)
	for key := range table {
		lists[0] = append(lists[0], key)
	}
	for i := 0; i < num; i++ {
		lists[i+1] = make([]string, len(lists[0]))
		for j, key := range lists[0] {
			lists[i+1][j] = table[key][i]
		}
	}
	return lists
}

func toCsv(lists [][]string) []byte {
	var buf bytes.Buffer
	for _, list := range lists {
		var bufline bytes.Buffer
		emptyNum := 0
		for _, item := range list {
			if len(item) == 0 {
				emptyNum++
			}
			bufline.WriteString(item)
			bufline.WriteString(",")
		}
		bufline.WriteString("\r\n")
		if emptyNum < len(list) {
			buf.Write(bufline.Bytes())
		}
	}
	return buf.Bytes()
}

func isBegin(str string) bool {
	chars := []rune(str)
	if len(chars) < 2 {
		return false
	}
	if chars[0] >= '0' && chars[0] <= '9' && chars[1] == '.' {
		return true
	}
	if strings.ContainsRune(strange, chars[0]) && chars[1] == '．' {
		return true
	}
	return false
}

func main() {
	files := []string{
		"JTT1022-2016.txt",
		"JTT1055-2016.txt",
		"JTT1057-2016.txt",
		"JTT1075-2016.txt",
		"JTT7353-2009.txt",
		"JTT9792-2015.txt",
	}
	for _, input := range files {
		output := strings.Split(input, ".")[0] + ".csv"
		err := parse(input, output)
		if err != nil {
			fmt.Printf("Parse file err, %+v \n", err)
		}
	}
}
