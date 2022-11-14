// Copyright  The OpenTelemetry Authors
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

package arrow

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"

	arrowpb "github.com/f5/otel-arrow-adapter/api/collector/arrow/v1"
	arrowCollectorMock "github.com/f5/otel-arrow-adapter/api/collector/arrow/v1/mock"
	arrowRecord "github.com/f5/otel-arrow-adapter/pkg/otel/arrow_record"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/internal/testdata"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/receiver/otlpreceiver/internal/arrow/mock"
)

type commonTestCase struct {
	ctrl      *gomock.Controller
	cancel    context.CancelFunc
	telset    component.TelemetrySettings
	consumers mockConsumers
	stream    *arrowCollectorMock.MockArrowStreamService_ArrowStreamServer
	producer  *arrowRecord.Producer

	ctxCall  *gomock.Call
	sendCall *gomock.Call
	recvCall *gomock.Call
}

type testChannel struct {
	ch chan recvResult

	rlock      sync.Mutex
	sentStatus []*arrowpb.BatchStatus
	recvTraces []ptrace.Traces
	recvLogs   []plog.Logs
}

type noisyTest bool

const Noisy noisyTest = true
const NotNoisy noisyTest = false

type recvResult struct {
	payload *arrowpb.BatchArrowRecords
	err     error
}

type mockConsumers struct {
	traces  *mock.MockTraces
	logs    *mock.MockLogs
	metrics *mock.MockMetrics

	tracesCall  *gomock.Call
	logsCall    *gomock.Call
	metricsCall *gomock.Call
}

func newTestTelemetry(t *testing.T, noisy noisyTest) component.TelemetrySettings {
	telset := componenttest.NewNopTelemetrySettings()
	if !noisy {
		telset.Logger = zaptest.NewLogger(t)
	}
	return telset
}

func newTestChannel() *testChannel {
	return &testChannel{
		ch: make(chan recvResult),
	}
}

func (tc *testChannel) putBatch(payload *arrowpb.BatchArrowRecords, err error) {
	tc.ch <- recvResult{
		payload: payload,
		err:     err,
	}
}

func (tc *testChannel) getBatch() (*arrowpb.BatchArrowRecords, error) {
	r, ok := <-tc.ch
	if !ok {
		return nil, io.EOF
	}
	return r.payload, r.err
}

func (tc *testChannel) doAndReturnSentStatus(normalReturn error) func(arg *arrowpb.BatchStatus) error {
	return func(arg *arrowpb.BatchStatus) error {
		tc.rlock.Lock()
		defer tc.rlock.Unlock()
		copy := &arrowpb.BatchStatus{}
		data, err := proto.Marshal(arg)
		if err != nil {
			return err
		}
		if err := proto.Unmarshal(data, copy); err != nil {
			return err
		}

		tc.sentStatus = append(tc.sentStatus, copy)
		return normalReturn
	}
}

func (tc *testChannel) doAndReturnConsumeTraces(normalReturn error) func(ctx context.Context, traces ptrace.Traces) error {
	return func(ctx context.Context, traces ptrace.Traces) error {
		tc.rlock.Lock()
		defer tc.rlock.Unlock()
		copy := ptrace.NewTraces()
		traces.CopyTo(copy)
		tc.recvTraces = append(tc.recvTraces, traces)
		return normalReturn
	}
}

func (tc *testChannel) doAndReturnConsumeLogs(normalReturn error) func(ctx context.Context, logs plog.Logs) error {
	return func(ctx context.Context, logs plog.Logs) error {
		tc.rlock.Lock()
		defer tc.rlock.Unlock()
		copy := plog.NewLogs()
		logs.CopyTo(copy)
		tc.recvLogs = append(tc.recvLogs, logs)
		return normalReturn
	}
}

func newMockConsumers(ctrl *gomock.Controller) mockConsumers {
	mc := mockConsumers{
		traces:  mock.NewMockTraces(ctrl),
		logs:    mock.NewMockLogs(ctrl),
		metrics: mock.NewMockMetrics(ctrl),
	}
	mc.traces.EXPECT().Capabilities().Times(0)
	mc.tracesCall = mc.traces.EXPECT().ConsumeTraces(
		gomock.Any(),
		gomock.Any(),
	).Times(0)
	mc.logs.EXPECT().Capabilities().Times(0)
	mc.logsCall = mc.logs.EXPECT().ConsumeLogs(
		gomock.Any(),
		gomock.Any(),
	).Times(0)
	mc.metrics.EXPECT().Capabilities().Times(0)
	mc.metricsCall = mc.metrics.EXPECT().ConsumeMetrics(
		gomock.Any(),
		gomock.Any(),
	).Times(0)
	return mc
}

func (m mockConsumers) Traces() consumer.Traces {
	return m.traces
}

func (m mockConsumers) Logs() consumer.Logs {
	return m.logs
}
func (m mockConsumers) Metrics() consumer.Metrics {
	return m.metrics
}

var _ Consumers = mockConsumers{}

func newCommonTestCase(t *testing.T, tc *testChannel, noisy noisyTest) *commonTestCase {
	ctrl := gomock.NewController(t)
	stream := arrowCollectorMock.NewMockArrowStreamService_ArrowStreamServer(ctrl)

	ctx, cancel := context.WithCancel(context.Background())

	ctc := &commonTestCase{
		ctrl:      ctrl,
		cancel:    cancel,
		telset:    newTestTelemetry(t, noisy),
		consumers: newMockConsumers(ctrl),
		stream:    stream,
		producer:  arrowRecord.NewProducer(),
		ctxCall:   stream.EXPECT().Context().Times(0),
		recvCall:  stream.EXPECT().Recv().Times(0),
		sendCall:  stream.EXPECT().Send(gomock.Any()).Times(0),
	}

	ctc.ctxCall.AnyTimes().Return(ctx)
	ctc.recvCall.AnyTimes().DoAndReturn(tc.getBatch)

	return ctc
}

func statusOKFor(batchID string) *arrowpb.BatchStatus {
	return &arrowpb.BatchStatus{
		Statuses: []*arrowpb.StatusMessage{
			{
				BatchId:    batchID,
				StatusCode: arrowpb.StatusCode_OK,
			},
		},
	}
}

func (ctc *commonTestCase) start() func() error {
	streamErr := make(chan error)

	rcvr := New(
		config.NewComponentID("arrowtest"),
		ctc.consumers,
		component.ReceiverCreateSettings{
			TelemetrySettings: ctc.telset,
			BuildInfo:         component.NewDefaultBuildInfo(),
		})

	go func() {
		streamErr <- rcvr.ArrowStream(ctc.stream)
	}()

	return func() error {
		ctc.cancel()
		return <-streamErr
	}
}

func TestReceiverTraces(t *testing.T) {
	tc := newTestChannel()
	ctc := newCommonTestCase(t, tc, NotNoisy)

	ctc.sendCall.AnyTimes().DoAndReturn(tc.doAndReturnSentStatus(nil))
	ctc.consumers.tracesCall.AnyTimes().DoAndReturn(tc.doAndReturnConsumeTraces(nil))

	stop := ctc.start()

	td := testdata.GenerateTraces(2)
	batch, err := ctc.producer.BatchArrowRecordsFromTraces(td)
	require.NoError(t, err)

	tc.putBatch(batch, nil)

	err = stop()
	require.Error(t, err)
	require.True(t, errors.Is(err, context.Canceled))

	// EqualValues works for the underlying gogo protos.
	assert.EqualValues(t, tc.recvTraces, []ptrace.Traces{td})

	// cmp.Diff works for new-style google protobuf protos.
	require.Equal(t, "", cmp.Diff(tc.sentStatus, []*arrowpb.BatchStatus{
		statusOKFor(batch.BatchId),
	}, protocmp.Transform()))
}

func TestReceiverLogs(t *testing.T) {
	tc := newTestChannel()
	ctc := newCommonTestCase(t, tc, NotNoisy)

	ctc.sendCall.AnyTimes().DoAndReturn(tc.doAndReturnSentStatus(nil))
	ctc.consumers.logsCall.AnyTimes().DoAndReturn(tc.doAndReturnConsumeLogs(nil))

	stop := ctc.start()

	td := testdata.GenerateLogs(2)
	batch, err := ctc.producer.BatchArrowRecordsFromLogs(td)
	require.NoError(t, err)

	tc.putBatch(batch, nil)

	err = stop()
	require.Error(t, err)
	require.True(t, errors.Is(err, context.Canceled), "for %v", err)

	// EqualValues works for the underlying gogo protos.
	assert.EqualValues(t, tc.recvLogs, []plog.Logs{td})

	// cmp.Diff works for new-style google protobuf protos.
	require.Equal(t, "", cmp.Diff(tc.sentStatus, []*arrowpb.BatchStatus{
		statusOKFor(batch.BatchId),
	}, protocmp.Transform()))
}
