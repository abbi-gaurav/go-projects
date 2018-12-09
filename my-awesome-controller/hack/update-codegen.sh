#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

echo "$SCRIPT_ROOT"
echo "$CODEGEN_PKG"

${CODEGEN_PKG}/generate-groups.sh "all" \
	      github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client \
	      github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/apis \
	      awesome.controller.io:v1 \
	      --go-header-file ${SCRIPT_ROOT}/hack/boilerplate.go.txt
