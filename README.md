# ctl

## Usage

API is still WIP and subject to change, but here's the general gist:

```go
// instantiate your gRPC client
c, _ := client.New( /* ... */ )
defer c.Close()

// create a controller
controller, _ := ctl.New(c)

// show help if you want
if args != nil && len(args) == 1 && args[0] == "help" {
  controller.Help(os.Stdout)
  return
}

// execute command
out, _ := controller.Exec(args, os.Stdout)

// print the output
fmt.Printf("%v\n", out)
```

In a command line application, you can then run:

```sh
$> my-command SomeRPCCall RequestParam=blah AnotherRequestParam=wow
```
