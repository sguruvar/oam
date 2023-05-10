/*
Copyright 2023.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OamTraitSpec defines the desired state of OamTrait
type OamTraitSpec struct {
	Name        string `json:"name"`
	Replicas    int32  `json:"replicas"`
	Image       string `json:"image"`
	Port        int32  `json:"port"`
	Namespace   string `json:"namespace,omitempty"`
	MaxReplicas int32  `json:"max_replicas"`
	CpuTarget   int32  `json:"cpu_target"`
	IngressReq  string `json:"expose_ingress"`
}

// OamTraitStatus defines the observed state of OamTrait
type OamTraitStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OamTrait is the Schema for the oamtraits API
type OamTrait struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OamTraitSpec   `json:"spec,omitempty"`
	Status OamTraitStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OamTraitList contains a list of OamTrait
type OamTraitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OamTrait `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OamTrait{}, &OamTraitList{})
}
