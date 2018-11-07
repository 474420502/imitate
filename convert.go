package main

import (
	py "github.com/sbinet/go-python"
)

// PyDict 简化Python的一堆字符的命令生成一个Dict
type PyDict struct {
	object *py.PyObject
}

// NewPyDict new pydict
func NewPyDict() *PyDict {
	return &PyDict{object: py.PyDict_New()}
}

// UpdateStrInt dict[key] = value value: int
func (dict *PyDict) UpdateStrInt(key string, value int) {
	py.PyDict_SetItem(dict.object, py.PyString_FromString(key), py.PyInt_FromLong(value))
}

// UpdateStrStr dict[key] = value value: int
func (dict *PyDict) UpdateStrStr(key string, value string) {
	py.PyDict_SetItem(dict.object, py.PyString_FromString(key), py.PyString_FromString(value))
}

// UpdateIntInt dict[key] = value value: int
func (dict *PyDict) UpdateIntInt(key int, value int) {
	py.PyDict_SetItem(dict.object, py.PyInt_FromLong(key), py.PyInt_FromLong(value))
}

// UpdateIntStr dict[key] = value value: int
func (dict *PyDict) UpdateIntStr(key int, value string) {
	py.PyDict_SetItem(dict.object, py.PyInt_FromLong(key), py.PyString_FromString(value))
}

// PyObject 返回python的object
func (dict *PyDict) PyObject() *py.PyObject {
	return dict.object
}

// UpdateMapFromPyDict 从Python的dict 更新至 golang map[string]interface{}
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

// GoResponseToPy Go Response 转 Python格式
// func GoResponseToPy(gresp *requests.Response) *py.PyObject {
// 	obj := py.PyDict_New()
// 	py.PyDict_SetItem(obj, py.PyString_FromString("status"), py.PyInt_FromLong(gresp.GResponse.StatusCode))
// 	py.PyDict_SetItem(obj, py.PyString_FromString("content"), py.PyString_FromString(gresp.Content()))
// 	return obj
// }
