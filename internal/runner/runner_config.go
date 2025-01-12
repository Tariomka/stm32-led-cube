package runner

import (
	"github.com/Tariomka/stm32-led-cube/internal/controller"
	"github.com/Tariomka/stm32-led-cube/internal/show"
)

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
	BaseSize   Size
	Height     Size
	LedType    LedType
	LightShows []controller.LightShow
}

func NewConfig() RunnerConfic {
	return RunnerConfic{
		BaseSize:   Size8,
		Height:     Size8,
		LedType:    RGB,
		LightShows: show.NewLedShowList(),
	}
}
