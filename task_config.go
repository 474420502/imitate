package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"regexp"

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

// TaskInfo 任务url 配置相关属性
type TaskInfo struct {
	BaseURL string
	Method  string
	Header  http.Header
	Cookies []*http.Cookie
	Query   url.Values
	Body    map[string]interface{}
}

// TaskSetting 任务基本设置
type TaskSetting struct {
	Name      string
	GroupName string

	Session  int
	Retry    int
	Priority int

	Plan ExecutePlan

	Proxies []string

	ResultProcessing string

	Device         string
	Platform       string
	AreaCC         int
	Channel        int
	Media          int
	SpiderID       int
	CatchAccountID string
}

// TaskConfig 任务配置相关结构
type TaskConfig struct {
	conf *py.PyObject
	info *py.PyObject

	loadedFilename string

	Info    *TaskInfo
	Setting *TaskSetting
}

// NewTaskConfig New一个TaskConfig 个例
func NewTaskConfig(filename string) *TaskConfig {
	tf := &TaskConfig{}
	tf.Load(filename)
	return tf
}

// Reload 重新加载配置文件
func (tf *TaskConfig) Reload() bool {
	if tf.loadedFilename == "" {
		log.Println("filename is null")
		return false
	}

	tf.Load(tf.loadedFilename)

	return true
}

// UpdateMapFromPyDict 从Python的dict 更新 golang map[string]interface{}
func UpdateMapFromPyDict(gomap map[string]interface{}, pydict *py.PyObject) {
	headersItemList := py.PyDict_Items(pydict)

	if headersItemList != nil {
		l := py.PyList_GET_SIZE(headersItemList)
		for i := 0; i < l; i++ {
			pyitem1 := py.PyList_GetItem(headersItemList, i)
			key := py.PyString_AsString(py.PyTuple_GetItem(pyitem1, 0))
			value := py.PyTuple_GetItem(pyitem1, 1)

			if py.PyString_Check(value) {
				gomap[key] = py.PyString_AsString(value)
			} else if py.PyList_Check(value) {
				ll := py.PyList_GET_SIZE(value)

				var vlist []string
				for ii := 0; ii < ll; ii++ {
					lvalue := py.PyString_AsString(py.PyList_GetItem(value, ii))
					vlist = append(vlist, lvalue)
				}
				gomap[key] = vlist
			}
		}
	}
}

// UpdateMapListFromPyDict 更新 golang map[string]string
func UpdateMapListFromPyDict(gomap map[string][]string, pydict *py.PyObject) {
	headersItemList := py.PyDict_Items(pydict)

	if headersItemList != nil {
		l := py.PyList_GET_SIZE(headersItemList)
		for i := 0; i < l; i++ {
			pyitem1 := py.PyList_GetItem(headersItemList, i)
			key := py.PyString_AsString(py.PyTuple_GetItem(pyitem1, 0))
			value := py.PyTuple_GetItem(pyitem1, 1)

			if py.PyString_Check(value) {
				gomap[key] = []string{py.PyString_AsString(value)}
			} else if py.PyList_Check(value) {
				ll := py.PyList_GET_SIZE(value)
				var vlist []string
				for ii := 0; ii < ll; ii++ {
					lvalue := py.PyString_AsString(py.PyList_GetItem(value, ii))
					vlist = append(vlist, lvalue)
				}
				gomap[key] = vlist
			}
		}
	}
}

// UpdateStringKVFromPyDict 更新 golang map[string]string
func UpdateStringKVFromPyDict(gomap map[string]string, pydict *py.PyObject) {
	headersItemList := py.PyDict_Items(pydict)

	if headersItemList != nil {
		l := py.PyList_GET_SIZE(headersItemList)
		for i := 0; i < l; i++ {
			pyitem1 := py.PyList_GetItem(headersItemList, i)
			key := py.PyString_AsString(py.PyTuple_GetItem(pyitem1, 0))
			value := py.PyString_AsString(py.PyTuple_GetItem(pyitem1, 1))

			gomap[key] = value
		}
	}
}

func (tf *TaskConfig) infoFromImportPythonScript() {
	tf.Info = &TaskInfo{}
	tf.Info.Header = make(http.Header)
	tf.Info.Query = make(url.Values)
	tf.Info.Body = make(map[string]interface{})

	tf.Info.Cookies = nil
	tempCookies := make(map[string]string)

	headers := tf.info.GetAttrString("headers")
	query := tf.info.GetAttrString("query_params")
	data := tf.info.GetAttrString("data")
	cookies := tf.info.GetAttrString("cookies")

	UpdateMapListFromPyDict(tf.Info.Header, headers)
	UpdateMapListFromPyDict(tf.Info.Query, query)
	UpdateMapFromPyDict(tf.Info.Body, data)
	UpdateStringKVFromPyDict(tempCookies, cookies)

	for k, v := range tempCookies {
		tf.Info.Cookies = append(tf.Info.Cookies, &http.Cookie{
			Name:     k,
			Value:    v,
			HttpOnly: true,
			Secure:   false,
		})
	}

	tf.Info.Method = strings.ToUpper(py.PyString_AsString(tf.info.GetAttrString("method")))
	tf.Info.BaseURL = py.PyString_AsString(tf.info.GetAttrString("base_url"))
}

func (tf *TaskConfig) confFromImportPythonScript() {
	var attr *py.PyObject

	tf.Setting = &TaskSetting{}

	attr = tf.conf.GetAttrString("name")
	if attr != nil {
		tf.Setting.Name = py.PyString_AsString(attr)
	} else {
		tf.Setting.Name = ""
	}

	attr = tf.conf.GetAttrString("session")
	if attr != nil {
		tf.Setting.Session = py.PyInt_AsLong(attr)
	} else {
		tf.Setting.Session = 0
	}

	attr = tf.conf.GetAttrString("retry")
	if attr != nil {
		tf.Setting.Retry = py.PyInt_AsLong(attr)
	} else {
		tf.Setting.Retry = 0
	}

	attr = tf.conf.GetAttrString("priority")
	if attr != nil {
		tf.Setting.Priority = py.PyInt_AsLong(attr)
	} else {
		tf.Setting.Priority = 10000
	}

	attr = tf.conf.GetAttrString("group_name")
	if attr != nil {
		tf.Setting.GroupName = py.PyString_AsString(attr)
	} else {
		tf.Setting.GroupName = "p"
	}

	attr = tf.conf.GetAttrString("device")
	if attr != nil {
		tf.Setting.Device = py.PyString_AsString(attr)
	} else {
		tf.Setting.Device = ""
	}

	attr = tf.conf.GetAttrString("platform")
	if attr != nil {
		tf.Setting.Platform = py.PyString_AsString(attr)
	} else {
		tf.Setting.Platform = ""
	}

	attr = tf.conf.GetAttrString("area_cc")
	if attr != nil {
		tf.Setting.AreaCC = py.PyInt_AsLong(attr)
	} else {
		tf.Setting.AreaCC = -1
	}

	attr = tf.conf.GetAttrString("channel")
	if attr != nil {
		tf.Setting.Channel = py.PyInt_AsLong(attr)
	} else {
		tf.Setting.Channel = -1
	}

	attr = tf.conf.GetAttrString("media")
	if attr != nil {
		tf.Setting.Media = py.PyInt_AsLong(attr)
	} else {
		tf.Setting.Media = -1
	}

	attr = tf.conf.GetAttrString("spider_id")
	if attr != nil {
		tf.Setting.SpiderID = py.PyInt_AsLong(attr)
	} else {
		tf.Setting.SpiderID = -1
	}

	attr = tf.conf.GetAttrString("catch_account_id")
	if attr != nil {
		tf.Setting.CatchAccountID = py.PyString_AsString(attr)
	} else {
		tf.Setting.CatchAccountID = ""
	}

	attr = tf.conf.GetAttrString("result_processing")
	if attr != nil {
		tf.Setting.ResultProcessing = py.PyString_AsString(attr)
	} else {
		tf.Setting.ResultProcessing = "save"
	}

	tf.Setting.Proxies = nil
	attr = tf.conf.GetAttrString("proxies")
	if attr != nil {
		l := py.PyList_GET_SIZE(attr)
		for i := 0; i < l; i++ {
			tf.Setting.Proxies = append(tf.Setting.Proxies, py.PyString_AsString(py.PyList_GetItem(attr, i)))
		}
	}

	tf.Setting.Plan.ClearIExecute()

	// execute_at = (-1, -1, -1, 12, 30, 12)+
	attr = tf.conf.GetAttrString("execute_at")
	if attr != nil {
		if py.PyList_Check(attr) {
			l := py.PyList_GET_SIZE(attr)
			for i := 0; i < l; i++ {
				ea := &ExecuteAt{}
				ea.SetStartStatus(true)
				ea.FromPyObject(py.PyList_GetItem(attr, i))
				tf.Setting.Plan.AppendIExecute(ea)
			}
		} else {
			ea := &ExecuteAt{}
			ea.SetStartStatus(true)
			ea.FromPyObject(attr)
			tf.Setting.Plan.AppendIExecute(ea)
		}

	}

	attr = tf.conf.GetAttrString("interval")
	if attr != nil {

		if py.PyList_Check(attr) {
			l := py.PyList_GET_SIZE(attr)
			for i := 0; i < l; i++ {
				ei := &ExecuteInterval{}
				ei.SetStartStatus(true)
				ei.FromPyObject(py.PyList_GetItem(attr, i))
				tf.Setting.Plan.AppendIExecute(ei)
			}
		} else {
			ei := &ExecuteInterval{}
			ei.SetStartStatus(true)
			ei.FromPyObject(attr)
			tf.Setting.Plan.AppendIExecute(ei)
		}
	}

}

// Load 加载配置文件
func (tf *TaskConfig) Load(filename string) {

	filename = strings.Replace(filename, "_config", "", -1)

	rec, err := regexp.Compile("/|\\\\|\\..*")
	if err != nil {
		panic(err)
	}
	filename = rec.ReplaceAllString(filename, ".")
	lastidx := len(filename) - 1

	var importpathTURL string
	var importpathConfig string
	if filename[lastidx] == '.' {
		importpathTURL = filename[0 : len(filename)-1]
	} else {
		importpathTURL = filename
	}

	importpathConfig = importpathTURL + "_config"
	log.Println("importpathTURL:", importpathTURL)
	log.Println("importpathConfig:", importpathConfig)

	tf.loadedFilename = filename
	tf.info = py.PyImport_ImportModule(importpathTURL)
	tf.conf = py.PyImport_ImportModule(importpathConfig)

	tf.confFromImportPythonScript()
	tf.infoFromImportPythonScript()
}

func (tf *TaskConfig) String() string {
	res := fmt.Sprintf("name: %s\ngroup_name: %s\nsession: %d\nretry: %d\npriority: %d\nexecute_at: %s\nproxies: %s\nresult_processing: %s\ndevice: %s\nplatform: %s\narea_cc: %d\nchannel: %d\nmedia: %d\nspider_id: %d\ncatch_account_id: %s\n",
		tf.Setting.Name, tf.Setting.GroupName, tf.Setting.Session,
		tf.Setting.Retry, tf.Setting.Priority, spew.Sdump(tf.Setting.Plan),
		spew.Sdump(tf.Setting.Proxies), tf.Setting.ResultProcessing,
		tf.Setting.Device, tf.Setting.Platform, tf.Setting.AreaCC,
		tf.Setting.Channel, tf.Setting.Media, tf.Setting.SpiderID,
		tf.Setting.CatchAccountID,
	)
	return res
}
