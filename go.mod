module github.com/fairwindsops/pluto/v3

go 1.16

require (
	github.com/gobuffalo/here v0.6.2 // indirect
	github.com/markbates/pkger v0.17.1
	github.com/mattn/go-runewidth v0.0.12 // indirect
	github.com/olekukonko/tablewriter v0.0.5
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.4.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	golang.org/x/mod v0.3.0
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	helm.sh/helm/v3 v3.5.1
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.8.3
)

replace github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
