package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

func main() {

	stdin := bufio.NewScanner(os.Stdin)
	fmt.Print("Input URL: ")
	stdin.Scan()

	text := stdin.Text()
	u, err := url.ParseRequestURI(strings.TrimSpace(text))
	if err != nil {
		panic("url is invalid")
	}
	doc, err := goquery.NewDocument(u.String())
	if err != nil {
		panic(err)
	}

	fmt.Print("■ ", doc.Url)
	fmt.Printf(":(%s)\n", doc.Find("title").Text())
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		u, _ := s.Attr("href")
		if err != nil {
			return
		}
		targetAttr, _ := s.Attr("target")
		if targetAttr == "_blank" {
			doc2, err := goquery.NewDocument(u)
			if err != nil {
				fmt.Println(err)
				return
			}

			doc2.Find("p").Each(func(_ int, s2 *goquery.Selection) {
				pText := sjis2utf8(s2.Text())
				if pText == "広告が見つかりません。" {
					fmt.Printf("- [ %s ]が無効です\n", u)
				}
			})
		}
	})
}

func sjis2utf8(str string) string {
	reader := strings.NewReader(str)
	transformReader := transform.NewReader(reader, japanese.ShiftJIS.NewDecoder())
	rtn, _ := ioutil.ReadAll(transformReader)
	return string(rtn)
}
