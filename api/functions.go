package api

import (
	"bytes"
	"fmt"
	"github.com/RulzUrLife/lasagna/common"
	"math"
	"reflect"
)

func div(divisor, dividend int) int {
	return int(math.Floor((float64(dividend) / float64(divisor)) + .5))
}

func slice(start, end int, item interface{}) (interface{}, error) {
	v := reflect.ValueOf(item)
	if !v.IsValid() {
		return nil, fmt.Errorf("slice of untyped nil")
	}
	if start > end {
		return nil, fmt.Errorf("slice indexes must be: %d < %d", start, end)
	}
	if start < 0 || end > int(v.Len()) {
		return nil, fmt.Errorf("slice indexes out of range: %d, %d", start, end)
	}
	switch k := v.Kind(); k {
	case reflect.Array, reflect.Slice, reflect.String:
		return v.Slice(start, end).Interface(), nil
	default:
		return nil, fmt.Errorf("cannot slice type %s, need array, slice or string", k)
	}
}

func url(path string, items ...interface{}) (string, error) {
	b := bytes.NewBuffer(nil)
	_, err := fmt.Fprintf(b, common.Url(common.Config.Url, path), items...)
	return b.String(), err
}
