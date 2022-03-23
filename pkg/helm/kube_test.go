package helm

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	testclient "k8s.io/client-go/kubernetes/fake"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func getMockConfigInstance() *kube {
	kubeClient = &kube{
		Client: testclient.NewSimpleClientset(),
	}
	return kubeClient
}

func Test_getKubeClient(t *testing.T) {
	tests := []struct {
		name        string
		kubeContext string
		kubeConfig  string
		wantErr     bool
	}{
		{
			name:        "context does not exist",
			kubeContext: "farglebargle",
			kubeConfig:  "testdata/kubeconfig",
			wantErr:     true,
		},
		{
			name:        "context exists",
			kubeContext: "kind-kind",
			kubeConfig:  "testdata/kubeconfig",
			wantErr:     false,
		},
		{
			name:        "invalid kubeconfig",
			kubeContext: "kind-kind",
			kubeConfig:  "testdata/kubeconfig_invalid",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("KUBECONFIG", tt.kubeConfig)
			_, err := getKubeClient(tt.kubeContext)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
