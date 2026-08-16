package transfer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	pptoutput "github.com/pandorafms/pandoraplugintools-go/pkg/output"
	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

// Mode selects the transport mode used by Send.
type Mode string

const (
	// ModeTentacle uses the Tentacle client binary.
	ModeTentacle Mode = "tentacle"
	// ModeLocal moves the generated file into a local data directory.
	ModeLocal Mode = "local"
)

const defaultTentacleBinary = "tentacle_client"
const defaultTentaclePort = 41121
const defaultDataDir = "/var/spool/pandora/data_in"

var defaultStagingDir = os.TempDir()

// Options defines transfer behavior for Phase 1.
type Options struct {
	Mode            Mode
	TentacleBinary  string
	Address         string
	Port            int
	ExtraArgs       []string
	DataDir         string
	RemoveOnSuccess bool
	// Retries is the number of additional attempts after the first failed
	// tentacle send (tentacle mode only). Zero means a single attempt, ported
	// from discovery.py's tentacle_xml(retry=False) default.
	Retries int
}

// Validate verifies the options for the selected mode.
func (o Options) Validate() error {
	normalized := normalize(o)

	switch normalized.Mode {
	case ModeTentacle:
		if strings.TrimSpace(normalized.Address) == "" {
			return errors.New("tentacle address is required")
		}
		if normalized.Port <= 0 {
			return errors.New("tentacle port must be positive")
		}
		if normalized.Retries < 0 {
			return errors.New("retries must not be negative")
		}
	case ModeLocal:
		if strings.TrimSpace(normalized.DataDir) == "" {
			return errors.New("local data directory is required")
		}
	default:
		return fmt.Errorf("unsupported transfer mode %q", normalized.Mode)
	}

	return nil
}

// WriteXML writes a Pandora payload to a timestamped .data file.
func WriteXML(content []byte, agentName string, dir string) (string, error) {
	if len(content) == 0 {
		return "", errors.New("xml content is required")
	}

	if strings.TrimSpace(agentName) == "" {
		return "", errors.New("agent name is required")
	}

	if strings.TrimSpace(dir) == "" {
		dir = defaultStagingDir
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	name := pptutil.GenerateMD5(agentName)

	for range 10 {
		path := filepath.Join(dir, fmt.Sprintf("%s.%d.data", name, time.Now().Unix()))

		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				now := time.Now()
				nextSecond := now.Truncate(time.Second).Add(time.Second)
				time.Sleep(time.Until(nextSecond))
				continue
			}
			return "", err
		}

		if _, err := file.Write(content); err != nil {
			file.Close()
			return "", err
		}

		if err := file.Close(); err != nil {
			return "", err
		}

		return path, nil
	}

	return "", errors.New("could not allocate a unique data file name")
}

// Send delivers a generated payload using the selected mode.
func Send(ctx context.Context, file string, opts Options) error {
	if strings.TrimSpace(file) == "" {
		return errors.New("file path is required")
	}

	normalized := normalize(opts)
	if err := normalized.Validate(); err != nil {
		return err
	}

	switch normalized.Mode {
	case ModeLocal:
		destination := filepath.Join(normalized.DataDir, filepath.Base(file))
		if err := os.MkdirAll(normalized.DataDir, 0o755); err != nil {
			return err
		}
		return moveFile(file, destination)
	case ModeTentacle:
		args := []string{"-v", "-a", normalized.Address, "-p", fmt.Sprintf("%d", normalized.Port)}
		args = append(args, normalized.ExtraArgs...)
		args = append(args, file)

		attempts := 1 + normalized.Retries

		var lastErr error
		for attempt := 1; attempt <= attempts; attempt++ {
			cmd := exec.CommandContext(ctx, normalized.TentacleBinary, args...)
			output, err := cmd.CombinedOutput()
			if err == nil {
				lastErr = nil
				break
			}

			lastErr = fmt.Errorf("tentacle send failed (attempt %d/%d): %w: %s", attempt, attempts, err, strings.TrimSpace(string(output)))

			if attempt < attempts {
				pptoutput.PrintStderr("%s", lastErr.Error())
			}

			if ctx.Err() != nil {
				break
			}
		}

		if lastErr != nil {
			return lastErr
		}

		if normalized.RemoveOnSuccess {
			if err := os.Remove(file); err != nil {
				return err
			}
		}

		return nil
	default:
		return fmt.Errorf("unsupported transfer mode %q", normalized.Mode)
	}
}

func moveFile(src string, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	} else if !errors.Is(err, syscall.EXDEV) {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}

	if err := out.Close(); err != nil {
		return err
	}

	return os.Remove(src)
}

func normalize(o Options) Options {
	if o.Mode == "" {
		o.Mode = ModeTentacle
	}

	if strings.TrimSpace(o.TentacleBinary) == "" {
		o.TentacleBinary = defaultTentacleBinary
	}

	if o.Port == 0 {
		o.Port = defaultTentaclePort
	}

	if o.Mode == ModeLocal && strings.TrimSpace(o.DataDir) == "" {
		o.DataDir = defaultDataDir
	}

	return o
}
