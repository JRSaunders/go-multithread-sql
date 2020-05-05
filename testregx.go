package main

import (
	"fmt"
	"regexp"
)

func main() {
	sql := "select * from users where created > ? and created < ? limit 10;"
	var re = regexp.MustCompile(` \? `)
	i := 0
	sql = re.ReplaceAllStringFunc(sql, func(s string) string {
		i++

		string := fmt.Sprintf(" $%d ", i)
		return string
	})

	fmt.Println(sql)
}
