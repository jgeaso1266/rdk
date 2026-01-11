// Package lerobot contains a gRPC based lerobot service server.
package lerobot

import (
	"context"

	commonpb "go.viam.com/api/common/v1"
	pb "go.viam.com/api/service/lerobot/v1"
	"go.viam.com/utils/trace"

	"go.viam.com/rdk/protoutils"
	"go.viam.com/rdk/resource"
)

// serviceServer implements the LeRobotService gRPC service.
type serviceServer struct {
	pb.UnimplementedLeRobotServiceServer
	coll resource.APIResourceGetter[Service]
}

// NewRPCServiceServer constructs a lerobot gRPC service server.
func NewRPCServiceServer(coll resource.APIResourceGetter[Service]) interface{} {
	return &serviceServer{coll: coll}
}

func (s *serviceServer) StartRecording(ctx context.Context, req *pb.StartRecordingRequest) (*pb.StartRecordingResponse, error) {
	svc, err := s.coll.Resource(req.Name)
	if err != nil {
		return nil, err
	}
	var extra map[string]interface{}
	if req.Extra != nil {
		extra = req.Extra.AsMap()
	}
	sessionID, err := svc.StartRecording(ctx, req.DatasetName, extra)
	if err != nil {
		return nil, err
	}
	return &pb.StartRecordingResponse{SessionId: sessionID}, nil
}

func (s *serviceServer) StopRecording(ctx context.Context, req *pb.StopRecordingRequest) (*pb.StopRecordingResponse, error) {
	svc, err := s.coll.Resource(req.Name)
	if err != nil {
		return nil, err
	}
	var extra map[string]interface{}
	if req.Extra != nil {
		extra = req.Extra.AsMap()
	}
	durationS, err := svc.StopRecording(ctx, req.SessionId, extra)
	if err != nil {
		return nil, err
	}
	return &pb.StopRecordingResponse{DurationS: durationS}, nil
}

func (s *serviceServer) RecordEpisode(ctx context.Context, req *pb.RecordEpisodeRequest) (*pb.RecordEpisodeResponse, error) {
	svc, err := s.coll.Resource(req.Name)
	if err != nil {
		return nil, err
	}
	var extra map[string]interface{}
	if req.Extra != nil {
		extra = req.Extra.AsMap()
	}
	resp, err := svc.RecordEpisode(ctx, RecordEpisodeRequest{
		DatasetName:  req.DatasetName,
		EpisodeIndex: req.EpisodeIndex,
		Source:       RecordingSource(req.Source),
		WarmupTimeS:  req.WarmupTimeS,
		EpisodeTimeS: req.EpisodeTimeS,
		ResetTimeS:   req.ResetTimeS,
		Fps:          req.Fps,
		Tags:         req.Tags,
		Extra:        extra,
	})
	if err != nil {
		return nil, err
	}
	return &pb.RecordEpisodeResponse{
		NumFrames:       resp.NumFrames,
		ActualDurationS: resp.ActualDurationS,
		EpisodePath:     resp.EpisodePath,
	}, nil
}

func (s *serviceServer) ReplayEpisode(ctx context.Context, req *pb.ReplayEpisodeRequest) (*pb.ReplayEpisodeResponse, error) {
	svc, err := s.coll.Resource(req.Name)
	if err != nil {
		return nil, err
	}
	var extra map[string]interface{}
	if req.Extra != nil {
		extra = req.Extra.AsMap()
	}
	resp, err := svc.ReplayEpisode(ctx, ReplayEpisodeRequest{
		DatasetName:  req.DatasetName,
		EpisodeIndex: req.EpisodeIndex,
		Fps:          req.Fps,
		Extra:        extra,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ReplayEpisodeResponse{
		NumFramesReplayed: resp.NumFramesReplayed,
		DurationS:         resp.DurationS,
	}, nil
}

func (s *serviceServer) StartTeleoperation(ctx context.Context, req *pb.StartTeleoperationRequest) (*pb.StartTeleoperationResponse, error) {
	svc, err := s.coll.Resource(req.Name)
	if err != nil {
		return nil, err
	}
	var extra map[string]interface{}
	if req.Extra != nil {
		extra = req.Extra.AsMap()
	}
	sessionID, err := svc.StartTeleoperation(ctx, StartTeleoperationRequest{
		TeleopDeviceType: req.TeleopDeviceType,
		Fps:              req.Fps,
		Extra:            extra,
	})
	if err != nil {
		return nil, err
	}
	return &pb.StartTeleoperationResponse{SessionId: sessionID}, nil
}

func (s *serviceServer) StopTeleoperation(ctx context.Context, req *pb.StopTeleoperationRequest) (*pb.StopTeleoperationResponse, error) {
	svc, err := s.coll.Resource(req.Name)
	if err != nil {
		return nil, err
	}
	var extra map[string]interface{}
	if req.Extra != nil {
		extra = req.Extra.AsMap()
	}
	durationS, err := svc.StopTeleoperation(ctx, req.SessionId, extra)
	if err != nil {
		return nil, err
	}
	return &pb.StopTeleoperationResponse{DurationS: durationS}, nil
}

func (s *serviceServer) LoadPolicy(ctx context.Context, req *pb.LoadPolicyRequest) (*pb.LoadPolicyResponse, error) {
	svc, err := s.coll.Resource(req.Name)
	if err != nil {
		return nil, err
	}
	var extra map[string]interface{}
	if req.Extra != nil {
		extra = req.Extra.AsMap()
	}
	resp, err := svc.LoadPolicy(ctx, LoadPolicyRequest{
		PolicyRepoID: req.PolicyRepoId,
		Extra:        extra,
	})
	if err != nil {
		return nil, err
	}
	return &pb.LoadPolicyResponse{
		PolicyId:   resp.PolicyID,
		PolicyType: resp.PolicyType,
	}, nil
}

func (s *serviceServer) RunPolicyEpisode(ctx context.Context, req *pb.RunPolicyEpisodeRequest) (*pb.RunPolicyEpisodeResponse, error) {
	svc, err := s.coll.Resource(req.Name)
	if err != nil {
		return nil, err
	}
	var extra map[string]interface{}
	if req.Extra != nil {
		extra = req.Extra.AsMap()
	}
	resp, err := svc.RunPolicyEpisode(ctx, RunPolicyEpisodeRequest{
		PolicyID:        req.PolicyId,
		MaxSteps:        req.MaxSteps,
		Fps:             req.Fps,
		RecordToDataset: req.RecordToDataset,
		DatasetName:     req.DatasetName,
		EpisodeIndex:    req.EpisodeIndex,
		Extra:           extra,
	})
	if err != nil {
		return nil, err
	}
	return &pb.RunPolicyEpisodeResponse{
		NumSteps:    resp.NumSteps,
		DurationS:   resp.DurationS,
		EpisodePath: resp.EpisodePath,
	}, nil
}

// DoCommand receives arbitrary commands.
func (server *serviceServer) DoCommand(ctx context.Context,
	req *commonpb.DoCommandRequest,
) (*commonpb.DoCommandResponse, error) {
	ctx, span := trace.StartSpan(ctx, "discovery::server::DoCommand")
	defer span.End()

	svc, err := server.coll.Resource(req.Name)
	if err != nil {
		return nil, err
	}
	return protoutils.DoFromResourceServer(ctx, svc, req)
}
