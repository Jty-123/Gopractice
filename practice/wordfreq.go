// Findlinks1 prints the links in an HTML document read from standard input.
package main

import (
    "fmt"
    "strings"
    "os"
    "bufio"
    "io"
    "unicode"
)

func main() {
    //open file
    file,err := os.Open("./pg-being_ernest.txt")
    mp := make(map[string]int)
    if err != nil {
        fmt.Println("open file err", err)
        return 
    }
    defer file.Close()
    reader := bufio.NewReader(file)
    for {
        line,_,err := reader.ReadLine()
        if err == io.EOF {
            break
        }
        if err != nil {
            return 
        }
        words:= splitWords(string(line))
        for _,v := range words {
            mp[v]++
        }
    }
    for k,v := range mp {
        fmt.Println(k,v)
    }
}

func splitWords(text string) []string {
	// 将所有非字母的字符替换为空格
	replaceFunc := func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		return ' '
	}
    text = strings.Map(replaceFunc, text)
    res := strings.Fields(text)
    return res
}