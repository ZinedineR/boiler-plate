package cmd

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "Go Simple API",
	Short: "Go Simple API / Service Demo",
	Long:  "Go Simple API / Service Demo HTTP & GRPC API & Kafka",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// register command
func init() {
	rootCmd.AddCommand(HttpCmd)

	// load environment variable
	if err := godotenv.Load(); err != nil {
		if os.Getenv("APP_ENV") == "development" {
			logrus.Println("unable to load environment variable", err.Error())
		} else {
			// fin, err := os.Open("./.env.example")
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// defer fin.Close()

			fout, err := os.Create("./.env")
			if err != nil {
				log.Fatal(err)
			}
			defer fout.Close()

			// _, err = io.Copy(fout, fin)

			// if err != nil {
			// 	logrus.Warningln("Can't find env.file. To use system's env vars for now")s
			// }
		}
	}
}
func Execute() error {
	cmd, _, err := rootCmd.Find(os.Args[1:])

	if err == nil && cmd.Use == rootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		args := append([]string{"http"}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	return rootCmd.Execute()
}
