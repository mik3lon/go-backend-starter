package kernel

import (
	"github.com/mik3lon/starter-template/pkg/bus"
	"github.com/mik3lon/starter-template/pkg/bus/command"
	"github.com/mik3lon/starter-template/pkg/bus/query"
)

type Modules []Module

// Module represents a module that can register routes.
type Module interface {
	RegisterRoutes(c *Kernel)
	Name() string
	Commands() map[bus.Dto]command.CommandHandler
	Queries() map[bus.Dto]query.QueryHandler
}

type BaseModule struct {
	commands map[bus.Dto]command.CommandHandler
	queries  map[bus.Dto]query.QueryHandler
}

// AddCommand adds a command to the module
func (bm *BaseModule) AddCommand(c bus.Dto, commandHandler command.CommandHandler) {
	if bm.commands == nil {
		bm.commands = make(map[bus.Dto]command.CommandHandler)
	}
	bm.commands[c] = commandHandler
}

// AddQuery adds a query to the module
func (bm *BaseModule) AddQuery(c bus.Dto, queryHandler query.QueryHandler) {
	if bm.queries == nil {
		bm.queries = make(map[bus.Dto]query.QueryHandler)
	}
	bm.queries[c] = queryHandler
}

// Commands returns all commands registered in the module
func (bm *BaseModule) Commands() map[bus.Dto]command.CommandHandler {
	return bm.commands
}

// Queries returns all commands registered in the module
func (bm *BaseModule) Queries() map[bus.Dto]query.QueryHandler {
	return bm.queries
}
