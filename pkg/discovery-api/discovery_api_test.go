// Copyright 2022 FairwindsOps Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

package discoveryapi

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	discoveryFake "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/dynamic/fake"
)

func TestNewDiscoveryAPIClientValidEmpty(t *testing.T) {

	scheme := runtime.NewScheme()
	clientset := fake.NewSimpleDynamicClient(scheme)
	discoveryClient := discoveryFake.FakeDiscovery{}
	testOpts := DiscoveryClient{
		ClientSet:       clientset,
		DiscoveryClient: &discoveryClient,
	}

	err := testOpts.GetApiResources()
	if err != nil {
		t.Errorf("Unable to fetch API Resources")
	}

}
