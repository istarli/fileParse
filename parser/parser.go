package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/istarli/fileParse/util"
)

const (
	strange = "０１２３４５６７８９Ｊ"
	path    = "./file/"
)

type Parser struct {
	CheckBegin func(string) (bool, bool)
}

func NewParser() *Parser {
	p := &Parser{}
	p.CheckBegin = p.defaultCheckBegin
	return p
}

func NewParserWithCheckFunc(fn func(string) (bool, bool)) *Parser {
	return &Parser{
		CheckBegin: fn,
	}
}

// Parse ...
func (p *Parser) Parse(input, output string) error {
	buf, err := ioutil.ReadFile(path + input)
	if err != nil {
		return fmt.Errorf("Open file err, %+v", err)
	}
	body := string(buf)
	table := make(map[string][]string)
	blockNum, err := p.parseBlock(body, table)
	if err != nil {
		return fmt.Errorf("ParseBlock err, %+v", err)
	}
	writeBuf, err := util.TableToCsv(blockNum, table)
	if err != nil {
		return fmt.Errorf("Transe to csv err, %+v", err)
	}
	err = ioutil.WriteFile(path+output, writeBuf, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("Write to file err, %+v", err)
	}
	return nil
}

func (p *Parser) parseBlock(body string, table map[string][]string) (int, error) {
	blockNum := 0
	lines := strings.Split(body, "\r\n")
	for _, line := range lines {
		line := strings.Trim(line, " ")
		if p.CheckBegin == nil {
			p.CheckBegin = p.defaultCheckBegin
		}
		isBegin, isJump := p.CheckBegin(line)
		if isBegin {
			blockNum++
		}
		if isJump {
			continue
		}
		items := strings.SplitN(line, "：", 2)
		if len(items) < 2 {
			fmt.Printf("invalid line:%s\n", line)
			continue
		}
		key, value := items[0], strings.Trim(items[1], "。")
		for len(table[key]) < blockNum {
			table[key] = append(table[key], "")
		}
		if blockNum < 1 {
			return 0, fmt.Errorf("blockNum error")
		}
		if len(table[key][blockNum-1]) == 0 {
			table[key][blockNum-1] = value
		} else {
			table[key][blockNum-1] += (":" + value)
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

func (p *Parser) defaultCheckBegin(str string) (bool, bool) {
	chars := []rune(str)
	if len(chars) < 2 {
		return false, false
	}
	if (chars[0] >= '0' && chars[0] <= '9' && chars[1] == '.') ||
		(chars[0] >= '０' && chars[0] <= '９' && chars[1] == '．') {
		return true, true
	}
	return false, false
}
