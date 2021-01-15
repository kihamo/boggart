package generators

type Step struct {
	Description string
	FilePath    string
	Content     string
}

type HasGeneratorOpenHab interface {
	GenerateConfigOpenHab() []Step
}
