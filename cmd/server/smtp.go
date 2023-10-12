package main

import (
	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/command"
	"github.com/Haya372/smtp-server/internal/connection"
	"github.com/Haya372/smtp-server/internal/server"
	"github.com/Haya372/smtp-server/internal/session"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			// TODO: 設定からロガーを生成できるようにする
			func() hlog.Config {
				return hlog.Config{
					LogLevel: hlog.Debug,
					Stdout:   true,
				}
			},
			hlog.NewLogger,
			command.AsCommandHandler(command.NewHeloHandler),
			command.AsCommandHandler(command.NewEhloHandler),
			command.AsCommandHandler(command.NewMailHandler),
			command.AsCommandHandler(command.NewRcptHandler),
			command.AsCommandHandler(command.NewDataHandler),
			command.AsCommandHandler(command.NewNoopHandler),
			command.AsCommandHandler(command.NewRsetHandler),
			command.AsCommandHandler(command.NewQuitHandler),
			command.AsCommandHandler(command.NewHelpHandler),
			session.NewSessionFactory,
			fx.Annotate(
				connection.NewSessionHandler,
				fx.ParamTags(``, `group:"commandhandler"`),
			),
			server.NewServer,
		),
		fx.Invoke(func(svr server.Server) {
			svr.ListenSmtp()
		}),
	)
	app.Run()
}
