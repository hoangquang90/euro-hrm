package try

//O Try Object
func O(o interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return o
}

//S Try Object
func S(s string, err error) string {
	return O(s, err).(string)
}

//I Try Object
func I(i int, err error) int {
	return O(i, err).(int)
}

//L Try Object
func L(b bool, err error) bool {
	return O(b, err).(bool)
}

//Ba Try Object
func Ba(b []byte, err error) []byte {
	return O(b, err).([]byte)
}

//V Try Void
func V(err error) {
	if err != nil {
		panic(err)
	}
}
