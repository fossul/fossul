module github.com/fossul/fossul

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/aws/aws-sdk-go v1.42.10
	github.com/fossul/fossul/src/client/k8s/virtctrl/client/versioned v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/kubernetes-csi/external-snapshotter/client/v4 v4.2.0
	github.com/lib/pq v1.10.4
	github.com/openshift/api v0.0.0-20211028023115-7224b732cc14
	github.com/openshift/client-go v0.0.0-20210831095141-e19a065e79f7
	github.com/pborman/getopt/v2 v2.1.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/swaggo/http-swagger v1.1.2
	github.com/swaggo/swag v1.7.4
	go.mongodb.org/mongo-driver v1.7.4
	k8s.io/api v0.22.4
	k8s.io/apimachinery v0.22.4
	k8s.io/client-go v0.22.4
)

replace (
	github.com/fossul/fossul/src/client/k8s => ./src/client/k8s
	github.com/fossul/fossul/src/client/k8s/virtctrl/client/versioned => ./src/client/k8s/virtctrl/client/versioned
	k8s.io/api => k8s.io/api v0.22.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.22.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.4
	k8s.io/apiserver => k8s.io/apiserver v0.22.4
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.22.4
	k8s.io/client-go => k8s.io/client-go v0.22.4
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.22.4
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.22.4
	k8s.io/code-generator => k8s.io/code-generator v0.22.4
	k8s.io/component-base => k8s.io/component-base v0.22.4
	k8s.io/cri-api => k8s.io/cri-api v0.22.4
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.22.4
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.22.4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.22.4
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20201113171705-d219536bb9fd
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.22.4
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.22.4
	k8s.io/kubectl => k8s.io/kubectl v0.22.4
	k8s.io/kubelet => k8s.io/kubelet v0.22.4
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.22.4
	k8s.io/metrics => k8s.io/metrics v0.22.4
	k8s.io/node-api => k8s.io/node-api v0.22.4
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.22.4
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.22.4
	k8s.io/sample-controller => k8s.io/sample-controller v0.22.4
)
