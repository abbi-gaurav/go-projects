package main

import (
	"encoding/json"
	"reflect"
	"testing"

	jsonpatch "github.com/evanphx/json-patch"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestJSONPatchForConfigMap(t *testing.T) {
	cm := corev1.ConfigMap{
		Data: map[string]string{
			"mutation-start": "yes",
		},
	}
	cmJS, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	patchObj, err := jsonpatch.DecodePatch([]byte(patch1))
	if err != nil {
		t.Fatal(err)
	}
	patchedJS, err := patchObj.Apply(cmJS)
	patchedObj := corev1.ConfigMap{}
	err = json.Unmarshal(patchedJS, &patchedObj)
	if err != nil {
		t.Fatal(err)
	}
	expected := corev1.ConfigMap{
		Data: map[string]string{
			"mutation-start":   "yes",
			"mutation-stage-1": "yes",
		},
	}

	if !reflect.DeepEqual(patchedObj, expected) {
		t.Errorf("\nexpected %#v\n, got %#v", expected, patchedObj)
	}
}

func TestJSONPatchForUnstructured(t *testing.T) {
	cr := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Something",
			"apiVersion": "somegroup/v1",
			"data": map[string]interface{}{
				"mutation-start": "yes",
			},
		},
	}
	crJS, err := json.Marshal(cr)
	if err != nil {
		t.Fatal(err)
	}

	patchObj, err := jsonpatch.DecodePatch([]byte(patch1))
	if err != nil {
		t.Fatal(err)
	}
	patchedJS, err := patchObj.Apply(crJS)
	patchedObj := unstructured.Unstructured{}
	err = json.Unmarshal(patchedJS, &patchedObj)
	if err != nil {
		t.Fatal(err)
	}
	expectedData := map[string]interface{}{
		"mutation-start":   "yes",
		"mutation-stage-1": "yes",
	}

	if !reflect.DeepEqual(patchedObj.Object["data"], expectedData) {
		t.Errorf("\nexpected %#v\n, got %#v", expectedData, patchedObj.Object["data"])
	}
}
