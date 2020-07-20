package broker

import (
	"code.cloudfoundry.org/lager"
	"context"
	"github.com/abbi-gaurav/go-projects/sample-broker/internal/config"
	"github.com/abbi-gaurav/go-projects/sample-broker/internal/constants"
	"github.com/abbi-gaurav/go-projects/sample-broker/internal/middleware"
	"github.com/abbi-gaurav/go-projects/sample-broker/internal/model"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-cf/brokerapi/domain"
)

type K8SServiceBroker struct {
	logger        lager.Logger
	availableSvcs []domain.Service
	service       *middleware.Service
}

func NewBroker(logger lager.Logger, services model.Services, service *middleware.Service) *K8SServiceBroker {
	availableSvcList := to(services)
	return &K8SServiceBroker{
		logger:        logger,
		availableSvcs: availableSvcList,
		service:       service,
	}
}

func to(services model.Services) []domain.Service {
	brokerSvcs := make([]domain.Service, len(services.Catalog))

	for i, svc := range services.Catalog {
		domainSvc := domain.Service{
			ID:                   svc.ServiceId,
			Name:                 svc.Name,
			Description:          svc.Description,
			Bindable:             false,
			InstancesRetrievable: false,
			BindingsRetrievable:  false,
			PlanUpdatable:        false,
			Plans:                []domain.ServicePlan{{ID: svc.PlanId, Name: "default", Description: "Default Plan"}},
			Metadata: &domain.ServiceMetadata{
				ImageUrl: svc.Image,
				AdditionalMetadata: map[string]interface{}{
					constants.ExposedPortMetadataField: svc.ExposedPort,
					constants.Namespace:                config.AppConfig().Namespace,
				},
			},
		}
		brokerSvcs[i] = domainSvc
	}
	return brokerSvcs
}

func (k *K8SServiceBroker) Services(ctx context.Context) ([]domain.Service, error) {
	k.logger.Info("list-services")
	return k.availableSvcs, nil
}

func (k *K8SServiceBroker) Provision(ctx context.Context, instanceID string, details domain.ProvisionDetails, asyncAllowed bool) (domain.ProvisionedServiceSpec, error) {
	k.logger.Info("provision", lager.Data{"instanceId": instanceID, "details": details, "asyncAllowed": asyncAllowed})
	params, err := model.Marshal(details.RawParameters)
	if err != nil {
		return domain.ProvisionedServiceSpec{}, err
	}
	service := brokerapi.RetrieveServiceFromContext(ctx)

	provisionedSpec, err := k.service.ProvisionService(service, params)
	if err != nil {
		return domain.ProvisionedServiceSpec{}, err
	}

	return provisionedSpec, nil
}

func (k *K8SServiceBroker) Deprovision(ctx context.Context, instanceID string, details domain.DeprovisionDetails, asyncAllowed bool) (domain.DeprovisionServiceSpec, error) {
	return domain.DeprovisionServiceSpec{}, nil
}

func (k *K8SServiceBroker) GetInstance(ctx context.Context, instanceID string) (domain.GetInstanceDetailsSpec, error) {
	return domain.GetInstanceDetailsSpec{}, nil
}

func (k *K8SServiceBroker) Update(ctx context.Context, instanceID string, details domain.UpdateDetails, asyncAllowed bool) (domain.UpdateServiceSpec, error) {
	return domain.UpdateServiceSpec{}, nil
}

func (k *K8SServiceBroker) LastOperation(ctx context.Context, instanceID string, details domain.PollDetails) (domain.LastOperation, error) {
	return domain.LastOperation{}, nil
}

func (k *K8SServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details domain.BindDetails, asyncAllowed bool) (domain.Binding, error) {
	return domain.Binding{}, nil
}

func (k *K8SServiceBroker) Unbind(ctx context.Context, instanceID, bindingID string, details domain.UnbindDetails, asyncAllowed bool) (domain.UnbindSpec, error) {
	return domain.UnbindSpec{}, nil
}

func (k *K8SServiceBroker) GetBinding(ctx context.Context, instanceID, bindingID string) (domain.GetBindingSpec, error) {
	return domain.GetBindingSpec{}, nil
}

func (k *K8SServiceBroker) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details domain.PollDetails) (domain.LastOperation, error) {
	return domain.LastOperation{}, nil
}
