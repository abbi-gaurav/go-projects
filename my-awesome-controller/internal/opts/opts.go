package opts

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type Options struct {
	Port           int
	DbType         string
	ResyncDuration time.Duration
}

var DefaultOptions = Options{
	Port:           8080,
	DbType:         "memory",
	ResyncDuration: 1 * time.Minute,
}

var version = os.Getenv("APP_VERSION")

var usageStr = `
	Usage: my-awesome-controller
`

func usage() {
	fmt.Printf("%s\n", usageStr)
}

func ParseFlags() *Options {
	fs := flag.NewFlagSet("my-awesome-controller", flag.ExitOnError)
	fs.Usage = usage

	opts, err := configureOptions(fs, os.Args[1:],
		func() {
			fmt.Printf("push version %s, ", version)
			os.Exit(0)
		},
		fs.Usage)
	if err != nil {
		log.Fatalf("failed to parse command line flags: %v", err.Error()+"\n"+usageStr)
	}

	return opts
}

func configureOptions(fs *flag.FlagSet, args []string, printVersion, printHelp func()) (*Options, error) {
	opts := &DefaultOptions

	var (
		showVersion bool
		showHelp    bool
	)

	fs.BoolVar(&showHelp, "h", false, "show this message")
	fs.BoolVar(&showHelp, "help", false, "show this message")
	fs.BoolVar(&showVersion, "v", false, "print version information")
	fs.BoolVar(&showVersion, "version", false, "print version information")

	fs.IntVar(&opts.Port, "port", opts.Port, "HTTP Server Port")
	fs.StringVar(&opts.DbType, "db_type", opts.DbType, "DB type to use")
	fs.DurationVar(&opts.ResyncDuration, "resync_duration", opts.ResyncDuration, "Duration to resyn custom resource")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if showVersion {
		printVersion()
		return nil, nil
	}

	if showHelp {
		printHelp()
		return nil, nil
	}

	return opts, nil
}
