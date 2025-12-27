package auth

import (
	"bufio"
	_ "embed"
	"io"
	"strings"

	"github.com/casbin/casbin/v3"
	casbinmodel "github.com/casbin/casbin/v3/model"
	casbinpersist "github.com/casbin/casbin/v3/persist"
	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
)

//go:embed model.conf
var modelConfig string

//go:embed policy.csv
var builtinPolicy string

type Enforcer struct {
	enforcer *casbin.Enforcer
}

func NewEnforcer(cfg config.AdminAuthz) (*Enforcer, error) {
	model, err := casbinmodel.NewModelFromString(modelConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create model")
	}

	adapter := &adapter{
		Policy: cfg.Policy,
	}

	e, err := casbin.NewEnforcer(model, adapter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create enforcer")
	}
	return &Enforcer{enforcer: e}, nil
}

func (e *Enforcer) Enforce(sub Viewer, obj, act string) (bool, error) {
	for _, group := range sub.Groups {
		if ok, err := e.enforcer.Enforce(group, obj, act); err != nil {
			return false, errors.Wrap(err, "failed to enforce")
		} else if ok {
			return true, nil
		}
	}
	return false, nil
}

type adapter struct {
	Policy string
}

var _ casbinpersist.Adapter = (*adapter)(nil)

func (a *adapter) LoadPolicy(model casbinmodel.Model) error {
	if err := loadPolicyFromReader(strings.NewReader(builtinPolicy), model); err != nil {
		return errors.Wrap(err, "failed to load builtin policy")
	}
	if a.Policy != "" {
		if err := loadPolicyFromReader(strings.NewReader(a.Policy), model); err != nil {
			return errors.Wrap(err, "failed to load policy")
		}
	}

	return nil
}

func loadPolicyFromReader(r io.Reader, model casbinmodel.Model) error {
	var lineno int
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lineno++
		text := scanner.Text()
		text = strings.TrimSpace(text)
		if err := casbinpersist.LoadPolicyLine(text, model); err != nil {
			return errors.Wrapf(err, "failed to load policy line %d", lineno)
		}
	}
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "failed to scan policy")
	}
	return nil
}

func (a *adapter) AddPolicy(sec string, ptype string, rule []string) error {
	panic("unimplemented")
}

func (a *adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	panic("unimplemented")
}

func (a *adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	panic("unimplemented")
}

func (a *adapter) SavePolicy(model casbinmodel.Model) error {
	panic("unimplemented")
}
