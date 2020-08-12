module github.com/fairwindsops/pluto

go 1.14

require (
	github.com/DATA-DOG/go-sqlmock v1.4.1 // indirect
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/lib/pq v1.3.0 // indirect
	github.com/markbates/pkger v0.17.0
	github.com/rogpeppe/go-internal v1.4.0
	github.com/rubenv/sql-migrate v0.0.0-20200212082348-64f95ea68aa3 // indirect
	github.com/spf13/cobra v0.0.6
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	golang.org/x/mod v0.2.0
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
	helm.sh/helm v2.16.6+incompatible
	helm.sh/helm/v3 v3.1.2
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v0.17.4
	k8s.io/helm v2.16.5+incompatible // indirect
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.5.2
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
)
