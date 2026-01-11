// Package lerobot defines a service for robot learning and teleoperation using LeRobot.
package lerobot

import (
	"context"

	pb "go.viam.com/api/service/lerobot/v1"

	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/robot"
)

func init() {
	resource.RegisterAPI(API, resource.APIRegistration[Service]{
		RPCServiceServerConstructor: NewRPCServiceServer,
		RPCServiceHandler:           pb.RegisterLeRobotServiceHandlerFromEndpoint,
		RPCServiceDesc:              &pb.LeRobotService_ServiceDesc,
		RPCClient:                   NewClientFromConn,
	})
}

// SubtypeName is a constant that identifies the service resource API string "lerobot".
const SubtypeName = "lerobot"

// API is a variable that identifies the service resource API.
var API = resource.APINamespaceRDK.WithServiceType(SubtypeName)

// Named is a helper for getting the named LeRobot's typed resource name.
func Named(name string) resource.Name {
	return resource.NewName(API, name)
}

// FromRobot is a helper for getting the named LeRobot from the given Robot.
// Deprecated: Use FromProvider instead.
func FromRobot(r robot.Robot, name string) (Service, error) {
	return robot.ResourceFromRobot[Service](r, Named(name))
}

// FromProvider is a helper for getting the named LeRobot
// from a resource Provider (collection of Dependencies or a Robot).
func FromProvider(provider resource.Provider, name string) (Service, error) {
	return resource.FromProvider[Service](provider, Named(name))
}

// NamesFromRobot is a helper for getting all lerobot names from the given Robot.
func NamesFromRobot(r robot.Robot) []string {
	return robot.NamesByAPI(r, API)
}

// RecordingSource represents the source of recording data.
type RecordingSource int32

const (
	RecordingSourceUnspecified   RecordingSource = 0
	RecordingSourceTeleoperation RecordingSource = 1
	RecordingSourcePolicy        RecordingSource = 2
)

// RecordEpisodeRequest contains parameters for recording an episode.
type RecordEpisodeRequest struct {
	DatasetName  string
	EpisodeIndex int32
	Source       RecordingSource
	WarmupTimeS  int32
	EpisodeTimeS int32
	ResetTimeS   int32
	Fps          int32
	Tags         []string
	Extra        map[string]interface{}
}

// RecordEpisodeResponse contains the result of recording an episode.
type RecordEpisodeResponse struct {
	NumFrames       int32
	ActualDurationS float32
	EpisodePath     string
}

// ReplayEpisodeRequest contains parameters for replaying an episode.
type ReplayEpisodeRequest struct {
	DatasetName  string
	EpisodeIndex int32
	Fps          int32
	Extra        map[string]interface{}
}

// ReplayEpisodeResponse contains the result of replaying an episode.
type ReplayEpisodeResponse struct {
	NumFramesReplayed int32
	DurationS         float32
}

// StartTeleoperationRequest contains parameters for starting teleoperation.
type StartTeleoperationRequest struct {
	TeleopDeviceType string
	Fps              int32
	DisplayCameras   bool
	Extra            map[string]interface{}
}

// LoadPolicyRequest contains parameters for loading a policy.
type LoadPolicyRequest struct {
	PolicyRepoID string
	Extra        map[string]interface{}
}

// LoadPolicyResponse contains the result of loading a policy.
type LoadPolicyResponse struct {
	PolicyID   string
	PolicyType string
}

// RunPolicyEpisodeRequest contains parameters for running a policy episode.
type RunPolicyEpisodeRequest struct {
	PolicyID        string
	MaxSteps        int32
	Fps             int32
	RecordToDataset bool
	DatasetName     string
	EpisodeIndex    int32
	Extra           map[string]interface{}
}

// RunPolicyEpisodeResponse contains the result of running a policy episode.
type RunPolicyEpisodeResponse struct {
	NumSteps    int32
	DurationS   float32
	Success     bool
	EpisodePath string
}

// Service defines the LeRobot service interface for robot learning and teleoperation.
type Service interface {
	resource.Resource

	// StartRecording begins a new recording session for the specified dataset.
	// Returns a session ID that can be used to stop the recording.
	StartRecording(ctx context.Context, datasetName string, extra map[string]interface{}) (string, error)

	// StopRecording ends the current recording session and saves the data.
	// Returns the duration of the recording in seconds.
	StopRecording(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error)

	// RecordEpisode records a single episode with the specified parameters.
	RecordEpisode(ctx context.Context, req RecordEpisodeRequest) (*RecordEpisodeResponse, error)

	// ReplayEpisode plays back a previously recorded episode.
	ReplayEpisode(ctx context.Context, req ReplayEpisodeRequest) (*ReplayEpisodeResponse, error)

	// StartTeleoperation begins a teleoperation session with the specified device.
	// Returns a session ID that can be used to stop the teleoperation.
	StartTeleoperation(ctx context.Context, req StartTeleoperationRequest) (string, error)

	// StopTeleoperation ends the current teleoperation session.
	// Returns the duration of the session in seconds.
	StopTeleoperation(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error)

	// LoadPolicy loads a policy from a HuggingFace repo or local path.
	LoadPolicy(ctx context.Context, req LoadPolicyRequest) (*LoadPolicyResponse, error)

	// RunPolicyEpisode executes a loaded policy for a single episode.
	RunPolicyEpisode(ctx context.Context, req RunPolicyEpisodeRequest) (*RunPolicyEpisodeResponse, error)
}
