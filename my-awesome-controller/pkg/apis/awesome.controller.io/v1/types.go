package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ADDED   string = "added"
	UPDATED string = "updated"
	DELETED string = "deleted"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Cake struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CakeSpec   `json:"spec"`
	Status            CakeStatus `json:"status"`
}

type CakeSpec struct {
	Type string `json:"type"`
}

type CakeStatus struct {
	State string `json:"cooked"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList is a list of Foo resources
type CakeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Cake `json:"items"`
}
