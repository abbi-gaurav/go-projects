package storage

import (
	apiRules "github.com/kyma-incubator/api-gateway/api/v1alpha1"
	"sync"
)

type Storage interface {
	DeleteInstance(instanceId string)
	GetInstance(instanceId string) *apiRules.APIRule
	AddInstance(instanceId string, rule *apiRules.APIRule)
}

type inMemStorage struct {
	instanceData sync.Map
}

func NewInMemory() *inMemStorage {
	return &inMemStorage{}
}

func (i *inMemStorage) AddInstance(instanceId string, rule *apiRules.APIRule) {
	i.instanceData.Store(instanceId, rule)
}

func (i *inMemStorage) DeleteInstance(instanceId string) {
	i.instanceData.Delete(instanceId)
}

func (i *inMemStorage) GetInstance(instanceId string) *apiRules.APIRule {
	value, ok := i.instanceData.Load(instanceId)
	if !ok {
		return nil
	}
	return value.(*apiRules.APIRule)
}
