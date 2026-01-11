// Package inject provides dependency injected structures for mocking interfaces.
package inject

import (
	"context"

	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/lerobot"
)

// LeRobotService is an injectable lerobot service.
type LeRobotService struct {
	resource.Resource
	name                    resource.Name
	StartRecordingFunc      func(ctx context.Context, datasetName string, extra map[string]interface{}) (string, error)
	StopRecordingFunc       func(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error)
	RecordEpisodeFunc       func(ctx context.Context, req lerobot.RecordEpisodeRequest) (*lerobot.RecordEpisodeResponse, error)
	ReplayEpisodeFunc       func(ctx context.Context, req lerobot.ReplayEpisodeRequest) (*lerobot.ReplayEpisodeResponse, error)
	StartTeleoperationFunc  func(ctx context.Context, req lerobot.StartTeleoperationRequest) (string, error)
	StopTeleoperationFunc   func(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error)
	LoadPolicyFunc          func(ctx context.Context, req lerobot.LoadPolicyRequest) (*lerobot.LoadPolicyResponse, error)
	RunPolicyEpisodeFunc    func(ctx context.Context, req lerobot.RunPolicyEpisodeRequest) (*lerobot.RunPolicyEpisodeResponse, error)
	DoFunc                  func(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error)
	CloseFunc               func(ctx context.Context) error
}

// NewLeRobotService returns a new injected lerobot service.
func NewLeRobotService(name string) *LeRobotService {
	return &LeRobotService{name: lerobot.Named(name)}
}

// Name returns the name of the resource.
func (s *LeRobotService) Name() resource.Name {
	return s.name
}

// StartRecording calls the injected function or returns an error.
func (s *LeRobotService) StartRecording(ctx context.Context, datasetName string, extra map[string]interface{}) (string, error) {
	if s.StartRecordingFunc == nil {
		return "", nil
	}
	return s.StartRecordingFunc(ctx, datasetName, extra)
}

// StopRecording calls the injected function or returns an error.
func (s *LeRobotService) StopRecording(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error) {
	if s.StopRecordingFunc == nil {
		return 0, nil
	}
	return s.StopRecordingFunc(ctx, sessionID, extra)
}

// RecordEpisode calls the injected function or returns an error.
func (s *LeRobotService) RecordEpisode(ctx context.Context, req lerobot.RecordEpisodeRequest) (*lerobot.RecordEpisodeResponse, error) {
	if s.RecordEpisodeFunc == nil {
		return &lerobot.RecordEpisodeResponse{}, nil
	}
	return s.RecordEpisodeFunc(ctx, req)
}

// ReplayEpisode calls the injected function or returns an error.
func (s *LeRobotService) ReplayEpisode(ctx context.Context, req lerobot.ReplayEpisodeRequest) (*lerobot.ReplayEpisodeResponse, error) {
	if s.ReplayEpisodeFunc == nil {
		return &lerobot.ReplayEpisodeResponse{}, nil
	}
	return s.ReplayEpisodeFunc(ctx, req)
}

// StartTeleoperation calls the injected function or returns an error.
func (s *LeRobotService) StartTeleoperation(ctx context.Context, req lerobot.StartTeleoperationRequest) (string, error) {
	if s.StartTeleoperationFunc == nil {
		return "", nil
	}
	return s.StartTeleoperationFunc(ctx, req)
}

// StopTeleoperation calls the injected function or returns an error.
func (s *LeRobotService) StopTeleoperation(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error) {
	if s.StopTeleoperationFunc == nil {
		return 0, nil
	}
	return s.StopTeleoperationFunc(ctx, sessionID, extra)
}

// LoadPolicy calls the injected function or returns an error.
func (s *LeRobotService) LoadPolicy(ctx context.Context, req lerobot.LoadPolicyRequest) (*lerobot.LoadPolicyResponse, error) {
	if s.LoadPolicyFunc == nil {
		return &lerobot.LoadPolicyResponse{}, nil
	}
	return s.LoadPolicyFunc(ctx, req)
}

// RunPolicyEpisode calls the injected function or returns an error.
func (s *LeRobotService) RunPolicyEpisode(ctx context.Context, req lerobot.RunPolicyEpisodeRequest) (*lerobot.RunPolicyEpisodeResponse, error) {
	if s.RunPolicyEpisodeFunc == nil {
		return &lerobot.RunPolicyEpisodeResponse{}, nil
	}
	return s.RunPolicyEpisodeFunc(ctx, req)
}

// DoCommand calls the injected DoCommand or the real version.
func (s *LeRobotService) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	if s.DoFunc == nil {
		if s.Resource == nil {
			return nil, nil
		}
		return s.Resource.DoCommand(ctx, cmd)
	}
	return s.DoFunc(ctx, cmd)
}

// Close calls the injected Close or the real version.
func (s *LeRobotService) Close(ctx context.Context) error {
	if s.CloseFunc == nil {
		if s.Resource == nil {
			return nil
		}
		return s.Resource.Close(ctx)
	}
	return s.CloseFunc(ctx)
}
