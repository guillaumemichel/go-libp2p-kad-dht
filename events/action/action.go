package action

import "context"

type Action interface {
	Run(context.Context)
}
