# ctl

Package ctl provides a library that makes it easy to embed a simple, minimal
client into your Go program to interact with your gRPC server.

## Usage

API is still WIP and subject to change, but here's the general gist:

```go
import "github.com/bobheadxi/ctl"

func main() {
  // instantiate your gRPC client
  c, _ := client.New( /* ... */ )

  // create a controller
  controller, _ := ctl.New(c)

  // show help if you want
  if os.Args != nil && len(os.Args) == 1 && os.Args[0] == "help" {
    controller.Help(os.Stdout)
    return
  }

  // execute command
  out, _ := controller.Exec(os.Args[0:], os.Stdout)

  // print the output
  fmt.Printf("%v\n", out)
}
```

In a command line application, you can then run:

```sh
$> my-command SomeRPCCall RequestParam=blah AnotherRequestParam=wow
```
