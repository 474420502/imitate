package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"regexp"

	py "github.com/sbinet/go-python"
)

// NewTaskConfig New一个TaskConfig 个例
func NewTaskConfig(filename string) *TaskConfig {
	tf := &TaskConfig{}
	tf.Load(filename)
	return tf
}

// TaskURL 任务url 配置相关属性
type TaskURL struct {
	BaseURL string
	Method  string
	Headers map[string]string
	Cookies []*http.Cookie
	Params  map[string]interface{}
	Data    map[string]interface{}
}

// TaskConfig 任务配置相关结构
type TaskConfig struct {
	conf *py.PyObject
	turl *py.PyObject

	loadedFilename string

	TURL *TaskURL

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

func (tf *TaskConfig) turlFromImportPythonScript() {
	tf.TURL = &TaskURL{}
	tf.TURL.Headers = make(map[string]string)
	tf.TURL.Params = make(map[string]interface{})
	tf.TURL.Data = make(map[string]interface{})

	tf.TURL.Cookies = nil
	tempCookies := make(map[string]string)

	headers := tf.turl.GetAttrString("headers")
	params := tf.turl.GetAttrString("query_params")
	data := tf.turl.GetAttrString("data")
	cookies := tf.turl.GetAttrString("cookies")

	UpdateStringKVFromPyDict(tf.TURL.Headers, headers)
	UpdateMapFromPyDict(tf.TURL.Params, params)
	UpdateMapFromPyDict(tf.TURL.Data, data)
	UpdateStringKVFromPyDict(tempCookies, cookies)

	for k, v := range tempCookies {
		tf.TURL.Cookies = append(tf.TURL.Cookies, &http.Cookie{
			Name:     k,
			Value:    v,
			HttpOnly: true,
			Secure:   false,
		})
	}

	tf.TURL.Method = strings.ToUpper(py.PyString_AsString(tf.turl.GetAttrString("method")))
	tf.TURL.BaseURL = py.PyString_AsString(tf.turl.GetAttrString("base_url"))
}

func (tf *TaskConfig) confFromImportPythonScript() {
	var attr *py.PyObject

	attr = tf.conf.GetAttrString("name")
	if attr != nil {
		tf.Name = py.PyString_AsString(attr)
	} else {
		tf.Name = ""
	}

	attr = tf.conf.GetAttrString("session")
	if attr != nil {
		tf.Session = py.PyInt_AsLong(attr)
	} else {
		tf.Session = 0
	}

	attr = tf.conf.GetAttrString("retry")
	if attr != nil {
		tf.Retry = py.PyInt_AsLong(attr)
	} else {
		tf.Retry = 0
	}

	attr = tf.conf.GetAttrString("priority")
	if attr != nil {
		tf.Priority = py.PyInt_AsLong(attr)
	} else {
		tf.Priority = 10000
	}

	attr = tf.conf.GetAttrString("group_name")
	if attr != nil {
		tf.GroupName = py.PyString_AsString(attr)
	} else {
		tf.GroupName = "p"
	}

	attr = tf.conf.GetAttrString("device")
	if attr != nil {
		tf.Device = py.PyString_AsString(attr)
	} else {
		tf.Device = ""
	}

	attr = tf.conf.GetAttrString("platform")
	if attr != nil {
		tf.Platform = py.PyString_AsString(attr)
	} else {
		tf.Platform = ""
	}

	attr = tf.conf.GetAttrString("area_cc")
	if attr != nil {
		tf.AreaCC = py.PyInt_AsLong(attr)
	} else {
		tf.AreaCC = -1
	}

	attr = tf.conf.GetAttrString("channel")
	if attr != nil {
		tf.Channel = py.PyInt_AsLong(attr)
	} else {
		tf.Channel = -1
	}

	attr = tf.conf.GetAttrString("media")
	if attr != nil {
		tf.Media = py.PyInt_AsLong(attr)
	} else {
		tf.Media = -1
	}

	attr = tf.conf.GetAttrString("spider_id")
	if attr != nil {
		tf.SpiderID = py.PyInt_AsLong(attr)
	} else {
		tf.SpiderID = -1
	}

	attr = tf.conf.GetAttrString("catch_account_id")
	if attr != nil {
		tf.CatchAccountID = py.PyString_AsString(attr)
	} else {
		tf.CatchAccountID = ""
	}

	attr = tf.conf.GetAttrString("result_processing")
	if attr != nil {
		tf.ResultProcessing = py.PyString_AsString(attr)
	} else {
		tf.ResultProcessing = "save"
	}

	tf.Proxies = nil
	attr = tf.conf.GetAttrString("proxies")
	if attr != nil {
		l := py.PyList_GET_SIZE(attr)
		for i := 0; i < l; i++ {
			tf.Proxies = append(tf.Proxies, py.PyString_AsString(py.PyList_GetItem(attr, i)))
		}
	}

	tf.Plan.ClearIExecute()

	// execute_at = (-1, -1, -1, 12, 30, 12)+
	attr = tf.conf.GetAttrString("execute_at")
	if attr != nil {
		if py.PyList_Check(attr) {
			l := py.PyList_GET_SIZE(attr)
			for i := 0; i < l; i++ {
				ea := &ExecuteAt{}
				ea.SetStartStatus(true)
				ea.FromPyObject(py.PyList_GetItem(attr, i))
				tf.Plan.AppendIExecute(ea)
			}
		} else {
			ea := &ExecuteAt{}
			ea.SetStartStatus(true)
			ea.FromPyObject(attr)
			tf.Plan.AppendIExecute(ea)
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
				tf.Plan.AppendIExecute(ei)
			}
		}
	} else {
		ei := &ExecuteInterval{}
		ei.SetStartStatus(true)
		ei.FromPyObject(attr)
		tf.Plan.AppendIExecute(ei)
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
	tf.turl = py.PyImport_ImportModule(importpathTURL)
	tf.conf = py.PyImport_ImportModule(importpathConfig)

	tf.confFromImportPythonScript()
	tf.turlFromImportPythonScript()
}

func (tf *TaskConfig) String() string {
	res := fmt.Sprintf("name: %s\ngroup_name: %s\nsession: %d\nretry: %d\npriority: %d\nexecute_at: %s\nproxies: %s\nresult_processing: %s\ndevice: %s\nplatform: %s\narea_cc: %d\nchannel: %d\nmedia: %d\nspider_id: %d\ncatch_account_id: %s\n",
		tf.Name, tf.GroupName, tf.Session,
		tf.Retry, tf.Priority, spew.Sdump(tf.Plan),
		spew.Sdump(tf.Proxies), tf.ResultProcessing,
		tf.Device, tf.Platform, tf.AreaCC,
		tf.Channel, tf.Media, tf.SpiderID,
		tf.CatchAccountID,
	)
	return res
}
