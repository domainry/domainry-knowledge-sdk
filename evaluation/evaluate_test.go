package evaluation

import (
	"math"
	"testing"
)

func boolean(value bool) *bool      { return &value }
func number(value float64) *float64 { return &value }
func integer(value int) *int        { return &value }

func TestEvaluateRetrievalCitationAndPermissionMetrics(t *testing.T) {
	dataset := Dataset{TopK: 4, Cases: []Case{
		{ID: "pricing", Query: "价格", ExpectedDocuments: []ExpectedDocument{{DocumentID: "a", Relevance: 3}, {DocumentID: "b", Relevance: 1}}, AllowedDocumentIDs: []string{"a", "b", "c"}},
		{ID: "policy", Query: "规则", ExpectedDocuments: []ExpectedDocument{{DocumentID: "d", Relevance: 2}}, AllowedDocumentIDs: []string{"d"}},
	}, Observations: []Observation{
		{CaseID: "pricing", Hits: []Hit{{DocumentID: "c", CitationCorrect: boolean(false)}, {DocumentID: "a", CitationCorrect: boolean(true)}, {DocumentID: "a", CitationCorrect: boolean(true)}, {DocumentID: "x", CitationCorrect: boolean(false)}, {DocumentID: "b", CitationCorrect: boolean(true)}}},
		{CaseID: "policy", Hits: []Hit{{DocumentID: "d", CitationCorrect: boolean(true)}}},
	}}
	metrics, err := Evaluate(dataset)
	if err != nil {
		t.Fatal(err)
	}
	if metrics.Cases != 2 || metrics.RetrievalCases != 2 || metrics.RecallAtK != 1 || math.Abs(metrics.MeanReciprocalRank-0.75) > 1e-9 || metrics.CitationEvaluated != 6 || metrics.CitationCorrect != 4 || math.Abs(metrics.CitationCorrectness-2.0/3.0) > 1e-9 || metrics.PermissionDocuments != 5 || metrics.UnauthorizedDocuments != 1 || math.Abs(metrics.PermissionLeakageRate-0.2) > 1e-9 || metrics.PermissionLeakCases != 1 || math.Abs(metrics.PermissionLeakageCaseRate-0.5) > 1e-9 || metrics.NormalizedDiscountedGain <= 0 || metrics.NormalizedDiscountedGain >= 1 {
		t.Fatalf("metrics=%+v", metrics)
	}
}

func TestEvaluateRejectsGroundTruthOutsideAllowedScope(t *testing.T) {
	_, err := Evaluate(Dataset{TopK: 5, Cases: []Case{{ID: "case", Query: "q", ExpectedDocuments: []ExpectedDocument{{DocumentID: "secret", Relevance: 1}}, AllowedDocumentIDs: []string{}}}, Observations: []Observation{{CaseID: "case"}}})
	if err == nil {
		t.Fatal("expected invalid ground truth")
	}
}

func TestEvaluateReportFailsConfiguredQualityAndPermissionGates(t *testing.T) {
	dataset := Dataset{TopK: 2, Cases: []Case{
		{ID: "policy", Query: "policy", ExpectedDocuments: []ExpectedDocument{{DocumentID: "allowed", Relevance: 3}}, AllowedDocumentIDs: []string{"allowed"}},
	}, Observations: []Observation{{CaseID: "policy", Hits: []Hit{
		{DocumentID: "forbidden", CitationCorrect: boolean(false)},
		{DocumentID: "allowed", CitationCorrect: boolean(true)},
	}}}, Thresholds: &QualityThresholds{
		MinimumRecallAtK:                 number(1),
		MinimumMeanReciprocalRank:        number(0.75),
		MinimumCitationCorrectness:       number(0.9),
		MinimumCitationsEvaluated:        integer(3),
		MaximumPermissionLeakageRate:     number(0),
		MaximumPermissionLeakageCaseRate: number(0),
		MaximumUnauthorizedDocuments:     integer(0),
	}}
	report, err := EvaluateReport(dataset)
	if err != nil {
		t.Fatal(err)
	}
	if report.Passed || len(report.Failures) != 6 {
		t.Fatalf("report=%+v", report)
	}
	want := []string{"mrr", "citation_correctness", "permission_leakage_rate", "permission_leakage_case_rate", "citation_evaluated", "unauthorized_documents"}
	for index, failure := range report.Failures {
		if failure.Metric != want[index] {
			t.Fatalf("failure[%d]=%+v", index, failure)
		}
	}
}

func TestEvaluateReportRejectsInvalidThreshold(t *testing.T) {
	_, err := EvaluateReport(Dataset{TopK: 1, Cases: []Case{{ID: "case", Query: "q", AllowedDocumentIDs: []string{}}}, Observations: []Observation{{CaseID: "case"}}, Thresholds: &QualityThresholds{MinimumRecallAtK: number(1.01)}})
	if err == nil {
		t.Fatal("expected invalid threshold")
	}
}
