package main

import (
	"log"
	"reflect"
)

type Student struct {
	StuName string
	StuAge  int
}
type Class struct {
	ClassName string
	Students  Student
}

func setNestedStructValue(val reflect.Value, fieldName string, setV any) {
	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		setNestedStructValue(val.Elem(), fieldName, setV)
	case reflect.Struct:
		fieldVal := val.FieldByName(fieldName)
		if fieldVal.IsValid() {
			fieldVal.Set(reflect.ValueOf(setV).Convert(fieldVal.Type()))
			return
		} else {
			for i := 0; i < val.NumField(); i++ {
				field := val.Field(i)
				setNestedStructValue(field, fieldName, setV)
			}
		}

	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			setNestedStructValue(elem, fieldName, setV)
		}
	case reflect.String:
		//_, b := val.Type().FieldByName(fieldName)
		log.Println(val.String())

	}
}

func setT[T any]() {
	t := new(T)
	//val := reflect.ValueOf(&class)
	val := reflect.ValueOf(t)
	setNestedStructValue(val.Elem(), "StuName", "linx")

	log.Printf("%#v\n", *t)

	setNestedStructValue(val.Elem(), "StuAge", 12)
	log.Printf("%#v\n", *t)

	setNestedStructValue(val.Elem(), "ClassName", "学前班")
	log.Printf("%#v\n", *t)
}

func main() {
	setT[Class]()

	//class := Class{}
	////val := reflect.ValueOf(&class)
	//val := reflect.ValueOf(new(Class))
	//setNestedStructValue(val.Elem(), "StuName", "linx")
	//
	//log.Printf("%#v\n", class)
	//
	//setNestedStructValue(val.Elem(), "StuAge", "linx")
	//log.Printf("%#v\n", class)
	//
	//setNestedStructValue(val.Elem(), "ClassName", "学前班")
	//log.Printf("%#v\n", class)
}
