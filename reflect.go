package dsy

import "reflect"

// satisfyFields checks that the map[string]interface data satisfies the field names and types of the struct.
func satisfyFields(m map[string]interface{}, s interface{}) bool {
	switch reflect.TypeOf(s).Kind() {
	case reflect.Struct:
		rv := reflect.ValueOf(s)
		rt := rv.Type()

		if rt.NumField() != len(m) {
			return false
		}

		for i := 0; i < rt.NumField(); i++ {
			rf := rt.Field(i)
			// check field name
			if _, ok := m[rf.Name]; !ok {
				return false
			}
			// check type
			if rf.Type != reflect.TypeOf(m[rf.Name]) {
				return false
			}
		}
		return true

	default:
		return false
	}
}
