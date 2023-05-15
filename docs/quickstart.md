---
meta:
  - name: description
    content: "Fairwinds Pluto | Quickstart Documentation"
---
# QuickStart

First, follow the install instructions to install pluto.

## File Detection in a Directory

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
$ pluto detect-helm -owide
NAME                                         NAMESPACE               KIND                           VERSION                                REPLACEMENT                       DEPRECATED   DEPRECATED IN   REMOVED   REMOVED IN
cert-manager/cert-manager-webhook            cert-manager            MutatingWebhookConfiguration   admissionregistration.k8s.io/v1beta1   admissionregistration.k8s.io/v1   true         v1.16.0         false     v1.19.0
```

This indicates that the StatefulSet audit-dashboard-prod-rabbitmq-ha was deployed with apps/v1beta1 which is deprecated in 1.16


If you want to see information for a single namespace, you can pass the `--namespace` or `-n` flag to restrict the output.

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

### API resources (in-cluster)
```
$ pluto detect-api-resources -owide
NAME                  NAMESPACE     KIND                VERSION          REPLACEMENT   DEPRECATED   DEPRECATED IN   REMOVED   REMOVED IN     
psp                   <UNKNOWN>     PodSecurityPolicy   policy/v1beta1                 true         v1.21.0         false     v1.25.0 
```

This indicates that the PodSecurityPolicy  was deployed with apps/v1beta1 which is deprecated in 1.21

### helm and API resources (in-cluster)

```
$ pluto detect-all-in-cluster -o wide 2>/dev/null
NAME              NAMESPACE   KIND                VERSION                     REPLACEMENT            DEPRECATED   DEPRECATED IN   REMOVED   REMOVED IN  
testing/viahelm   viahelm     Ingress             networking.k8s.io/v1beta1   networking.k8s.io/v1   true         v1.19.0         true      v1.22.0     
webapp            default     Ingress             networking.k8s.io/v1beta1   networking.k8s.io/v1   true         v1.19.0         true      v1.22.0     
eks.privileged    <UNKNOWN>   PodSecurityPolicy   policy/v1beta1                                     true         v1.21.0         false     v1.25.0     
```

This combines all available in-cluster detections, showing results from Helm releases and API resources.
