package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LabelSpec defines the desired state of Label
type LabelSpec struct {
	Projects []string `json:"projects"`
}

// LabelStatus defines the observed state of Label
type LabelStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Label is the Schema for the labels API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=labels,scope=Namespaced
// +genclient:nonNamespaced
type Label struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LabelSpec   `json:"spec,omitempty"`
	Status LabelStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LabelList contains a list of Label
type LabelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Label `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Label{}, &LabelList{})
}
