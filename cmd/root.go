package cmd

import (
	_ "embed"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/rfancn/protorefine/config"
	"github.com/rfancn/protorefine/pkg/finder"
	"github.com/rfancn/protorefine/pkg/reader"
	"github.com/rfancn/protorefine/pkg/render"
	"github.com/rfancn/protorefine/utils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type argument struct {
	projectName  string
	projectDir   string   // project path
	protoDir     string   // proto path
	outputDir    string   // outputDir filename
	skipDirs     []string // skip dirs when find protobuf types in project
	configFile   string   // config file
	pbImportPath string   // protobuf package import path in source code
}

var (
	arg = &argument{
		skipDirs: make([]string, 0),
	}
)

var rootCmd = &cobra.Command{
	Long:    "generate new proto source file that only contains pb types referenced in project source codes from proto files repository",
	PostRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		err := validateArgument(arg)
		if err != nil {
			utils.Fatalf(err, "validate argument")
		}

		err = genProtoFile(arg)
		if err != nil {
			utils.Fatalf(err, "generate proto file")
		}
	},
	Short: "generate proto file",
}

//go:embed configdata/config.toml
var defaultConfigContent string

func init() {
	cobra.OnInitialize(loadConfig)

	rootCmd.PersistentFlags().StringVarP(&arg.projectDir, "project-dir", "", "", "project souce code directory")
	_ = rootCmd.MarkPersistentFlagRequired("project-dir")
	rootCmd.PersistentFlags().StringVarP(&arg.protoDir, "proto-dir", "", "", "proto repository directory")
	_ = rootCmd.MarkPersistentFlagRequired("proto-dir")
	rootCmd.PersistentFlags().StringVarP(&arg.pbImportPath, "pb-import-path", "", "", "protobuf package import path in source code")
	_ = rootCmd.MarkPersistentFlagRequired("pb-import-path")

	rootCmd.PersistentFlags().StringVarP(&arg.outputDir, "output-dir", "", "", "output directory")
	rootCmd.PersistentFlags().StringVarP(&arg.outputDir, "config", "", "", "config file path")
	rootCmd.PersistentFlags().StringSliceVarP(&arg.skipDirs, "skip-dirs", "", nil, "skip directories when finding in project, separated by comma, e,g: autogen,xxx...")

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.Fatalf(err, "execute command")
	}
}

func validateArgument(arg *argument) error {
	if arg.projectName == "" {
		arg.projectName = filepath.Base(arg.projectDir)
	}

	var err error
	if arg.outputDir == "" {
		arg.outputDir, err = os.MkdirTemp("", "protorefine_*")
		if err != nil {
			utils.Fatalf(err, "generate temporary output directory")
		}
	}

	return nil
}

func genProtoFile(arg *argument) error {
	fmt.Printf("%-18s%s\n%-18s%s\n%-18s%s\n%-18s%s\n",
		"project-dir:", arg.projectDir,
		"proto-dir:", arg.protoDir,
		"pb-import-path:", arg.pbImportPath,
		"output-dir:", arg.outputDir)

	pbTypeNames, err := finder.New(arg.projectDir, arg.pbImportPath).Find(arg.skipDirs...)
	if err != nil {
		return errors.Wrap(err, "get pb types from project")
	}

	pbPkgName := filepath.Base(arg.pbImportPath)
	pbTypeDefs, err := reader.New(pbPkgName).ExtractTypeDefs(arg.protoDir, pbTypeNames)
	if err != nil {
		return errors.Wrap(err, "get proto pb type definitions")
	}

	err = render.New(arg.protoDir, arg.outputDir, arg.projectName).Render(pbTypeDefs)
	if err != nil {
		return errors.Wrap(err, "render proto file content")
	}

	return nil
}

func loadConfig() {
	configContent := defaultConfigContent
	if arg.configFile != "" {
		data, err := os.ReadFile(arg.configFile)
		if err != nil {
			utils.Fatalf(err, "read config file: %s", arg.configFile)
		}
		configContent = string(data)
	}

	_, err := toml.Decode(configContent, &config.Config)
	if err != nil {
		utils.Fatal(err)
	}
}
