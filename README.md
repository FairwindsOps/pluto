[![CircleCI](https://circleci.com/gh/FairwindsOps/pluto.svg?style=svg)](https://circleci.com/gh/FairwindsOps/pluto)[![codecov](https://codecov.io/gh/FairwindsOps/pluto/branch/master/graph/badge.svg?token=A23F79JTNA)](https://codecov.io/gh/FairwindsOps/pluto)

# pluto

This is a very simple utility to help users find deprecated Kubernetes apiVersions in their code repositories and their helm releases.

**Want to learn more?** Reach out on [the Slack channel](https://fairwindscommunity.slack.com/messages/goldilocks), send an email to `opensource@fairwinds.com`, or join us for [office hours on Zoom](https://fairwindscommunity.slack.com/messages/office-hours)

## QuickStart

Install the binary from our [releases](https://github.com/FairwindsOps/pluto/releases) page.

### File Detection

Run `pluto detect-files -d <DIRECTORY YOU WANT TO SCAN>`

You should see an output something like:

```
$ pluto detect-files -d pkg/finder/testdata
KIND         VERSION              DEPRECATED   FILE
Deployment   extensions/v1beta1   true         pkg/finder/testdata/deployment-extensions-v1beta1.json
Deployment   extensions/v1beta1   true         pkg/finder/testdata/deployment-extensions-v1beta1.yaml
```

This indicates that we have two files in our directory that have deprecated apiVersions. This will need to be fixed prior to a 1.16 upgrade.

### Helm Detection

```
$ pluto detect-helm --helm-version 3
KIND          VERSION        DEPRECATED   RESOURCE NAME
StatefulSet   apps/v1beta1   true         audit-dashboard-prod-rabbitmq-ha
```

This indicates that the StatefulSet audit-dashboard-prod-rabbitmq-ha was deployed with apps/v1beta1 which is deprecated in 1.16

## Usage

```
A tool to detect Kubernetes apiVersions

Usage:
  pluto [flags]
  pluto [command]

Available Commands:
  detect-files detect-files
  detect-helm  detect-helm
  help         Help about any command
  version      Prints the current version of the tool.

Flags:
      --add_dir_header                   If true, adds the file directory to the header
      --alsologtostderr                  log to standard error as well as files
  -h, --help                             help for pluto
      --kubeconfig string                Paths to a kubeconfig. Only required if out-of-cluster.
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --log_file string                  If non-empty, use this log file
      --log_file_max_size uint           Defines the maximum size a log file can grow to. Unit is megabytes. If the value is 0, the maximum file size is unlimited. (default 1800)
      --logtostderr                      log to standard error instead of files (default true)
      --master --kubeconfig              (Deprecated: switch to --kubeconfig) The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.
      --skip_headers                     If true, avoid header prefixes in the log messages
      --skip_log_headers                 If true, avoid headers when opening log files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          number for the log level verbosity
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging

Use "pluto [command] --help" for more information about a command.
```

## Detect Files Options

```
Usage:
  pluto detect-files [flags]

Flags:
  -d, --directory string      The directory to scan. If blank, defaults to current workding dir.
  -h, --help                  help for detect-files
  -o, --output string         The output format to use. (tabular|json|yaml) (default "tabular")
      --show-non-deprecated   If enabled, will show files that have non-deprecated apiVersion. Only applies to tabular output.
```

## Detect Helm Options

NOTE: Only helm 3 is currently supported

```
Detect Kubernetes apiVersions in a helm release (in cluster)

Usage:
  pluto detect-helm [flags]

Flags:
      --helm-version string   Helm version in current cluster (2|3) (default "3")
  -h, --help                  help for detect-helm
      --show-non-deprecated   If enabled, will show files that have non-deprecated apiVersion. Only applies to tabular output.
```
