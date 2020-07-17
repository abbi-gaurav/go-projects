package middleware

import (
	"code.cloudfoundry.org/lager"
	"context"
	"github.com/abbi-gaurav/go-projects/sample-broker/internal/model"
	apiRules "github.com/kyma-incubator/api-gateway/api/v1alpha1"
	rulev1alpha1 "github.com/ory/oathkeeper-maester/api/v1alpha1"
	"github.com/pivotal-cf/brokerapi/domain"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/deprecated/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type Service struct {
	k8sClient client.Client
	logger    lager.Logger
}

func New(logger lager.Logger) (*Service, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	sch, err := createScheme()

	if err != nil {
		return nil, err
	}

	k8sClient, err := client.New(cfg, client.Options{Scheme: sch})
	if err != nil {
		return nil, err
	}

	return &Service{
		k8sClient: k8sClient,
		logger:    logger,
	}, nil
}

func createScheme() (*runtime.Scheme, error) {
	sch := scheme.Scheme
	var addToSchemes runtime.SchemeBuilder
	addToSchemes = append(addToSchemes, apiRules.AddToScheme)
	addToSchemes = append(addToSchemes, rulev1alpha1.AddToScheme)
	err := addToSchemes.AddToScheme(sch)
	if err != nil {
		return nil, err
	}
	return sch, nil
}

func (k *Service) ProvisionService(service *domain.Service, params *model.ServiceParams) error {
	labels := map[string]string{"app": service.Name, "created-by": "sample-broker"}

	if err := k.provisionDeployment(service, params, labels); err != nil {
		return err
	}

	if err := k.provisionK8sService(service, params, labels); err != nil {
		return err
	}

	if err := k.provisionAPIRule(service, params, labels); err != nil {
		return err
	}

	return nil
}

func (k *Service) provisionDeployment(service *domain.Service, params *model.ServiceParams, labels map[string]string) error {
	deployment := k.createK8sDeploymentObject(service, params, labels)
	return k.k8sClient.Create(context.TODO(), deployment)
}

func (k *Service) provisionK8sService(service *domain.Service, params *model.ServiceParams, labels map[string]string) error {
	k8sService := k.createK8SServiceObject(service, params, labels)
	return k.k8sClient.Create(context.TODO(), k8sService)
}

func (k *Service) provisionAPIRule(service *domain.Service, params *model.ServiceParams, labels map[string]string) error {
	kymaApiRule := k.createAPIRuleObject(service, params, labels)
	return k.k8sClient.Create(context.TODO(), kymaApiRule)
}

func (k *Service) createK8sDeploymentObject(service *domain.Service, params *model.ServiceParams, labels map[string]string) *appsv1.Deployment {
	replicas := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: params.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: service.Metadata.ImageUrl,
						Name:  service.Name,
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
						}},
					}},
				},
			},
		},
	}
	k.logger.Info("create-deployment", lager.Data{"deployment": deployment})
	return deployment
}

func (k *Service) createK8SServiceObject(service *domain.Service, params *model.ServiceParams, labels map[string]string) *corev1.Service {
	k8sSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: params.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				TargetPort: intstr.IntOrString{
					IntVal: 8080,
				},
				Port: 8080,
			}},
			Selector: labels,
		},
	}
	return k8sSvc
}

func (k *Service) createAPIRuleObject(service *domain.Service, params *model.ServiceParams, labels map[string]string) *apiRules.APIRule {
	fqdn := service.Name + params.Namespace
	port := uint32(8080)
	gateway := "kyma-gateway.kyma-system.svc.cluster.local"
	apiRule := apiRules.APIRule{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: params.Namespace,
			Name:      service.Name,
			Labels:    labels,
		},
		Spec: apiRules.APIRuleSpec{
			Rules: []apiRules.Rule{
				{
					Methods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
					Path: "/.*",
					AccessStrategies: []*rulev1alpha1.Authenticator{
						{
							Handler: &rulev1alpha1.Handler{
								Name:   "noop",
								Config: nil,
							},
						},
					},
				},
			},
			Service: &apiRules.Service{
				Name: &service.Name,
				Port: &port,
				Host: &fqdn,
			},
			Gateway: &gateway,
		},
	}
	return &apiRule
}
