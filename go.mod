module github.com/oracle/coherence-operator

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/coreos/prometheus-operator v0.31.1 // indirect
	github.com/docker/docker v0.0.0-20180612054059-a9fbbdc8dd87 // indirect
	github.com/elastic/go-elasticsearch/v6 v6.8.3-0.20190731061920-efbed2e4c2f8
	github.com/emicklei/go-restful v2.9.3+incompatible // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.1.0
	github.com/go-openapi/spec v0.19.0
	github.com/go-openapi/validate v0.18.0
	github.com/go-test/deep v1.0.3
	github.com/google/btree v1.0.0 // indirect
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/gregjones/httpcache v0.0.0-20190212212710-3befbb6ad0cc // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/operator-framework/operator-sdk v0.11.1-0.20191014155558-888dde512025
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/spf13/pflag v1.0.3
	github.com/tebeka/go2xunit v1.4.10
	github.com/ugorji/go/codec v0.0.0-20190320090025-2dc34c0b8780 // indirect
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092
	google.golang.org/appengine v1.5.0 // indirect
	gotest.tools v2.2.0+incompatible // indirect
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apiextensions-apiserver v0.0.0-20190918161926-8f644eb6e783
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/apiserver v0.0.0-20181213151703-3ccfe8365421 // indirect
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/helm v2.14.1+incompatible
	k8s.io/kube-openapi v0.0.0-20190603182131-db7b694dc208
	k8s.io/kubernetes v1.14.2
	k8s.io/utils v0.0.0-20190506122338-8fab8cb257d5
	sigs.k8s.io/controller-runtime v0.2.0
	sigs.k8s.io/yaml v1.1.0

)

// Pinned to kubernetes-1.14.1
replace (
	k8s.io/api => k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190409022649-727a075fdec8
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190409021813-1ec86e4da56c
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190409023024-d644b00f3b79
	k8s.io/client-go => k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190409023720-1bc0c81fa51d
	k8s.io/kubernetes => k8s.io/kubernetes v1.14.1
)

replace (
	// Indirect operator-sdk dependencies use git.apache.org, which is frequently
	// down. The github mirror should be used instead.
	// Locking to a specific version (from 'go mod graph'):
	git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
	github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.31.1
	// Pinned to v2.10.0 (kubernetes-1.14.1) so https://proxy.golang.org can
	// resolve it correctly.
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v0.0.0-20190525122359-d20e84d0fb64
)

replace github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.11.0
