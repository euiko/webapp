package api

import "context"

type (
	App interface {
		Register(ModuleFactory)
		Run(context.Context) error
		Start(context.Context) error
	}
)
