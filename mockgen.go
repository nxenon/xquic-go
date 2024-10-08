//go:build gomock || generate

package quic

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_send_conn_test.go github.com/nxenon/xquic-go SendConn"
type SendConn = sendConn

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_raw_conn_test.go github.com/nxenon/xquic-go RawConn"
type RawConn = rawConn

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_sender_test.go github.com/nxenon/xquic-go Sender"
type Sender = sender

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_stream_internal_test.go github.com/nxenon/xquic-go StreamI"
type StreamI = streamI

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_crypto_stream_test.go github.com/nxenon/xquic-go CryptoStream"
type CryptoStream = cryptoStream

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_receive_stream_internal_test.go github.com/nxenon/xquic-go ReceiveStreamI"
type ReceiveStreamI = receiveStreamI

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_send_stream_internal_test.go github.com/nxenon/xquic-go SendStreamI"
type SendStreamI = sendStreamI

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_stream_getter_test.go github.com/nxenon/xquic-go StreamGetter"
type StreamGetter = streamGetter

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_stream_sender_test.go github.com/nxenon/xquic-go StreamSender"
type StreamSender = streamSender

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_crypto_data_handler_test.go github.com/nxenon/xquic-go CryptoDataHandler"
type CryptoDataHandler = cryptoDataHandler

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_frame_source_test.go github.com/nxenon/xquic-go FrameSource"
type FrameSource = frameSource

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_ack_frame_source_test.go github.com/nxenon/xquic-go AckFrameSource"
type AckFrameSource = ackFrameSource

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_stream_manager_test.go github.com/nxenon/xquic-go StreamManager"
type StreamManager = streamManager

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_sealing_manager_test.go github.com/nxenon/xquic-go SealingManager"
type SealingManager = sealingManager

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_unpacker_test.go github.com/nxenon/xquic-go Unpacker"
type Unpacker = unpacker

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_packer_test.go github.com/nxenon/xquic-go Packer"
type Packer = packer

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_mtu_discoverer_test.go github.com/nxenon/xquic-go MTUDiscoverer"
type MTUDiscoverer = mtuDiscoverer

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_conn_runner_test.go github.com/nxenon/xquic-go ConnRunner"
type ConnRunner = connRunner

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_quic_conn_test.go github.com/nxenon/xquic-go QUICConn"
type QUICConn = quicConn

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_packet_handler_test.go github.com/nxenon/xquic-go PacketHandler"
type PacketHandler = packetHandler

//go:generate sh -c "go run go.uber.org/mock/mockgen -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_packet_handler_manager_test.go github.com/nxenon/xquic-go PacketHandlerManager"

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/nxenon/xquic-go -destination mock_packet_handler_manager_test.go github.com/nxenon/xquic-go PacketHandlerManager"
type PacketHandlerManager = packetHandlerManager

// Need to use source mode for the batchConn, since reflect mode follows type aliases.
// See https://github.com/golang/mock/issues/244 for details.
//
//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -package quic -self_package github.com/nxenon/xquic-go -source sys_conn_oob.go -destination mock_batch_conn_test.go -mock_names batchConn=MockBatchConn"

//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -package quic -self_package github.com/nxenon/xquic-go -self_package github.com/nxenon/xquic-go -destination mock_token_store_test.go github.com/nxenon/xquic-go TokenStore"
//go:generate sh -c "go run go.uber.org/mock/mockgen -typed -package quic -self_package github.com/nxenon/xquic-go -self_package github.com/nxenon/xquic-go -destination mock_packetconn_test.go net PacketConn"
