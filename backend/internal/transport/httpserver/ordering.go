package httpserver

import (
	"sort"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

// orderedContentProblems expands the manifest's section references instead of
// relying on the storage order of the top-level problem objects. This preserves
// the author-controlled order inside each section while keeping sections in
// chronological order.
func orderedContentProblems(snapshot core.ContentSnapshot) []core.Problem {
	problemByCode := make(map[string]core.Problem, len(snapshot.Manifest.Problems))
	for _, problem := range snapshot.Manifest.Problems {
		problemByCode[problem.Code] = problem
	}
	sections := append([]core.Section(nil), snapshot.Manifest.Sections...)
	sort.SliceStable(sections, func(i, j int) bool {
		if sections[i].Beginning.Equal(sections[j].Beginning) {
			return sections[i].Slug < sections[j].Slug
		}
		return sections[i].Beginning.Before(sections[j].Beginning)
	})
	problems := make([]core.Problem, 0, len(snapshot.Manifest.Problems))
	for _, section := range sections {
		for _, code := range section.ProblemIDs {
			if problem, ok := problemByCode[code]; ok {
				problems = append(problems, problem)
			}
		}
	}
	return problems
}

func sortAdminDeployments(deployments []api.AdminDeployment) {
	sort.SliceStable(deployments, func(i, j int) bool {
		if deployments[i].TeamCode != deployments[j].TeamCode {
			return deployments[i].TeamCode < deployments[j].TeamCode
		}
		if deployments[i].ProblemCode != deployments[j].ProblemCode {
			return deployments[i].ProblemCode < deployments[j].ProblemCode
		}
		return deployments[i].Revision > deployments[j].Revision
	})
}

func sortContestantDeployments(deployments []api.ContestantDeployment) {
	sort.SliceStable(deployments, func(i, j int) bool {
		return deployments[i].Revision > deployments[j].Revision
	})
}
