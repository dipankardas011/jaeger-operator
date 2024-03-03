/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v2alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JaegerSpec defines the desired state of Jaeger
type JaegerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Jaeger. Edit jaeger_types.go to remove/update
	Config string `json:"config,omitempty"`
}

// JaegerStatus defines the observed state of Jaeger
type JaegerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Version string `json:"version"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Jaeger is the Schema for the jaegers API
type Jaeger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JaegerSpec   `json:"spec,omitempty"`
	Status JaegerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// JaegerList contains a list of Jaeger
type JaegerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Jaeger `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Jaeger{}, &JaegerList{})
}
