package runner

import "github.com/Tariomka/stm32-led-cube/internal/show/original"

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

type RunnerConfig struct {
	BaseSize Size
	Height   Size
	LedType  LedType
	Shows    string // Show data blob in the iCubeSmart PC software format
}

func NewConfig() RunnerConfig {
	return RunnerConfig{
		BaseSize: Size8,
		Height:   Size8,
		LedType:  RGB,
		Shows:    original.Blob,
	}
}
