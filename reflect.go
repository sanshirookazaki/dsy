package dsy

import "reflect"

func satisfyFields(m map[string]interface{}, s interface{}) bool {
	switch reflect.TypeOf(s).Kind() {
	case reflect.Struct:
		rv := reflect.ValueOf(s)
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			rf := rt.Field(i)
			// check field name
			if _, ok := m[rf.Name]; !ok {
				return false
			}
			// check type
			//t := fmt.Sprintf("%T", v)
			//if t != rf.Type.Kind().String() {
			//	return false
			//}
		}
		return true

	default:
		return false
	}
}
