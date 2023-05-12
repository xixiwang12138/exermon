package exermon

import "exermon/gateway"

func main() {
	gateway.Instance.RegisterGlobalMiddleWare()
}
