package main

import (
	"errors"
	"reflect"
)

type fieldType interface{}

type defaultFieldType struct {
	fieldType
}

func (_this defaultFieldType) DString(Default string) string {

	if _this.fieldType == nil {
		return Default
	}

	err, str := _this.String()

	if err != nil {
		return Default
	}

	return str

}

func (_this defaultFieldType) ArrayByte() (error, [][]byte) {

	var result [][]byte

	s := reflect.ValueOf(_this.fieldType)

	if s.Kind() != reflect.Slice {
		return errors.New("type não suportado"), nil
	}

	if s.IsNil() || s.IsZero() {
		return nil, nil
	}

	for i := 0; i < s.Len(); i++ {

		var aux defaultFieldType

		aux = defaultFieldType{fieldType: s.Index(i).Interface()}

		err, str := aux.String()

		if err != nil {
			continue
		}

		result = append(result, []byte(str))

	}

	return nil, result

}

func (_this defaultFieldType) ArrayString() (error, []string) {

	var result []string

	s := reflect.ValueOf(_this.fieldType)

	if s.Kind() != reflect.Slice {
		return errors.New("type não suportado"), nil
	}

	if s.IsNil() || s.IsZero() {
		return nil, nil
	}

	for i := 0; i < s.Len(); i++ {

		var aux defaultFieldType

		aux = defaultFieldType{fieldType: s.Index(i).Interface()}

		err, str := aux.String()

		if err != nil {
			continue
		}

		result = append(result, str)

	}

	return nil, result

}

func (_this defaultFieldType) String() (error, string) {

	s := reflect.ValueOf(_this.fieldType)

	if s.Kind() != reflect.String {
		return errors.New("type não suportado"), ""
	}

	return nil, _this.fieldType.(string)

}

func (_this defaultFieldType) DInteger(Default int64) int64 {

	if _this.fieldType == nil {
		return Default
	}

	err, i64 := _this.Integer()

	if err != nil {
		return Default
	}

	return i64

}

func (_this defaultFieldType) Integer() (error, int64) {

	s := reflect.ValueOf(_this.fieldType)

	if s.Kind() != reflect.Float64 {
		return errors.New("type não suportado"), 0
	}

	i64 := int64(_this.fieldType.(float64))

	return nil, i64

}

func (_this defaultFieldType) DBool(Default bool) bool {

	if _this.fieldType == nil {
		return Default
	}

	err, res := _this.Bool()

	if err != nil {
		return Default
	}

	return res

}

func (_this defaultFieldType) Bool() (error, bool) {

	s := reflect.ValueOf(_this.fieldType)

	if s.Kind() != reflect.Bool {
		return errors.New("type não suportado"), false
	}

	return nil, _this.fieldType.(bool)

}

func (_this defaultFieldType) Array() (error, []interface{}) {

	var result []interface{}

	s := reflect.ValueOf(_this.fieldType)

	if s.Kind() != reflect.Slice {
		return errors.New("type não suportado"), nil
	}

	if s.IsNil() || s.IsZero() {
		return nil, nil
	}

	for i := 0; i < s.Len(); i++ {

		result = append(result, s.Index(i).Interface())

	}

	return nil, result

}
