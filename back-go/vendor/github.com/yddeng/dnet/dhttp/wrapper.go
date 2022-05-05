package dhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

type Handle struct {
}

func httpServeError(w http.ResponseWriter, status int, txt string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintln(w, txt)
}

func httpHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	(*w).Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	(*w).Header().Set("content-type", "application/json")             //返回数据格式是json
}

func HandlerFuncJsonJson(fn interface{}) http.HandlerFunc {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	if typ.NumIn() != 1 || typ.NumOut() != 1 {
		panic("func symbol error")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		httpHeader(&w)
		inType := typ.In(0)
		var inValue reflect.Value
		if inType.Kind() == reflect.Ptr {
			inValue = reflect.New(inType.Elem())
		} else {
			inValue = reflect.New(inType)
		}

		if err := json.NewDecoder(r.Body).Decode(inValue.Interface()); err != nil {
			httpServeError(w, 404, err.Error())
			return
		}

		if inType.Kind() != reflect.Ptr {
			inValue = inValue.Elem()
		}

		outValue := val.Call([]reflect.Value{inValue})[0]
		json.NewEncoder(w).Encode(outValue.Interface())
	}
}

func HandlerFuncParamJson(fn interface{}) http.HandlerFunc {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	if typ.NumIn() != 1 || typ.NumOut() != 1 {
		panic("func symbol error")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()

		inType := typ.In(0)
		var inValue reflect.Value
		if inType.Kind() == reflect.Ptr {
			inValue = reflect.New(inType.Elem())
		} else {
			inValue = reflect.New(inType)
		}

		if err := json.NewDecoder(r.Body).Decode(inValue.Interface()); err != nil {
			httpServeError(w, 404, err.Error())
			return
		}

		if inType.Kind() != reflect.Ptr {
			inValue = inValue.Elem()
		}

		outValue := val.Call([]reflect.Value{inValue})[0]
		json.NewEncoder(w).Encode(outValue.Interface())
	}
}
