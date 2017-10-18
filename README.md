# Decision Maker

[ ![Codeship Status for jaypeng2015/go-decision-maker](https://app.codeship.com/projects/9ffed2c0-9513-0135-cd90-467ad8efcfd1/status?branch=master)](https://app.codeship.com/projects/251157)

The [decision maker](https://github.com/jaypeng2015/decision-maker) in Go.

Instead of using AWS Step Functions, this one is just AWS API Gateway + Lambda Function.

## Step by Step

### Install Go

Mac:

```
brew install go
```

The latest version of go will be installed into `~/go`.

For non Mac users, please read [this](https://golang.org/doc/install).

### Quick Start

```
cd ~/go
mkdir -p src/jaypeng2015/
git clone https://github.com/jaypeng2015/decision-maker
cd decision-maker
go get -v ./...
go test -v
```

### Setup IDE

 - Install [Visual Studio Code](https://code.visualstudio.com/)
 - Install [Go Extension](https://github.com/Microsoft/vscode-go/)
 - Install [Go Tools](https://github.com/Microsoft/vscode-go/wiki/Go-tools-that-the-Go-extension-depends-on) from Go Extension

 ### Provision

 #### Configure your AWS credentials

  - Install the [AWS CLI](http://docs.aws.amazon.com/cli/latest/userguide/installing.html) for your operating system.
  - Configure your AWS security access keys

    ```
      aws --profile decisionmaker configure
      AWS Access Key ID [None]: xxxxxxxxxxxxxxxx
      AWS Secret Access Key [None]: xxxxxxxxxxxxxxxxxxxx
      Default region name [None]: ap-southeast-2
      Default output format [None]: json
    ```

  - Activate the decisionmaker profile

    ```
      export AWS_PROFILE=decisionmaker
    ```

#### Command Lines

A compiled application provides several command line options which are available by providing the -h/--help option as in:

```
$ go run application.go --help

The Decision Maker in Go.

Usage:
  application [command]

Available Commands:
  delete      Delete service
  describe    Describe service
  execute     Execute
  explore     Interactively explore service
  help        Help about any command
  provision   Provision service
  version     Sparta framework version

Flags:
  -f, --format string    Log format [text, json] (default "text")
  -h, --help             help for application
      --ldflags string   Go linker string definition flags (https://golang.org/cmd/link/)
  -l, --level string     Log level [panic, fatal, error, warn, info, debug] (default "info")
  -n, --noop             Dry-run behavior only (do not perform mutations)
  -t, --tags string      Optional build tags for conditional compilation

Use "application [command] --help" for more information about a command.

```

More information can be found [here](http://gosparta.io/docs/application/commandline/).
