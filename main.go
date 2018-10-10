package main

import (
	"fmt"
	"log"
	"path"
	"runtime"

	py "github.com/sbinet/go-python"
)

// PythonInit 初始化
func init() {
	err := py.Initialize()
	if err != nil {
		panic(err.Error())
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	py.PyRun_SimpleString("import sys")
	py.PyRun_SimpleString(fmt.Sprintf("sys.path.append(\"%s\")", path.Dir(filename)))
	log.Println("python init success!", path.Dir(filename))
}

// PersonGroup 私人组
type PersonGroup interface {
}

func main() {
	// ses := grequests.NewSession(nil)

	// burl := "https://www.baidu.com"

	// ses.Get(burl)
	// log.Println(ses.RequestOptions.DomainCookies)
}
