package manager

type BindItemsList []*BindItem

func (l BindItemsList) MarshalYAML() (interface{}, error) {
	return struct {
		Devices []*BindItem
	}{
		Devices: l,
	}, nil
}
