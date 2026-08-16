package module

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const defaultType = "generic_data_string"
const defaultValue = "0"
const pandoraTimestampLayout = "2006/01/02 15:04:05"

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

// XMLOptions controls optional XML output behavior for module fragments.
type XMLOptions struct {
	LogEncoding string
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

// XML renders only module and log_module nodes without an agent wrapper.
func XML(modules []Module, logModules []LogModule) ([]byte, error) {
	return XMLWithOptions(modules, logModules, XMLOptions{})
}

// XMLWithOptions renders only module and log_module nodes without an agent wrapper.
func XMLWithOptions(modules []Module, logModules []LogModule, opts XMLOptions) ([]byte, error) {
	for _, m := range modules {
		if err := m.Validate(); err != nil {
			return nil, err
		}
	}

	for _, m := range logModules {
		if err := m.Validate(); err != nil {
			return nil, err
		}
	}

	var buf bytes.Buffer

	for _, m := range modules {
		body, err := xml.MarshalIndent(moduleXMLPayload(m), "", "  ")
		if err != nil {
			return nil, err
		}
		buf.Write(body)
		buf.WriteByte('\n')
	}

	for _, m := range logModules {
		body, err := xml.MarshalIndent(logModuleXMLPayload(m, opts), "", "  ")
		if err != nil {
			return nil, err
		}
		buf.Write(body)
		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
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

type cdataText struct {
	Text string `xml:",cdata"`
}

type dataListXML struct {
	Items []dataPointXML `xml:"data"`
}

type dataPointXML struct {
	Value     cdataText  `xml:"value"`
	Timestamp *cdataText `xml:"timestamp,omitempty"`
}

type moduleXML struct {
	XMLName              xml.Name     `xml:"module"`
	Name                 cdataText    `xml:"name"`
	Type                 string       `xml:"type"`
	Data                 *cdataText   `xml:"data,omitempty"`
	DataList             *dataListXML `xml:"datalist,omitempty"`
	Description          *cdataText   `xml:"description,omitempty"`
	Unit                 *cdataText   `xml:"unit,omitempty"`
	ModuleInterval       *cdataText   `xml:"module_interval,omitempty"`
	Tags                 string       `xml:"tags,omitempty"`
	ModuleGroup          string       `xml:"module_group,omitempty"`
	ModuleParent         string       `xml:"module_parent,omitempty"`
	MinWarning           *cdataText   `xml:"min_warning,omitempty"`
	MinWarningForced     *cdataText   `xml:"min_warning_forced,omitempty"`
	MaxWarning           *cdataText   `xml:"max_warning,omitempty"`
	MaxWarningForced     *cdataText   `xml:"max_warning_forced,omitempty"`
	MinCritical          *cdataText   `xml:"min_critical,omitempty"`
	MinCriticalForced    *cdataText   `xml:"min_critical_forced,omitempty"`
	MaxCritical          *cdataText   `xml:"max_critical,omitempty"`
	MaxCriticalForced    *cdataText   `xml:"max_critical_forced,omitempty"`
	StrWarning           *cdataText   `xml:"str_warning,omitempty"`
	StrWarningForced     *cdataText   `xml:"str_warning_forced,omitempty"`
	StrCritical          *cdataText   `xml:"str_critical,omitempty"`
	StrCriticalForced    *cdataText   `xml:"str_critical_forced,omitempty"`
	CriticalInverse      *cdataText   `xml:"critical_inverse,omitempty"`
	WarningInverse       *cdataText   `xml:"warning_inverse,omitempty"`
	Max                  *cdataText   `xml:"max,omitempty"`
	Min                  *cdataText   `xml:"min,omitempty"`
	PostProcess          *cdataText   `xml:"post_process,omitempty"`
	Disabled             *cdataText   `xml:"disabled,omitempty"`
	MinFFEvent           *cdataText   `xml:"min_ff_event,omitempty"`
	Status               *cdataText   `xml:"status,omitempty"`
	Timestamp            *cdataText   `xml:"timestamp,omitempty"`
	CustomID             *cdataText   `xml:"custom_id,omitempty"`
	CriticalInstructions *cdataText   `xml:"critical_instructions,omitempty"`
	WarningInstructions  *cdataText   `xml:"warning_instructions,omitempty"`
	UnknownInstructions  *cdataText   `xml:"unknown_instructions,omitempty"`
	Quiet                *cdataText   `xml:"quiet,omitempty"`
	ModuleFFInterval     *cdataText   `xml:"module_ff_interval,omitempty"`
	CronTab              *cdataText   `xml:"crontab,omitempty"`
	MinFFEventNormal     *cdataText   `xml:"min_ff_event_normal,omitempty"`
	MinFFEventWarning    *cdataText   `xml:"min_ff_event_warning,omitempty"`
	MinFFEventCritical   *cdataText   `xml:"min_ff_event_critical,omitempty"`
	FFType               *cdataText   `xml:"ff_type,omitempty"`
	FFTimeout            *cdataText   `xml:"ff_timeout,omitempty"`
	EachFF               *cdataText   `xml:"each_ff,omitempty"`
	ModuleParentUnlink   *cdataText   `xml:"module_parent_unlink,omitempty"`
	ExtraData            *cdataText   `xml:"extra_data,omitempty"`
	AlertTemplates       []cdataText  `xml:"alert_template,omitempty"`
}

type logModuleXML struct {
	XMLName  xml.Name  `xml:"log_module"`
	Source   cdataText `xml:"source"`
	Data     string    `xml:"data"`
	Encoding string    `xml:"encoding,omitempty"`
}

func moduleXMLPayload(m Module) moduleXML {
	return moduleXML{
		Name:                 cdataText{Text: m.Config.Name},
		Type:                 m.Config.Type,
		Data:                 scalarData(m.Config.Value, m.Config.DataList),
		DataList:             dataList(m.Config.DataList),
		Description:          optionalCDATA(m.Config.Description),
		Unit:                 optionalCDATA(m.Config.Unit),
		ModuleInterval:       optionalCDATA(m.Config.Interval),
		Tags:                 m.Config.Tags,
		ModuleGroup:          m.Config.ModuleGroup,
		ModuleParent:         m.Config.ModuleParent,
		MinWarning:           optionalCDATA(m.Config.MinWarning),
		MinWarningForced:     optionalCDATA(m.Config.MinWarningForced),
		MaxWarning:           optionalCDATA(m.Config.MaxWarning),
		MaxWarningForced:     optionalCDATA(m.Config.MaxWarningForced),
		MinCritical:          optionalCDATA(m.Config.MinCritical),
		MinCriticalForced:    optionalCDATA(m.Config.MinCriticalForced),
		MaxCritical:          optionalCDATA(m.Config.MaxCritical),
		MaxCriticalForced:    optionalCDATA(m.Config.MaxCriticalForced),
		StrWarning:           optionalCDATA(m.Config.StrWarning),
		StrWarningForced:     optionalCDATA(m.Config.StrWarningForced),
		StrCritical:          optionalCDATA(m.Config.StrCritical),
		StrCriticalForced:    optionalCDATA(m.Config.StrCriticalForced),
		CriticalInverse:      optionalCDATA(m.Config.CriticalInverse),
		WarningInverse:       optionalCDATA(m.Config.WarningInverse),
		Min:                  optionalCDATA(m.Config.Min),
		Max:                  optionalCDATA(m.Config.Max),
		PostProcess:          optionalCDATA(m.Config.PostProcess),
		Disabled:             optionalCDATA(m.Config.Disabled),
		MinFFEvent:           optionalCDATA(m.Config.MinFFEvent),
		Status:               optionalCDATA(m.Config.Status),
		Timestamp:            optionalCDATA(m.Config.Timestamp),
		CustomID:             optionalCDATA(m.Config.CustomID),
		CriticalInstructions: optionalCDATA(m.Config.CriticalInstructions),
		WarningInstructions:  optionalCDATA(m.Config.WarningInstructions),
		UnknownInstructions:  optionalCDATA(m.Config.UnknownInstructions),
		Quiet:                optionalCDATA(m.Config.Quiet),
		ModuleFFInterval:     optionalCDATA(m.Config.ModuleFFInterval),
		CronTab:              optionalCDATA(m.Config.CronTab),
		MinFFEventNormal:     optionalCDATA(m.Config.MinFFEventNormal),
		MinFFEventWarning:    optionalCDATA(m.Config.MinFFEventWarning),
		MinFFEventCritical:   optionalCDATA(m.Config.MinFFEventCritical),
		FFType:               optionalCDATA(m.Config.FFType),
		FFTimeout:            optionalCDATA(m.Config.FFTimeout),
		EachFF:               optionalCDATA(m.Config.EachFF),
		ModuleParentUnlink:   optionalCDATA(m.Config.ModuleParentUnlink),
		ExtraData:            optionalCDATA(m.Config.ExtraData),
		AlertTemplates:       toCDATAList(m.Config.AlertTemplates),
	}
}

func logModuleXMLPayload(m LogModule, opts XMLOptions) logModuleXML {
	return logModuleXML{
		Source:   cdataText{Text: m.Config.Source},
		Data:     fmt.Sprintf("\"%s\"", m.Config.Value),
		Encoding: opts.LogEncoding,
	}
}

func normalizePandoraTimestamp(ts string) string {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts
	}
	return t.Format(pandoraTimestampLayout)
}

func scalarData(value string, points []DataPoint) *cdataText {
	if len(points) > 0 {
		return nil
	}

	return &cdataText{Text: value}
}

func dataList(points []DataPoint) *dataListXML {
	if len(points) == 0 {
		return nil
	}

	items := make([]dataPointXML, 0, len(points))
	for _, point := range points {
		items = append(items, dataPointXML{
			Value:     cdataText{Text: point.Value},
			Timestamp: optionalCDATA(normalizePandoraTimestamp(point.Timestamp)),
		})
	}

	return &dataListXML{Items: items}
}

func optionalCDATA(value string) *cdataText {
	if value == "" {
		return nil
	}

	return &cdataText{Text: value}
}

func toCDATAList(values []string) []cdataText {
	if len(values) == 0 {
		return nil
	}

	result := make([]cdataText, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		result = append(result, cdataText{Text: value})
	}

	if len(result) == 0 {
		return nil
	}

	return result
}
