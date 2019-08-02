package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/elastic/beats/libbeat/logp"
	collector "github.com/elastic/ece-support-diagnostics/collectors"
	"github.com/elastic/ece-support-diagnostics/collectors/docker"
	"github.com/elastic/ece-support-diagnostics/collectors/eceAPI"
	"github.com/elastic/ece-support-diagnostics/collectors/eceMetrics"
	"github.com/elastic/ece-support-diagnostics/collectors/systemInfo"
	"github.com/elastic/ece-support-diagnostics/collectors/systemLogs"
	"github.com/elastic/ece-support-diagnostics/config"
	"github.com/elastic/ece-support-diagnostics/helpers"
	"github.com/elastic/ece-support-diagnostics/store/tar"
	"github.com/elastic/ece-support-diagnostics/uploader"
	"github.com/spf13/cobra"
)

var (
	cfg        *config.Config
	olderThan  string
	cpuprofile string
	VERSION    string
)

// Execute is the main entry point
func Execute(version string) {
	VERSION = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "ece-support-diagnostics",
	Short: "ece-support-diagnostics collects support data for an ECE deployment",
	// Long: `A Fast and Flexible Static Site Generator built with
	// 			  love by spf13 and friends in Go.
	// 			  Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		if cpuprofile != "" {
			f, err := os.Create(cpuprofile)
			if err != nil {
				log.Fatal("could not create CPU profile: ", err)
			}
			defer f.Close()
			if err := pprof.StartCPUProfile(f); err != nil {
				log.Fatal("could not start CPU profile: ", err)
			}
			defer pprof.StopCPUProfile()
		}

		// fmt.Printf("%+v\n", cfg)
		cfg.Initalize()

		fmt.Printf("Using %s for ECE install location\n", cfg.ElasticFolder)
		fmt.Printf("\tECE Runner Name: %s\n", cfg.RunnerName())

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

		messages := make(chan string) // channel used to print status from each of collectors

		spinner := helpers.NewSpinner("%s Working...")
		spinner.Start()
		defer spinner.Stop()

		go startCollectors(messages)

		for message := range messages {
			fmt.Println(helpers.ClearLine + message)
		}
		spinner.Stop()

		logTarPath := filepath.Join(cfg.DiagnosticFilename(), "diagnostic.log")
		tar.Finalize(cfg.DiagnosticLogFilePath(), logTarPath)

		fmt.Printf("Finished creating file: %s (total: %s)\n",
			cfg.DiagnosticTarFilePath(),
			time.Since(cfg.StartTime).Truncate(time.Millisecond),
		)

		if cfg.UploadUID != "" {
			uploader.RunUpload(tar.Filepath(), cfg.UploadUID)
		}
	},
}

func startCollectors(returnCh chan<- string) {

	var wg sync.WaitGroup

	wg.Add(5)
	go collector.StartCollector(systemLogs.Run, returnCh, cfg, &wg)
	go collector.StartCollector(systemInfo.Run, returnCh, cfg, &wg)
	go collector.StartCollector(eceMetrics.Run, returnCh, cfg, &wg)
	go collector.StartCollector(docker.Run, returnCh, cfg, &wg)
	go collector.StartCollector(eceAPI.Run, returnCh, cfg, &wg)

	wg.Wait()       // wait for all tasks
	close(returnCh) // printing each item using range over the channel. Closing to end range.

}

func init() {
	cobra.OnInitialize(initConfig)

	cfg = config.New()

	rootCmd.PersistentFlags().StringVarP(
		&cfg.Auth.User,
		"username",
		"u",
		"",
		`<user:password>
		
	Specify the user name and password to use to authenticate to ECE

	If you only specify the user name, you will be prompted for a password`,
	)
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
		"disableApiCalls",
		false,
		"Disable API calls",
	)
	rootCmd.PersistentFlags().StringVarP(
		&cfg.UploadUID,
		"uploadID",
		"U",
		"",
		"Elastic Upload ID",
	)
	rootCmd.Flags().StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to `file`")
	rootCmd.Flags().MarkHidden("cpuprofile")
}

func initConfig() {
	duration, err := time.ParseDuration(olderThan)
	if err != nil {
		log.Fatal("Could not parse -i / --ignoreOlderThan duration")
	}
	cfg.OlderThan = duration
}
