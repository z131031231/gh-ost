/*
   Copyright 2016 GitHub Inc.
	 See https://github.com/github/gh-osc/blob/master/LICENSE
*/

package sql

import (
	"fmt"
	"reflect"
	"strings"
)

// ColumnList makes for a named list of columns
type ColumnList []string

// ParseColumnList parses a comma delimited list of column names
func ParseColumnList(columns string) *ColumnList {
	result := ColumnList(strings.Split(columns, ","))
	return &result
}

func (this *ColumnList) String() string {
	return strings.Join(*this, ",")
}

func (this *ColumnList) Equals(other *ColumnList) bool {
	return reflect.DeepEqual(*this, *other)
}

// ColumnsMap maps a column onto its ordinal position
type ColumnsMap map[string]int

func NewColumnsMap(columnList ColumnList) ColumnsMap {
	columnsMap := make(map[string]int)
	for i, column := range columnList {
		columnsMap[column] = i
	}
	return ColumnsMap(columnsMap)
}

// UniqueKey is the combination of a key's name and columns
type UniqueKey struct {
	Name        string
	Columns     ColumnList
	HasNullable bool
}

// IsPrimary checks if this unique key is primary
func (this *UniqueKey) IsPrimary() bool {
	return this.Name == "PRIMARY"
}

func (this *UniqueKey) String() string {
	return fmt.Sprintf("%s: %s; has nullable: %+v", this.Name, this.Columns, this.HasNullable)
}

type ColumnValues struct {
	abstractValues []interface{}
	ValuesPointers []interface{}
}

func NewColumnValues(length int) *ColumnValues {
	result := &ColumnValues{
		abstractValues: make([]interface{}, length),
		ValuesPointers: make([]interface{}, length),
	}
	for i := 0; i < length; i++ {
		result.ValuesPointers[i] = &result.abstractValues[i]
	}

	return result
}

func ToColumnValues(abstractValues []interface{}) *ColumnValues {
	result := &ColumnValues{
		abstractValues: abstractValues,
		ValuesPointers: make([]interface{}, len(abstractValues)),
	}
	for i := 0; i < len(abstractValues); i++ {
		result.ValuesPointers[i] = &result.abstractValues[i]
	}

	return result
}

func (this *ColumnValues) AbstractValues() []interface{} {
	return this.abstractValues
}

func (this *ColumnValues) StringColumn(index int) string {
	val := this.AbstractValues()[index]
	if ints, ok := val.([]uint8); ok {
		return string(ints)
	}
	return fmt.Sprintf("%+v", val)
}

func (this *ColumnValues) String() string {
	stringValues := []string{}
	for i := range this.AbstractValues() {
		stringValues = append(stringValues, this.StringColumn(i))
	}
	return strings.Join(stringValues, ",")
}
