package auth_test

import (
	"testing"

	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
)

func Test_Enforcer(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		cfg      config.AdminAuthz
		viewer   auth.Viewer
		obj, act string

		can bool
	}{
		"role:admin": {
			viewer: auth.Viewer{Groups: []string{"role:admin"}},
			obj:    "teams",
			act:    "create",

			can: true,
		},
		"unauthenticated": {
			viewer: auth.Viewer{Groups: []string{"system:unauthenticated"}},
			obj:    "teams",
			act:    "create",

			can: false,
		},
		"custom group": {
			cfg: config.AdminAuthz{
				Policy: "g, ictsc:Admin, role:admin",
			},
			viewer: auth.Viewer{Groups: []string{"system:authenticated", "ictsc:Admin"}},
			obj:    "teams",
			act:    "create",

			can: true,
		},
		"custom policy": {
			cfg: config.AdminAuthz{
				Policy: "p, system:unauthenticated, teams, get",
			},
			viewer: auth.Viewer{Groups: []string{"system:unauthenticated"}},
			obj:    "teams",
			act:    "get",

			can: true,
		},
		"custom policy (deny)": {
			cfg: config.AdminAuthz{
				Policy: "p, system:unauthenticated, teams, get",
			},
			viewer: auth.Viewer{Groups: []string{"system:unauthenticated"}},
			obj:    "teams",
			act:    "create",

			can: false,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			enforcer, err := auth.NewEnforcer(tt.cfg)
			if err != nil {
				t.Fatalf("Failed to create enforcer: %v", err)
			}

			can, err := enforcer.Enforce(tt.viewer, tt.obj, tt.act)
			if err != nil {
				t.Fatalf("Failed to enforce: %v", err)
			}

			if can != tt.can {
				t.Errorf("Unexpected result: got %v, want %v", can, tt.can)
			}
		})
	}
}
