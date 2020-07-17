module github.com/abbi-gaurav/go-projects/sample-broker

go 1.14

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/drewolson/testflight v1.0.0 // indirect
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/kyma-incubator/api-gateway v0.0.0-20200619075331-7f6fd0cfdac9
	github.com/onsi/ginkgo v1.14.0 // indirect
	github.com/ory/oathkeeper-maester v0.0.2-beta.1
	github.com/pivotal-cf/brokerapi v6.4.2+incompatible
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/utils v0.0.0-20200716102541-988ee3149bb2 // indirect
	sigs.k8s.io/controller-runtime v0.6.1
)

replace k8s.io/client-go => k8s.io/client-go v0.18.6 // Required by prometheus-operator
