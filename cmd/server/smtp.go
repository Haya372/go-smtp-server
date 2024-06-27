package main

import (
	"context"

	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/command"
	"github.com/Haya372/smtp-server/internal/config"
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
			// TODO: 設定ファイルから読み込んだものを使用する
			config.NewDefaultConfig,
			config.NewServerConfig,
			config.NewSmtpConfig,
			config.NewTlsConfig,
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
			command.AsCommandHandler(command.NewStartTlsHandler),
			session.NewSessionFactory,
			fx.Annotate(
				connection.NewSessionHandler,
				fx.ParamTags(``, `group:"commandhandler"`),
			),
			func(lc fx.Lifecycle, log hlog.Logger, conf *config.ServerConfig, factory session.SessionFactory, handler connection.SessionHandler) *server.Server {
				s := server.NewServer(log, conf, factory, handler)
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						return s.ListenSmtp(ctx)
					},
				})
				return &s
			},
		),
		fx.Invoke(func(s *server.Server) {}),
	)
	app.Run()
}
