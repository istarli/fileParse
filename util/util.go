package util

import (
	"bytes"
	"fmt"
)

// TableToCsv ...
func TableToCsv(num int, table map[string][]string) ([]byte, error) {
	if table == nil {
		return nil, fmt.Errorf("Table is nil")
	}
	return listToCsv(tableToList(num, table)), nil
}

func tableToList(num int, table map[string][]string) [][]string {
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

func listToCsv(lists [][]string) []byte {
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
