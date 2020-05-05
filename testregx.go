package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	sql := "select * from users where created > :uuid and created < :uuid2 limit 10;"
	re := regexp.MustCompile(`(<|and|>|>=|<=|like|between|^) :\w+ (,|and|or|limit|\))`)

	sql = re.ReplaceAllString(sql, "$1 ? $2")

	re = regexp.MustCompile(`(<|and|>|>=|<=|like|between|^) \? (,|and|or|limit|\))`)
	i := 0
	sql = re.ReplaceAllStringFunc(sql, func(s string) string {
		i++
		qmark := fmt.Sprintf("$%d", i)
		s = strings.Replace(s, "?", qmark, -1)
		return s
	})

	fmt.Println(sql)
}
