package core

var _ Module = (*moduleImpl)(nil)

type ModuleInfo struct {
	Name string

	Hidden bool
}

type Module interface {
	Info() ModuleInfo
}

type moduleImpl struct {
	info ModuleInfo
}

func (m *moduleImpl) Info() ModuleInfo {
	return m.info
}

func NewModule(info ModuleInfo) Module {
	return &moduleImpl{
		info: info,
	}
}
