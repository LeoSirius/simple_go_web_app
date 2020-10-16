package main

import (
	"fmt"
	"io/ioutil"
)

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"

	// 0 0o 0O 开头的字面量都是八进制，600表示自己有读写权限
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title:title, Body: body}, nil
}

func main() {
	// 创建一个测试的Page p1，然后调用save方法保存到磁盘
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page.")}
	p1.save()

	// 用loadPage读取刚才保存的page，然后打印内容
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}
