package server

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fanfeilong/dot_20_arch_draft/internal/analyzer"
	"github.com/fanfeilong/dot_20_arch_draft/internal/deriver"
	"github.com/fanfeilong/dot_20_arch_draft/internal/installer"
	"github.com/fanfeilong/dot_20_arch_draft/internal/reporter"
	"github.com/fanfeilong/dot_20_arch_draft/internal/state"
	"github.com/fanfeilong/dot_20_arch_draft/internal/tester"
)

func TestReportHandler(t *testing.T) {
	repo := t.TempDir()
	if _, err := installer.Install(repo); err != nil {
		t.Fatalf("Install returned error: %v", err)
	}

	target := t.TempDir()
	if _, err := analyzer.Analyze(target, repo); err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}
	if _, err := state.RecordSkill(repo, state.SkillUpdate{
		Skill:          "d2a-challenge-architecture",
		Status:         "completed",
		Stage:          state.StageArchitectureChallengeComplete,
		Phase:          "challenge-dialogue",
		QuestionIndex:  6,
		QuestionTotal:  6,
		Recommendation: "proceed",
		Summary:        "Challenge phase complete.",
	}); err != nil {
		t.Fatalf("RecordSkill returned error: %v", err)
	}
	if _, err := deriver.DeriveMini(repo); err != nil {
		t.Fatalf("DeriveMini returned error: %v", err)
	}
	if _, err := tester.PrepareTests(repo); err != nil {
		t.Fatalf("PrepareTests returned error: %v", err)
	}
	if _, err := reporter.BuildReport(repo); err != nil {
		t.Fatalf("BuildReport returned error: %v", err)
	}

	handler, err := ReportHandler(repo)
	if err != nil {
		t.Fatalf("ReportHandler returned error: %v", err)
	}

	t.Run("index html", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("unexpected status: %d", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "d2a Report") {
			t.Fatalf("unexpected body: %q", rec.Body.String())
		}
	})

	t.Run("summary json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example/data/summary.json", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("unexpected status: %d", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), filepath.Base(target)) && !strings.Contains(rec.Body.String(), target) {
			t.Fatalf("summary.json does not mention target repo: %q", rec.Body.String())
		}
	})
}

func TestReportHandlerRequiresReport(t *testing.T) {
	repo := t.TempDir()
	if _, err := installer.Install(repo); err != nil {
		t.Fatalf("Install returned error: %v", err)
	}

	if _, err := ReportHandler(repo); err == nil {
		t.Fatalf("expected error when report html is missing")
	}
}
