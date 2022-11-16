---
meta:
  - name: description
    content: "Fairwinds Pluto | Contribution Guidelines"
---
# Contributing Guide

Issues, whether bugs, tasks, or feature requests are essential for keeping pluto great. We believe it should be as easy as possible to contribute changes that get things working in your environment. There are a few guidelines that we need contributors to follow so that we can keep on top of things.

## Code of Conduct

This project adheres to a [code of conduct](code-of-conduct.md). Please review this document before contributing to this project.

## Sign the CLA
Before you can contribute, you will need to sign the [Contributor License Agreement](https://cla-assistant.io/fairwindsops/pluto).

## Project Structure

### CLI
pluto is a relatively simple cobra cli tool that helps deal with deprecated api versions in Kubernetes. The [/cmd](https://github.com/FairwindsOps/pluto/tree/master/cmd) folder contains the flags and other cobra config, while the [/pkg](https://github.com/FairwindsOps/pluto/tree/master/pkg) folder has the various packages.

### API

This contains the various structs and helper functions to deal with Kubernetes objects and their apiVersions. It assumes that any file we care about will have a `Kind` and an `apiVersion`.

### Finder

This package is for dealing with a set of static files and analyzing the apiVersions in them. It can search through a directory and find any files that conform to the specifications of the versions package.

### Helm

This package deals with finding helm release manifests that are in your running Kubernetes cluster.

## Versions Updates

The versions.yaml file contains the source of truth for the standard list of deprecated versions that Pluto can deal with. If you wish to update it or change the default targetVersions, it is easiest to use Pluto to do that.

Just create a deprecated-versions yaml file that has the additional versions and/or target versions you wish to edit. Here's an example:

```
target-versions:
  other: v1.0.0
deprecated-versions:
- version: someother/v1beta1
  kind: AnotherCRD
  deprecated-in: v1.9.0
  removed-in: v1.16.0
  replacement-api: apps/v1
  component: new-component
```

Put this somewhere temporary, like /tmp/new.yaml. Then you can run:

```
pluto list-versions -f /tmp/new.yaml -oyaml > versions.yaml
```

This will ensure that you have parsable valid versions, and it should make patching easier.

## Getting Started

We label issues with the ["good first issue" tag](https://github.com/FairwindsOps/pluto/labels/good%20first%20issue) if we believe they'll be a good starting point for new contributors. If you're interested in working on an issue, please start a conversation on that issue, and we can help answer any questions as they come up.

## Setting Up Your Development Environment
### Prerequisites
* A properly configured Golang environment with Go 1.13 or higher
* Install `pkger` - [documentation](https://github.com/markbates/pkger#cli)

### Installation
* Clone the project with `go get github.com/fairwindsops/pluto`
* Change into the pluto directory which is installed at `$GOPATH/src/github.com/fairwindsops/pluto`
* Use `make build` to build the binary locally.
* Use `make test` to run the tests and generate a coverage report.

## Creating a New Issue

If you've encountered an issue that is not already reported, please create an issue that contains the following:

- Clear description of the issue
- Steps to reproduce it
- Appropriate labels

## Creating a Pull Request

Each new pull request should:

- Reference any related issues
- Add tests that show the issues have been solved
- Pass existing tests and linting
- Contain a clear indication of if they're ready for review or a work in progress
- Be up to date and/or rebased on the master branch

## Orb development

There is a Makefile that can assist in validating and testing the orb locally.  See the commands there for more info.

## Creating a new release

Push a new annotated tag.  This tag should contain a changelog of pertinent changes. Goreleaser will take care of the rest.

## Pre-commit

This repo contains a pre-commit file for use with [pre-commit](https://pre-commit.com/). Just run `pre-commit install` and you will have the hooks.
