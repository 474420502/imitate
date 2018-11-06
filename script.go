package main

import (
	"errors"
	"log"

	py "github.com/sbinet/go-python"
)

// ScriptResult 固定的脚本结构
type ScriptResult struct {
	NextDo string
	Result interface{}
}

func callScript(sresult *ScriptResult) {
	if sresult.NextDo != "" {
		if method, ok := ScriptBook[sresult.NextDo]; ok {
			switch m := method.(type) {
			case *py.PyObject:
				result := m.CallFunction(sresult.Result) // GoResponseToPy(PResult.Resp)
				if result != nil {
					if py.PyTuple_Check(result) {
						l := py.PyTuple_GET_SIZE(result)
						nextDo := py.PyString_AsString(py.PyList_GetItem(result, 0))
						sr := &ScriptResult{NextDo: nextDo}
						if l >= 2 {
							params := py.PyList_New(l - 1)
							for i := 1; i < l; i++ {
								item := py.PyList_GetItem(result, i)
								py.PyList_Append(params, item)
							}
							sr.Result = params
						}
						callScript(sr)
					} else if py.PyString_Check(result) {
						sr := &ScriptResult{NextDo: py.PyString_AS_STRING(result)}
						callScript(sr)
					}
				}
			}

		} else {
			panic(errors.New("method is error, key is " + sresult.NextDo))
		}

	}
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
	log.Println(ScriptBook)
}
