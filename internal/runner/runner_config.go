package runner

type Size uint8

const (
	Size8 Size = iota
	Size16
	Size32
)

type LedType uint8

const (
	Monochrome LedType = iota
	RGB
)

type RunnerConfic struct {
	BaseSize Size
	Height   Size
	LedType  LedType
}

func NewConfig() RunnerConfic {
	return RunnerConfic{
		BaseSize: Size8,
		Height:   Size8,
		LedType:  RGB,
	}
}
