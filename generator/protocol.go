package main

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
)

type TypeBrowserProtocol map[string]interface{}

func (_this *TypeBrowserProtocol) Load(fileName string) error {

	fileTest, err := os.Stat(fileName)

	if !((err == nil) && !fileTest.IsDir()) {
		return errors.New("File config.json not found")
	}

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&_this)

	if err != nil {
		return err
	}

	return nil

}

func (_this *TypeBrowserProtocol) Exist(Keys ...string) bool {

	err, _ := _this.Value(Keys...)

	return err == nil

}

func (_this *TypeBrowserProtocol) FValue(Keys ...string) defaultFieldType {

	err, val := _this.Value(Keys...)

	if err != nil {
		return defaultFieldType{nil}
	}

	return val

}

func (_this *TypeBrowserProtocol) Value(Keys ...string) (error, defaultFieldType) {

	var err error
	var val interface{}
	var ok bool

	if _this == nil {
		return errors.New("nil"), defaultFieldType{nil}
	}

	auxObj := *_this

	for idx, key := range Keys {

		val, ok = auxObj[key]

		if !ok {
			return errors.New("key: " + key + " not found"), defaultFieldType{nil}
		}

		if idx != len(Keys)-1 {

			s := reflect.ValueOf(val)

			if s.Kind() != reflect.Map {
				return errors.New("Tipo não suportado"), defaultFieldType{nil}
			}

			err, auxObj = getMap(val)

			if err != nil {
				return err, defaultFieldType{nil}
			}

		}

	}

	return nil, defaultFieldType{fieldType: val}

}

func getMap(obj interface{}) (error, map[string]interface{}) {

	myMap := make(map[string]interface{})

	if obj == nil {
		return errors.New("nil"), nil
	}

	s := reflect.ValueOf(obj)

	if s.Kind() != reflect.Map {
		return errors.New("Tipo não suportado"), nil
	}

	for _, key := range s.MapKeys() {

		strct := s.MapIndex(key)
		myMap[key.String()] = strct.Interface()

	}

	return nil, myMap
}
