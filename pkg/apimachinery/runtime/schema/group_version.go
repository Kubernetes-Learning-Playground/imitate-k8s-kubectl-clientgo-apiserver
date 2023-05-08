package schema

// @ See https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/schema/group_version.go

type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

type GroupResource struct {
	Group    string
	Resource string
}

type GroupVersionResource struct {
	Group    string
	Version  string
	Resource string
}
