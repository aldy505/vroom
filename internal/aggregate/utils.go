package aggregate

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"strings"

	"github.com/getsentry/vroom/internal/calltree"
	"github.com/getsentry/vroom/internal/quantile"
)

type (
	DisplayModeType int
)

const (
	DisplayModeIOS DisplayModeType = iota
	DisplayModeAndroid
)

func newCallTreeFrameP(root *calltree.AggregateCallTree, hashOfParents []byte, displayMode DisplayModeType) Frame {
	// Compute the hash of the current frame, then append the hash of the parents
	// so that we get the hash of all of the nodes on the path to this node in the tree.
	h := computeFunctionHash(root.Image, root.Symbol)
	h.Write(hashOfParents)
	currentHash := h.Sum(nil)

	children := make([]Frame, 0, len(root.Children))
	for _, child := range root.Children {
		// Now further sub-compute the hashes of its children.
		children = append(children, newCallTreeFrameP(child, currentHash, displayMode))
	}

	var image, symbol string
	var isApplicationSymbol bool

	switch displayMode {
	case DisplayModeIOS:
		image = root.Image
		symbol = root.Symbol
		isApplicationSymbol = IsIOSApplicationImage(root.Image)
	case DisplayModeAndroid:
		image = root.Image
		symbol = root.Symbol
		isApplicationSymbol = !IsAndroidSystemPackage(root.Image)
	default:
		image = root.Image
		symbol = root.Symbol
	}

	return Frame{
		Children:              children,
		ID:                    fmt.Sprintf("%x", currentHash),
		Image:                 image,
		IsApplicationSymbol:   isApplicationSymbol,
		Line:                  root.Line,
		Path:                  root.Path,
		SelfDurationNs:        quantileToAggQuantiles(quantile.Quantile{Xs: root.SelfDurationsNs}),
		SelfDurationNsValues:  root.SelfDurationsNs,
		Symbol:                symbol,
		TotalDurationNs:       quantileToAggQuantiles(quantile.Quantile{Xs: root.TotalDurationsNs}),
		TotalDurationNsValues: root.TotalDurationsNs,
	}
}

var (
	androidPackagePrefixes = []string{
		"android.",
		"androidx.",
		"com.android.",
		"com.google.android.",
		"com.motorola.",
		"java.",
		"kotlin.",
		"kotlinx.",
	}
)

// Checking if synmbol belongs to an Android system package
func IsAndroidSystemPackage(packageName string) bool {
	for _, p := range androidPackagePrefixes {
		if strings.HasPrefix(packageName, p) {
			return true
		}
	}
	return false
}

// isApplicationSymbol determines whether the image represents that of the application
// binary (or a binary embedded in the application binary) by checking its path.
func IsIOSApplicationImage(image string) bool {
	// These are the path patterns that iOS uses for applications, system
	// libraries are stored elsewhere.
	//
	// Must be kept in sync with the corresponding Python implementation of
	// this function in python/symbolicate/__init__.py
	return strings.HasPrefix(image, "/private/var/containers") ||
		strings.HasPrefix(image, "/var/containers") ||
		strings.Contains(image, "/Developer/Xcode/DerivedData") ||
		strings.Contains(image, "/data/Containers/Bundle/Application")
}

func computeFunctionHash(image, symbol string) hash.Hash {
	h := md5.New()
	_, _ = io.WriteString(h, image)
	_, _ = io.WriteString(h, symbol)
	return h
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func removeDurationValuesFromFrameP(f Frame) {
	f.TotalDurationNsValues = nil
	f.SelfDurationNsValues = nil
	for _, c := range f.Children {
		removeDurationValuesFromFrameP(c)
	}
}

func RemoveDurationValuesFromCallTreesP(callTrees []AggregateCallTree) {
	for _, ct := range callTrees {
		removeDurationValuesFromFrameP(ct.RootFrame)
	}
}

// Structs

/** A container type to group aggregations by app_version, interaction and
 * entry_type. **/
type AggregationResult struct {
	AppVersion      string             `json:"app_version,omitempty"`
	TransactionName string             `json:"interaction,omitempty"`
	EntityType      string             `json:"entry_type,omitempty"`
	RowCount        uint32             `json:"row_count"`
	Aggregation     BacktraceAggregate `json:"aggregation"`
}

type CallTreeData struct {
	FunctionCall BacktraceAggregateFunctionCall `json:"function_call"`
	CallTrees    []AggregateCallTree            `json:"call_trees"`
}

type BacktraceAggregate struct {
	FunctionCalls       []BacktraceAggregateFunctionCall `json:"function_call"`
	FunctionToCallTrees map[string][]AggregateCallTree   `json:"function_to_call_trees"`
}

type BacktraceAggregateFunctionCall struct {
	// The name of the binary/library that the function is in
	Image string `json:"image"`
	// String representation of the function name
	Symbol string `json:"symbol"`
	// Wall time duration for the execution of the function
	DurationNs       Quantiles `json:"duration_ns"`
	DurationNsValues []float64 `json:"duration_ns_values"`

	// How frequently the function is called within a single transaction
	Frequency       Quantiles `json:"frequency"`
	FrequencyValues []float64 `json:"frequency_values"`

	// Percentage of how frequently the function is called on the main thread,
	// from [0, 1]
	MainThreadPercent float32 `json:"main_thread_percent"`
	// Map from thread name to the % of the time that the function is called
	// on that thread, from [0, 1], includes the main thread.

	ThreadNameToPercent map[string]float32 `json:"thread_name_to_percent"`
	// Line is the line number for the function in its original source file,
	// if that information is available, otherwise 0.
	Line int `json:"line"`

	// Path is the path to the original source file that contains the function,
	// if that information is available, otherwise "".
	Path string `json:"path"`

	// ProfileIDs is a unique list of the profile identifiers that this function
	// appears in.
	ProfileIDs []string `json:"profile_ids"`

	ProfileIDToThreadID map[string]uint64 `json:"profile_id_to_thread_id"`

	// The key that can be used to look up this function in the
	// function_to_call_trees map.
	Key string `json:"key"`

	// List of unique transaction names where this function is found
	TransactionNames []string `json:"transaction_names"`
}

type AggregateCallTree struct {
	// An identifier that uniquely identifies the pattern for this call tree,
	// which is computed as an MD5 hash over the image & symbol for each of
	// the frames, recursively.
	ID string `json:"id"`

	// The number of times this call tree pattern was recorded.
	Count uint64 `json:"count"`

	// Map from thread name to the number of times this call tree pattern was
	// recorded for that thread.
	ThreadNameToCount map[string]uint64 `json:"thread_name_to_count"`

	// A unique list of the trace identifiers that this call tree appears in.
	ProfileIDs []string `json:"profile_ids"`

	// The root frame of this call tree.
	RootFrame Frame `json:"root_frame"`
}

type Frame struct {
	// A stable identifier that uniquely identifies this frame within the
	// tree. This identifier is an MD5 hash of image/symbol of this node and
	// all of its parents nodes, so even when the frame for a function appears
	// multiple times in the tree, each node will have a unique ID.
	ID string `json:"id"`

	// The name of the binary/library that the function is in.
	Image string `json:"image"`

	// String representation of the function name.
	Symbol string `json:"symbol"`

	// Whether the symbol exists in application code (as opposed to system/SDK
	// code)
	IsApplicationSymbol bool `json:"is_application_symbol"`

	// Line is the line number for the function in its original source file,
	// if that information is available, otherwise 0.
	Line uint32 `json:"line"`

	// Path is the path to the original source file that contains the
	// function, if that information is available, otherwise "".
	Path string `json:"path"`

	// Wall time duration for the execution of the function and its children.
	TotalDurationNs       Quantiles `json:"total_duration_ns"`
	TotalDurationNsValues []float64 `json:"total_duration_ns_values"`

	// Child frames of this frame.
	Children []Frame `json:"children"`

	// Wall time duration for the execution of the function, excluding its
	// children.
	SelfDurationNs       Quantiles `json:"self_duration_ns"`
	SelfDurationNsValues []float64 `json:"self_duration_ns_values"`
}

type Quantiles struct {
	P50 float64 `json:"p50"`
	P75 float64 `json:"p75"`
	P90 float64 `json:"p90"`
	P95 float64 `json:"p95"`
	P99 float64 `json:"p99"`
}

// Ios & Android Profiles

type IosFrame struct {
	AbsPath         string `json:"abs_path,omitempty"`
	Filename        string `json:"filename,omitempty"`
	Function        string `json:"function,omitempty"`
	InstructionAddr string `json:"instruction_addr,omitempty"`
	Lang            string `json:"lang,omitempty"`
	LineNo          int    `json:"lineno,omitempty"`
	OriginalIndex   int    `json:"original_index,omitempty"`
	Package         string `json:"package"`
	Status          string `json:"status,omitempty"`
	SymAddr         string `json:"sym_addr,omitempty"`
	Symbol          string `json:"symbol,omitempty"`
}

type Sample struct {
	Frames              []IosFrame  `json:"frames,omitempty"`
	Priority            int         `json:"priority,omitempty"`
	QueueAddress        string      `json:"queue_address,omitempty"`
	RelativeTimestampNS interface{} `json:"relative_timestamp_ns,omitempty"`
	ThreadID            interface{} `json:"thread_id,omitempty"`
}

type IosProfile struct {
	QueueMetadata  map[string]QueueMetadata `json:"queue_metadata"`
	Samples        []Sample                 `json:"samples"`
	ThreadMetadata map[string]ThreadMedata  `json:"thread_metadata"`
}

type ThreadMedata struct {
	Name     string `json:"name"`
	Priority uint64 `json:"priority"`
}

type QueueMetadata struct {
	Label string `json:"label"`
}

type Symbol struct {
	Image    string `json:"image"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
	Line     int    `json:"line"`
}

// Path() returns (line, path, ok) where ok indicates whether the
// values are valid and can be used.
func (s *Symbol) GetPath() (int, string, bool) {
	if s.Filename != "" && s.Filename != "<compiler-generated>" {
		return s.Line, s.Path, true
	}
	return 0, "", false
}

func quantileToAggQuantiles(q quantile.Quantile) Quantiles {
	return Quantiles{
		P50: q.Percentile(0.5),
		P75: q.Percentile(0.75),
		P90: q.Percentile(0.90),
		P95: q.Percentile(0.95),
		P99: q.Percentile(0.99),
	}
}