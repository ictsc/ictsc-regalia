package domain_test

import (
	"testing"

	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

func TestRedeployRuleManual_String(t *testing.T) {
	t.Parallel()

	if got := domain.RedeployRuleManual.String(); got != "Manual" {
		t.Errorf("String() = %q, want %q", got, "Manual")
	}
}

func TestRedeployRuleManual_MarshalJSON(t *testing.T) {
	t.Parallel()

	b, err := domain.RedeployRuleManual.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	if got := string(b); got != `"Manual"` {
		t.Errorf("MarshalJSON() = %q, want %q", got, `"Manual"`)
	}
}

func TestRedeployRuleManual_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	var rule domain.RedeployRule
	if err := rule.UnmarshalJSON([]byte(`"Manual"`)); err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if rule != domain.RedeployRuleManual {
		t.Errorf("UnmarshalJSON() = %v, want %v", rule, domain.RedeployRuleManual)
	}
}

func TestProblem_Manual_Redeployable(t *testing.T) {
	t.Parallel()

	problem := domain.FixDescriptiveProblemManual(t).Problem()

	if !problem.Redeployable() {
		t.Errorf("Redeployable() = false, want true")
	}
}

func TestProblem_Manual_RemainingDeployments(t *testing.T) {
	t.Parallel()

	problem := domain.FixDescriptiveProblemManual(t).Problem()

	for i := uint32(0); i <= 10; i++ {
		if got := problem.RemainingDeployments(i); got != 0 {
			t.Errorf("RemainingDeployments(%d) = %d, want 0", i, got)
		}
	}
}

func TestProblem_Manual_Penalty(t *testing.T) {
	t.Parallel()

	problem := domain.FixDescriptiveProblemManual(t).Problem()

	for i := uint32(0); i <= 10; i++ {
		if got := problem.Penalty(i); got != 0 {
			t.Errorf("Penalty(%d) = %d, want 0", i, got)
		}
	}
}
