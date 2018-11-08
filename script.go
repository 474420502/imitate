package main

import (
	"errors"

	py "github.com/sbinet/go-python"
)

// ScriptResult 固定的脚本结构
type ScriptResult struct {
	NextDo string
	Result interface{}
}

// callScript 递归脚本 直到 完成
func callScript(sresult *ScriptResult) error {
	if sresult.NextDo != "" {
		if method, ok := ScriptBook[sresult.NextDo]; ok {
			switch m := method.(type) {
			case *py.PyObject:
				result := m.CallFunction(sresult.Result) // GoResponseToPy(PResult.Resp)
				// log.Println(sresult.NextDo)
				if result != nil {
					if py.PyTuple_Check(result) {
						l := py.PyTuple_GET_SIZE(result)
						nextDo := py.PyString_AsString(py.PyTuple_GetItem(result, 0))
						sr := &ScriptResult{NextDo: nextDo}
						if l >= 2 {
							item := py.PyTuple_GetItem(result, 1)
							sr.Result = item
						}
						return callScript(sr)
					} else if py.PyString_Check(result) {
						sr := &ScriptResult{NextDo: py.PyString_AS_STRING(result)}
						return callScript(sr)
					}
				}
			}
		} else {
			return errors.New("method is error, key is " + sresult.NextDo)
		}
	}

	return nil
}

// LoadScript script.py加载所有script文件夹下的脚本
func LoadScript(spath string) {

	ScriptBook = make(map[string]interface{})

	sbook := py.PyImport_ImportModule(spath).CallMethod("load_script")
	sbookItems := py.PyDict_Items(sbook)
	l := py.PyList_GET_SIZE(sbookItems)
	for i := 0; i < l; i++ {
		item := py.PyList_GetItem(sbookItems, i)
		key := py.PyString_AS_STRING(py.PyTuple_GET_ITEM(item, 0))
		value := py.PyTuple_GET_ITEM(item, 1)
		ScriptBook[key] = value
	}
}
