package vm

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var (
	httpRequestClass  *RClass
	httpResponseClass *RClass
)

func initHTTPClass(vm *VM) {
	net := vm.loadConstant("Net", true)
	http := vm.initializeClass("HTTP", false)
	http.setBuiltInMethods(builtinHTTPClassMethods(), true)
	initRequestClass(vm, http)
	initResponseClass(vm, http)

	net.setClassConstant(http)

	// Use Goby code to extend request and response classes.
	vm.execGobyLib("net/http/response.gb")
	vm.execGobyLib("net/http/request.gb")
}

func initRequestClass(vm *VM, hc *RClass) *RClass {
	requestClass := vm.initializeClass("Request", false)
	hc.setClassConstant(requestClass)
	builtinHTTPRequestInstanceMethods := []*BuiltInMethodObject{}

	requestClass.setBuiltInMethods(builtinHTTPRequestInstanceMethods, false)

	httpRequestClass = requestClass
	return requestClass
}

func initResponseClass(vm *VM, hc *RClass) *RClass {
	responseClass := vm.initializeClass("Response", false)
	hc.setClassConstant(responseClass)
	builtinHTTPResponseInstanceMethods := []*BuiltInMethodObject{}

	responseClass.setBuiltInMethods(builtinHTTPResponseInstanceMethods, false)

	httpResponseClass = responseClass
	return responseClass
}

func builtinHTTPClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			// Sends a GET request to the target and returns the HTTP response as a string.
			Name: "get",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					var path string

					domain := args[0].(*StringObject).value

					if len(args) > 1 {
						path = args[1].(*StringObject).value
					}

					if !strings.HasPrefix(path, "/") {
						path = "/" + path
					}

					resp, err := http.Get(domain + path)

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					content, err := ioutil.ReadAll(resp.Body)
					resp.Body.Close()

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					return t.vm.initStringObject(string(content))
				}
			},
		}, {
			// Sends a GET request to the target and returns the HTTP response as a string.
			Name: "post",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 3 {
						return t.vm.initErrorObject(ArgumentError, "Expect 3 arguments. got=%v", strconv.Itoa(len(args)))
					}

					url := args[0].(*StringObject).value
					contentType := args[1].(*StringObject).value
					body := args[2].(*StringObject).value

					resp, err := http.Post(url, contentType, strings.NewReader(body))

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					content, err := ioutil.ReadAll(resp.Body)
					resp.Body.Close()

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					return t.vm.initStringObject(string(content))
				}
			},
		},
	}
}
