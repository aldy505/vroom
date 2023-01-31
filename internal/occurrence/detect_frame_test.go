package occurrence

import (
	"testing"
	"time"

	"github.com/getsentry/vroom/internal/frame"
	"github.com/getsentry/vroom/internal/nodetree"
	"github.com/getsentry/vroom/internal/testutil"
)

func TestDetectFrameInCallTree(t *testing.T) {
	trueValue := true
	falseValue := false
	tests := []struct {
		job  DetectExactFrameOptions
		name string
		node *nodetree.Node
		want map[nodeKey]nodeInfo
	}{
		{
			job: DetectExactFrameOptions{
				DurationThreshold: 16 * time.Millisecond,
				FunctionsByPackage: map[string]map[string]struct{}{
					"CoreFoundation": map[string]struct{}{
						"CFReadStreamRead": {},
					},
				},
			},
			name: "Detect frame in call tree",
			node: &nodetree.Node{
				DurationNS:    uint64(30 * time.Millisecond),
				EndNS:         uint64(30 * time.Millisecond),
				Fingerprint:   0,
				IsApplication: true,
				Line:          0,
				Name:          "root",
				Package:       "package",
				Path:          "path",
				StartNS:       0,
				Children: []*nodetree.Node{
					&nodetree.Node{
						DurationNS:    uint64(20 * time.Millisecond),
						EndNS:         uint64(20 * time.Millisecond),
						Fingerprint:   0,
						IsApplication: false,
						Line:          0,
						Name:          "child1-1",
						Package:       "package",
						Path:          "path",
						StartNS:       0,
						Children: []*nodetree.Node{
							&nodetree.Node{
								DurationNS:    uint64(20 * time.Millisecond),
								EndNS:         uint64(20 * time.Millisecond),
								Fingerprint:   0,
								IsApplication: true,
								Line:          0,
								Name:          "child2-1",
								Package:       "package",
								Path:          "path",
								StartNS:       0,
								Children: []*nodetree.Node{
									&nodetree.Node{
										DurationNS:    uint64(20 * time.Millisecond),
										EndNS:         uint64(20 * time.Millisecond),
										Fingerprint:   0,
										IsApplication: false,
										Line:          0,
										Name:          "CFReadStreamRead",
										Package:       "CoreFoundation",
										Path:          "path",
										StartNS:       0,
										Children:      []*nodetree.Node{},
									},
								},
							},
						},
					},
					&nodetree.Node{
						DurationNS:    5,
						EndNS:         10,
						Fingerprint:   0,
						IsApplication: false,
						Line:          0,
						Name:          "child1-2",
						Package:       "package",
						Path:          "path",
						StartNS:       5,
						Children: []*nodetree.Node{
							&nodetree.Node{
								DurationNS:    5,
								EndNS:         10,
								Fingerprint:   0,
								IsApplication: true,
								Line:          0,
								Name:          "",
								Package:       "",
								Path:          "",
								StartNS:       5,
								Children: []*nodetree.Node{
									&nodetree.Node{
										DurationNS:    5,
										EndNS:         10,
										Fingerprint:   0,
										IsApplication: false,
										Line:          0,
										Name:          "child3-1",
										Package:       "package",
										Path:          "path",
										StartNS:       5,
										Children:      []*nodetree.Node{},
									},
								},
							},
						},
					},
				},
			},
			want: map[nodeKey]nodeInfo{
				nodeKey{
					Package:  "CoreFoundation",
					Function: "CFReadStreamRead",
				}: nodeInfo{
					Node: &nodetree.Node{
						DurationNS:    uint64(20 * time.Millisecond),
						EndNS:         uint64(20 * time.Millisecond),
						Fingerprint:   0,
						IsApplication: false,
						Line:          0,
						Name:          "CFReadStreamRead",
						Package:       "CoreFoundation",
						Path:          "path",
						StartNS:       0,
						Children:      []*nodetree.Node{},
					},
					StackTrace: []frame.Frame{
						frame.Frame{
							Function: "root",
							InApp:    &trueValue,
							Line:     0,
							Package:  "package",
							Path:     "path",
						},
						frame.Frame{
							Function: "child1-1",
							InApp:    &falseValue,
							Line:     0,
							Package:  "package",
							Path:     "path",
						},
						frame.Frame{
							Function: "child2-1",
							InApp:    &trueValue,
							Line:     0,
							Package:  "package",
							Path:     "path",
						},
						frame.Frame{
							Function: "CFReadStreamRead",
							InApp:    &falseValue,
							Line:     0,
							Package:  "CoreFoundation",
							Path:     "path",
						},
					},
				},
			},
		},
		{
			job: DetectExactFrameOptions{
				DurationThreshold: 16 * time.Millisecond,
				FunctionsByPackage: map[string]map[string]struct{}{
					"CoreFoundation": map[string]struct{}{
						"CFReadStreamRead": {},
					},
					"vroom": map[string]struct{}{
						"SuperShortFunction": {},
					},
				},
			},
			name: "Do not detect frame in call tree under threshold",
			node: &nodetree.Node{
				DurationNS:    uint64(30 * time.Millisecond),
				EndNS:         uint64(30 * time.Millisecond),
				Fingerprint:   0,
				IsApplication: true,
				Line:          0,
				Name:          "root",
				Package:       "package",
				Path:          "path",
				StartNS:       0,
				Children: []*nodetree.Node{
					&nodetree.Node{
						DurationNS:    uint64(20 * time.Millisecond),
						EndNS:         uint64(20 * time.Millisecond),
						Fingerprint:   0,
						IsApplication: false,
						Line:          0,
						Name:          "child1-1",
						Package:       "package",
						Path:          "path",
						StartNS:       0,
						Children: []*nodetree.Node{
							&nodetree.Node{
								DurationNS:    uint64(20 * time.Millisecond),
								EndNS:         uint64(20 * time.Millisecond),
								Fingerprint:   0,
								IsApplication: true,
								Line:          0,
								Name:          "child2-1",
								Package:       "package",
								Path:          "path",
								StartNS:       0,
								Children: []*nodetree.Node{
									&nodetree.Node{
										DurationNS:    uint64(10 * time.Millisecond),
										EndNS:         uint64(10 * time.Millisecond),
										Fingerprint:   0,
										IsApplication: false,
										Line:          0,
										Name:          "SuperShortFunction",
										Package:       "vroom",
										Path:          "path",
										StartNS:       0,
										Children:      []*nodetree.Node{},
									},
								},
							},
						},
					},
				},
			},
			want: map[nodeKey]nodeInfo{},
		},
		{
			job: DetectExactFrameOptions{
				DurationThreshold: 16 * time.Millisecond,
				FunctionsByPackage: map[string]map[string]struct{}{
					"vroom": map[string]struct{}{
						"FunctionWithOneSample":  {},
						"FunctionWithTwoSamples": {},
					},
				},
			},
			name: "Do not detect frame in call tree under threshold",
			node: &nodetree.Node{
				DurationNS:    uint64(30 * time.Millisecond),
				EndNS:         uint64(30 * time.Millisecond),
				Fingerprint:   0,
				IsApplication: true,
				Line:          0,
				Name:          "root",
				Package:       "package",
				Path:          "path",
				StartNS:       0,
				Children: []*nodetree.Node{
					&nodetree.Node{
						DurationNS:    uint64(20 * time.Millisecond),
						EndNS:         uint64(20 * time.Millisecond),
						Fingerprint:   0,
						IsApplication: false,
						Line:          0,
						Name:          "child1-1",
						Package:       "package",
						Path:          "path",
						StartNS:       0,
						Children: []*nodetree.Node{
							&nodetree.Node{
								DurationNS:    uint64(20 * time.Millisecond),
								EndNS:         uint64(20 * time.Millisecond),
								Fingerprint:   0,
								IsApplication: true,
								Line:          0,
								Name:          "child2-1",
								Package:       "package",
								Path:          "path",
								StartNS:       0,
								Children: []*nodetree.Node{
									&nodetree.Node{
										DurationNS:    uint64(20 * time.Millisecond),
										EndNS:         uint64(20 * time.Millisecond),
										Fingerprint:   0,
										IsApplication: false,
										Line:          0,
										Name:          "FunctionWithOneSample",
										Package:       "vroom",
										Path:          "path",
										SampleCount:   1,
										StartNS:       0,
										Children:      []*nodetree.Node{},
									},
									&nodetree.Node{
										DurationNS:    uint64(20 * time.Millisecond),
										EndNS:         uint64(20 * time.Millisecond),
										Fingerprint:   0,
										IsApplication: true,
										Line:          0,
										Name:          "child3-1",
										Package:       "package",
										Path:          "path",
										StartNS:       0,
										Children: []*nodetree.Node{
											&nodetree.Node{
												DurationNS:    uint64(20 * time.Millisecond),
												EndNS:         uint64(20 * time.Millisecond),
												Fingerprint:   0,
												IsApplication: false,
												Line:          0,
												Name:          "FunctionWithTwoSamples",
												Package:       "vroom",
												Path:          "path",
												SampleCount:   2,
												StartNS:       0,
												Children:      []*nodetree.Node{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: map[nodeKey]nodeInfo{
				nodeKey{
					Package:  "vroom",
					Function: "FunctionWithTwoSamples",
				}: nodeInfo{
					Node: &nodetree.Node{
						DurationNS:    uint64(20 * time.Millisecond),
						EndNS:         uint64(20 * time.Millisecond),
						Fingerprint:   0,
						IsApplication: false,
						Line:          0,
						Name:          "FunctionWithTwoSamples",
						Package:       "vroom",
						Path:          "path",
						SampleCount:   2,
						StartNS:       0,
						Children:      []*nodetree.Node{},
					},
					StackTrace: []frame.Frame{
						frame.Frame{
							Function: "root",
							InApp:    &trueValue,
							Line:     0,
							Package:  "package",
							Path:     "path",
						},
						frame.Frame{
							Function: "child1-1",
							InApp:    &falseValue,
							Line:     0,
							Package:  "package",
							Path:     "path",
						},
						frame.Frame{
							Function: "child2-1",
							InApp:    &trueValue,
							Line:     0,
							Package:  "package",
							Path:     "path",
						},
						frame.Frame{
							Function: "child3-1",
							InApp:    &trueValue,
							Line:     0,
							Package:  "package",
							Path:     "path",
						},
						frame.Frame{
							Function: "FunctionWithTwoSamples",
							InApp:    &falseValue,
							Line:     0,
							Package:  "vroom",
							Path:     "path",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes := make(map[nodeKey]nodeInfo)
			var stackTrace []frame.Frame
			detectFrameInCallTree(tt.node, tt.job, nodes, &stackTrace)
			if diff := testutil.Diff(nodes, tt.want); diff != "" {
				t.Fatalf("Result mismatch: got - want +\n%s", diff)
			}
		})
	}
}