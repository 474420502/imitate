package main

import (
	"errors"
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

// ScriptBook 程序的脚本加载总集合
var ScriptBook map[string]interface{}

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

	LoadScript("script") // script.py加载所有script文件夹下的脚本
	log.Println("python init success!", path.Dir(filename))
}

// TaskInfo 任务url 配置相关属性
type TaskInfo struct {
	BaseURL string
	Method  string
	Header  http.Header
	Cookies []*http.Cookie
	Query   url.Values
	Body    map[string][]string // 目前从Python data参数 转换 缺少data为 byte数据的例子 TODO: 待完善
}

// TaskSetting 任务基本设置
type TaskSetting struct {
	Name      string
	GroupName string

	Session  int
	Mode     TypeMode
	Retry    int
	Priority int

	Plan ExecutePlan

	Proxies []string

	NextDo string

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
	setting *py.PyObject
	info    *py.PyObject

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

func (tf *TaskConfig) infoFromImportPythonScript() {
	tf.Info = &TaskInfo{}
	tf.Info.Header = make(http.Header)
	tf.Info.Query = make(url.Values)
	tf.Info.Body = make(map[string][]string)

	tf.Info.Cookies = make([]*http.Cookie, 0)
	tempCookies := make(map[string]string)

	headers := tf.info.GetAttrString("headers")
	query := tf.info.GetAttrString("query_params")
	data := tf.info.GetAttrString("data")
	cookies := tf.info.GetAttrString("cookies")

	UpdateMapListFromPyDict(tf.Info.Header, headers)
	UpdateMapListFromPyDict(tf.Info.Query, query)
	UpdateMapListFromPyDict(tf.Info.Body, data)
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

func (tf *TaskConfig) settingFromImportPythonScript() {
	var attr *py.PyObject

	tf.Setting = &TaskSetting{}

	attr = tf.setting.GetAttrString("name")
	if attr != nil {
		tf.Setting.Name = py.PyString_AsString(attr)
	} else {
		tf.Setting.Name = ""
	}

	attr = tf.setting.GetAttrString("session")
	if attr != nil {
		tf.Setting.Session = py.PyInt_AsLong(attr)
	} else {
		tf.Setting.Session = 0
	}

	attr = tf.setting.GetAttrString("mode")
	if attr != nil {
		tf.Setting.Mode = TypeMode(py.PyInt_AsLong(attr))
	} else {
		tf.Setting.Mode = ModeCookie
	}

	attr = tf.setting.GetAttrString("retry")
	if attr != nil {
		tf.Setting.Retry = py.PyInt_AsLong(attr)
	} else {
		tf.Setting.Retry = 0
	}

	attr = tf.setting.GetAttrString("priority")
	if attr != nil {
		tf.Setting.Priority = py.PyInt_AsLong(attr)
	} else {
		tf.Setting.Priority = 10000
	}

	attr = tf.setting.GetAttrString("group_name")
	if attr != nil {
		tf.Setting.GroupName = py.PyString_AsString(attr)
	} else {
		tf.Setting.GroupName = "p"
	}

	attr = tf.setting.GetAttrString("device")
	if attr != nil {
		tf.Setting.Device = py.PyString_AsString(attr)
	} else {
		tf.Setting.Device = ""
	}

	attr = tf.setting.GetAttrString("platform")
	if attr != nil {
		tf.Setting.Platform = py.PyString_AsString(attr)
	} else {
		tf.Setting.Platform = ""
	}

	attr = tf.setting.GetAttrString("area_cc")
	if attr != nil {
		tf.Setting.AreaCC = py.PyInt_AsLong(attr)
	} else {
		tf.Setting.AreaCC = -1
	}

	attr = tf.setting.GetAttrString("channel")
	if attr != nil {
		tf.Setting.Channel = py.PyInt_AsLong(attr)
	} else {
		tf.Setting.Channel = -1
	}

	attr = tf.setting.GetAttrString("media")
	if attr != nil {
		tf.Setting.Media = py.PyInt_AsLong(attr)
	} else {
		tf.Setting.Media = -1
	}

	attr = tf.setting.GetAttrString("spider_id")
	if attr != nil {
		tf.Setting.SpiderID = py.PyInt_AsLong(attr)
	} else {
		tf.Setting.SpiderID = -1
	}

	attr = tf.setting.GetAttrString("catch_account_id")
	if attr != nil {
		tf.Setting.CatchAccountID = py.PyString_AsString(attr)
	} else {
		tf.Setting.CatchAccountID = ""
	}

	attr = tf.setting.GetAttrString("next_do")
	if attr != nil {
		tf.Setting.NextDo = py.PyString_AsString(attr)
	} else {
		tf.Setting.NextDo = "save"
	}

	tf.Setting.Proxies = nil
	attr = tf.setting.GetAttrString("proxies")
	if attr != nil {
		l := py.PyList_GET_SIZE(attr)
		for i := 0; i < l; i++ {
			tf.Setting.Proxies = append(tf.Setting.Proxies, py.PyString_AsString(py.PyList_GetItem(attr, i)))
		}
	}

	tf.Setting.Plan.ClearIExecute()

	// execute_at = (-1, -1, -1, 12, 30, 12)+
	attr = tf.setting.GetAttrString("execute_at")
	if attr != nil {
		if py.PyList_Check(attr) {
			l := py.PyList_GET_SIZE(attr)
			for i := 0; i < l; i++ {
				ea := &ExecuteAt{}

				ea.FromPyObject(py.PyList_GetItem(attr, i))
				ea.SetStartStatus(true)
				ea.CalculateTrigger()
				tf.Setting.Plan.AppendIExecute(ea)
			}
		} else {
			ea := &ExecuteAt{}
			ea.FromPyObject(attr)
			ea.SetStartStatus(true)
			ea.CalculateTrigger()
			tf.Setting.Plan.AppendIExecute(ea)
		}

	}

	attr = tf.setting.GetAttrString("interval")
	if attr != nil {

		if py.PyList_Check(attr) {
			l := py.PyList_GET_SIZE(attr)
			for i := 0; i < l; i++ {
				ei := &ExecuteInterval{}
				ei.FromPyObject(py.PyList_GetItem(attr, i))
				ei.SetStartStatus(true)
				ei.CalculateTrigger()
				tf.Setting.Plan.AppendIExecute(ei)
			}
		} else {
			ei := &ExecuteInterval{}
			ei.FromPyObject(attr)
			ei.SetStartStatus(true)
			ei.CalculateTrigger()
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

	tf.loadedFilename = filename
	tf.info = py.PyImport_ImportModule(importpathTURL)
	if tf.info == nil {
		log.Panic(errors.New("info is error, path is " + importpathTURL))
	}
	tf.setting = py.PyImport_ImportModule(importpathConfig)
	if tf.setting == nil {
		log.Panic(errors.New("setting is error, path is " + importpathConfig))
	}

	tf.settingFromImportPythonScript()
	tf.infoFromImportPythonScript()
}

func (tf *TaskConfig) String() string {
	res := fmt.Sprintf("name: %s\ngroup_name: %s\nsession: %d\nretry: %d\npriority: %d\nexecute_at: %s\nproxies: %s\nnext_do: %s\ndevice: %s\nplatform: %s\narea_cc: %d\nchannel: %d\nmedia: %d\nspider_id: %d\ncatch_account_id: %s\n",
		tf.Setting.Name, tf.Setting.GroupName, tf.Setting.Session,
		tf.Setting.Retry, tf.Setting.Priority, spew.Sdump(tf.Setting.Plan),
		spew.Sdump(tf.Setting.Proxies), tf.Setting.NextDo,
		tf.Setting.Device, tf.Setting.Platform, tf.Setting.AreaCC,
		tf.Setting.Channel, tf.Setting.Media, tf.Setting.SpiderID,
		tf.Setting.CatchAccountID,
	)
	return res
}
