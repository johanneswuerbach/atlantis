package events_test

import (
	"fmt"
	"testing"

	"github.com/runatlantis/atlantis/server/core/config/raw"
	"github.com/runatlantis/atlantis/server/core/config/valid"
	"github.com/runatlantis/atlantis/server/events"
	"github.com/runatlantis/atlantis/server/events/models"
	"github.com/runatlantis/atlantis/server/logging/mocks/matchers"

	. "github.com/petergtz/pegomock"

	"github.com/runatlantis/atlantis/server/events/command"
	"github.com/runatlantis/atlantis/server/events/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAggregateApplyRequirements_ValidateApplyProject(t *testing.T) {
	repoDir := "repoDir"
	fullRequirements := []string{
		raw.ApprovedRequirement,
		valid.PoliciesPassedCommandReq,
		raw.MergeableRequirement,
		raw.UnDivergedRequirement,
	}
	tests := []struct {
		name        string
		ctx         command.ProjectContext
		setup       func(workingDir *mocks.MockWorkingDir)
		wantFailure string
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:    "pass no requirements",
			ctx:     command.ProjectContext{},
			wantErr: assert.NoError,
		},
		{
			name: "pass full requirements",
			ctx: command.ProjectContext{
				ApplyRequirements: fullRequirements,
				PullReqStatus: models.PullReqStatus{
					ApprovalStatus: models.ApprovalStatus{IsApproved: true},
					Mergeable:      true,
				},
				ProjectPlanStatus: models.PassedPolicyCheckStatus,
			},
			setup: func(workingDir *mocks.MockWorkingDir) {
				When(workingDir.HasDiverged(matchers.AnyLoggingSimpleLogging(), AnyString())).ThenReturn(false)
			},
			wantErr: assert.NoError,
		},
		{
			name: "fail by no approved",
			ctx: command.ProjectContext{
				ApplyRequirements: []string{raw.ApprovedRequirement},
				PullReqStatus: models.PullReqStatus{
					ApprovalStatus: models.ApprovalStatus{IsApproved: false},
				},
			},
			wantFailure: "Pull request must be approved by at least one person other than the author before running apply.",
			wantErr:     assert.NoError,
		},
		{
			name: "fail by no policy passed",
			ctx: command.ProjectContext{
				ApplyRequirements: []string{valid.PoliciesPassedCommandReq},
				ProjectPlanStatus: models.ErroredPolicyCheckStatus,
			},
			wantFailure: "All policies must pass for project before running apply.",
			wantErr:     assert.NoError,
		},
		{
			name: "fail by no mergeable",
			ctx: command.ProjectContext{
				ApplyRequirements: []string{raw.MergeableRequirement},
				PullReqStatus:     models.PullReqStatus{Mergeable: false},
			},
			wantFailure: "Pull request must be mergeable before running apply.",
			wantErr:     assert.NoError,
		},
		{
			name: "fail by diverged",
			ctx: command.ProjectContext{
				ApplyRequirements: []string{raw.UnDivergedRequirement},
			},
			setup: func(workingDir *mocks.MockWorkingDir) {
				When(workingDir.HasDiverged(matchers.AnyLoggingSimpleLogging(), AnyString())).ThenReturn(true)
			},
			wantFailure: "Default branch must be rebased onto pull request before running apply.",
			wantErr:     assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterMockTestingT(t)
			workingDir := mocks.NewMockWorkingDir()
			a := &events.DefaultCommandRequirementHandler{WorkingDir: workingDir}
			if tt.setup != nil {
				tt.setup(workingDir)
			}
			gotFailure, err := a.ValidateApplyProject(repoDir, tt.ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("ValidateApplyProject(%v, %v)", repoDir, tt.ctx)) {
				return
			}
			assert.Equalf(t, tt.wantFailure, gotFailure, "ValidateApplyProject(%v, %v)", repoDir, tt.ctx)
		})
	}
}

func TestAggregateApplyRequirements_ValidateImportProject(t *testing.T) {
	repoDir := "repoDir"
	fullRequirements := []string{
		raw.ApprovedRequirement,
		raw.MergeableRequirement,
		raw.UnDivergedRequirement,
	}
	tests := []struct {
		name        string
		ctx         command.ProjectContext
		setup       func(workingDir *mocks.MockWorkingDir)
		wantFailure string
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:    "pass no requirements",
			ctx:     command.ProjectContext{},
			wantErr: assert.NoError,
		},
		{
			name: "pass full requirements",
			ctx: command.ProjectContext{
				ImportRequirements: fullRequirements,
				PullReqStatus: models.PullReqStatus{
					ApprovalStatus: models.ApprovalStatus{IsApproved: true},
					Mergeable:      true,
				},
				ProjectPlanStatus: models.PassedPolicyCheckStatus,
			},
			setup: func(workingDir *mocks.MockWorkingDir) {
				When(workingDir.HasDiverged(matchers.AnyLoggingSimpleLogging(), AnyString())).ThenReturn(false)
			},
			wantErr: assert.NoError,
		},
		{
			name: "fail by no approved",
			ctx: command.ProjectContext{
				ImportRequirements: []string{raw.ApprovedRequirement},
				PullReqStatus: models.PullReqStatus{
					ApprovalStatus: models.ApprovalStatus{IsApproved: false},
				},
			},
			wantFailure: "Pull request must be approved by at least one person other than the author before running import.",
			wantErr:     assert.NoError,
		},
		{
			name: "fail by no mergeable",
			ctx: command.ProjectContext{
				ImportRequirements: []string{raw.MergeableRequirement},
				PullReqStatus:      models.PullReqStatus{Mergeable: false},
			},
			wantFailure: "Pull request must be mergeable before running import.",
			wantErr:     assert.NoError,
		},
		{
			name: "fail by diverged",
			ctx: command.ProjectContext{
				ImportRequirements: []string{raw.UnDivergedRequirement},
			},
			setup: func(workingDir *mocks.MockWorkingDir) {
				When(workingDir.HasDiverged(matchers.AnyLoggingSimpleLogging(), AnyString())).ThenReturn(true)
			},
			wantFailure: "Default branch must be rebased onto pull request before running import.",
			wantErr:     assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterMockTestingT(t)
			workingDir := mocks.NewMockWorkingDir()
			a := &events.DefaultCommandRequirementHandler{WorkingDir: workingDir}
			if tt.setup != nil {
				tt.setup(workingDir)
			}
			gotFailure, err := a.ValidateImportProject(repoDir, tt.ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("ValidateImportProject(%v, %v)", repoDir, tt.ctx)) {
				return
			}
			assert.Equalf(t, tt.wantFailure, gotFailure, "ValidateImportProject(%v, %v)", repoDir, tt.ctx)
		})
	}
}
