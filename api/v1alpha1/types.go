package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type KeycloakGroupBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KeycloakGroupBindingSpec   `json:"spec,omitempty"`
	Status            KeycloakGroupBindingStatus `json:"status,omitempty"`
}

func (in *KeycloakGroupBinding) DeepCopyObject() runtime.Object {
	out := new(KeycloakGroupBinding)
	*out = *in
	return out
}

type KeycloakGroupBindingSpec struct {
	KeycloakGroup string      `json:"keycloakGroup"`
	ClusterRole   string      `json:"clusterRole"`
	Namespace     string      `json:"namespace,omitempty"`
	GroupPrefix   string      `json:"groupPrefix,omitempty"`
	KeycloakRef   KeycloakRef `json:"keycloakRef"`
}

type KeycloakRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type KeycloakGroupBindingStatus struct {
	Conditions     []metav1.Condition `json:"conditions,omitempty"`
	BoundGroupName string             `json:"boundGroupName,omitempty"`
	LastSyncTime   *metav1.Time       `json:"lastSyncTime,omitempty"`
}

type KeycloakGroupBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeycloakGroupBinding `json:"items"`
}

func (in *KeycloakGroupBindingList) DeepCopyObject() runtime.Object {
	out := new(KeycloakGroupBindingList)
	*out = *in
	return out
}

func init() {
	SchemeBuilder.Register(&KeycloakGroupBinding{}, &KeycloakGroupBindingList{})
}
