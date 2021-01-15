package generators

type Step struct {
	Description string
	Content     string
}

type HasGeneratorOpenHab interface {
	GenerateConfigOpenHab() []Step
}
