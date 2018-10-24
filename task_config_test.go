package main

import (
	"testing"
)

func TestTaskConfigLoad(t *testing.T) {
	tf := TaskConfig{}
	tf.Load("task/task_config.py")
	if tf.Setting.Channel != 105 {
		t.Error("config Channel is not 105, is ", tf.Setting.Channel)
	}

	if tf.Setting.Media != 55 {
		t.Error("config Media is not 55, is ", tf.Setting.Media)
	}

	if tf.Setting.SpiderID != 73 {
		t.Error("config SpiderID is not 73, is ", tf.Setting.SpiderID)
	}

	if tf.Setting.Platform != "Android" {
		t.Error("config Platform is not Android, is ", tf.Setting.Platform)
	}
}

func TestExecuteAt(t *testing.T) {
	tf := TaskConfig{}
	tf.Load("task/task_config.py")

	if tf.Setting.Plan.CountIExecute() != 2 {
		t.Error("iexecute size is not 2, is ", tf.Setting.Plan.CountIExecute())
	}
}
