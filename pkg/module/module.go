package module

import (
	"errors"
	"strconv"
	"strings"
)

const defaultType = "generic_data_string"
const defaultValue = "0"

// DataPoint represents a single value entry in a Pandora datalist payload.
type DataPoint struct {
	Value     string
	Timestamp string
}

// Config defines the public configuration used to create a Pandora module.
type Config struct {
	Name                 string
	Type                 string
	Value                string
	Data                 string
	DataList             []DataPoint
	Description          string
	Desc                 string
	Unit                 string
	Interval             string
	Tags                 string
	ModuleGroup          string
	ModuleParent         string
	MinWarning           string
	MinWarningForced     string
	MaxWarning           string
	MaxWarningForced     string
	MinCritical          string
	MinCriticalForced    string
	MaxCritical          string
	MaxCriticalForced    string
	StrWarning           string
	StrWarningForced     string
	StrCritical          string
	StrCriticalForced    string
	CriticalInverse      string
	WarningInverse       string
	Min                  string
	Max                  string
	PostProcess          string
	Disabled             string
	MinFFEvent           string
	Status               string
	Timestamp            string
	CustomID             string
	CriticalInstructions string
	WarningInstructions  string
	UnknownInstructions  string
	Quiet                string
	ModuleFFInterval     string
	CronTab              string
	MinFFEventNormal     string
	MinFFEventWarning    string
	MinFFEventCritical   string
	FFType               string
	FFTimeout            string
	EachFF               string
	ModuleParentUnlink   string
	ExtraData            string
	AlertTemplates       []string
	Alert                []string
}

// Module is the validated public module model used by the library.
type Module struct {
	Config Config
}

// LogConfig defines the public configuration used to create a Pandora log module.
type LogConfig struct {
	Source string
	Value  string
}

// LogModule is the validated public log-module model used by the library.
type LogModule struct {
	Config LogConfig
}

// New creates a new module with Phase 1 defaults applied.
func New(cfg Config) (Module, error) {
	cfg = applyDefaults(cfg)

	m := Module{Config: cfg}
	if err := m.Validate(); err != nil {
		return Module{}, err
	}

	return m, nil
}

// NewLog creates a new log module.
func NewLog(cfg LogConfig) (LogModule, error) {
	m := LogModule{Config: cfg}
	if err := m.Validate(); err != nil {
		return LogModule{}, err
	}

	return m, nil
}

// Validate verifies the minimum invariants for scalar and datalist module payloads.
func (m Module) Validate() error {
	if strings.TrimSpace(m.Config.Name) == "" {
		return errors.New("module name is required")
	}

	if strings.TrimSpace(m.Config.Type) == "" {
		return errors.New("module type is required")
	}

	for i, point := range m.Config.DataList {
		if strings.TrimSpace(point.Value) == "" {
			return errors.New("datalist point value is required at index " + strconv.Itoa(i))
		}
	}

	return nil
}

// Validate verifies the minimum invariants for log modules.
func (m LogModule) Validate() error {
	if strings.TrimSpace(m.Config.Source) == "" {
		return errors.New("log module source is required")
	}

	return nil
}

func applyDefaults(cfg Config) Config {
	if strings.TrimSpace(cfg.Type) == "" {
		cfg.Type = defaultType
	}

	if strings.TrimSpace(cfg.Data) != "" {
		cfg.Value = cfg.Data
	}

	if strings.TrimSpace(cfg.Description) == "" && strings.TrimSpace(cfg.Desc) != "" {
		cfg.Description = cfg.Desc
	}

	if len(cfg.AlertTemplates) == 0 && len(cfg.Alert) > 0 {
		cfg.AlertTemplates = cfg.Alert
	}

	if cfg.Value == "" && len(cfg.DataList) == 0 {
		cfg.Value = defaultValue
	}

	return cfg
}
