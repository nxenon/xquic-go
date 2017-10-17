package mocks

//go:generate sh -c "./mockgen_internal.sh mocks stream_flow_controller.go github.com/lucas-clemente/quic-go/internal/flowcontrol StreamFlowController"
//go:generate sh -c "./mockgen_internal.sh mocks connection_flow_controller.go github.com/lucas-clemente/quic-go/internal/flowcontrol ConnectionFlowController"
//go:generate sh -c "./mockgen_stream.sh mocks stream.go github.com/lucas-clemente/quic-go StreamI"
//go:generate sh -c "goimports -w ."
