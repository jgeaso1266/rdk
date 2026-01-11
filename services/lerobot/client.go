// Package lerobot contains a gRPC based lerobot client.
package lerobot

import (
	"context"

	commonpb "go.viam.com/api/common/v1"
	pb "go.viam.com/api/service/lerobot/v1"
	"go.viam.com/utils/protoutils"
	"go.viam.com/utils/rpc"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

// client implements Service.
type client struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable
	name   string
	client pb.LeRobotServiceClient
	logger logging.Logger
}

// NewClientFromConn constructs a new Client from connection passed in.
func NewClientFromConn(
	ctx context.Context,
	conn rpc.ClientConn,
	remoteName string,
	name resource.Name,
	logger logging.Logger,
) (Service, error) {
	c := pb.NewLeRobotServiceClient(conn)
	return &client{
		Named:  name.PrependRemote(remoteName).AsNamed(),
		name:   name.Name,
		client: c,
		logger: logger,
	}, nil
}

func (c *client) StartRecording(ctx context.Context, datasetName string, extra map[string]interface{}) (string, error) {
	ext, err := protoutils.StructToStructPb(extra)
	if err != nil {
		return "", err
	}
	resp, err := c.client.StartRecording(ctx, &pb.StartRecordingRequest{
		Name:        c.name,
		DatasetName: datasetName,
		Extra:       ext,
	})
	if err != nil {
		return "", err
	}
	return resp.SessionId, nil
}

func (c *client) StopRecording(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error) {
	ext, err := protoutils.StructToStructPb(extra)
	if err != nil {
		return 0, err
	}
	resp, err := c.client.StopRecording(ctx, &pb.StopRecordingRequest{
		Name:      c.name,
		SessionId: sessionID,
		Extra:     ext,
	})
	if err != nil {
		return 0, err
	}
	return resp.DurationS, nil
}

func (c *client) RecordEpisode(ctx context.Context, req RecordEpisodeRequest) (*RecordEpisodeResponse, error) {
	ext, err := protoutils.StructToStructPb(req.Extra)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.RecordEpisode(ctx, &pb.RecordEpisodeRequest{
		Name:         c.name,
		DatasetName:  req.DatasetName,
		EpisodeIndex: req.EpisodeIndex,
		Source:       pb.RecordingSource(req.Source),
		WarmupTimeS:  req.WarmupTimeS,
		EpisodeTimeS: req.EpisodeTimeS,
		ResetTimeS:   req.ResetTimeS,
		Fps:          req.Fps,
		Tags:         req.Tags,
		Extra:        ext,
	})
	if err != nil {
		return nil, err
	}
	return &RecordEpisodeResponse{
		NumFrames:       resp.NumFrames,
		ActualDurationS: resp.ActualDurationS,
		EpisodePath:     resp.EpisodePath,
	}, nil
}

func (c *client) ReplayEpisode(ctx context.Context, req ReplayEpisodeRequest) (*ReplayEpisodeResponse, error) {
	ext, err := protoutils.StructToStructPb(req.Extra)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.ReplayEpisode(ctx, &pb.ReplayEpisodeRequest{
		Name:         c.name,
		DatasetName:  req.DatasetName,
		EpisodeIndex: req.EpisodeIndex,
		Fps:          req.Fps,
		Extra:        ext,
	})
	if err != nil {
		return nil, err
	}
	return &ReplayEpisodeResponse{
		NumFramesReplayed: resp.NumFramesReplayed,
		DurationS:         resp.DurationS,
	}, nil
}

func (c *client) StartTeleoperation(ctx context.Context, req StartTeleoperationRequest) (string, error) {
	ext, err := protoutils.StructToStructPb(req.Extra)
	if err != nil {
		return "", err
	}
	resp, err := c.client.StartTeleoperation(ctx, &pb.StartTeleoperationRequest{
		Name:             c.name,
		TeleopDeviceType: req.TeleopDeviceType,
		Fps:              req.Fps,
		Extra:            ext,
	})
	if err != nil {
		return "", err
	}
	return resp.SessionId, nil
}

func (c *client) StopTeleoperation(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error) {
	ext, err := protoutils.StructToStructPb(extra)
	if err != nil {
		return 0, err
	}
	resp, err := c.client.StopTeleoperation(ctx, &pb.StopTeleoperationRequest{
		Name:      c.name,
		SessionId: sessionID,
		Extra:     ext,
	})
	if err != nil {
		return 0, err
	}
	return resp.DurationS, nil
}

func (c *client) LoadPolicy(ctx context.Context, req LoadPolicyRequest) (*LoadPolicyResponse, error) {
	ext, err := protoutils.StructToStructPb(req.Extra)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.LoadPolicy(ctx, &pb.LoadPolicyRequest{
		Name:         c.name,
		PolicyRepoId: req.PolicyRepoID,
		Extra:        ext,
	})
	if err != nil {
		return nil, err
	}
	return &LoadPolicyResponse{
		PolicyID:   resp.PolicyId,
		PolicyType: resp.PolicyType,
	}, nil
}

func (c *client) RunPolicyEpisode(ctx context.Context, req RunPolicyEpisodeRequest) (*RunPolicyEpisodeResponse, error) {
	ext, err := protoutils.StructToStructPb(req.Extra)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.RunPolicyEpisode(ctx, &pb.RunPolicyEpisodeRequest{
		Name:            c.name,
		PolicyId:        req.PolicyID,
		MaxSteps:        req.MaxSteps,
		Fps:             req.Fps,
		RecordToDataset: req.RecordToDataset,
		DatasetName:     req.DatasetName,
		EpisodeIndex:    req.EpisodeIndex,
		Extra:           ext,
	})
	if err != nil {
		return nil, err
	}
	return &RunPolicyEpisodeResponse{
		NumSteps:    resp.NumSteps,
		DurationS:   resp.DurationS,
		EpisodePath: resp.EpisodePath,
	}, nil
}

func (c *client) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	command, err := protoutils.StructToStructPb(cmd)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.DoCommand(ctx, &commonpb.DoCommandRequest{
		Name:    c.name,
		Command: command,
	})
	if err != nil {
		return nil, err
	}
	return resp.Result.AsMap(), nil
}
