package main

import (
	"bufio"
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
				fmt.Printf("Warning: Split line err, item num: %d", len(items))
				continue
			}
			key, value := items[0], strings.Trim(items[1], "。")
			if len(table[key]) >= blockNum {
				return 0, fmt.Errorf("key[%s], length err, blockNum=%d but len=%d", key, blockNum, len(table[key]))
			}
			for len(table[key]) < blockNum {
				table[key] = append(table[key], "")
			}
			table[key][blockNum-1] = value
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
	return strings.IndexRune(str, '5') == 0
}

func main() {
	fmt.Println("输入文件名:")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	inputFileName := input.Text()
	outputFileName := strings.Split(inputFileName, ".")[0] + ".csv"

	err := parse(inputFileName, outputFileName)
	if err != nil {
		fmt.Println("Parse file err, %+v", err)
	}
}
