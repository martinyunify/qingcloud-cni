package nicmanagr

const (
	Random = iota
	Rotate
)

type AllocationPolicy interface {
	GetResource() (ResourcePool, error)
}

type ResourcePool interface {
	GetName()
	GetCapability()
}
