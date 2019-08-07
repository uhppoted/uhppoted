package simulator

import (
	"uhppote-simulator/entities"
)

type Context struct {
	Simulators []*Simulator
	Directory  string
	TxQ        chan entities.Message // TODO: interim hack - replace when TxQ is a simulator property
}
