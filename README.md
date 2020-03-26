[![CircleCI](https://circleci.com/gh/FairwindsOps/pluto.svg?style=svg)](https://circleci.com/gh/FairwindsOps/pluto)[![codecov](https://codecov.io/gh/FairwindsOps/pluto/branch/master/graph/badge.svg?token=A23F79JTNA)](https://codecov.io/gh/FairwindsOps/pluto)

# pluto

This is a very simple utility to help users find deprecated Kubernetes apiVersions in their code repositories and their helm releases.

**Want to learn more?** Reach out on [the Slack channel](https://fairwindscommunity.slack.com/messages/goldilocks), send an email to `opensource@fairwinds.com`, or join us for [office hours on Zoom](https://fairwindscommunity.slack.com/messages/office-hours)

## QuickStart

Install the binary from our [releases](https://github.com/FairwindsOps/pluto/releases) page.

### File Detection in a Directory

Run `pluto detect-files -d <DIRECTORY YOU WANT TO SCAN>`

You should see an output something like:

```
$ pluto detect-files -d pkg/finder/testdata
KIND         VERSION              DEPRECATED   FILE
Deployment   extensions/v1beta1   true         pkg/finder/testdata/deployment-extensions-v1beta1.json
Deployment   extensions/v1beta1   true         pkg/finder/testdata/deployment-extensions-v1beta1.yaml
```

This indicates that we have two files in our directory that have deprecated apiVersions. This will need to be fixed prior to a 1.16 upgrade.

### Helm Detection (in-cluster)

```
$ pluto detect-helm --helm-version 3
KIND          VERSION        DEPRECATED   RESOURCE NAME
StatefulSet   apps/v1beta1   true         audit-dashboard-prod-rabbitmq-ha
```

This indicates that the StatefulSet audit-dashboard-prod-rabbitmq-ha was deployed with apps/v1beta1 which is deprecated in 1.16

### Helm Chart Checking (local files)

You can run `helm template <chart-dir> | pluto detect --show-non-deprecated -`

This will output something like so:

```
KIND         VERSION   DEPRECATED   RESOURCE NAME
Deployment   apps/v1   false        RELEASE-NAME-goldilocks-controller
Deployment   apps/v1   false        RELEASE-NAME-goldilocks-dashboard
```
