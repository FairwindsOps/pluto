---
meta:
  - name: description
    content: "Fairwinds Pluto | Documentation FAQ"
---
## Frequently Asked Questions

### I updated my deployment method to use the new API version and Pluto doesn't report anything but kubectl still shows the old API. What gives?

See above in the [Purpose](/#purpose) section of this doc. Kubectl is likely lying to you because it only tells you what the default is for the given kubernetes version even if an object was deployed with a newer API version.

### Why doesn't Pluto check the last-applied-configuration annotation?

If you see the annotation `kubectl.kubernetes.io/last-applied-configuration` on an object in your cluster it means that object was updated with `kubectl apply`. We don't consider this an entirely reliable solution for checking. In fact, others have pointed out that updating the same object with `kubectl patch` will **remove** the annotation. Due to the flaky behavior here, we will not plan on supporting this.

### I don't use helm, how can I do in cluster checks?

Currently, the only in-cluster check we are confident in supporting is helm. If your deployment method can generate yaml manifests for kubernetes, you should be able to use the `detect` or `detect-files` functionality described below after the manifest files have been generated.

### I updated the API version of an object, but pluto still reports that the apiVersion needs to be updated.

Pluto looks at the API Versions of objects in releases that are in a `Deployed` state, and Helm has an issue where it might list old revisions of a release as still being in a `Deployed` state. To fix this, look at the release revision history with `helm history <release name>`, and determine if older releases still show a `Deployed` state. If so, delete the Helm release secret(s) associated with the revision number(s). For example, `kubectl delete secret sh.helm.release.v1.my-release.v10` where `10` corresponds to the release number. Then run Pluto again to see if the object has been removed from the report.