package quic

//go:generate sh -c "./mockgen_private.sh quic mock_stream_internal_test.go github.com/lucas-clemente/quic-go streamI"
//go:generate sh -c "./mockgen_private.sh quic mock_receive_stream_internal_test.go github.com/lucas-clemente/quic-go receiveStreamI"
//go:generate sh -c "./mockgen_private.sh quic mock_send_stream_internal_test.go github.com/lucas-clemente/quic-go sendStreamI"
//go:generate sh -c "./mockgen_private.sh quic mock_stream_sender_test.go github.com/lucas-clemente/quic-go streamSender"
//go:generate sh -c "./mockgen_private.sh quic mock_stream_getter_test.go github.com/lucas-clemente/quic-go streamGetter"
//go:generate sh -c "./mockgen_private.sh quic mock_stream_frame_source_test.go github.com/lucas-clemente/quic-go streamFrameSource"
//go:generate sh -c "./mockgen_private.sh quic mock_crypto_stream_test.go github.com/lucas-clemente/quic-go cryptoStream"
//go:generate sh -c "./mockgen_private.sh quic mock_stream_manager_test.go github.com/lucas-clemente/quic-go streamManager"
//go:generate sh -c "./mockgen_private.sh quic mock_unpacker_test.go github.com/lucas-clemente/quic-go unpacker"
//go:generate sh -c "./mockgen_private.sh quic mock_quic_aead_test.go github.com/lucas-clemente/quic-go quicAEAD"
//go:generate sh -c "./mockgen_private.sh quic mock_gquic_aead_test.go github.com/lucas-clemente/quic-go gQUICAEAD"
//go:generate sh -c "./mockgen_private.sh quic mock_session_runner_test.go github.com/lucas-clemente/quic-go sessionRunner"
//go:generate sh -c "./mockgen_private.sh quic mock_quic_session_test.go github.com/lucas-clemente/quic-go quicSession"
//go:generate sh -c "./mockgen_private.sh quic mock_packet_handler_test.go github.com/lucas-clemente/quic-go packetHandler"
//go:generate sh -c "./mockgen_private.sh quic mock_unknown_packet_handler_test.go github.com/lucas-clemente/quic-go unknownPacketHandler"
//go:generate sh -c "./mockgen_private.sh quic mock_packet_handler_manager_test.go github.com/lucas-clemente/quic-go packetHandlerManager"
//go:generate sh -c "./mockgen_private.sh quic mock_multiplexer_test.go github.com/lucas-clemente/quic-go multiplexer"

//go:generate sh -c "find . -type f -name 'mock_*_test.go' | xargs sed -i '' 's/quic_go.//g'"
//go:generate sh -c "goimports -w mock*_test.go"
