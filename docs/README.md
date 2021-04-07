<div align="center" class="no-border">
  <img src="/img/pluto-logo.png" alt="Pluto Logo">
  <br>
  <h3>Find Kubernetes resources that have been deprecated</h3>
  <a href="https://github.com/FairwindsOps/pluto/releases">
    <img src="https://img.shields.io/github/v/release/FairwindsOps/pluto">
  </a>
  <a href="https://goreportcard.com/report/github.com/FairwindsOps/pluto">
    <img src="https://goreportcard.com/badge/github.com/FairwindsOps/pluto">
  </a>
  <a href="https://circleci.com/gh/FairwindsOps/pluto.svg">
    <img src="https://circleci.com/gh/FairwindsOps/pluto.svg?style=svg">
  </a>
  <a href="https://codecov.io/gh/FairwindsOps/pluto">
    <img src="https://codecov.io/gh/FairwindsOps/pluto/branch/master/graph/badge.svg">
  </a>
</div>

Pluto is a utility to help users find [deprecated Kubernetes apiVersions](https://k8s.io/docs/reference/using-api/deprecation-guide/) in their code repositories and their helm releases.

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
