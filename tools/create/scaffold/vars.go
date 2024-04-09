package scaffold

import "errors"

const (
	ManifestFileFlagName    = "gen-manifest"
	ManifestFileFlagDefault = "manifest.yaml"

	ModuleConfigFileFlagName    = "module-config"
	ModuleConfigFileFlagDefault = "scaffold-module-config.yaml"

	DefaultCRFlagName    = "gen-default-cr"
	DefaultCRFlagDefault = "default-cr.yaml"

	SecurityConfigFlagName    = "gen-security-config"
	SecurityConfigFlagDefault = "sec-scanners-config.yaml"
)

// Flags
var (
	ModuleName    string
	ModuleVersion string
	ModuleChannel string

	ModuleConfigFile   string
	ManifestFile       string
	SecurityConfigFile string
	DefaultCRFile      string

	Directory string
	Overwrite bool
)

// Errors
var (
	errModuleConfigExists           = errors.New("module config file already exists. use --overwrite flag to overwrite it")
	errModuleNameEmpty              = errors.New("--module-name flag must not be empty")
	errModuleVersionEmpty           = errors.New("--module-version flag must not be empty")
	errModuleChannelEmpty           = errors.New("--module-channel flag must not be empty")
	errManifestFileEmpty            = errors.New("--gen-manifest flag must not be empty")
	errModuleConfigEmpty            = errors.New("--module-config flag must not be empty")
	errManifestCreation             = errors.New("could not generate manifest")
	errDefaultCRCreationFailed      = errors.New("could not generate default CR")
	errModuleConfigCreationFailed   = errors.New("could not generate module config")
	errSecurityConfigCreationFailed = errors.New("could not generate security config")
)
