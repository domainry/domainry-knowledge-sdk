package evaluation

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type ExpectedDocument struct {
	DocumentID string  `json:"doc_id"`
	Relevance  float64 `json:"relevance"`
}

type Case struct {
	ID                 string             `json:"id"`
	Query              string             `json:"query"`
	ExpectedDocuments  []ExpectedDocument `json:"expected_documents"`
	AllowedDocumentIDs []string           `json:"allowed_document_ids"`
}

type Hit struct {
	DocumentID      string `json:"doc_id"`
	CitationCorrect *bool  `json:"citation_correct,omitempty"`
}

type Observation struct {
	CaseID string `json:"case_id"`
	Hits   []Hit  `json:"hits"`
}

type Dataset struct {
	TopK         int                `json:"top_k"`
	Cases        []Case             `json:"cases"`
	Observations []Observation      `json:"observations"`
	Thresholds   *QualityThresholds `json:"thresholds,omitempty"`
}

// QualityThresholds turns the same frozen evaluation dataset into a CI gate.
// Pointer fields distinguish an explicit zero (especially no permission
// leakage) from a metric that the caller chose not to gate.
type QualityThresholds struct {
	MinimumRecallAtK                 *float64 `json:"minimum_recall_at_k,omitempty"`
	MinimumMeanReciprocalRank        *float64 `json:"minimum_mrr,omitempty"`
	MinimumNormalizedDiscountedGain  *float64 `json:"minimum_ndcg_at_k,omitempty"`
	MinimumCitationCorrectness       *float64 `json:"minimum_citation_correctness,omitempty"`
	MinimumCitationsEvaluated        *int     `json:"minimum_citations_evaluated,omitempty"`
	MaximumPermissionLeakageRate     *float64 `json:"maximum_permission_leakage_rate,omitempty"`
	MaximumPermissionLeakageCaseRate *float64 `json:"maximum_permission_leakage_case_rate,omitempty"`
	MaximumUnauthorizedDocuments     *int     `json:"maximum_unauthorized_documents,omitempty"`
}

type GateFailure struct {
	Metric    string  `json:"metric"`
	Rule      string  `json:"rule"`
	Actual    float64 `json:"actual"`
	Threshold float64 `json:"threshold"`
}

type Report struct {
	Metrics  Metrics       `json:"metrics"`
	Passed   bool          `json:"passed"`
	Failures []GateFailure `json:"failures"`
}

type Metrics struct {
	TopK                      int     `json:"top_k"`
	Cases                     int     `json:"cases"`
	RetrievalCases            int     `json:"retrieval_cases"`
	RecallAtK                 float64 `json:"recall_at_k"`
	MeanReciprocalRank        float64 `json:"mrr"`
	NormalizedDiscountedGain  float64 `json:"ndcg_at_k"`
	CitationEvaluated         int     `json:"citation_evaluated"`
	CitationCorrect           int     `json:"citation_correct"`
	CitationCorrectness       float64 `json:"citation_correctness"`
	PermissionDocuments       int     `json:"permission_documents"`
	UnauthorizedDocuments     int     `json:"unauthorized_documents"`
	PermissionLeakageRate     float64 `json:"permission_leakage_rate"`
	PermissionLeakCases       int     `json:"permission_leak_cases"`
	PermissionLeakageCaseRate float64 `json:"permission_leakage_case_rate"`
}

func Evaluate(dataset Dataset) (Metrics, error) {
	metrics := Metrics{TopK: dataset.TopK, Cases: len(dataset.Cases)}
	if dataset.TopK < 1 || dataset.TopK > 100 || len(dataset.Cases) == 0 || len(dataset.Cases) > 100_000 || len(dataset.Observations) != len(dataset.Cases) {
		return Metrics{}, fmt.Errorf("invalid evaluation dataset dimensions")
	}
	cases := make(map[string]Case, len(dataset.Cases))
	for _, item := range dataset.Cases {
		if err := validateCase(item); err != nil {
			return Metrics{}, fmt.Errorf("case %q: %w", item.ID, err)
		}
		if _, exists := cases[item.ID]; exists {
			return Metrics{}, fmt.Errorf("duplicate evaluation case %q", item.ID)
		}
		cases[item.ID] = item
	}
	observed := make(map[string]bool, len(dataset.Observations))
	var recall, reciprocalRank, ndcg float64
	for _, observation := range dataset.Observations {
		item, exists := cases[observation.CaseID]
		if !exists || observed[observation.CaseID] || len(observation.Hits) > 10_000 {
			return Metrics{}, fmt.Errorf("invalid observation for case %q", observation.CaseID)
		}
		observed[observation.CaseID] = true
		expected := make(map[string]float64, len(item.ExpectedDocuments))
		for _, document := range item.ExpectedDocuments {
			expected[document.DocumentID] = document.Relevance
		}
		allowed := make(map[string]bool, len(item.AllowedDocumentIDs))
		for _, id := range item.AllowedDocumentIDs {
			allowed[id] = true
		}
		unique, seen := make([]Hit, 0, dataset.TopK), map[string]bool{}
		for _, hit := range observation.Hits {
			if !validText(hit.DocumentID, 4096) {
				return Metrics{}, fmt.Errorf("case %q contains an invalid hit document id", item.ID)
			}
			if hit.CitationCorrect != nil {
				metrics.CitationEvaluated++
				if *hit.CitationCorrect {
					metrics.CitationCorrect++
				}
			}
			if seen[hit.DocumentID] || len(unique) >= dataset.TopK {
				continue
			}
			seen[hit.DocumentID] = true
			unique = append(unique, hit)
		}
		for _, hit := range unique {
			metrics.PermissionDocuments++
			if !allowed[hit.DocumentID] {
				metrics.UnauthorizedDocuments++
			}
		}
		if hasUnauthorized(unique, allowed) {
			metrics.PermissionLeakCases++
		}
		if len(expected) == 0 {
			continue
		}
		metrics.RetrievalCases++
		relevant, firstRank, gain := 0, 0, 0.0
		for index, hit := range unique {
			grade := expected[hit.DocumentID]
			if grade <= 0 {
				continue
			}
			relevant++
			if firstRank == 0 {
				firstRank = index + 1
			}
			gain += (math.Pow(2, grade) - 1) / math.Log2(float64(index)+2)
		}
		recall += float64(relevant) / float64(len(expected))
		if firstRank > 0 {
			reciprocalRank += 1 / float64(firstRank)
		}
		ideal := make([]float64, 0, len(expected))
		for _, grade := range expected {
			ideal = append(ideal, grade)
		}
		sort.Slice(ideal, func(i, j int) bool { return ideal[i] > ideal[j] })
		idealGain := 0.0
		for index, grade := range ideal {
			if index >= dataset.TopK {
				break
			}
			idealGain += (math.Pow(2, grade) - 1) / math.Log2(float64(index)+2)
		}
		if idealGain > 0 {
			ndcg += gain / idealGain
		}
	}
	if len(observed) != len(cases) {
		return Metrics{}, fmt.Errorf("one or more evaluation cases have no observation")
	}
	if metrics.RetrievalCases > 0 {
		denominator := float64(metrics.RetrievalCases)
		metrics.RecallAtK = recall / denominator
		metrics.MeanReciprocalRank = reciprocalRank / denominator
		metrics.NormalizedDiscountedGain = ndcg / denominator
	}
	if metrics.CitationEvaluated > 0 {
		metrics.CitationCorrectness = float64(metrics.CitationCorrect) / float64(metrics.CitationEvaluated)
	}
	if metrics.PermissionDocuments > 0 {
		metrics.PermissionLeakageRate = float64(metrics.UnauthorizedDocuments) / float64(metrics.PermissionDocuments)
	}
	metrics.PermissionLeakageCaseRate = float64(metrics.PermissionLeakCases) / float64(metrics.Cases)
	return metrics, nil
}

// EvaluateReport calculates metrics and evaluates every explicitly configured
// threshold. A failed quality gate is data, not a malformed dataset error, so
// callers can always serialize the complete report before returning a failing
// process status to CI.
func EvaluateReport(dataset Dataset) (Report, error) {
	metrics, err := Evaluate(dataset)
	if err != nil {
		return Report{}, err
	}
	failures, err := evaluateThresholds(metrics, dataset.Thresholds)
	if err != nil {
		return Report{}, err
	}
	return Report{Metrics: metrics, Passed: len(failures) == 0, Failures: failures}, nil
}

func evaluateThresholds(metrics Metrics, thresholds *QualityThresholds) ([]GateFailure, error) {
	if thresholds == nil {
		return []GateFailure{}, nil
	}
	minimums := []struct {
		name   string
		actual float64
		value  *float64
	}{
		{name: "recall_at_k", actual: metrics.RecallAtK, value: thresholds.MinimumRecallAtK},
		{name: "mrr", actual: metrics.MeanReciprocalRank, value: thresholds.MinimumMeanReciprocalRank},
		{name: "ndcg_at_k", actual: metrics.NormalizedDiscountedGain, value: thresholds.MinimumNormalizedDiscountedGain},
		{name: "citation_correctness", actual: metrics.CitationCorrectness, value: thresholds.MinimumCitationCorrectness},
	}
	maximums := []struct {
		name   string
		actual float64
		value  *float64
	}{
		{name: "permission_leakage_rate", actual: metrics.PermissionLeakageRate, value: thresholds.MaximumPermissionLeakageRate},
		{name: "permission_leakage_case_rate", actual: metrics.PermissionLeakageCaseRate, value: thresholds.MaximumPermissionLeakageCaseRate},
	}
	failures := make([]GateFailure, 0, len(minimums)+len(maximums)+2)
	for _, item := range minimums {
		if err := validateRatioThreshold(item.name, item.value); err != nil {
			return nil, err
		}
		if item.value != nil && item.actual < *item.value {
			failures = append(failures, GateFailure{Metric: item.name, Rule: "minimum", Actual: item.actual, Threshold: *item.value})
		}
	}
	for _, item := range maximums {
		if err := validateRatioThreshold(item.name, item.value); err != nil {
			return nil, err
		}
		if item.value != nil && item.actual > *item.value {
			failures = append(failures, GateFailure{Metric: item.name, Rule: "maximum", Actual: item.actual, Threshold: *item.value})
		}
	}
	if thresholds.MinimumCitationsEvaluated != nil {
		if *thresholds.MinimumCitationsEvaluated < 0 {
			return nil, fmt.Errorf("minimum_citations_evaluated must be nonnegative")
		}
		if metrics.CitationEvaluated < *thresholds.MinimumCitationsEvaluated {
			failures = append(failures, GateFailure{Metric: "citation_evaluated", Rule: "minimum", Actual: float64(metrics.CitationEvaluated), Threshold: float64(*thresholds.MinimumCitationsEvaluated)})
		}
	}
	if thresholds.MaximumUnauthorizedDocuments != nil {
		if *thresholds.MaximumUnauthorizedDocuments < 0 {
			return nil, fmt.Errorf("maximum_unauthorized_documents must be nonnegative")
		}
		if metrics.UnauthorizedDocuments > *thresholds.MaximumUnauthorizedDocuments {
			failures = append(failures, GateFailure{Metric: "unauthorized_documents", Rule: "maximum", Actual: float64(metrics.UnauthorizedDocuments), Threshold: float64(*thresholds.MaximumUnauthorizedDocuments)})
		}
	}
	return failures, nil
}

func validateRatioThreshold(name string, value *float64) error {
	if value == nil {
		return nil
	}
	if math.IsNaN(*value) || math.IsInf(*value, 0) || *value < 0 || *value > 1 {
		return fmt.Errorf("%s threshold must be between 0 and 1", name)
	}
	return nil
}

func validateCase(item Case) error {
	if !validText(item.ID, 256) || !validText(item.Query, 4096) || item.AllowedDocumentIDs == nil || len(item.ExpectedDocuments) > 10_000 || len(item.AllowedDocumentIDs) > 100_000 {
		return fmt.Errorf("invalid identity, query, or document dimensions")
	}
	allowed := map[string]bool{}
	for _, id := range item.AllowedDocumentIDs {
		if !validText(id, 4096) || allowed[id] {
			return fmt.Errorf("invalid or duplicate allowed document id")
		}
		allowed[id] = true
	}
	expected := map[string]bool{}
	for _, document := range item.ExpectedDocuments {
		if !validText(document.DocumentID, 4096) || expected[document.DocumentID] || document.Relevance <= 0 || document.Relevance > 3 || math.IsNaN(document.Relevance) || math.IsInf(document.Relevance, 0) || !allowed[document.DocumentID] {
			return fmt.Errorf("invalid expected document")
		}
		expected[document.DocumentID] = true
	}
	return nil
}

func validText(value string, maximum int) bool {
	return value == strings.TrimSpace(value) && value != "" && len(value) <= maximum && !strings.ContainsRune(value, 0)
}

func hasUnauthorized(hits []Hit, allowed map[string]bool) bool {
	for _, hit := range hits {
		if !allowed[hit.DocumentID] {
			return true
		}
	}
	return false
}
