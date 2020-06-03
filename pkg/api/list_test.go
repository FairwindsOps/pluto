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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/mod/semver"
)

func Test_VersionListIsValid(t *testing.T) {
	// This test validates that all of the versions in VersionList are valid semVer
	// it should prevent us from putting bad values in that list in future development
	for _, version := range VersionList {
		if version.DeprecatedIn != "" {
			assert.True(t, semver.IsValid(version.DeprecatedIn), fmt.Sprintf("version %s is not valid semver", version.DeprecatedIn))
		}
		if version.RemovedIn != "" {
			assert.True(t, semver.IsValid(version.RemovedIn), fmt.Sprintf("version %s is not valid semver", version.RemovedIn))
		}
	}
}

var TestVersionList = []Version{
	{"apps/v1beta1", "Deployment", "v1.9.0", "v1.16.0", "apps/v1", "foo"},
}

func Example_printVersionsTabular() {
	padChar = '-'
	_ = printVersionsTabular()

	// Output:
	// KIND-------------------------- NAME---------------------------------- DEPRECATED IN-- REMOVED IN-- REPLACEMENT----------------------
	// Deployment-------------------- apps/v1------------------------------- n/a------------ n/a--------- n/a------------------------------
	// NetworkPolicy----------------- networking.k8s.io/v1------------------ n/a------------ n/a--------- n/a------------------------------
	// PodSecurityPolicy------------- policy/v1beta1------------------------ n/a------------ n/a--------- n/a------------------------------
	// Deployment-------------------- extensions/v1beta1-------------------- v1.9.0--------- v1.16.0----- apps/v1--------------------------
	// Deployment-------------------- apps/v1beta2-------------------------- v1.9.0--------- v1.16.0----- apps/v1--------------------------
	// Deployment-------------------- apps/v1beta1-------------------------- v1.9.0--------- v1.16.0----- apps/v1--------------------------
	// StatefulSet------------------- apps/v1beta1-------------------------- v1.9.0--------- v1.16.0----- apps/v1--------------------------
	// StatefulSet------------------- apps/v1beta2-------------------------- v1.9.0--------- v1.16.0----- apps/v1--------------------------
	// NetworkPolicy----------------- extensions/v1beta1-------------------- v1.9.0--------- v1.16.0----- networking.k8s.io/v1-------------
	// Ingress----------------------- extensions/v1beta1-------------------- v1.14.0-------- v1.22.0----- networking.k8s.io/v1beta1--------
	// DaemonSet--------------------- apps/v1beta2-------------------------- v1.9.0--------- v1.16.0----- apps/v1--------------------------
	// DaemonSet--------------------- extensions/v1beta1-------------------- v1.9.0--------- v1.16.0----- apps/v1--------------------------
	// PodSecurityPolicy------------- extensions/v1beta1-------------------- v1.10.0-------- v1.16.0----- policy/v1beta1-------------------
	// ReplicaSet-------------------- extensions/v1beta1-------------------- n/a------------ v1.16.0----- apps/v1--------------------------
	// ReplicaSet-------------------- apps/v1beta1-------------------------- n/a------------ v1.16.0----- apps/v1--------------------------
	// ReplicaSet-------------------- apps/v1beta2-------------------------- n/a------------ v1.16.0----- apps/v1--------------------------
	// PriorityClass----------------- scheduling.k8s.io/v1beta1------------- v1.14.0-------- v1.17.0----- scheduling.k8s.io/v1-------------
	// PriorityClass----------------- scheduling.k8s.io/v1alpha1------------ v1.14.0-------- v1.17.0----- scheduling.k8s.io/v1-------------
	// CustomResourceDefinition------ apiextensions.k8s.io/v1beta1---------- v1.16.0-------- v1.19.0----- apiextensions.k8s.io/v1----------
	// MutatingWebhookConfiguration-- admissionregistration.k8s.io/v1beta1-- v1.16.0-------- v1.19.0----- admissionregistration.k8s.io/v1--
	// ClusterRoleBinding------------ rbac.authorization.k8s.io/v1alpha1---- v1.17.0-------- v1.20.0----- rbac.authorization.k8s.io/v1-----
	// ClusterRole------------------- rbac.authorization.k8s.io/v1alpha1---- v1.17.0-------- v1.20.0----- rbac.authorization.k8s.io/v1-----
	// ClusterRoleBindingList-------- rbac.authorization.k8s.io/v1alpha1---- v1.17.0-------- v1.20.0----- rbac.authorization.k8s.io/v1-----
	// ClusterRoleList--------------- rbac.authorization.k8s.io/v1alpha1---- v1.17.0-------- v1.20.0----- rbac.authorization.k8s.io/v1-----
	// Role-------------------------- rbac.authorization.k8s.io/v1alpha1---- v1.17.0-------- v1.20.0----- rbac.authorization.k8s.io/v1-----
	// RoleBinding------------------- rbac.authorization.k8s.io/v1alpha1---- v1.17.0-------- v1.20.0----- rbac.authorization.k8s.io/v1-----
	// RoleList---------------------- rbac.authorization.k8s.io/v1alpha1---- v1.17.0-------- v1.20.0----- rbac.authorization.k8s.io/v1-----
	// RoleBindingList--------------- rbac.authorization.k8s.io/v1alpha1---- v1.17.0-------- v1.20.0----- rbac.authorization.k8s.io/v1-----
	// CSINode----------------------- storage.k8s.io/v1beta1---------------- v1.17.0-------- n/a--------- n/a------------------------------
	// AuthorizationPolicies--------- rbac.istio.io------------------------- v1.4.0--------- v1.4.0------ security.istio.io/v1beta1--------
	// ------------------------------ authentication.istio.io/v1alpha1------ v1.5.0--------- v1.6.0------ security.istio.io/v1beta1--------
	// ------------------------------ networking.istio.io/v1alpha3---------- v1.5.0--------- n/a--------- networking.istio.io/v1beta1------
	// Certificate------------------- certmanager.k8s.io/v1alpha1----------- v0.11.0-------- v0.11.0----- cert-manager.io/v1alpha2---------
	// Issuer------------------------ certmanager.k8s.io/v1alpha1----------- v0.11.0-------- v0.11.0----- cert-manager.io/v1alpha2---------
	// ClusterIssuer----------------- certmanager.k8s.io/v1alpha1----------- v0.11.0-------- v0.11.0----- cert-manager.io/v1alpha2---------
	// CertificateRequest------------ certmanager.k8s.io/v1alpha1----------- v0.11.0-------- v0.11.0----- cert-manager.io/v1alpha2---------
	// Order------------------------- certmanager.k8s.io/v1alpha1----------- v0.11.0-------- v0.11.0----- cert-manager.io/v1alpha2---------
	// Challenge--------------------- certmanager.k8s.io/v1alpha1----------- v0.11.0-------- v0.11.0----- cert-manager.io/v1alpha2---------
}

func ExamplePrintVersionList_json() {
	VersionList = TestVersionList
	_ = PrintVersionList("json")

	// Output:
	// [{"version":"apps/v1beta1","kind":"Deployment","deprecated-in":"v1.9.0","removed-in":"v1.16.0","replacement-api":"apps/v1"}]
}

func ExamplePrintVersionList_yaml() {
	VersionList = TestVersionList
	_ = PrintVersionList("yaml")

	// Output:
	// - version: apps/v1beta1
	//   kind: Deployment
	//   deprecated-in: v1.9.0
	//   removed-in: v1.16.0
	//   replacement-api: apps/v1
}

func ExamplePrintVersionList_normal() {
	VersionList = TestVersionList
	_ = PrintVersionList("normal")

	// Output:
	// KIND-------- NAME---------- DEPRECATED IN-- REMOVED IN-- REPLACEMENT--
	// Deployment-- apps/v1beta1-- v1.9.0--------- v1.16.0----- apps/v1------
}

func ExamplePrintVersionList_wide() {
	VersionList = TestVersionList
	_ = PrintVersionList("wide")

	// Output:
	// KIND-------- NAME---------- DEPRECATED IN-- REMOVED IN-- REPLACEMENT--
	// Deployment-- apps/v1beta1-- v1.9.0--------- v1.16.0----- apps/v1------
}

func ExamplePrintVersionList_badformat() {
	VersionList = TestVersionList
	_ = PrintVersionList("foo")

	// Output:
	// The output format must one of (normal|json|yaml)
}
