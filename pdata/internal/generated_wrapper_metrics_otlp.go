// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by "model/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "go run model/internal/cmd/pdatagen/main.go".

package internal

import (
	otlpcollectormetrics "go.opentelemetry.io/collector/pdata/internal/data/protogen/collector/metrics/v1"
)

type MetricsExportPartialSuccess struct {
	orig *otlpcollectormetrics.ExportMetricsPartialSuccess
}

func GetOrigMetricsExportPartialSuccess(ms MetricsExportPartialSuccess) *otlpcollectormetrics.ExportMetricsPartialSuccess {
	return ms.orig
}

func NewMetricsExportPartialSuccess(orig *otlpcollectormetrics.ExportMetricsPartialSuccess) MetricsExportPartialSuccess {
	return MetricsExportPartialSuccess{orig: orig}
}

func GenerateTestMetricsExportPartialSuccess() MetricsExportPartialSuccess {
	orig := otlpcollectormetrics.ExportMetricsPartialSuccess{}
	tv := NewMetricsExportPartialSuccess(&orig)
	FillTestMetricsExportPartialSuccess(tv)
	return tv
}

func FillTestMetricsExportPartialSuccess(tv MetricsExportPartialSuccess) {
	tv.orig.RejectedDataPoints = int64(13)
	tv.orig.ErrorMessage = "error message"
}
