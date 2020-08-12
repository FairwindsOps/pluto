module github.com/fairwindsops/pluto

go 1.14

require (
	github.com/markbates/pkger v0.17.0
	github.com/rogpeppe/go-internal v1.6.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	golang.org/x/mod v0.3.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	helm.sh/helm v2.16.9+incompatible
	helm.sh/helm/v3 v3.3.0
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/helm v2.16.9+incompatible // indirect
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.6.2
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
)
