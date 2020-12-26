// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nn

import (
	"fmt"
	"github.com/nlpodyssey/spago/pkg/utils"
	"reflect"
	"strings"
)

// paramsTraversal allows the traversal of Model parameters.
// The given callback is invoked for each parameter of the Model.
// If exploreSubModels is true, every nested Model and its parameters are
// also visited.
type paramsTraversal struct {
	callback         func(param Param)
	exploreSubModels bool
}

// newParamsTraversal returns a new paramsTraversal.
func newParamsTraversal(callback func(param Param), exploreSubModels bool) paramsTraversal {
	return paramsTraversal{
		callback:         callback,
		exploreSubModels: exploreSubModels,
	}
}

// walk iterates through all the parameters of m.
// TODO: don't loop the field every time, use a lazy initialized "params list" instead
func (pt paramsTraversal) walk(m interface{}) {
	utils.ForEachField(m, func(field interface{}, name string, tag reflect.StructTag) {
		v := reflect.ValueOf(field)
		switch v.Kind() {
		case reflect.Struct, reflect.Ptr, reflect.Interface:
			pt.walkStructOrPtr(field, name, tag)
		case reflect.Slice:
			pt.walkSlice(v, name, tag)
		case reflect.Map:
			pt.walkMap(v, name, tag)
		}
	})
}

func (pt paramsTraversal) walkStructOrPtr(item interface{}, name string, tag reflect.StructTag) {
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr && v.Elem().Kind() != reflect.Struct {
		return
	}
	switch itemT := item.(type) {
	case *param:
		pt.walkParam(itemT, name, tag)
	case Module:
		if pt.exploreSubModels {
			pt.walk(item)
		}
	default:
		if tag.Get("type") == "params" {
			pt.walk(item)
		}
	}
}

func (pt paramsTraversal) walkSlice(v reflect.Value, name string, tag reflect.StructTag) {
	length := v.Len()
	for i := 0; i < length; i++ {
		p := v.Index(i)
		switch p.Kind() {
		case reflect.Struct, reflect.Ptr, reflect.Interface:
			pt.walkStructOrPtr(p.Interface(), name, tag)
		default:
			return // skip
		}
	}
}

func (pt paramsTraversal) walkMap(v reflect.Value, name string, tag reflect.StructTag) {
	mapRange := v.MapRange()
	for mapRange.Next() {
		key := ""
		switch k := mapRange.Key().Interface().(type) {
		case string:
			key = k
		case int:
			key = fmt.Sprintf("%d", k)
		default:
			return // skip map if the key is not a string or an int
		}

		name := strings.ToLower(fmt.Sprintf("%s.%s", name, key))
		switch mapRange.Value().Kind() {
		case reflect.Struct, reflect.Ptr, reflect.Interface:
			pt.walkStructOrPtr(mapRange.Value().Interface(), name, tag)
		default:
			return // skip
		}
	}
}

func (pt paramsTraversal) walkParam(item *param, name string, tag reflect.StructTag) {
	if item.Name() == "" {
		item.SetName(strings.ToLower(name))
	}
	item.SetType(ToType(tag.Get("type")))
	pt.callback(item)
}
