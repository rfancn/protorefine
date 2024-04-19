package cmd

import (
	"github.com/rfancn/protorefine/utils"
	"testing"
)

func Test_genProtoFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test_genProtoFile",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg = &argument{
				projectDir:   "D:\\codes\\hdmall\\backend\\product",
				protoDir:     "D:\\codes\\hdmall\\backend\\proto",
				pbImportPath: "product/autogen/pb",
				outputDir:    "autogen/proto",
				skipDirs:     []string{"autogen"},
			}

			loadConfig()

			err := validateArgument(arg)
			if err != nil {
				utils.Fatalf(err, "validate argument")
			}

			_ = genProtoFile(arg)
		})
	}
}
