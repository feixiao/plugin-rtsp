package rtsp

import (
	"github.com/aler9/gortsplib"
	"github.com/aler9/gortsplib/pkg/mpeg4audio"
	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/codec"
	"m7s.live/engine/v4/track"
)

type RTSPSubscriber struct {
	Subscriber
	RTSPIO
}

func (s *RTSPSubscriber) OnEvent(event any) {
	switch v := event.(type) {
	case *track.Video:
		if s.Video != nil {
			return
		}
		switch v.CodecID {
		case codec.CodecID_H264:
			vtrack := &gortsplib.TrackH264{
				PayloadType: v.PayloadType, SPS: v.ParamaterSets[0], PPS: v.ParamaterSets[1],
			}
			s.videoTrackId = len(s.tracks)
			s.tracks = append(s.tracks, vtrack)
		case codec.CodecID_H265:
			vtrack := &gortsplib.TrackH265{
				PayloadType: v.PayloadType, VPS: v.ParamaterSets[0], SPS: v.ParamaterSets[1], PPS: v.ParamaterSets[2],
			}
			s.videoTrackId = len(s.tracks)
			s.tracks = append(s.tracks, vtrack)
		}
		s.AddTrack(v)
	case *track.Audio:
		if s.Audio != nil {
			return
		}
		switch v.CodecID {
		case codec.CodecID_AAC:
			var mpegConf mpeg4audio.Config
			mpegConf.Unmarshal(v.SequenceHead[2:])
			atrack := &gortsplib.TrackMPEG4Audio{
				PayloadType: v.PayloadType, Config: &mpegConf, SizeLength: 13, IndexLength: 3, IndexDeltaLength: 3,
			}
			s.audioTrackId = len(s.tracks)
			s.tracks = append(s.tracks, atrack)
		case codec.CodecID_PCMA:
			s.audioTrackId = len(s.tracks)
			s.tracks = append(s.tracks, &gortsplib.TrackPCMA{})
		case codec.CodecID_PCMU:
			s.audioTrackId = len(s.tracks)
			s.tracks = append(s.tracks, &gortsplib.TrackPCMU{})
		}
		s.AddTrack(v)
	case ISubscriber:
		s.stream = gortsplib.NewServerStream(s.tracks)
	case VideoRTP:
		s.stream.WritePacketRTP(s.videoTrackId, &v.Packet)
	case AudioRTP:
		s.stream.WritePacketRTP(s.audioTrackId, &v.Packet)
	default:
		s.Subscriber.OnEvent(event)
	}
}
