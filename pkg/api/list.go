// Copyright 2020 Fairwinds
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package api

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"gopkg.in/yaml.v2"
)

// VersionList is a set of apiVersions and if they are deprecated or not.
var VersionList = []Version{
	// Not Removed or Deprecated
	{"apps/v1", "Deployment", "", "", "", "k8s"},
	{"networking.k8s.io/v1", "NetworkPolicy", "", "", "", "k8s"},
	{"policy/v1beta1", "PodSecurityPolicy", "", "", "", "k8s"},

	// Removed in 1.16
	{"extensions/v1beta1", "Deployment", "v1.9.0", "v1.16.0", "apps/v1", "k8s"},
	{"apps/v1beta2", "Deployment", "v1.9.0", "v1.16.0", "apps/v1", "k8s"},
	{"apps/v1beta1", "Deployment", "v1.9.0", "v1.16.0", "apps/v1", "k8s"},

	{"apps/v1beta1", "StatefulSet", "v1.9.0", "v1.16.0", "apps/v1", "k8s"},
	{"apps/v1beta2", "StatefulSet", "v1.9.0", "v1.16.0", "apps/v1", "k8s"},

	{"extensions/v1beta1", "NetworkPolicy", "v1.9.0", "v1.16.0", "networking.k8s.io/v1", "k8s"},
	{"extensions/v1beta1", "Ingress", "v1.14.0", "v1.22.0", "networking.k8s.io/v1beta1", "k8s"},

	{"apps/v1beta2", "DaemonSet", "v1.9.0", "v1.16.0", "apps/v1", "k8s"},
	{"extensions/v1beta1", "DaemonSet", "v1.9.0", "v1.16.0", "apps/v1", "k8s"},

	{"extensions/v1beta1", "PodSecurityPolicy", "v1.10.0", "v1.16.0", "policy/v1beta1", "k8s"},

	{"extensions/v1beta1", "ReplicaSet", "", "v1.16.0", "apps/v1", "k8s"},
	{"apps/v1beta1", "ReplicaSet", "", "v1.16.0", "apps/v1", "k8s"},
	{"apps/v1beta2", "ReplicaSet", "", "v1.16.0", "apps/v1", "k8s"},

	// Removed in 1.17
	{"scheduling.k8s.io/v1beta1", "PriorityClass", "v1.14.0", "v1.17.0", "scheduling.k8s.io/v1", "k8s"},
	{"scheduling.k8s.io/v1alpha1", "PriorityClass", "v1.14.0", "v1.17.0", "scheduling.k8s.io/v1", "k8s"},
	{"apiextensions.k8s.io/v1beta1", "CustomResourceDefinition", "v1.16.0", "v1.19.0", "apiextensions.k8s.io/v1", "k8s"},

	// Removed in 1.19
	{"admissionregistration.k8s.io/v1beta1", "MutatingWebhookConfiguration", "v1.16.0", "v1.19.0", "admissionregistration.k8s.io/v1", "k8s"},

	// Removed in 1.20
	{"rbac.authorization.k8s.io/v1alpha1", "ClusterRoleBinding", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1", "k8s"},
	{"rbac.authorization.k8s.io/v1alpha1", "ClusterRole", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1", "k8s"},
	{"rbac.authorization.k8s.io/v1alpha1", "ClusterRoleBindingList", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1", "k8s"},
	{"rbac.authorization.k8s.io/v1alpha1", "ClusterRoleList", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1", "k8s"},
	{"rbac.authorization.k8s.io/v1alpha1", "Role", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1", "k8s"},
	{"rbac.authorization.k8s.io/v1alpha1", "RoleBinding", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1", "k8s"},
	{"rbac.authorization.k8s.io/v1alpha1", "RoleList", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1", "k8s"},
	{"rbac.authorization.k8s.io/v1alpha1", "RoleBindingList", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1", "k8s"},

	// Unknown Removal, but deprecated
	{"storage.k8s.io/v1beta1", "CSINode", "v1.17.0", "", "", "k8s"},

	// Istio API versions
	{"rbac.istio.io", "AuthorizationPolicies", "v1.4.0", "v1.4.0", "security.istio.io/v1beta1", "istio"},
	{"authentication.istio.io/v1alpha1", "", "v1.5.0", "v1.6.0", "security.istio.io/v1beta1", "istio"},
	{"networking.istio.io/v1alpha3", "", "v1.5.0", "", "networking.istio.io/v1beta1", "istio"},
}

// PrintVersionList prints out the list of versions
// in a specific format
func PrintVersionList(outputFormat string) error {
	switch outputFormat {
	case "normal":
		err := printVersionsTabular()
		if err != nil {
			return err
		}
	case "wide":
		err := printVersionsTabular()
		if err != nil {
			return err
		}
	case "json":
		data, err := json.Marshal(VersionList)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml":
		data, err := yaml.Marshal(VersionList)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	default:
		errText := "The output format must one of (normal|json|yaml)"
		fmt.Println(errText)
		return fmt.Errorf(errText)
	}
	return nil
}

func printVersionsTabular() error {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 15, 2, padChar, 0)

	_, _ = fmt.Fprintln(w, "KIND\t NAME\t DEPRECATED IN\t REMOVED IN\t REPLACEMENT\t")

	for _, version := range VersionList {
		deprecatedIn := version.DeprecatedIn
		if deprecatedIn == "" {
			deprecatedIn = "n/a"
		}
		removedIn := version.RemovedIn
		if removedIn == "" {
			removedIn = "n/a"
		}

		replacementAPI := version.ReplacementAPI
		if replacementAPI == "" {
			replacementAPI = "n/a"
		}

		_, _ = fmt.Fprintf(w, "%s\t %s\t %s\t %s\t %s\t\n", version.Kind, version.Name, deprecatedIn, removedIn, replacementAPI)
	}
	err := w.Flush()
	if err != nil {
		return err
	}
	return nil
}
