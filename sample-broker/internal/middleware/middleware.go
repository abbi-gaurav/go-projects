package middleware

import (
	"code.cloudfoundry.org/lager"
	"context"
	goerrors "errors"
	appconfig "github.com/abbi-gaurav/go-projects/sample-broker/internal/config"
	"github.com/abbi-gaurav/go-projects/sample-broker/internal/constants"
	"github.com/abbi-gaurav/go-projects/sample-broker/internal/model"
	"github.com/abbi-gaurav/go-projects/sample-broker/internal/storage"
	apiRules "github.com/kyma-incubator/api-gateway/api/v1alpha1"
	rulev1alpha1 "github.com/ory/oathkeeper-maester/api/v1alpha1"
	"github.com/pivotal-cf/brokerapi/domain"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/deprecated/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"strconv"
)

type Service struct {
	k8sClient   client.Client
	logger      lager.Logger
	persistence storage.Storage
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
		k8sClient:   k8sClient,
		logger:      logger,
		persistence: storage.NewInMemory(),
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

func getContainerPort(additionalMetadata map[string]interface{}) (int32, error) {
	i, exists := additionalMetadata[constants.ExposedPortMetadataField]

	if exists {
		containerPort, err := strconv.Atoi(i.(string))
		if err != nil {
			return -1, err
		}
		return int32(containerPort), nil
	} else {
		return int32(constants.DefaultServicePort), nil
	}
}

func (k *Service) ProvisionService(service *domain.Service, params *model.ServiceParams, instanceId string) (domain.ProvisionedServiceSpec, error) {
	fqdn := params.ServiceInstanceName + "-" + appconfig.AppConfig().Namespace

	existingInstance, err := k.tryGettingAPIRule(instanceId)
	if err != nil {
		return domain.ProvisionedServiceSpec{}, err
	}
	if existingInstance != nil {
		k.logger.Info("Provision Instance", lager.Data{"existing instance": existingInstance})
		return generateProvisionedSpec(fqdn), nil
	}
	labels := map[string]string{
		"app": params.ServiceInstanceName, "created-by": "sample-broker",
		"service": service.Name, constants.DemoOsbInstanceIdLabelName: instanceId,
	}

	if err := k.provisionAPIRule(params, labels, &fqdn, instanceId, service); err != nil {
		return domain.ProvisionedServiceSpec{}, err
	}

	return generateProvisionedSpec(fqdn), nil
}

func (k *Service) provisionAPIRule(params *model.ServiceParams, labels map[string]string, fqdn *string, instanceId string, service *domain.Service) error {
	kymaApiRule := k.createAPIRuleObject(params, labels, fqdn, service)
	err := k.k8sClient.Create(context.TODO(), kymaApiRule)
	k.persistence.AddInstance(instanceId, kymaApiRule)
	return err
}

func (k *Service) createAPIRuleObject(params *model.ServiceParams, labels map[string]string, fqdn *string, service *domain.Service) *apiRules.APIRule {
	port := uint32(constants.DefaultServicePort)
	gateway := constants.KymaGatewayDomain
	apiRule := apiRules.APIRule{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: appconfig.AppConfig().Namespace,
			Name:      params.ServiceInstanceName,
			Labels:    labels,
		},
		Spec: apiRules.APIRuleSpec{
			Rules: []apiRules.Rule{
				{
					Methods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
					Path:    "/.*",
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
				Host: fqdn,
			},
			Gateway: &gateway,
		},
	}
	return &apiRule
}

func (k *Service) tryGettingAPIRule(instanceId string) (*apiRules.APIRule, error) {
	apiRuleList := &apiRules.APIRuleList{}
	err := k.k8sClient.List(context.TODO(), apiRuleList, client.MatchingLabels{
		constants.DemoOsbInstanceIdLabelName: instanceId,
	}, client.InNamespace(appconfig.AppConfig().Namespace))
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	if len(apiRuleList.Items) > 1 {
		return nil, goerrors.New("More than one API Rule for instance Id - " + instanceId)
	} else if len(apiRuleList.Items) == 0 {
		return nil, nil
	}
	return &apiRuleList.Items[0], nil
}

func (k *Service) Deprovision(ctx context.Context, instanceID string, deprovisionDetails domain.DeprovisionDetails) (domain.DeprovisionServiceSpec, error) {
	err := k.k8sClient.DeleteAllOf(
		context.TODO(),
		&apiRules.APIRule{},
		client.MatchingLabels{constants.DemoOsbInstanceIdLabelName: instanceID},
		client.InNamespace(appconfig.AppConfig().Namespace),
	)
	if err != nil {
		return domain.DeprovisionServiceSpec{}, err
	}
	return domain.DeprovisionServiceSpec{
		OperationData: "api rule deleted",
	}, nil
}

func generateProvisionedSpec(fqdn string) domain.ProvisionedServiceSpec {
	serviceUrl := createServiceUrl(fqdn)
	return domain.ProvisionedServiceSpec{
		IsAsync:       false,
		AlreadyExists: false,
		DashboardURL:  serviceUrl,
		OperationData: serviceUrl,
	}
}

func createServiceUrl(fqdn string) string {
	return "https://" + fqdn + "." + appconfig.AppConfig().ClusterDomain
}

func (k *Service) provisionK8sService(params *model.ServiceParams, labels map[string]string, containerPort int32) error {
	k8sService := k.createK8SServiceObject(params, labels, containerPort)
	return k.k8sClient.Create(context.TODO(), k8sService)
}

func (k *Service) provisionDeployment(service *domain.Service, params *model.ServiceParams, labels map[string]string, containerPort int32) error {
	deployment := k.createK8sDeploymentObject(service, params, labels, containerPort)
	return k.k8sClient.Create(context.TODO(), deployment)
}

func (k *Service) createK8SServiceObject(params *model.ServiceParams, labels map[string]string, containerPort int32) *corev1.Service {
	k8sSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      params.ServiceInstanceName,
			Namespace: appconfig.AppConfig().Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				TargetPort: intstr.IntOrString{
					IntVal: containerPort,
				},
				Port: constants.DefaultServicePort,
			}},
			Selector: labels,
		},
	}
	return k8sSvc
}

func (k *Service) createK8sDeploymentObject(service *domain.Service, params *model.ServiceParams, labels map[string]string, containerPort int32) *appsv1.Deployment {
	replicas := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      params.ServiceInstanceName,
			Namespace: appconfig.AppConfig().Namespace,
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
						Name:  params.ServiceInstanceName,
						Ports: []corev1.ContainerPort{{
							ContainerPort: containerPort,
						}},
					}},
				},
			},
		},
	}
	k.logger.Info("create-deployment", lager.Data{"deployment": deployment})
	return deployment
}

func (k *Service) Bind(instanceID string, bindingID string, bindDetails domain.BindDetails) (domain.Binding, error) {
	apiRule, err := k.tryGettingAPIRule(instanceID)
	if err != nil {
		return domain.Binding{}, err
	}
	return domain.Binding{Credentials: dummyCredentials(apiRule)}, nil
}

func dummyCredentials(apiRule *apiRules.APIRule) map[string]interface{} {
	return map[string]interface{}{
		"url": createServiceUrl(*apiRule.Spec.Service.Host),
		"auth": model.Auth{
			BaUser:     "admin",
			BaPassword: "nimda123",
		},
	}
}

func (k *Service) GetBinding(instanceID string, bindingID string) (domain.GetBindingSpec, error) {
	apiRule, err := k.tryGettingAPIRule(instanceID)
	if err != nil {
		return domain.GetBindingSpec{}, err
	}
	return domain.GetBindingSpec{Credentials: dummyCredentials(apiRule)}, nil
}

func (k *Service) GetInstance(service *domain.Service, instanceID string) (domain.GetInstanceDetailsSpec, error) {
	apiRule, err := k.tryGettingAPIRule(instanceID)
	if err != nil {
		return domain.GetInstanceDetailsSpec{}, err
	}
	return domain.GetInstanceDetailsSpec{
		ServiceID:    service.ID,
		PlanID:       service.Plans[0].ID,
		DashboardURL: createServiceUrl(*apiRule.Spec.Service.Host),
	}, nil
}
