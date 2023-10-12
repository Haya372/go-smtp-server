package command

import (
	"context"

	"github.com/Haya372/smtp-server/internal/session"
	"go.uber.org/fx"
)

type CommandHandler interface {
	HandleCommand(ctx context.Context, s session.Session, arg []string) error
	Command() string
}

func AsCommandHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(CommandHandler)),
		fx.ResultTags(`group:"commandhandler"`),
	)
}
