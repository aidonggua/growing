package gorm

import (
	"reflect"
	"github.com/aidonggua/growing/gutils"
	"fmt"
)

var StructInfoMap = make(map[reflect.Type]*StructInfo)

//结构体信息
type StructInfo struct {
	FieldsMap map[string]*StructField //字段字典集合
	Name      string                  //类型名
	TableName string                  //表名
}

//结构体字段信息
type StructField struct {
	name           string        //字段名
	value          reflect.Value //字段值
	stringValue    string        //字符串值
	tableFieldName string        //表属性名
	tableFieldType string        //表属性类型
}


//获得结构体的信息
func GetStructInfo(target interface{}) *StructInfo {
	v := reflect.Indirect(reflect.ValueOf(target))
	t := v.Type()
	//判断target的类型
	if t.Kind() != reflect.Struct {
		fmt.Println(GET_STRUCTINFO_ERROR)
		return nil
	}
	return GetReflectInfo(t, v)
}

//获得结构体的反射的信息
func GetReflectInfo(t reflect.Type, v reflect.Value) *StructInfo {

	var (
		tableName string
		tName string
		structInfo *StructInfo
	)

	fieldsMap := make(map[string]*StructField)
	//从map里取结构体信息,如果map没有则新建一个然后存map
	if value, ok := StructInfoMap[t]; ok {
		structInfo = value

		//更新缓存的结构体字段的值
		for key, _ := range structInfo.FieldsMap {
			structInfo.FieldsMap[key].value.Set(v.FieldByName(key))
			structInfo.FieldsMap[key].stringValue = gutils.ParseValueToString(v.FieldByName(key))
		}

	} else {
		//遍历所有属性
		for index := 0; index < t.NumField(); index++ {
			structField := t.Field(index)
			structFieldValue := v.Field(index)

			sf := &StructField{
				name:structField.Name,
				tableFieldName:gutils.UnCamelCase(structField.Name),
				tableFieldType:gutils.GetDBType(t.Kind().String()),
				value:structFieldValue,
				stringValue:gutils.ParseValueToString(structFieldValue),
			}

			fieldsMap[sf.name] = sf

			tableName = gutils.UnCamelCase(t.Name())
			tName = t.Name()

			structInfo = &StructInfo{
				Name:tName,
				TableName:tableName,
				FieldsMap:fieldsMap,
			}
			StructInfoMap[t] = structInfo
		}
	}
	return structInfo
}


