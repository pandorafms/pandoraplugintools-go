package pandoraxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"time"

	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

const pandoraTimestampLayout = "2006/01/02 15:04:05"

// normalizePandoraTimestamp converts an RFC3339 timestamp to the Pandora
// server format (YYYY/MM/DD HH:MM:SS). If the input is not RFC3339 it is
// returned unchanged so callers that already use the Pandora format work too.
func normalizePandoraTimestamp(ts string) string {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts
	}
	return t.Format(pandoraTimestampLayout)
}

// AgentData is the internal agent payload consumed by the XML encoder.
type AgentData struct {
	AgentName       string
	AgentAlias      string
	ParentAgentName string
	Description     string
	Version         string
	OSName          string
	OSVersion       string
	Timestamp       string
	Address         string
	Group           string
	Interval        int
	AgentMode       string
}

// EncodeOptions controls optional XML output behavior.
type EncodeOptions struct {
	LogEncoding string
}

type agentXML struct {
	XMLName         xml.Name       `xml:"agent_data"`
	AgentName       string         `xml:"agent_name,attr,omitempty"`
	AgentAlias      string         `xml:"agent_alias,attr,omitempty"`
	ParentAgentName string         `xml:"parent_agent_name,attr,omitempty"`
	Description     string         `xml:"description,attr,omitempty"`
	Version         string         `xml:"version,attr,omitempty"`
	OSName          string         `xml:"os_name,attr,omitempty"`
	OSVersion       string         `xml:"os_version,attr,omitempty"`
	Timestamp       string         `xml:"timestamp,attr,omitempty"`
	Address         string         `xml:"address,attr,omitempty"`
	Group           string         `xml:"group,attr,omitempty"`
	Interval        int            `xml:"interval,attr,omitempty"`
	AgentMode       string         `xml:"agent_mode,attr,omitempty"`
	Modules         []moduleXML    `xml:"module"`
	LogModules      []logModuleXML `xml:"log_module"`
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
	Min                  *cdataText   `xml:"min,omitempty"`
	Max                  *cdataText   `xml:"max,omitempty"`
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
	Source   cdataText `xml:"source"`
	Data     string    `xml:"data"`
	Encoding string    `xml:"encoding,omitempty"`
}

// Encode converts the public model into Pandora XML.
func Encode(agent AgentData, modules []pptmodule.Module, logModules []pptmodule.LogModule, opts EncodeOptions) ([]byte, error) {
	payload := agentXML{
		AgentName:       agent.AgentName,
		AgentAlias:      agent.AgentAlias,
		ParentAgentName: agent.ParentAgentName,
		Description:     agent.Description,
		Version:         agent.Version,
		OSName:          agent.OSName,
		OSVersion:       agent.OSVersion,
		Timestamp:       normalizePandoraTimestamp(agent.Timestamp),
		Address:         agent.Address,
		Group:           agent.Group,
		Interval:        agent.Interval,
		AgentMode:       agent.AgentMode,
		Modules:         make([]moduleXML, 0, len(modules)),
		LogModules:      make([]logModuleXML, 0, len(logModules)),
	}

	for _, m := range modules {
		payload.Modules = append(payload.Modules, moduleXML{
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
		})
	}

	for _, m := range logModules {
		payload.LogModules = append(payload.LogModules, logModuleXML{
			Source:   cdataText{Text: m.Config.Source},
			Data:     fmt.Sprintf("\"%s\"", m.Config.Value),
			Encoding: opts.LogEncoding,
		})
	}

	body, err := xml.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	buf.Write(body)
	buf.WriteByte('\n')

	return buf.Bytes(), nil
}

func scalarData(value string, points []pptmodule.DataPoint) *cdataText {
	if len(points) > 0 {
		return nil
	}

	return &cdataText{Text: value}
}

func dataList(points []pptmodule.DataPoint) *dataListXML {
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
