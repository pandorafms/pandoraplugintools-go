package agent

import (
	"errors"
	"strings"
	"time"

	"github.com/pandorafms/pandoraPlugintoolsGo/internal/pandoraxml"
	pptmodule "github.com/pandorafms/pandoraPlugintoolsGo/pkg/module"
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

// Agent aggregates configuration, standard modules, and log modules for XML generation.
type Agent struct {
	Config     Config
	Modules    []pptmodule.Module
	LogModules []pptmodule.LogModule
}

// New creates a new agent with Phase 1 defaults applied.
func New(cfg Config) (*Agent, error) {
	cfg = applyDefaults(cfg)

	a := &Agent{
		Config:     cfg,
		Modules:    []pptmodule.Module{},
		LogModules: []pptmodule.LogModule{},
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
		pandoraxml.EncodeOptions{LogEncoding: opts.LogEncoding},
	)
}

func applyDefaults(cfg Config) Config {
	if cfg.Interval == 0 {
		cfg.Interval = defaultInterval
	}

	if strings.TrimSpace(cfg.AgentMode) == "" {
		cfg.AgentMode = defaultAgentMode
	}

	if strings.TrimSpace(cfg.Timestamp) == "" {
		cfg.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	return cfg
}
