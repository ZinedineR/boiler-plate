package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"boiler-plate/app/api"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var HttpCmd = &cobra.Command{
	Use:   "http serve",
	Short: "Run Http API",
	Long:  "Run Http API",
	RunE: func(cmd *cobra.Command, args []string) error {
		initHTTP()

		// running open telemetry
		// cleanup := initTracer()
		// defer cleanup(context.Background())
		app := api.New(appConf.AppEnvConfig.AppName, baseHandler, URLHandler)

		echan := make(chan error)
		go func() {
			echan <- app.Run(appConf)
		}()

		//go func() {
		//	if err := serverio.Serve(); err != nil {
		//		logrus.Fatalf("socketio listen error: %s\n", err)
		//	}
		//}()
		//defer serverio.Close()

		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)

		select {
		case <-term:
			logrus.Infoln("signal terminated detected")
			return nil
		case err := <-echan:
			return errors.Wrap(err, "service runtime error")
		}
	},
}
