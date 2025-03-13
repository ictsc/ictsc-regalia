package domain

import (
	"context"
	"slices"
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	Deployment struct {
		id          uuid.UUID
		teamCode    TeamCode
		problemCode ProblemCode
		revision    uint32
		events      []*DeploymentEvent
	}
	DeploymentEvent struct {
		status     DeploymentStatus
		occurredAt time.Time
	}
)

func (d *Deployment) ID() uuid.UUID {
	return d.id
}

func (d *Deployment) TeamCode() TeamCode {
	return d.teamCode
}

func (d *Deployment) ProblemCode() ProblemCode {
	return d.problemCode
}

func (d *Deployment) Revision() uint32 {
	return d.revision
}

func (d *Deployment) Status() DeploymentStatus {
	return d.events[len(d.events)-1].status
}

func (d *Deployment) CreatedAt() time.Time {
	return d.events[0].occurredAt
}

func (d *Deployment) Events() []*DeploymentEvent {
	return d.events
}

func (e *DeploymentEvent) Status() DeploymentStatus {
	return e.status
}

func (e *DeploymentEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func ListDeployments(ctx context.Context, eff DeploymentReader) ([]*Deployment, error) {
	deployments, err := eff.ListDeployments(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*Deployment, 0, len(deployments))
	for _, d := range deployments {
		deployment, err := d.parse()
		if err != nil {
			return nil, err
		}
		results = append(results, deployment)
	}
	return results, nil
}

func (tp *TeamProblem) Deployments(ctx context.Context, eff DeploymentReader) ([]*Deployment, error) {
	list, err := ListDeployments(ctx, eff)
	if err != nil {
		return nil, err
	}

	forTeam := make([]*Deployment, 0, len(list))
	for _, d := range list {
		if d.TeamCode() == tp.Team().Code() && d.ProblemCode() == tp.Problem().Code() {
			forTeam = append(forTeam, d)
		}
	}
	return forTeam, nil
}

func (tp *TeamProblem) latestDeployment(ctx context.Context, eff DeploymentReader) (*Deployment, error) {
	list, err := ListDeployments(ctx, eff)
	if err != nil {
		return nil, err
	}
	idx := slices.IndexFunc(list, func(d *Deployment) bool {
		return d.TeamCode() == tp.Team().Code() && d.ProblemCode() == tp.Problem().Code()
	})
	if idx < 0 {
		return nil, NewNotFoundError("deployment not found", nil)
	}
	return list[idx], nil
}

func (d *Deployment) UpdateStatus(ctx context.Context, eff DeploymentWriter, status DeploymentStatus, occurredAt time.Time) error {
	if !slices.Contains(d.Status().Next(), status) {
		return NewInvalidArgumentError("invalid status transition", nil)
	}
	if err := eff.UpdateDeploymentStatus(ctx, &UpdateDeploymentStatusInput{
		ID:         d.id,
		Status:     status,
		OccurredAt: occurredAt,
	}); err != nil {
		return WrapAsInternal(err, "failed to update deployment status")
	}
	d.events = append(d.events, &DeploymentEvent{
		status:     status,
		occurredAt: occurredAt,
	})
	return nil
}

func (tp *TeamProblem) Deploy(ctx context.Context, eff DeploymentWriter, occurredAt time.Time) (*Deployment, error) {
	latestDeploy, err := tp.latestDeployment(ctx, eff)
	if err != nil {
		return nil, err
	}
	if !latestDeploy.Status().IsFinished() {
		return nil, NewInvalidArgumentError("the latest deployment is not finished", nil) // InvaidArgument ではない気もする
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, WrapAsInternal(err, "failed to generate UUID")
	}
	revision := latestDeploy.Revision() + 1
	if err := eff.CreateDeployment(ctx, &CreateDeploymentInput{
		ID:         id,
		TeamID:     uuid.UUID(tp.Team().teamID),
		ProblemID:  uuid.UUID(tp.Problem().problemID),
		Revision:   revision,
		Status:     DeploymentStatusQueued,
		OccurredAt: occurredAt,
	}); err != nil {
		return nil, WrapAsInternal(err, "failed to create deployment")
	}

	return &Deployment{
		id:          id,
		teamCode:    tp.Team().Code(),
		problemCode: tp.Problem().Code(),
		revision:    revision,
		events: []*DeploymentEvent{
			{status: DeploymentStatusQueued, occurredAt: occurredAt},
		},
	}, nil
}

type (
	DeploymentData struct {
		ID          uuid.UUID              `json:"id"`
		TeamCode    int64                  `json:"team_code"`
		ProblemCode string                 `json:"problem_code"`
		Revision    uint32                 `json:"revision"`
		Events      []*DeploymentEventData `json:"events"`
	}
	DeploymentEventData struct {
		Status     DeploymentStatus `json:"status"`
		OccurredAt time.Time        `json:"occurred_at"`
	}
	DeploymentReader interface {
		ListDeployments(ctx context.Context) ([]*DeploymentData, error)
	}
	CreateDeploymentInput struct {
		ID         uuid.UUID
		TeamID     uuid.UUID
		ProblemID  uuid.UUID
		Revision   uint32
		Status     DeploymentStatus
		OccurredAt time.Time
	}
	UpdateDeploymentStatusInput struct {
		ID         uuid.UUID
		Status     DeploymentStatus
		OccurredAt time.Time
	}
	DeploymentWriter interface {
		DeploymentReader
		CreateDeployment(ctx context.Context, input *CreateDeploymentInput) error
		UpdateDeploymentStatus(ctx context.Context, input *UpdateDeploymentStatusInput) error
	}
)

func (d *Deployment) Data() *DeploymentData {
	events := make([]*DeploymentEventData, 0, len(d.events))
	for _, e := range d.events {
		events = append(events, e.Data())
	}
	return &DeploymentData{
		ID:          d.id,
		TeamCode:    int64(d.teamCode),
		ProblemCode: string(d.problemCode),
		Revision:    d.revision,
		Events:      events,
	}
}

func (e *DeploymentEvent) Data() *DeploymentEventData {
	return &DeploymentEventData{
		Status:     e.status,
		OccurredAt: e.occurredAt,
	}
}

func (d *DeploymentData) parse() (*Deployment, error) {
	teamCode, err := NewTeamCode(d.TeamCode)
	if err != nil {
		return nil, err
	}
	problemCode, err := NewProblemCode(d.ProblemCode)
	if err != nil {
		return nil, err
	}

	events := make([]*DeploymentEvent, 0, len(d.Events))
	for _, e := range d.Events {
		events = append(events, e.parse())
	}
	slices.SortFunc(events, func(i, j *DeploymentEvent) int {
		return i.occurredAt.Compare(j.occurredAt)
	})
	if len(events) == 0 {
		return nil, NewInvalidArgumentError("events must not be empty", nil)
	}
	if events[0].status != DeploymentStatusQueued {
		return nil, NewInvalidArgumentError("the first event status must be Queued", nil)
	}
	for i := range len(events) - 1 {
		if !slices.Contains(events[i].status.Next(), events[i+1].status) {
			return nil, NewInvalidArgumentError("invalid status transition", nil)
		}
	}

	return &Deployment{
		id:          d.ID,
		teamCode:    teamCode,
		problemCode: problemCode,
		revision:    d.Revision,
		events:      events,
	}, nil
}

func (e *DeploymentEventData) parse() *DeploymentEvent {
	return &DeploymentEvent{
		status:     e.Status,
		occurredAt: e.OccurredAt,
	}
}
