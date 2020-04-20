package api

// VersionList is a set of apiVersions and if they are deprecated or not.
var VersionList = []Version{
	// Not Removed or Deprecated
	{"apps/v1", "Deployment", "", "", ""},
	{"networking.k8s.io/v1", "NetworkPolicy", "", "", ""},
	{"policy/v1beta1", "PodSecurityPolicy", "", "", ""},

	// Removed in 1.16
	{"extensions/v1beta1", "Deployment", "v1.9.0", "v1.16.0", "apps/v1"},
	{"apps/v1beta2", "Deployment", "v1.9.0", "v1.16.0", "apps/v1"},
	{"apps/v1beta1", "Deployment", "v1.9.0", "v1.16.0", "apps/v1"},

	{"apps/v1beta1", "StatefulSet", "v1.9.0", "v1.16.0", "apps/v1"},
	{"apps/v1beta2", "StatefulSet", "v1.9.0", "v1.16.0", "apps/v1"},

	{"extensions/v1beta1", "NetworkPolicy", "v1.9.0", "v1.16.0", "networking.k8s.io/v1"},

	{"apps/v1beta2", "DaemonSet", "v1.9.0", "v1.16.0", "apps/v1"},
	{"extensions/v1beta1", "DaemonSet", "v1.9.0", "v1.16.0", "apps/v1"},

	{"extensions/v1beta1", "PodSecurityPolicy", "v1.10.0", "v1.16.0", "policy/v1beta1"},

	{"extensions/v1beta1", "ReplicaSet", "", "v1.16.0", "apps/v1"},
	{"apps/v1beta1", "ReplicaSet", "", "v1.16.0", "apps/v1"},
	{"apps/v1beta2", "ReplicaSet", "", "v1.16.0", "apps/v1"},

	// Removed in 1.17
	{"scheduling.k8s.io/v1beta1", "PriorityClass", "v1.14.0", "v1.17.0", "scheduling.k8s.io/v1"},
	{"scheduling.k8s.io/v1alpha1", "PriorityClass", "v1.14.0", "v1.17.0", "scheduling.k8s.io/v1"},
	{"apiextensions.k8s.io/v1beta1", "CustomResourceDefinition", "v1.16.0", "v1.19.0", "apiextensions.k8s.io/v1"},

	// Removed in 1.19
	{"admissionregistration.k8s.io/v1beta1", "MutatingWebhookConfiguration", "v1.16.0", "v1.19.0", "admissionregistration.k8s.io/v1"},

	// Removed in 1.20
	{"rbac.authorization.k8s.io/v1alpha1", "ClusterRoleBinding", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},
	{"rbac.authorization.k8s.io/v1alpha1", "ClusterRole", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},
	{"rbac.authorization.k8s.io/v1alpha1", "ClusterRoleBindingList", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},
	{"rbac.authorization.k8s.io/v1alpha1", "ClusterRoleList", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},
	{"rbac.authorization.k8s.io/v1alpha1", "Role", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},
	{"rbac.authorization.k8s.io/v1alpha1", "RoleBinding", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},
	{"rbac.authorization.k8s.io/v1alpha1", "RoleList", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},
	{"rbac.authorization.k8s.io/v1alpha1", "RoleBindingList", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},

	// Unknown Removal, but deprecated
	{"storage.k8s.io/v1beta1", "CSINode", "v1.17.0", "", ""},
}
