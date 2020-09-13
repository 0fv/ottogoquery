package ottogoquery

import (
	"errors"
	"reflect"

	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
)

func functionInit(vm *otto.Otto, doc *goquery.Document) {
	setDocExec(vm, doc)
}

func executeDoc(doc *goquery.Document, vm *otto.Otto, method string, param ...interface{}) (func(call otto.FunctionCall) otto.Value, error) {
	var rvalues []reflect.Value = make([]reflect.Value, 0, len(param))
	for _, v := range param {
		rvalues = append(rvalues, reflect.ValueOf(v))
	}
	m := reflect.ValueOf(doc).MethodByName(method)
	if m.IsZero() {
		return nil, errors.New("function not found")
	}
	result := m.Call(rvalues)
	if len(result) != 1 {
		return nil, errors.New("function not support")
	}
	d := result[0]
	selector, ok := d.Interface().(*goquery.Selection)
	if !ok {
		return nil, errors.New("function not support")
	}
	return func(call otto.FunctionCall) otto.Value {
		if len(selector.Nodes) == 0 {
			return otto.NullValue()
		}
		method2, err := call.Argument(0).ToString()
		if err != nil {
			return otto.NullValue()
		}
		t := call.Argument(1)
		var param2 reflect.Value
		flag := false
		switch {
		case t.IsUndefined():
			flag = true
		case t.IsNumber():
			p, err := t.ToInteger()
			if err != nil {
				return otto.NullValue()
			}
			param2 = reflect.ValueOf(int(p))
		case t.IsString():
			p, err := t.ToString()
			if err != nil {
				return otto.NullValue()
			}
			param2 = reflect.ValueOf(p)
		default:
			return otto.NullValue()
		}
		m2 := reflect.ValueOf(selector).MethodByName(method2)
		if m2.IsZero() {
			return otto.NullValue()
		}
		var result2 []reflect.Value
		if flag {
			result2 = m2.Call([]reflect.Value{})
		} else {
			result2 = m2.Call([]reflect.Value{param2})
		}
		if len(result2) != 1 {
			return otto.NullValue()
		}
		d := result2[0]
		vv := d.Interface()
		s2, ok := vv.(*goquery.Selection)
		if ok {
			selector = s2
			return otto.NullValue()
		}
		val, err := vm.ToValue(vv)
		if err != nil {
			return otto.NullValue()
		}
		return val
	}, nil
}

func setDocExec(vm *otto.Otto, doc *goquery.Document) {
	vm.Set("docExec", func(call otto.FunctionCall) otto.Value {
		method, err := call.Argument(0).ToString()
		if err != nil {
			return otto.NullValue()
		}
		selector, err := call.Argument(1).ToString()
		if err != nil {
			return otto.NullValue()
		}
		f, err := executeDoc(doc, vm, method, selector)
		if err != nil {
			return otto.NullValue()
		}
		of, err := vm.ToValue(f)
		if err != nil {
			return otto.NullValue()
		}
		return of
	})
}
