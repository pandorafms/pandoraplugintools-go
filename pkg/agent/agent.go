package agent

import (
	"bytes"
	"errors"
	"strings"

	"github.com/pandorafms/pandoraplugintools-go/internal/pandoraxml"
	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

const defaultInterval = 300
const defaultAgentMode = "1"

// Config defines the public agent configuration for Phase 1.
type Config struct {
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

// XMLOptions controls optional XML output behavior.
type XMLOptions struct {
	LogEncoding string
}

// Agent aggregates configuration, standard modules, log modules, and image
// modules for XML generation.
type Agent struct {
	Config       Config
	Modules      []pptmodule.Module
	LogModules   []pptmodule.LogModule
	ImageModules []pptmodule.Module
}

// New creates a new agent with Phase 1 defaults applied.
func New(cfg Config) (*Agent, error) {
	cfg = applyDefaults(cfg)

	a := &Agent{
		Config:       cfg,
		Modules:      []pptmodule.Module{},
		LogModules:   []pptmodule.LogModule{},
		ImageModules: []pptmodule.Module{},
	}
	if err := a.Validate(); err != nil {
		return nil, err
	}

	return a, nil
}

// AddModule validates and appends a module to the agent.
func (a *Agent) AddModule(m pptmodule.Module) error {
	if a == nil {
		return errors.New("agent is nil")
	}

	if err := m.Validate(); err != nil {
		return err
	}

	a.Modules = append(a.Modules, m)
	return nil
}

// AddLogModule validates and appends a log module to the agent.
func (a *Agent) AddLogModule(m pptmodule.LogModule) error {
	if a == nil {
		return errors.New("agent is nil")
	}

	if err := m.Validate(); err != nil {
		return err
	}

	a.LogModules = append(a.LogModules, m)
	return nil
}

// AddImageModule validates and appends an image module to the agent. Value
// is expected to already contain base64-encoded image data; it is rendered
// with a PNG data URI prefix, same as (pptmodule.Module).ImageXML.
func (a *Agent) AddImageModule(m pptmodule.Module) error {
	if a == nil {
		return errors.New("agent is nil")
	}

	if err := m.Validate(); err != nil {
		return err
	}

	a.ImageModules = append(a.ImageModules, m)
	return nil
}

// Validate verifies the minimum Phase 1 agent invariants.
func (a *Agent) Validate() error {
	if a == nil {
		return errors.New("agent is nil")
	}

	if strings.TrimSpace(a.Config.AgentName) == "" {
		return errors.New("agent name is required")
	}

	for _, m := range a.Modules {
		if err := m.Validate(); err != nil {
			return err
		}
	}

	for _, m := range a.LogModules {
		if err := m.Validate(); err != nil {
			return err
		}
	}

	for _, m := range a.ImageModules {
		if err := m.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// XML serializes the agent and attached modules to Pandora XML.
func (a *Agent) XML() ([]byte, error) {
	return a.XMLWithOptions(XMLOptions{})
}

// XMLWithOptions serializes the agent using optional XML output behavior.
func (a *Agent) XMLWithOptions(opts XMLOptions) ([]byte, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}

	return pandoraxml.Encode(
		pandoraxml.AgentData{
			AgentName:       a.Config.AgentName,
			AgentAlias:      a.Config.AgentAlias,
			ParentAgentName: a.Config.ParentAgentName,
			Description:     a.Config.Description,
			Version:         a.Config.Version,
			OSName:          a.Config.OSName,
			OSVersion:       a.Config.OSVersion,
			Timestamp:       a.Config.Timestamp,
			Address:         a.Config.Address,
			Group:           a.Config.Group,
			Interval:        a.Config.Interval,
			AgentMode:       a.Config.AgentMode,
		},
		a.Modules,
		a.LogModules,
		a.ImageModules,
		pandoraxml.EncodeOptions{LogEncoding: opts.LogEncoding},
	)
}

// ModulesXML serializes only the agent modules, which is useful for agent plugins
// that render module fragments directly to stdout.
func (a *Agent) ModulesXML() ([]byte, error) {
	return a.ModulesXMLWithOptions(XMLOptions{})
}

// ModulesXMLWithOptions serializes only the attached module nodes, including
// image modules (rendered via (pptmodule.Module).ImageXML) after regular and
// log modules, matching XMLWithOptions' element order.
func (a *Agent) ModulesXMLWithOptions(opts XMLOptions) ([]byte, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}

	body, err := pptmodule.XMLWithOptions(
		a.Modules,
		a.LogModules,
		pptmodule.XMLOptions{LogEncoding: opts.LogEncoding},
	)
	if err != nil {
		return nil, err
	}

	if len(a.ImageModules) == 0 {
		return body, nil
	}

	var buf bytes.Buffer
	buf.Write(body)

	for _, m := range a.ImageModules {
		imgXML, err := pptmodule.ImageXML(m)
		if err != nil {
			return nil, err
		}
		buf.Write(imgXML)
		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}

func applyDefaults(cfg Config) Config {
	if cfg.Interval == 0 {
		cfg.Interval = defaultInterval
	}

	if strings.TrimSpace(cfg.AgentMode) == "" {
		cfg.AgentMode = defaultAgentMode
	}

	if strings.TrimSpace(cfg.Timestamp) == "" {
		cfg.Timestamp = pptutil.Now()
	}

	return cfg
}
