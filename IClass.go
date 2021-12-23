package gmoon

type IClass interface {
	Build(moon *GMoon)
	Name() string
}
