package helm

import (
	testclient "k8s.io/client-go/kubernetes/fake"
)

func getMockConfigInstance() *kube {
	kubeClient = &kube{
		Client: testclient.NewSimpleClientset(),
	}
	return kubeClient
}
