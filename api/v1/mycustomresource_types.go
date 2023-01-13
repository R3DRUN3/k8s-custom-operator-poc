package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MyCustomResourceSpec defines the desired state of MyCustomResource
type MyCustomResourceSpec struct {
	// Name of the friend MyCustomResource is looking for
	Name string `json:"name"`
}

// MyCustomResourceStatus defines the observed state of MyCustomResource
type MyCustomResourceStatus struct {
	// Healthy will be set to true if MyCustomResource found a friend
	Healthy bool `json:"Healthy,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MyCustomResource is the Schema for the MyCustomResources API
type MyCustomResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyCustomResourceSpec   `json:"spec,omitempty"`
	Status MyCustomResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MyCustomResourceList contains a list of MyCustomResource
type MyCustomResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyCustomResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyCustomResource{}, &MyCustomResourceList{})
}
