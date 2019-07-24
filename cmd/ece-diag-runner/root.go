package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/ece-support-diagnostics/collectors/docker"
	"github.com/elastic/ece-support-diagnostics/collectors/eceAPI"
	"github.com/elastic/ece-support-diagnostics/collectors/eceMetrics"
	"github.com/elastic/ece-support-diagnostics/collectors/systemInfo"
	"github.com/elastic/ece-support-diagnostics/collectors/systemLogs"
	"github.com/elastic/ece-support-diagnostics/config"
	"github.com/elastic/ece-support-diagnostics/store/tar"
	"github.com/elastic/ece-support-diagnostics/uploader"
	"github.com/spf13/cobra"
)

var (
	cfg       *config.Config
	olderThan string
)

var rootCmd = &cobra.Command{
	Use:   "ece-diag-runner",
	Short: "ece-diag-runner collects support data for an ECE deployment",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("%+v\n", cfg)
		cfg.Initalize()

		fmt.Printf("Using %s for ECE install location\n", cfg.ElasticFolder)

		l := logp.NewLogger("Main")
		l.Infof("Using %s as temporary storage location", cfg.Basepath)

		tar, err := tar.Create(cfg.DiagnosticTarFilePath())
		defer tar.Close()
		if err != nil {
			// Exit here because we could not create the tar file
			log.Fatal(err)
		}
		// set Store interface in the config
		cfg.Store = tar

		docker.Run(cfg)
		eceAPI.Run(eceAPI.NewRestCalls(), cfg)
		systemInfo.Run(systemInfo.NewSystemCmdTasks(), systemInfo.NewSystemFileTasks(), cfg)
		systemLogs.Run(cfg)
		eceMetrics.Run(cfg)

		logTarPath := filepath.Join(cfg.DiagnosticFilename(), "diagnostic.log")
		tar.Finalize(cfg.DiagnosticLogFilePath(), logTarPath)

		if cfg.UploadUID != "" {
			uploader.RunUpload(tar.Filepath(), cfg.UploadUID)
		}

		// add an empty line
		println()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	cfg = config.New()

	rootCmd.PersistentFlags().StringVarP(
		&cfg.Basepath,
		"tmpFolder",
		"t",
		"/tmp",
		"specify a temporary location to write the output to",
	)
	rootCmd.PersistentFlags().StringVarP(
		&cfg.ElasticFolder,
		"eceHomeFolder",
		"f",
		"/mnt/data/elastic",
		"this should point to where you installed ECE",
	)
	rootCmd.PersistentFlags().StringVarP(
		&olderThan,
		"ignoreOlderThan",
		"i",
		"72h",
		"This is a cutoff, any log file with a modification time older than this duration will be ignored",
	)
	rootCmd.PersistentFlags().BoolVar(
		&cfg.DisableRest,
		"disableRest",
		false,
		"Disable Rest calls",
	)
	rootCmd.PersistentFlags().StringVarP(
		&cfg.UploadUID,
		"uploadID",
		"u",
		"",
		"Elastic Upload ID",
	)
}

func initConfig() {
	duration, err := time.ParseDuration(olderThan)
	if err != nil {
		log.Fatal("Could not parse -i / --ignoreOlderThan duration")
	}
	cfg.OlderThan = duration
}
