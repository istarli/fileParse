package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	fileName       = "./input.txt"
	outputFileName = "./result.csv"
)

var (
	table map[string][]string
)

func init() {
	table = make(map[string][]string)
}

func parse(fileName string) error {
	buf, err := ioutil.ReadFile(fileName)
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
	err = ioutil.WriteFile(outputFileName, writeBuf, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("Write to file err, %+v", err)
	}
	return nil
}

func parseBlock(body string) (int, error) {
	blockNum := 0
	lines := strings.Split(body, "\r\n")
	for _, line := range lines {
		if isBegin(line) {
			blockNum++
		} else {
			items := strings.Split(line, "：")
			if len(items) != 2 {
				return 0, fmt.Errorf("Split line err, item num:%d", len(items))
			}
			if len(table[items[0]]) >= blockNum {
				return 0, fmt.Errorf("key[%s], length err, blockNum=%d but len=%d", items[0], blockNum, len(table[items[0]]))
			}
			for len(table[items[0]]) < blockNum {
				table[items[0]] = append(table[items[0]], "")
			}
			table[items[0]][blockNum-1] = items[1]
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
		for _, item := range list {
			buf.WriteString(item)
			buf.WriteString(",")
		}
		buf.WriteString("\r\n")
	}
	return buf.Bytes()
}

func isBegin(str string) bool {
	return strings.Contains(str, "５")
}

func main() {
	err := parse(fileName)
	if err != nil {
		fmt.Println("Parse file err, %+v", err)
	}
}
