module github.com/alicefr/kubevirt-hook

go 1.16

replace (
	github.com/go-kit/kit => github.com/go-kit/kit v0.3.0
	github.com/openshift/api => github.com/openshift/api v0.0.0-20191219222812-2987a591a72c
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20210112165513-ebc401615f47
	github.com/operator-framework/operator-lifecycle-manager => github.com/operator-framework/operator-lifecycle-manager v0.0.0-20190128024246-5eb7ae5bdb7a
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/client-go => k8s.io/client-go v0.20.2
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.20.2
	kubevirt.io/client-go => kubevirt.io/client-go v0.44.1
)

require (
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/spf13/pflag v1.0.5
	google.golang.org/grpc v1.40.0
	kubevirt.io/client-go v0.0.0-00010101000000-000000000000
	kubevirt.io/kubevirt v0.44.1
)
