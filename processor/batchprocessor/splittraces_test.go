// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package batchprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/internal/testdata"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func TestSplitTraces_noop(t *testing.T) {
	td := testdata.GenerateTraces(20)
	splitSize := 40
	split := splitTraces(splitSize, td)
	assert.Equal(t, td, split)

	i := 0
	td.ResourceSpans().At(0).ScopeSpans().At(0).Spans().RemoveIf(func(_ ptrace.Span) bool {
		i++
		return i > 5
	})
	assert.EqualValues(t, td, split)
}

func TestSplitTraces(t *testing.T) {
	td := testdata.GenerateTraces(20)
	spans := td.ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	for i := 0; i < spans.Len(); i++ {
		spans.At(i).SetName(getTestSpanName(0, i))
	}
	cp := ptrace.NewTraces()
	cpSpans := cp.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans()
	cpSpans.EnsureCapacity(5)
	td.ResourceSpans().At(0).Resource().CopyTo(
		cp.ResourceSpans().At(0).Resource())
	td.ResourceSpans().At(0).ScopeSpans().At(0).Scope().CopyTo(
		cp.ResourceSpans().At(0).ScopeSpans().At(0).Scope())
	spans.At(0).CopyTo(cpSpans.AppendEmpty())
	spans.At(1).CopyTo(cpSpans.AppendEmpty())
	spans.At(2).CopyTo(cpSpans.AppendEmpty())
	spans.At(3).CopyTo(cpSpans.AppendEmpty())
	spans.At(4).CopyTo(cpSpans.AppendEmpty())

	splitSize := 5
	split := splitTraces(splitSize, td)
	assert.Equal(t, splitSize, split.SpanCount())
	assert.Equal(t, cp, split)
	assert.Equal(t, 15, td.SpanCount())
	assert.Equal(t, "test-span-0-0", split.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Name())
	assert.Equal(t, "test-span-0-4", split.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(4).Name())

	split = splitTraces(splitSize, td)
	assert.Equal(t, 10, td.SpanCount())
	assert.Equal(t, "test-span-0-5", split.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Name())
	assert.Equal(t, "test-span-0-9", split.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(4).Name())

	split = splitTraces(splitSize, td)
	assert.Equal(t, 5, td.SpanCount())
	assert.Equal(t, "test-span-0-10", split.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Name())
	assert.Equal(t, "test-span-0-14", split.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(4).Name())

	split = splitTraces(splitSize, td)
	assert.Equal(t, 5, td.SpanCount())
	assert.Equal(t, "test-span-0-15", split.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Name())
	assert.Equal(t, "test-span-0-19", split.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(4).Name())
}

func TestSplitTracesMultipleResourceSpans(t *testing.T) {
	td := testdata.GenerateTraces(20)
	spans := td.ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	for i := 0; i < spans.Len(); i++ {
		spans.At(i).SetName(getTestSpanName(0, i))
	}
	// add second index to resource spans
	testdata.GenerateTraces(20).
		ResourceSpans().At(0).CopyTo(td.ResourceSpans().AppendEmpty())
	spans = td.ResourceSpans().At(1).ScopeSpans().At(0).Spans()
	for i := 0; i < spans.Len(); i++ {
		spans.At(i).SetName(getTestSpanName(1, i))
	}

	splitSize := 5
	split := splitTraces(splitSize, td)
	assert.Equal(t, splitSize, split.SpanCount())
	assert.Equal(t, 35, td.SpanCount())
	assert.Equal(t, "test-span-0-0", split.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Name())
	assert.Equal(t, "test-span-0-4", split.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(4).Name())
}

func TestSplitTracesMultipleResourceSpans_SplitSizeGreaterThanSpanSize(t *testing.T) {
	td := testdata.GenerateTraces(20)
	spans := td.ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	for i := 0; i < spans.Len(); i++ {
		spans.At(i).SetName(getTestSpanName(0, i))
	}
	// add second index to resource spans
	testdata.GenerateTraces(20).
		ResourceSpans().At(0).CopyTo(td.ResourceSpans().AppendEmpty())
	spans = td.ResourceSpans().At(1).ScopeSpans().At(0).Spans()
	for i := 0; i < spans.Len(); i++ {
		spans.At(i).SetName(getTestSpanName(1, i))
	}

	splitSize := 25
	split := splitTraces(splitSize, td)
	assert.Equal(t, splitSize, split.SpanCount())
	assert.Equal(t, 40-splitSize, td.SpanCount())
	assert.Equal(t, 1, td.ResourceSpans().Len())
	assert.Equal(t, "test-span-0-0", split.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Name())
	assert.Equal(t, "test-span-0-19", split.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(19).Name())
	assert.Equal(t, "test-span-1-0", split.ResourceSpans().At(1).ScopeSpans().At(0).Spans().At(0).Name())
	assert.Equal(t, "test-span-1-4", split.ResourceSpans().At(1).ScopeSpans().At(0).Spans().At(4).Name())
}

func TestSplitTracesMultipleILS(t *testing.T) {
	td := testdata.GenerateTraces(20)
	spans := td.ResourceSpans().At(0).ScopeSpans().At(0).Spans()
	for i := 0; i < spans.Len(); i++ {
		spans.At(i).SetName(getTestSpanName(0, i))
	}
	// add second index to ILS
	td.ResourceSpans().At(0).ScopeSpans().At(0).
		CopyTo(td.ResourceSpans().At(0).ScopeSpans().AppendEmpty())
	spans = td.ResourceSpans().At(0).ScopeSpans().At(1).Spans()
	for i := 0; i < spans.Len(); i++ {
		spans.At(i).SetName(getTestSpanName(1, i))
	}

	// add third index to ILS
	td.ResourceSpans().At(0).ScopeSpans().At(0).
		CopyTo(td.ResourceSpans().At(0).ScopeSpans().AppendEmpty())
	spans = td.ResourceSpans().At(0).ScopeSpans().At(2).Spans()
	for i := 0; i < spans.Len(); i++ {
		spans.At(i).SetName(getTestSpanName(2, i))
	}

	splitSize := 40
	split := splitTraces(splitSize, td)
	assert.Equal(t, splitSize, split.SpanCount())
	assert.Equal(t, 20, td.SpanCount())
	assert.Equal(t, "test-span-0-0", split.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Name())
	assert.Equal(t, "test-span-0-4", split.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(4).Name())
}
