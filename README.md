# Simple go web app

这篇文章是官方教程[Writing Web Applications](https://golang.org/doc/articles/wiki/)的意译（说是意译，因为不是逐字逐句翻译的，且加入了我自己的理解）

[导言](#导言)

[开始](#开始)

[数据结构](#数据结构)

[`net/http`介绍](#nethttp介绍)

[使用`net/http`展示wiki页面](#使用nethttp展示wiki页面)

[编辑Page](#编辑Page)

[`html/template`包](#htmltemplate包)

[处理不存在的页面](#处理不存在的页面)

[保存页面](#保存页面)

[错误处理](#错误处理)

[模板缓存](#模板缓存)

[验证](#验证)

[函数字面值和闭包](#函数字面值和闭包)

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

## `net/http`介绍

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

url设计：

- `/view/`展示wikis
- `/edit/`编辑wiki
- `/save/`保存编辑好的wiki

我们会为这三个路径各自添加一个handler。这里先写`viewHandler`

回到`wiki.go`中。在import中加上`net/http`

```go
import (
	"fmt"
	"io/ioutil"
	"net/http"
)
```

然后创建一个视图handler，处理url为`/view/`的情况。

```go
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, _ := loadPage(title)
    fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}
```

在main中注册这个handler

```go
func main() {
    http.HandleFunc("/view/", viewHandler)
    log.Fatal(http.ListenAndServe(":8888", nil))
}
```

然后在浏览器访问`http://39.106.229.224:8888/view/TestPage`。就可以看到我们在第一步中添加的测试Page。

## 编辑Page

为edit添加handler

```go
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    fmt.Fprintf(w, "<h1>Editing %s</h1>"+
        "<form action=\"/save/%s\" method=\"POST\">"+
        "<textarea name=\"body\">%s</textarea><br>"+
        "<input type=\"submit\" value=\"Save\">"+
        "</form>",
        p.Title, p.Title, p.Body)
}
```

然后在main中注册`http.HandleFunc("/edit/", editHandler)`

## `html/template`包

前面我们的html是直接用字符串写在go代码中的，这样修改模板后需要重新编译一遍go代码。使用`html/template`可以把html文件和go源码分开。

首先import中加入`"html/template"`。然后创建一个`edit.html`，把之前写在go源码中的内容写到html里

```html
<h1>Editing {{.Title}}</h1>

<form action="/save/{{.Title}}" method="POST">

<!-- Body是[]byte类型，这里用printf "%s"将其转换成字符串 -->
<div><textarea name="body" rows="20" cols="80">{{printf "%s" .Body}}</textarea></div>
<div><input type="submit" value="Save"></div>
</form>
```

修改editHandler函数

```go
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)

	// 如果用户输入的是新的title，loadPage会返回err，这里用用户新输入的title创建一个新的Page
	if err != nil {
		p = &Page{Title: title}
	}

    // ParseFiles读取”edit.html“文件并返回一个模板
    t, _ := template.ParseFiles("edit.html")
    // Execute会把模板写到ResponseWriter里
	t.Execute(w, p)
}
```

然后把之前viewHandler也改成使用模板，创建`view.html`

```html
<h1>{{.Title}}</h1>

<p>[<a href="/edit/{{.Title}}">edit</a>]</p>

<div>{{printf "%s" .Body}}</div>
```

然后修改`viewHandler`

```go
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	t, _ := template.ParseFiles("view.html")
	t.Execute(w, p)
}
```

`viewHandler`和`editHandler`都有处理模板的共同逻辑，可以把这块抽出来单独写个函数`renderTemplate`

```go
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	// 如果用户输入的是新的title，loadPage会返回err，这里用用户新输入的title创建一个新的Page
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}
```

## 处理不存在的页面

如果用户输入`view/ANonExistsPage`，我们将其重定向到edit页面让用户编辑新的Page

```go
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}
```

## 保存页面

类似前面的`editHandler`和`viewHandler`，我们为保存添加一个处理函数。

```go
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}
```

然后在main中注册`http.HandleFunc("/save/", saveHandler)`

## 错误处理

在`renderTemplate`和`saveHandler`都加上对错误的处理

http.Error()会返回一个response

```go
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
```

```go
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}
```

## 模板缓存

目前为止，我们每刷新一次view或edit页面，对应的模板都要重新加载。`template.Must`可以在程序开始执行时一次性的把模板都加载出来。

在全局位置加上Must

```go
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
```

再修改一下`renderTemplate`函数

```go
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
```

## 验证

到目前位置，我们直接把用户在url中输入的路径来创建或修改文件，这是非常危险的。现在用`regexp`对用户输入的url进行验证。

先创建一个全局变量

```go
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
```

正则表达式不在这里详细介绍。大概解释一下，`"^/(edit|save|view)/([a-zA-Z0-9]+)$"`的意思就是url的第一部分必须是`edit|save|view`其实中之一，url的第二部分必须是`[a-zA-Z0-9]`中的多个字符。这样就把path限制在了当前目录下。

在`viewHandler`, `editHandler`和`saveHandler`中都有从URL中提取title的部分，把他们做成一个公共函数。然后在三个handlers中使用这个函数。

```go
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	// url输入localhost:8888/view/hello
	// m的内容为下面，可见m[2]就是title
	// [/view/hello view hello]
	return m[2], nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}

	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}

	p, err := loadPage(title)
	// 如果用户输入的是新的title，loadPage会返回err，这里用用户新输入的title创建一个新的Page
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}

	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}
```

## 函数字面值和闭包

`viewHandler`, `editHandler`和`saveHandler`都有共同的逻辑，这里再抽象一下，用高阶函数和闭包

首先写一个创建handler的函数。这个函数的参数是一个函数，并且返回这个函数。在返回的函数里，会调用传进来的函数。

```go
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
```
