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
  <a href="https://insights.fairwinds.com/gh/FairwindsOps/pluto">
    <img src="https://insights.fairwinds.com/v0/gh/FairwindsOps/pluto/badge.svg">
  </a>
</div>

Pluto is a utility to help users find [deprecated Kubernetes apiVersions](https://k8s.io/docs/reference/using-api/deprecation-guide/) in their code repositories and their helm releases.

## Documentation
Check out the [documentation at docs.fairwinds.com](https://pluto.docs.fairwinds.com)

## Purpose

Kubernetes sometimes deprecates apiVersions. Most notably, a large number of deprecations happened in the [1.16 release](https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/). This is fine, and it's a fairly easy thing to deal with. However, it can be difficult to find all the places where you might have used a version that will be deprecated in your next upgrade.

You might think, "I'll just ask the api-server to tell me!", but this is fraught with danger. If you ask the api-server to give you `deployments.v1.apps`, and the deployment was deployed as `deployments.v1beta1.extensions`, the api-server will quite happily convert the api version and return a manifest with `apps/v1`. This is fairly well outlined in the discussion in [this issue](https://github.com/kubernetes/kubernetes/issues/58131#issuecomment-356823588).

So, long story short, finding the places where you have deployed a deprecated apiVersion can be challenging. This is where `pluto` comes in. You can use pluto to check a couple different places where you might have placed a deprecated version:
* Infrastructure-as-Code repos: Pluto can check both static manifests and Helm charts for deprecated apiVersions
* Live Helm releases: Pluto can check both Helm 2 and Helm 3 releases running in your cluster for deprecated apiVersions

## Kubernetes Deprecation Policy

You can read the full policy [here](https://kubernetes.io/docs/reference/using-api/deprecation-policy/)

Long story short, apiVersions get deprecated, and then they eventually get removed entirely. Pluto differentiates between these two, and will tell you if a version is `DEPRECATED` or `REMOVED`

## GitHub Action Usage
Want to use pluto within your GitHub workflows?

```yaml
- name: Download Pluto
  uses: FairwindsOps/pluto/github-action@master

- name: Use pluto
  run: |
    pluto detect-files -d pkg/finder/testdata
```

<!-- Begin boilerplate -->
## Join the Fairwinds Open Source Community

The goal of the Fairwinds Community is to exchange ideas, influence the open source roadmap,
and network with fellow Kubernetes users.
[Chat with us on Slack](https://join.slack.com/t/fairwindscommunity/shared_invite/zt-e3c6vj4l-3lIH6dvKqzWII5fSSFDi1g)
or
[join the user group](https://www.fairwinds.com/open-source-software-user-group) to get involved!

<a href="https://www.fairwinds.com/t-shirt-offer?utm_source=pluto&utm_medium=pluto&utm_campaign=pluto-tshirt">
  <img src="https://www.fairwinds.com/hubfs/Doc_Banners/Fairwinds_OSS_User_Group_740x125_v6.png" alt="Love Fairwinds Open Source? Share your business email and job title and we'll send you a free Fairwinds t-shirt!" />
</a>

## Other Projects from Fairwinds

Enjoying Pluto? Check out some of our other projects:
* [Polaris](https://github.com/FairwindsOps/Polaris) - Audit, enforce, and build policies for Kubernetes resources, including over 20 built-in checks for best practices
* [Goldilocks](https://github.com/FairwindsOps/Goldilocks) - Right-size your Kubernetes Deployments by compare your memory and CPU settings against actual usage
* [Nova](https://github.com/FairwindsOps/Nova) - Check to see if any of your Helm charts have updates available
* [rbac-manager](https://github.com/FairwindsOps/rbac-manager) - Simplify the management of RBAC in your Kubernetes clusters

Or [check out the full list](https://www.fairwinds.com/open-source-software?utm_source=pluto&utm_medium=pluto&utm_campaign=pluto)
## Fairwinds Insights
If you're interested in running Pluto in multiple clusters,
tracking the results over time, integrating with Slack, Datadog, and Jira,
or unlocking other functionality, check out
[Fairwinds Insights](https://fairwinds.com/pricing),
a platform for auditing and enforcing policy in Kubernetes clusters.

<a href="https://fairwinds.com/pricing">
  <img src="https://www.fairwinds.com/hubfs/Doc_Banners/Fairwinds_Pluto_Ad.png" alt="Fairwinds Insights" />
</a>
