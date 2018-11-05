# ctl

## Usage

API is still WIP and subject to change, but here's the general gist:

```go
// instantiate your gRPC client
c, err := client.New( /* ... */ )
if err != nil {
	fatal(err.Error())
}
defer c.Close()

// create a controller
controller, err := ctl.New(c)
if err != nil {
  fatal(err.Error())
}

// show help if you want
if args != nil && len(args) == 1 && args[0] == "help" {
  controller.Help(os.Stdout)
  return
}

// execute command
out, err := controller.Exec(args, os.Stdout)
if err != nil {
  fatal(err.Error())
}

// print the output
fmt.Printf("%v\n", out)
```
