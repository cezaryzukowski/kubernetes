package topologymanager

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
)

func TestCtnCalculateAffinity(t *testing.T) {
	tcases := []struct {
		name     string
		hp       []HintProvider
		expected []map[string][]TopologyHint
	}{
		{
			name:     "No hint providers",
			hp:       []HintProvider{},
			expected: ([]map[string][]TopologyHint)(nil),
		},
		{
			name: "HintProvider returns empty non-nil map[string][]TopologyHint",
			hp: []HintProvider{
				&mockHintProvider{
					map[string][]TopologyHint{},
				},
			},
			expected: []map[string][]TopologyHint{
				{},
			},
		},
		{
			name: "HintProvider returns -nil map[string][]TopologyHint from provider",
			hp: []HintProvider{
				&mockHintProvider{
					map[string][]TopologyHint{
						"resource": nil,
					},
				},
			},
			expected: []map[string][]TopologyHint{
				{
					"resource": nil,
				},
			},
		},
		{
			name: "Assorted HintProviders",
			hp: []HintProvider{
				&mockHintProvider{
					map[string][]TopologyHint{
						"resource-1/A": {
							{NUMANodeAffinity: NewTestBitMask(0), Preferred: true},
							{NUMANodeAffinity: NewTestBitMask(0, 1), Preferred: false},
						},
						"resource-1/B": {
							{NUMANodeAffinity: NewTestBitMask(1), Preferred: true},
							{NUMANodeAffinity: NewTestBitMask(1, 2), Preferred: false},
						},
					},
				},
				&mockHintProvider{
					map[string][]TopologyHint{
						"resource-2/A": {
							{NUMANodeAffinity: NewTestBitMask(2), Preferred: true},
							{NUMANodeAffinity: NewTestBitMask(3, 4), Preferred: false},
						},
						"resource-2/B": {
							{NUMANodeAffinity: NewTestBitMask(2), Preferred: true},
							{NUMANodeAffinity: NewTestBitMask(3, 4), Preferred: false},
						},
					},
				},
				&mockHintProvider{
					map[string][]TopologyHint{
						"resource-3": nil,
					},
				},
			},
			expected: []map[string][]TopologyHint{
				{
					"resource-1/A": {
						{NUMANodeAffinity: NewTestBitMask(0), Preferred: true},
						{NUMANodeAffinity: NewTestBitMask(0, 1), Preferred: false},
					},
					"resource-1/B": {
						{NUMANodeAffinity: NewTestBitMask(1), Preferred: true},
						{NUMANodeAffinity: NewTestBitMask(1, 2), Preferred: false},
					},
				},
				{
					"resource-2/A": {
						{NUMANodeAffinity: NewTestBitMask(2), Preferred: true},
						{NUMANodeAffinity: NewTestBitMask(3, 4), Preferred: false},
					},
					"resource-2/B": {
						{NUMANodeAffinity: NewTestBitMask(2), Preferred: true},
						{NUMANodeAffinity: NewTestBitMask(3, 4), Preferred: false},
					},
				},
				{
					"resource-3": nil,
				},
			},
		},
	}

	for _, tc := range tcases {
		ctnScope := &containerScope{
			scope{
				hintProviders: tc.hp,
				policy:        &mockPolicy{},
				name:          podTopologyScope,
			},
		}

		ctnScope.calculateAffinity(&v1.Pod{}, &v1.Container{})
		actual := ctnScope.policy.(*mockPolicy).ph
		if !reflect.DeepEqual(tc.expected, actual) {
			t.Errorf("Test Case: %s", tc.name)
			t.Errorf("Expected result to be %v, got %v", tc.expected, actual)
		}
	}
}

func TestCtnAccumulateProvidersHints(t *testing.T) {
	tcases := []struct {
		name     string
		hp       []HintProvider
		expected []map[string][]TopologyHint
	}{
		{
			name:     "TopologyHint not set",
			hp:       []HintProvider{},
			expected: nil,
		},
		{
			name: "HintProvider returns empty non-nil map[string][]TopologyHint",
			hp: []HintProvider{
				&mockHintProvider{
					map[string][]TopologyHint{},
				},
			},
			expected: []map[string][]TopologyHint{
				{},
			},
		},
		{
			name: "HintProvider returns - nil map[string][]TopologyHint from provider",
			hp: []HintProvider{
				&mockHintProvider{
					map[string][]TopologyHint{
						"resource": nil,
					},
				},
			},
			expected: []map[string][]TopologyHint{
				{
					"resource": nil,
				},
			},
		},
		{
			name: "2 HintProviders with 1 resource returns hints",
			hp: []HintProvider{
				&mockHintProvider{
					map[string][]TopologyHint{
						"resource1": {TopologyHint{}},
					},
				},
				&mockHintProvider{
					map[string][]TopologyHint{
						"resource2": {TopologyHint{}},
					},
				},
			},
			expected: []map[string][]TopologyHint{
				{
					"resource1": {TopologyHint{}},
				},
				{
					"resource2": {TopologyHint{}},
				},
			},
		},
		{
			name: "2 HintProviders 1 with 1 resource 1 with nil hints",
			hp: []HintProvider{
				&mockHintProvider{
					map[string][]TopologyHint{
						"resource1": {TopologyHint{}},
					},
				},
				&mockHintProvider{nil},
			},
			expected: []map[string][]TopologyHint{
				{
					"resource1": {TopologyHint{}},
				},
				nil,
			},
		},
		{
			name: "2 HintProviders 1 with 1 resource 1 empty hints",
			hp: []HintProvider{
				&mockHintProvider{
					map[string][]TopologyHint{
						"resource1": {TopologyHint{}},
					},
				},
				&mockHintProvider{
					map[string][]TopologyHint{},
				},
			},
			expected: []map[string][]TopologyHint{
				{
					"resource1": {TopologyHint{}},
				},
				{},
			},
		},
		{
			name: "HintProvider with 2 resources returns hints",
			hp: []HintProvider{
				&mockHintProvider{
					map[string][]TopologyHint{
						"resource1": {TopologyHint{}},
						"resource2": {TopologyHint{}},
					},
				},
			},
			expected: []map[string][]TopologyHint{
				{
					"resource1": {TopologyHint{}},
					"resource2": {TopologyHint{}},
				},
			},
		},
	}

	for _, tc := range tcases {
		ctnScope := containerScope{
			scope{
				hintProviders: tc.hp,
			},
		}
		actual := ctnScope.accumulateProvidersHints(&v1.Pod{}, &v1.Container{})
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("Test Case %s: Expected NUMANodeAffinity in result to be %v, got %v", tc.name, tc.expected, actual)
		}
	}
}
