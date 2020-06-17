package main

import (
	"os"
	"strings"

	"github.com/bearcherian/rollzap"
	"github.com/mkideal/cli"
	"github.com/rollbar/rollbar-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Args represents the command line argument.
type Args struct {
	cli.Helper
	Token string `cli:"*t" usage:"Rollbar Token"`
	Env   string `cli:"*e" usage:"Rollbar Environment"`
	Level string `cli:"l" usage:"Log Level"`
	Msg   string `cli:"*m" usage:"Message"`
}

func main() {
	os.Exit(cli.Run(&Args{}, func(ctx *cli.Context) error {
		args := ctx.Argv().(*Args)
		rollbar.SetToken(args.Token)
		rollbar.SetEnvironment(args.Env)
		rollbarCore := rollzap.NewRollbarCore(zapcore.WarnLevel)
		logger, err := zap.NewProduction()
		if err != nil {
			return err
		}
		logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, rollbarCore)
		}))
		switch strings.ToLower(args.Level) {
		case "info":
			logger.Info(args.Msg, zap.String("test", "test"))
			break
		case "warn":
			logger.Warn(args.Msg, zap.String("testwarn", "testwarn"))
			break
		case "error":
			logger.Error(args.Msg, zap.String("testerr", "testerr"))
			break
		case "fatal":
			logger.Fatal(args.Msg, zap.String("testfatal", "testfatal"))
			break
		}
		rollbar.Wait()
		return nil
	}))
}
