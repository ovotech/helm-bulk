package utils

import "log"

//PanicCheck panics if error is not nil
func PanicCheck(e error) {
	if e != nil {
		log.Panic(e.Error())
	}
}
