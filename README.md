# Pluto
<div align="center">
  <img src="/img/logo.png" alt="Pluto Logo" height="200"/>
  <br>

  <h3>Find Kubernetes resources that have been deprecated</h3>

  [![GitHub release (latest SemVer)][release-image]][release-link] [![GitHub go.mod Go version][version-image]][version-link] [![CircleCI][circleci-image]][circleci-link] [![Code Coverage][codecov-image]][codecov-link] [![Go Report Card][goreport-image]][goreport-link]
</div>

[version-image]: https://img.shields.io/github/go-mod/go-version/FairwindsOps/pluto
[version-link]: https://github.com/FairwindsOps/pluto

[release-image]: https://img.shields.io/github/v/release/FairwindsOps/pluto
[release-link]: https://github.com/FairwindsOps/pluto

[goreport-image]: https://goreportcard.com/badge/github.com/FairwindsOps/pluto
[goreport-link]: https://goreportcard.com/report/github.com/FairwindsOps/pluto

[circleci-image]: https://circleci.com/gh/FairwindsOps/pluto/tree/master.svg?style=svg
[circleci-link]: https://circleci.com/gh/FairwindsOps/pluto

[codecov-image]: https://codecov.io/gh/FairwindsOps/pluto/branch/master/graph/badge.svg
[codecov-link]: https://codecov.io/gh/FairwindsOps/pluto

This is a very simple utility to help users find deprecated Kubernetes apiVersions in their code repositories and their helm releases.

**Want to learn more?** Reach out on [the Slack channel](https://fairwindscommunity.slack.com/messages/pluto) ([request invite](https://join.slack.com/t/fairwindscommunity/shared_invite/zt-e3c6vj4l-3lIH6dvKqzWII5fSSFDi1g)), send an email to `opensource@fairwinds.com`, or join us for [office hours on Zoom](https://fairwindscommunity.slack.com/messages/office-hours)

## Purpose

Kubernetes sometimes deprecates apiVersions. Most notably, a large number of deprecations happened in the [1.16 release](https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/). This is fine, and it's a fairly easy thing to deal with. However, it can be difficult to find all the places where you might have used a version that will be deprecated in your next upgrade.

You might think, "I'll just ask the api-server to tell me!", but this is fraught with danger. If you ask the api-server to give you `deployments.v1.apps`, and the deployment was deployed as `deployments.v1beta1.extensions`, the api-server will quite happily convert the api version and return a manifest with `apps/v1`. This is fairly well outlined in the discussion in [this issue](https://github.com/kubernetes/kubernetes/issues/58131#issuecomment-356823588).

So, long story short, finding the places where you have deployed a deprecated apiVersion can be challenging. This is where `pluto` comes in. You can use pluto to check a couple different places where you might have placed a deprecated version:
* Infrastructure-as-Code repos: Pluto can check both static manifests and Helm charts for deprecated apiVersions
* Live Helm releases: Pluto can check both Helm 2 and Helm 3 releases running in your cluster for deprecated apiVersions

## Kubernetes Deprecation Policy

You can read the full policy [here](https://kubernetes.io/docs/reference/using-api/deprecation-policy/)

Long story short, apiVersions get deprecated, and then they eventually get removed entirely. Pluto differentiates between these two, and will tell you if a version is `DEPRECATED` or `REMOVED`

## Installation

### asdf

We have an asdf plugin [here](https://github.com/FairwindsOps/asdf-pluto). You can install with:

```
asdf plugin-add pluto
asdf list-all pluto
asdf install pluto <latest version>
```

### Binary

Install the binary from our [releases](https://github.com/FairwindsOps/pluto/releases) page.

### Homebrew Tap

```
brew install FairwindsOps/tap/pluto
```

## QuickStart

First, follow the install instructions to install pluto.

### File Detection in a Directory

Run `pluto detect-files -d <DIRECTORY YOU WANT TO SCAN>`

You should see an output something like:

```
$ pluto detect-files -d pkg/finder/testdata
NAME        KIND         VERSION              REPLACEMENT   REMOVED   DEPRECATED
utilities   Deployment   extensions/v1beta1   apps/v1       true      true
utilities   Deployment   extensions/v1beta1   apps/v1       true      true
```

This indicates that we have two files in our directory that have deprecated apiVersions. This will need to be fixed prior to a 1.16 upgrade.

### Helm Detection (in-cluster)

```
$ pluto detect-helm --helm-version 3 -owide
NAME                                         NAMESPACE               KIND                           VERSION                                REPLACEMENT                       DEPRECATED   DEPRECATED IN   REMOVED   REMOVED IN
cert-manager/cert-manager-webhook            cert-manager            MutatingWebhookConfiguration   admissionregistration.k8s.io/v1beta1   admissionregistration.k8s.io/v1   true         v1.16.0         false     v1.19.0
```

This indicates that the StatefulSet audit-dashboard-prod-rabbitmq-ha was deployed with apps/v1beta1 which is deprecated in 1.16

You can also use Pluto with helm 2:

```
$ pluto detect-helm --helm-version=2 -A
NAME                             KIND         VERSION              REPLACEMENT   REMOVED   DEPRECATED
invincible-zebu-metrics-server   Deployment   extensions/v1beta1   apps/v1       true      true
lunging-bat-metrics-server       Deployment   apps/v1                            false     false
```

In case you have configured helm 2 with `--storage=secret`, you can add the optional `--helm-store secrets`. Default configuration for helm 2 is configmaps.

If you only have access to a certain namespace (helm 3), or if you only want to see information for a single namespace, you can pass the `--namespace` or `-n` flag to restrict the output. Due to the way that helm2 stores configmaps, you will still need access to the tiller namespace if you are using helm 2.

```
$ pluto detect-helm -n cert-manager -owide
NAME                                NAMESPACE      KIND                           VERSION                                REPLACEMENT                       DEPRECATED   DEPRECATED IN   REMOVED   REMOVED IN
cert-manager/cert-manager-webhook   cert-manager   MutatingWebhookConfiguration   admissionregistration.k8s.io/v1beta1   admissionregistration.k8s.io/v1   true         v1.16.0         false     v1.19.0
```

### Helm Chart Checking (local files)

You can run `helm template <chart-dir> | pluto detect -`

This will output something like so:

```
$ helm template e2e/tests/assets/helm3chart | pluto detect -
KIND         VERSION              DEPRECATED   DEPRECATED IN   RESOURCE NAME
Deployment   extensions/v1beta1   true         v1.16.0         RELEASE-NAME-helm3chart-v1beta1
```

## Other Usage Options

### Display Options

In addition to the standard output, Pluto can output yaml, json, or wide.

#### Wide

The wide output provides more information about when an apiVersion was removed or deprecated.

```
$ pluto detect-helm --helm-version 3 -owide
└─ pluto detect-helm -owide
NAME                                         NAMESPACE               KIND                           VERSION                                REPLACEMENT                       DEPRECATED   DEPRECATED IN   REMOVED   REMOVED IN
cert-manager/cert-manager-webhook            cert-manager            MutatingWebhookConfiguration   admissionregistration.k8s.io/v1beta1   admissionregistration.k8s.io/v1   true         v1.16.0         false     v1.19.0
```

#### JSON

```
$ pluto detect-helm --helm-version 3 -ojson | jq .
{
  "items": [
    {
      "name": "cert-manager/cert-manager-webhook",
      "namespace": "cert-manager",
      "api": {
        "version": "admissionregistration.k8s.io/v1beta1",
        "kind": "MutatingWebhookConfiguration",
        "deprecated-in": "v1.16.0",
        "removed-in": "v1.19.0",
        "replacement-api": "admissionregistration.k8s.io/v1",
        "component": "k8s"
      },
      "deprecated": true,
      "removed": false
    }
  ],
  "target-versions": {
    "cert-manager": "v0.15.1",
    "istio": "v1.6.0",
    "k8s": "v1.16.0"
  }
}

```

#### YAML

```yaml
items:
- name: cert-manager/cert-manager-webhook
  namespace: cert-manager
  api:
    version: admissionregistration.k8s.io/v1beta1
    kind: MutatingWebhookConfiguration
    deprecated-in: v1.16.0
    removed-in: v1.19.0
    replacement-api: admissionregistration.k8s.io/v1
    component: k8s
  deprecated: true
  removed: false
target-versions:
  cert-manager: v0.15.1
  istio: v1.6.0
  k8s: v1.16.0
```

### CI Pipelines

Pluto has specific exit codes that is uses to indicate certain results:

- Exit Code 1 - An error. A message will be displayed
- Exit Code 2 - A deprecated apiVersion has been found.
- Exit Code 3 - A removed apiVersion has been found.

If you wish to bypass the generation of exit codes 2 and 3, you may do so with two different flags:

```
--ignore-deprecations              Ignore the default behavior to exit 2 if deprecated apiVersions are found.
--ignore-removals                  Ignore the default behavior to exit 3 if removed apiVersions are found.
```

### Target Versions

By default, Pluto was designed with deprecations related to Kubernetes v1.16.0. However, as more deprecations are introduced, we will try to keep it updated.

You can target the version you are concerned with by using the `--target-versions` or `-t` flag. You must pass the `component=version`, and the version must begin with a `v` (this is a limitation of the semver library we are using to verify).

For example:

```
$ pluto detect-helm --target-version k8s=v1.15.0
No output to display

$ echo $?
0
```

Notice that there is no output, despite the fact that we might have recognized apiVersions present in the cluster that are not yet deprecated or removed in v1.15.0. This particular run exited 0.

### Components

By default Pluto will scan for all components in the versionsList that it can find. If you wish to only see deprecations for a specific component, you can use the `--components` flag to specify a list.

### Adding Custom Version Checks

If you want to check additional apiVersions and/or types, you can pass an additional file with the `--additional-versions` or `-f` flag.

The file should look something like this:

```
target-versions:
  custom: v1.0.0
deprecated-versions:
- version: someother/v1beta1
  kind: AnotherCRD
  deprecated-in: v1.9.0
  removed-in: v1.16.0
  replacement-api: apps/v1
  component: custom
```

You can test that it's working by using `list-versions`:

```
$ pluto list-versions -f new.yaml
KIND                           NAME                                   DEPRECATED IN   REMOVED IN   REPLACEMENT   COMPONENT
AnotherCRD                     someother/v1beta1                      v1.9.0          v1.16.0      apps/v1       custom
```
_NOTE: This output is truncated to show only the additional version. Normally this will include the defaults as well_

The `target-versions` field in this custom file will set the default target version for that component. You can still override this with `--target-versions custom=vX.X.X` when you run Pluto.

Please note that we do not allow overriding anything contained in the default [versions.yaml](versions.yaml) that Pluto uses.
