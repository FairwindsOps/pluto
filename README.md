# api-version-finder

This is a very simple utility to help users find deprecated Kubernetes apiVersions in their code repositories and their helm releases.

## QuickStart

Install the binary from our [releases](https://github.com/FairwindsOps/api-version-finder/releases) page.

Run `api-version-finder detect-files -d <DIRECTORY YOU WANT TO SCAN>`

You should see an output something like:

```
$ api-version-finder detect-files -d pkg/finder/testdata
KIND         VERSION              DEPRECATED   FILE
Deployment   extensions/v1beta1   true         pkg/finder/testdata/deployment-extensions-v1beta1.json
Deployment   extensions/v1beta1   true         pkg/finder/testdata/deployment-extensions-v1beta1.yaml
```

This indicates that we have two files in our directory that have deprecated apiVersions. We should fix this.

## Usage

```
A tool to detect Kubernetes apiVersions

Usage:
  api-version-finder [flags]
  api-version-finder [command]

Available Commands:
  detect-files detect-files
  help         Help about any command
  version      Prints the current version of the tool.

Flags:
      --add_dir_header                   If true, adds the file directory to the header
      --alsologtostderr                  log to standard error as well as files
  -h, --help                             help for api-version-finder
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --log_file string                  If non-empty, use this log file
      --log_file_max_size uint           Defines the maximum size a log file can grow to. Unit is megabytes. If the value is 0, the maximum file size is unlimited. (default 1800)
      --logtostderr                      log to standard error instead of files (default true)
      --skip_headers                     If true, avoid header prefixes in the log messages
      --skip_log_headers                 If true, avoid headers when opening log files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          number for the log level verbosity
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging

Use "api-version-finder [command] --help" for more information about a command.
```

## Detect Files Options

```
Usage:
  api-version-finder detect-files [flags]

Flags:
  -d, --directory string      The directory to scan. If blank, defaults to current workding dir.
  -h, --help                  help for detect-files
  -o, --output string         The output format to use. (tabular|json|yaml) (default "tabular")
      --show-non-deprecated   If enabled, will show files that have non-deprecated apiVersion. Only applies to tabular output.
```
