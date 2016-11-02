package dt

import (
	"database/sql"
	"fmt"
	"reflect"
)

type DatabaseTemplate interface {
	// Query(sql string, mapRow MapRow, params ...interface{}) (object interface{}, err error)
	ExecDDL(sql string, params ...interface{}) (err error)
	Exec(sql string, params ...interface{}) (err error)
	ExecForResult(sql string, params ...interface{}) (sql.Result, error)
	QueryArray(sql string, mapRow MapRow, params ...interface{}) ([]interface{}, error)
	// QueryIntoArray(resultList interface{}, sql string, mapRow MapRow, params ...interface{}) error
	QueryObject(sql string, mapRow MapRow, params ...interface{}) (interface{}, error)
	Close() error
}

type DatabaseTemplateImpl struct {
	Conn *sql.DB
}

type MapRow func(resultSet *sql.Rows) (object interface{}, err error)

func (this *DatabaseTemplateImpl) Close() (err error) {
	if this.Conn == nil {
		return nil
	}
	err = this.Conn.Close()
	this.Conn = nil
	return
}

func (this *DatabaseTemplateImpl) QueryArray(sql string, mapRow MapRow, params ...interface{}) ([]interface{}, error) {
	result, err := this.Conn.Query(sql, params...)
	d := func() {
		if result != nil {
			result.Close()
		}
	}
	defer d()
	if err != nil {
		return nil, err
	}
	var resArray []interface{}
	if result == nil {
		return nil, nil
	}
	for result.Next() {
		obj, err := mapRow(result)
		if err != nil {
			return nil, err
		}
		resArray = append(resArray, obj)
	}
	return resArray, nil
}

func (this *DatabaseTemplateImpl) QueryObject(sql string, mapRow MapRow, params ...interface{}) (interface{}, error) {
	result, error := this.Conn.Query(sql, params...)
	d := func() {
		if result != nil {
			result.Close()
		}
	}
	defer d()
	if error != nil {
		return nil, error
	}
	if result == nil {
		return nil, nil
	}
	if result.Next() {
		object, err := mapRow(result)
		return object, err
	}
	return nil, nil
}

func (this *DatabaseTemplateImpl) Exec(sql string, params ...interface{}) (err error) {
	_, error := this.Conn.Exec(sql, params...)
	if error != nil {
		err = error
	}
	return
}

func (this *DatabaseTemplateImpl) ExecDDL(sql string, params ...interface{}) (err error) {
	_, error := this.Conn.Exec(sql, params...)
	if error != nil {
		err = error
	}

	return
}

func (this *DatabaseTemplateImpl) ExecForResult(sql string, params ...interface{}) (sql.Result, error) {
	result, error := this.Conn.Exec(sql, params...)
	return result, error
}

// toSliceType returns the element type of the given object, if the object is a
// "*[]*Element" or "*[]Element". If not, returns nil.
// err is returned if the user was trying to pass a pointer-to-slice but failed.
func toSliceType(i interface{}) (reflect.Type, error) {
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr {
		// If it's a slice, return a more helpful error message
		if t.Kind() == reflect.Slice {
			return nil, fmt.Errorf("database_template: Cannot SELECT into a non-pointer slice: %v", t)
		}
		return nil, nil
	}
	if t = t.Elem(); t.Kind() != reflect.Slice {
		return nil, nil
	}
	return t.Elem(), nil
}

func toType(i interface{}) (reflect.Type, error) {
	t := reflect.TypeOf(i)

	// If a Pointer to a type, follow
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("database_template: Cannot SELECT into this type: %v", reflect.TypeOf(i))
	}
	return t, nil
}

func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
