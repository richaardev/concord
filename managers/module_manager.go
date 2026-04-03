package managers

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/richaardev/concord/discord/core"
)

type ModuleManager interface {
	Register(modules ...core.Module)
	Modules() []core.Module

	Load()
}

type moduleManagerImpl struct {
	modules []core.Module
}

func NewModuleManager() ModuleManager {
	return &moduleManagerImpl{}
}

func (m *moduleManagerImpl) Load() {
	var wg sync.WaitGroup

	for _, m := range m.modules {
		wg.Add(1)
		go func(m core.Module) {
			defer wg.Done()

			slog.Info(fmt.Sprintf("Plugin %s loaded successufully", m.Info().Name))
		}(m)
	}
	wg.Wait()
}

func (m *moduleManagerImpl) Modules() []core.Module {
	return m.modules
}

func (m *moduleManagerImpl) Register(modules ...core.Module) {
	m.modules = append(m.modules, modules...)
}
