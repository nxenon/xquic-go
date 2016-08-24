package quic

import (
	"errors"

	"github.com/lucas-clemente/quic-go/protocol"
	"github.com/lucas-clemente/quic-go/qerr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Streams Map", func() {
	var (
		m *streamsMap
	)

	BeforeEach(func() {
		m = newStreamsMap(nil)
	})

	Context("getting and creating streams", func() {
		BeforeEach(func() {
			m.newStream = func(id protocol.StreamID) (*stream, error) {
				return &stream{streamID: id}, nil
			}
		})

		It("gets new streams", func() {
			s, err := m.GetOrOpenStream(5)
			Expect(err).NotTo(HaveOccurred())
			Expect(s.StreamID()).To(Equal(protocol.StreamID(5)))
		})

		It("rejects streams with even IDs", func() {
			_, err := m.GetOrOpenStream(6)
			Expect(err).To(MatchError("InvalidStreamID: attempted to open stream 6 from client-side"))
		})

		It("gets existing streams", func() {
			s, err := m.GetOrOpenStream(5)
			Expect(err).NotTo(HaveOccurred())
			s, err = m.GetOrOpenStream(5)
			Expect(err).NotTo(HaveOccurred())
			Expect(s.StreamID()).To(Equal(protocol.StreamID(5)))
		})

		It("returns nil for closed streams", func() {
			s, err := m.GetOrOpenStream(5)
			Expect(err).NotTo(HaveOccurred())
			err = m.RemoveStream(5)
			Expect(err).NotTo(HaveOccurred())
			s, err = m.GetOrOpenStream(5)
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(BeNil())
		})

		It("panics on OpenStream", func() {
			Expect(func() { m.OpenStream(0) }).To(Panic())
		})

		Context("counting streams", func() {
			It("errors when too many streams are opened", func() {
				for i := 0; i < m.maxNumStreams; i++ {
					_, err := m.GetOrOpenStream(protocol.StreamID(i*2 + 1))
					Expect(err).NotTo(HaveOccurred())
				}
				_, err := m.GetOrOpenStream(protocol.StreamID(m.maxNumStreams))
				Expect(err).To(MatchError(qerr.TooManyOpenStreams))
			})

			It("does not error when many streams are opened and closed", func() {
				for i := 2; i < 10*m.maxNumStreams; i++ {
					_, err := m.GetOrOpenStream(protocol.StreamID(i*2 + 1))
					Expect(err).NotTo(HaveOccurred())
					m.RemoveStream(protocol.StreamID(i*2 + 1))
				}
			})
		})

		Context("DoS mitigation", func() {
			It("opens and closes a lot of streams", func() {
				for i := 1; i < 2*protocol.MaxNewStreamIDDelta; i += 2 {
					streamID := protocol.StreamID(i)
					_, err := m.GetOrOpenStream(streamID)
					Expect(m.highestStreamOpenedByClient).To(Equal(streamID))
					Expect(err).NotTo(HaveOccurred())
					err = m.RemoveStream(streamID)
					Expect(err).NotTo(HaveOccurred())
				}
			})

			It("prevents opening of streams with very low StreamIDs, if higher streams have already been opened", func() {
				for i := 1; i < protocol.MaxNewStreamIDDelta+14; i += 2 {
					if i == 11 || i == 13 {
						continue
					}
					streamID := protocol.StreamID(i)
					_, err := m.GetOrOpenStream(streamID)
					Expect(err).NotTo(HaveOccurred())
					err = m.RemoveStream(streamID)
					Expect(err).NotTo(HaveOccurred())
				}
				Expect(m.highestStreamOpenedByClient).To(Equal(protocol.StreamID(protocol.MaxNewStreamIDDelta + 13)))
				_, err := m.GetOrOpenStream(11)
				Expect(err).To(MatchError("InvalidStreamID: attempted to open stream 11, which is a lot smaller than the highest opened stream, 413"))
				_, err = m.GetOrOpenStream(13)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("deleting streams", func() {
		BeforeEach(func() {
			for i := 1; i <= 5; i++ {
				err := m.putStream(&stream{streamID: protocol.StreamID(i)})
				Expect(err).ToNot(HaveOccurred())
			}
			Expect(m.openStreams).To(Equal([]protocol.StreamID{1, 2, 3, 4, 5}))
		})

		It("errors when removing non-existing stream", func() {
			err := m.RemoveStream(1337)
			Expect(err).To(MatchError("attempted to remove non-existing stream: 1337"))
		})

		It("removes the first stream", func() {
			err := m.RemoveStream(1)
			Expect(err).ToNot(HaveOccurred())
			Expect(m.openStreams).To(HaveLen(4))
			Expect(m.openStreams).To(Equal([]protocol.StreamID{2, 3, 4, 5}))
		})

		It("removes a stream in the middle", func() {
			err := m.RemoveStream(3)
			Expect(err).ToNot(HaveOccurred())
			Expect(m.openStreams).To(HaveLen(4))
			Expect(m.openStreams).To(Equal([]protocol.StreamID{1, 2, 4, 5}))
		})

		It("removes a stream at the end", func() {
			err := m.RemoveStream(5)
			Expect(err).ToNot(HaveOccurred())
			Expect(m.openStreams).To(HaveLen(4))
			Expect(m.openStreams).To(Equal([]protocol.StreamID{1, 2, 3, 4}))
		})

		It("removes all streams", func() {
			for i := 1; i <= 5; i++ {
				err := m.RemoveStream(protocol.StreamID(i))
				Expect(err).ToNot(HaveOccurred())
			}
			Expect(m.openStreams).To(BeEmpty())
		})
	})

	Context("number of streams", func() {
		It("returns 0 in the beginning", func() {
			Expect(m.NumberOfStreams()).To(Equal(0))
		})

		It("increases the counter when a new stream is added", func() {
			err := m.putStream(&stream{streamID: 5})
			Expect(err).ToNot(HaveOccurred())
			Expect(m.NumberOfStreams()).To(Equal(1))
		})

		It("decreases the counter when removing a stream", func() {
			err := m.putStream(&stream{streamID: 5})
			Expect(err).ToNot(HaveOccurred())
			err = m.RemoveStream(5)
			Expect(err).ToNot(HaveOccurred())
			Expect(m.NumberOfStreams()).To(BeZero())
		})
	})

	Context("Iterate", func() {
		// create 3 streams, ids 1 to 3
		BeforeEach(func() {
			for i := 1; i <= 3; i++ {
				err := m.putStream(&stream{streamID: protocol.StreamID(i)})
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("executes the lambda exactly once for every stream", func() {
			var numIterations int
			callbackCalled := make(map[protocol.StreamID]bool)
			fn := func(str *stream) (bool, error) {
				callbackCalled[str.StreamID()] = true
				numIterations++
				return true, nil
			}
			err := m.Iterate(fn)
			Expect(err).ToNot(HaveOccurred())
			Expect(callbackCalled).To(HaveKey(protocol.StreamID(1)))
			Expect(callbackCalled).To(HaveKey(protocol.StreamID(2)))
			Expect(callbackCalled).To(HaveKey(protocol.StreamID(3)))
			Expect(numIterations).To(Equal(3))
		})

		It("stops iterating when the callback returns false", func() {
			var numIterations int
			fn := func(str *stream) (bool, error) {
				numIterations++
				return false, nil
			}
			err := m.Iterate(fn)
			Expect(err).ToNot(HaveOccurred())
			// due to map access randomization, we don't know for which stream the callback was executed
			// but it must only be executed once
			Expect(numIterations).To(Equal(1))
		})

		It("returns the error, if the lambda returns one", func() {
			var numIterations int
			expectedError := errors.New("test")
			fn := func(str *stream) (bool, error) {
				numIterations++
				return true, expectedError
			}
			err := m.Iterate(fn)
			Expect(err).To(MatchError(expectedError))
			Expect(numIterations).To(Equal(1))
		})
	})

	Context("RoundRobinIterate", func() {
		// create 5 streams, ids 1 to 5
		var lambdaCalledForStream []protocol.StreamID
		var numIterations int

		BeforeEach(func() {
			lambdaCalledForStream = lambdaCalledForStream[:0]
			numIterations = 0
			for i := 1; i <= 5; i++ {
				err := m.putStream(&stream{streamID: protocol.StreamID(i)})
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("executes the lambda exactly once for every stream", func() {
			fn := func(str *stream) (bool, error) {
				lambdaCalledForStream = append(lambdaCalledForStream, str.StreamID())
				numIterations++
				return true, nil
			}
			err := m.RoundRobinIterate(fn)
			Expect(err).ToNot(HaveOccurred())
			Expect(numIterations).To(Equal(5))
			Expect(lambdaCalledForStream).To(Equal([]protocol.StreamID{1, 2, 3, 4, 5}))
			Expect(m.roundRobinIndex).To(BeZero())
		})

		It("goes around once when starting in the middle", func() {
			fn := func(str *stream) (bool, error) {
				lambdaCalledForStream = append(lambdaCalledForStream, str.StreamID())
				numIterations++
				return true, nil
			}
			m.roundRobinIndex = 3
			err := m.RoundRobinIterate(fn)
			Expect(err).ToNot(HaveOccurred())
			Expect(numIterations).To(Equal(5))
			Expect(lambdaCalledForStream).To(Equal([]protocol.StreamID{4, 5, 1, 2, 3}))
			Expect(m.roundRobinIndex).To(Equal(3))
		})

		It("picks up at the index where it last stopped", func() {
			fn := func(str *stream) (bool, error) {
				lambdaCalledForStream = append(lambdaCalledForStream, str.StreamID())
				numIterations++
				if str.StreamID() == 2 {
					return false, nil
				}
				return true, nil
			}
			err := m.RoundRobinIterate(fn)
			Expect(err).ToNot(HaveOccurred())
			Expect(numIterations).To(Equal(2))
			Expect(lambdaCalledForStream).To(Equal([]protocol.StreamID{1, 2}))
			Expect(m.roundRobinIndex).To(Equal(1))
			numIterations = 0
			lambdaCalledForStream = lambdaCalledForStream[:0]
			fn2 := func(str *stream) (bool, error) {
				lambdaCalledForStream = append(lambdaCalledForStream, str.StreamID())
				numIterations++
				if str.StreamID() == 4 {
					return false, nil
				}
				return true, nil
			}
			err = m.RoundRobinIterate(fn2)
			Expect(err).ToNot(HaveOccurred())
			Expect(numIterations).To(Equal(3))
			Expect(lambdaCalledForStream).To(Equal([]protocol.StreamID{2, 3, 4}))
		})

		It("adjust the RoundRobinIndex when deleting an element in front", func() {
			m.roundRobinIndex = 3 // stream 4
			m.RemoveStream(2)
			Expect(m.roundRobinIndex).To(Equal(2))
		})

		It("doesn't adjust the RoundRobinIndex when deleting an element at the back", func() {
			m.roundRobinIndex = 1 // stream 2
			m.RemoveStream(4)
			Expect(m.roundRobinIndex).To(Equal(1))
		})

		It("doesn't adjust the RoundRobinIndex when deleting the element it is pointing to", func() {
			m.roundRobinIndex = 3 // stream 4
			m.RemoveStream(4)
			Expect(m.roundRobinIndex).To(Equal(3))
		})
	})
})
