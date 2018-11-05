package ctl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// CTL is a client controller
type CTL struct {
	Registry map[string]reflect.Type

	client interface{}
}

// New instantiates a new controller
func New(client interface{}) (CTL, error) {
	r := make(map[string]reflect.Type)

	// generate type dictionary based on arguments of given client's methods
	cType := reflect.TypeOf(client)
	for i := 0; i < cType.NumMethod(); i++ {
		t := cType.Method(i).Type
		if t.NumIn() > 2 {
			arg := t.In(2).Elem()
			r[arg.PkgPath()+"."+arg.Name()] = arg
		}
	}

	return CTL{
		Registry: r,
		client:   client,
	}, nil
}

// Exec takes command line args and maps them to a client call
func (c *CTL) Exec(ctx context.Context, args []string, out io.Writer) (interface{}, error) {
	if args == nil || len(args) < 1 {
		return nil, errors.New("insufficient arguments provided")
	}

	fn := reflect.ValueOf(c.client).MethodByName(args[0])
	fmt.Fprintf(out, "function %s found\n", args[0])

	// go-grpc interfaces should have 2 arguments
	in := make([]reflect.Value, 2)

	// first argument is context
	in[0] = reflect.ValueOf(ctx)

	// second argument is a struct
	fnType, _ := reflect.TypeOf(c.client).MethodByName(args[0])
	argType := fnType.Type.In(2).Elem()
	arg := c.instantiate(argType.PkgPath() + "." + argType.Name())
	for _, v := range args[0:] {
		split := strings.Split(v, "=")
		if len(split) > 1 {
			setProperty(split[0], split[1], arg)
		}
	}
	in[1] = reflect.ValueOf(arg)

	// output generated call
	fmt.Fprintf(out, "generated function call: \n%s(ctx, { %v})\n", args[0], arg)

	// execute and get results of call: [interface{}, error]
	result := fn.Call(in)

	return result[0].Interface(), result[1].Interface().(error)
}

// Help lists functions and arguments of the client
func (c *CTL) Help(out io.Writer) {
	cType := reflect.TypeOf(c.client)
	for i := 0; i < cType.NumMethod(); i++ {
		method := cType.Method(i)

		fmt.Fprintf(out, "%s:\n", method.Name)
		t := method.Type
		if t.NumIn() > 2 {
			arg := t.In(2)
			fmt.Fprintf(out, "  %s\n", arg.String())
			fmt.Fprintln(out, "  arguments:")
			elem := arg.Elem()
			hasArgs := false
			for j := 0; j < elem.NumField(); j++ {
				if !strings.Contains(elem.Field(j).Name, "XXX") {
					hasArgs = true
					fmt.Fprintf(out, "    %s=%s\n", elem.Field(j).Name, elem.Field(j).Type.Name())
				}
			}
			if !hasArgs {
				fmt.Fprintln(out, "    none")
			}
		}
		fmt.Fprintln(out, "")
	}
}

func (c *CTL) instantiate(t string) interface{} {
	v := reflect.New(c.Registry[t]).Elem()
	// Maybe fill in fields here if necessary
	return v.Addr().Interface()
}
