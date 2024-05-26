package component

import "fmt"

type ComponentError struct{}

func RaiseIfComponentError(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}
	panic("[Component] " + err.Error() + fmt.Sprintf(msg, args...))
}
