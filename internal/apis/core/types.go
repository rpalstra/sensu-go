package core

import metav1 "github.com/sensu/sensu-go/apis/meta/v1"

// A Namespace is a resource that defines where other resources are located.
type Namespace struct {
	metav1.TypeMeta   `json:",inline" protobuf:"bytes,1,opt,name=typeMeta"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,2,opt,name=metadata"`
}

type Subscription struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:",inline"`

	CheckRefs []CheckRef `json:"checkRefs"`
}

type CheckRef struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type Check struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec CheckSpec
}

type CheckSpec struct {
	Command  string       `json:"command"`
	Handlers []HandlerRef `json:"handlers,omitempty"`
	Assets   []AssetRef   `json:"assets,omitempty"`

	Interval uint32     `json:"interval"`
	Cron     string     `json:"cron"`
	Publish  bool       `json:"active"`
	Subdue   SubdueSpec `json:"subdue"`

	LowFlapThreshold  uint32 `json:"lowFlapThreshold"`
	HighFlapThreshold uint32 `json:"highFlapThreshold"`

	Subscriptions []string `json:"subscriptions"`

	Stdin bool `json:"stdin"`
}

type SubdueSpec struct {
	Start string   `json:"start"`
	End   string   `json:"end"`
	Days  []string `json:"days"`
}

type HandlerRef struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type AssetRef struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}
