package vcs

import "sync"

type RepositoryService interface {
	GetRepository(user, project string) error
	Name() string
}

var services = struct {
	sync.RWMutex
	s map[string]RepositoryService
}{s: map[string]RepositoryService{}}

func RegisterService(service RepositoryService) {
	services.Lock()
	defer services.Unlock()

	services.s[service.Name()] = service
}

func GetServices() []RepositoryService {
	services.RLock()
	defer services.RUnlock()

	srvcs := make([]RepositoryService, len(services.s))

	for _, srvc := range services.s {
		srvcs = append(srvcs, srvc)
	}

	return srvcs
}
