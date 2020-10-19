# Simple go web app

这篇文章是官方教程[Writing Web Applications](https://golang.org/doc/articles/wiki/)的意译（说是意译，因为不是逐字逐句翻译的，且加入了我自己的理解）

[导言](#导言)

[开始](#开始)

[数据结构](#数据结构)

[net/http介绍](#net%2Fhttp介绍)

## 导言

本教程覆盖的内容：

- 创建一个有save和load方法的数据结构
- 使用`net/http`创建web app
- 使用`html/template`处理html模板
- 使用`regexp`验证用户输入
- 使用闭包

要求的你掌握的知识：

- 一点go基础
- 一点web技术，如http，html
- 一点unix-like命令行知识

## 开始

首先需要的你的系统中安装好go。开启`go mod`

```bash
mkdir simple_go_web_app
cd simlple_go_web_app
go mod init simlple_go_web_app
```

项目中会自动生成`go.mod`。然后创建`wiki.go`，开始写代码。首先是main包和import

```go
package main

import (
    "fmt"
    "io/ioutil"
)
```

## 数据结构

wiki是由一个个page组成的，每个page包括体格title和body。这里body的类型是`byte slice`而不是string，这是为了方便io库使用

```go
type Page struct {
	Title string
	Body []byte
}
```

然后实现`save`和`loadPage`两个方法（save是Page的方法，loadPage是函数，这篇教程里不再严格区分）。save方法把内存中的Page结构体分别持久化到磁盘，loadPage则反过来从磁盘读取文件，并返回Page指针。

```go
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
```

在main函数中测试一下这两个方法。

```go
func main() {
	// 创建一个测试的Page p1，然后调用save方法保存到磁盘
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page.")}
	p1.save()

	// 用loadPage读取刚才保存的page，然后打印内容
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}
```

先跑一下程序，可以看到打印的p2.Body。本地磁盘可以看到`TestPage.txt`

```bash
$ go run wiki.go
This is a sample page.
$ cat TestPage.txt
This is a sample page.
```

## net/http介绍

先来看一个简单的demo

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// 如果浏览器输入127.0.0.1:8888/go
	// 那么r.URL的值是 “/go”
	// 这里切片去掉了最前面的"/"
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)

	// ListenAndServe返回的话，必然是error
	log.Fatal(http.ListenAndServe(":8888", nil))
}
```

访问`127.0.0.1:8888/go`，可以看到返回的`Hi there, I love go!`

## 使用net/http展示wiki页面

回到`wiki.go`中。在import中加上`net/http`

```go
import (
	"fmt"
	"io/ioutil"
	"net/http"
)
```


