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
	client   interface{}
	registry map[string]reflect.Type
}

// New instantiates a new controller
func New(client interface{}) (CTL, error) {
	r := make(map[string]reflect.Type)

	cType := reflect.TypeOf(client)
	for i := 0; i < cType.NumMethod(); i++ {
		t := cType.Method(i).Type
		if t.NumIn() > 2 {
			arg := t.In(2).Elem()
			r[arg.PkgPath()+"."+arg.Name()] = arg
		}
	}

	fmt.Printf("registered: %v\n", r)

	return CTL{
		client:   client,
		registry: r,
	}, nil
}

// Exec takes command line args and maps them to a client call
func (c *CTL) Exec(args []string, out io.Writer) (interface{}, error) {
	if args == nil || len(args) < 1 {
		return nil, errors.New("insufficient arguments provided")
	}

	fn := reflect.ValueOf(c.client).MethodByName(args[0])
	fmt.Fprintf(out, "function %s found\n", args[0])

	in := make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(context.Background())

	// construct input
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

	fmt.Fprintf(out, "generated argument: \n{ %v}\n", arg)

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
	fmt.Printf("%s: %v\n", t, c.registry[t])
	v := reflect.New(c.registry[t]).Elem()
	// Maybe fill in fields here if necessary
	return v.Addr().Interface()
}

func setProperty(name string, value string, obj interface{}) bool {
	val := reflect.ValueOf(obj)

	if val.Kind() != reflect.Ptr {
		return false
	}
	structVal := val.Elem()
	for i := 0; i < structVal.NumField(); i++ {
		valueField := structVal.Field(i)
		typeField := structVal.Type().Field(i)
		if strings.ToLower(typeField.Name) == strings.ToLower(name) {
			if valueField.IsValid() && valueField.CanSet() && valueField.Kind() == reflect.String {
				valueField.SetString(value)
				return true
			}
		}
	}
	return false
}
