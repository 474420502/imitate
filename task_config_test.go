package main

import (
	"testing"

	"github.com/sbinet/go-python"
)

// TestLoadConfig 加载配置函数功能
func TestLoadConfig(t *testing.T) {

	m := python.PyImport_ImportModule("task_config")
	if m == nil {
		t.Error("can not load module task_config, pointer is", m)
	}

	attrSession := m.GetAttrString("session")
	v1 := python.PyInt_AsLong(attrSession)
	if v1 != 1 {
		t.Error("error task_config.session is not 1, and session is", v1)
	}

	attrDevice := m.GetAttrString("device")
	v2 := python.PyString_AsString(attrDevice)
	if v2 != "eson-OnePlus" {
		t.Error("error task_config.session is not eson-OnePlus, and session is", v2)
	}
}

func TestTaskConfigLoad(t *testing.T) {
	tf := TaskConfig{}
	tf.Load("task/task_config.py")
	if tf.Channel != 105 {
		t.Error("config Channel is not 105, is ", tf.Channel)
	}

	if tf.Media != 55 {
		t.Error("config Media is not 55, is ", tf.Media)
	}

	if tf.SpiderID != 73 {
		t.Error("config SpiderID is not 73, is ", tf.SpiderID)
	}

	if tf.Platform != "Android" {
		t.Error("config Platform is not Android, is ", tf.Platform)
	}
}

func TestExecuteAt(t *testing.T) {
	tf := TaskConfig{}
	tf.Load("task/task_config.py")

	if tf.Plan.CountIExecute() != 2 {
		t.Error("iexecute size is not 2, is ", tf.Plan.CountIExecute())
	}
}
