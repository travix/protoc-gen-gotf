package terraform

var _ TypeValue = &typeValue{}

type TypeValue interface {
}

type typeValue struct {
}
