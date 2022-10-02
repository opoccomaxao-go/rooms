package engine

type Factory interface {
	New() Engine
}

type Engine interface{}
