package config

import (
	"fmt"
	"github.com/prometheus/common/version"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
)

const defaultConfigFile = "config.toml"

type Loader struct {
	serviceName string
}

func NewLoader(serviceName string) *Loader {
	return &Loader{
		serviceName: serviceName,
	}
}

type CliFlags struct {
	configFile  string
	dumpConfig  bool
	showHelp    bool
	showVersion bool
}

func (f *CliFlags) Parse(serviceName string) (*pflag.FlagSet, error) {
	commandLine := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	// use standalone instead of shared default pflag.CommandLine avoid "pflag redefined: config" error when unit tests
	commandLine.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s %s\n\n", serviceName, version.Version)
		fmt.Fprintln(os.Stderr, commandLine.FlagUsages())
	}
	commandLine.SortFlags = false

	commandLine.StringVarP(&f.configFile, "config", "c", "", "config file path")
	commandLine.BoolVar(&f.dumpConfig, "dump", false, "dump config to toml")
	commandLine.BoolVarP(&f.showVersion, "version", "v", false, "display the current version of this CLI")
	commandLine.BoolVarP(&f.showHelp, "help", "h", false, "show help")

	return commandLine, commandLine.Parse(os.Args[1:])
}

func (l *Loader) Load(cfg *Config) error {
	flag := &CliFlags{}
	commandLine, err := flag.Parse(l.serviceName)
	if err != nil {
		return err
	}

	switch {
	case flag.showHelp:
		commandLine.Usage()
		os.Exit(0)
	case flag.showVersion:
		// nolint: forbidigo
		fmt.Println(version.Print(l.serviceName))
		os.Exit(0)
	case flag.dumpConfig:
		text, err := TomlMarshalIndent(cfg)
		if err != nil {
			panic(err)
		}
		fmt.Println(text)
		os.Exit(0)
	}

	if flag.configFile == "" {
		flag.configFile = defaultConfigFile
	}

	content, err := ioutil.ReadFile(flag.configFile)
	if err != nil {
		return fmt.Errorf("read config failed, config_file=%s err=%w", flag.configFile, err)
	}
	if len(content) == 0 {
		return fmt.Errorf("empty config file, config_file=%s", flag.configFile)
	}

	if err := TomlUnmarshaler(content, cfg); err != nil {
		return fmt.Errorf("unmarshal config failed, config_file=%s err=%w", flag.configFile, err)
	}
	return nil
}
