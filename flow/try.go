package flow

func Try(fn func(), handler func(interface{})) {
	defer func() {
		if r := recover(); r != nil {
			handler(r)
		}
	}()
	fn()
	return
}
