package config

var (
	Config ProtoRefineConfig
)

type ProtoRefineConfig struct {
	Import struct {
		Rules []ConfImportRule `toml:"rules"`
	} `toml:"import"`
}

type ConfImportRule struct {
	Match      string   `toml:"match"`
	ProtoFile  string   `toml:"file"`
	Dependents []string `toml:"dependents"`
}
