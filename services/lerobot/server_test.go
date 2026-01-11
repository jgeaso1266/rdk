package lerobot_test

import (
	"context"
	"errors"
	"testing"

	pb "go.viam.com/api/service/lerobot/v1"
	"go.viam.com/test"

	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/lerobot"
	"go.viam.com/rdk/testutils/inject"
)

var (
	testLeRobotName = "lerobot1"
	failLeRobotName = "lerobot2"
	errTestFailed   = errors.New("test failed")
)

func newServer() (pb.LeRobotServiceServer, *inject.LeRobotService, *inject.LeRobotService, error) {
	injectLeRobot := &inject.LeRobotService{}
	injectLeRobot2 := &inject.LeRobotService{}
	resourceMap := map[resource.Name]lerobot.Service{
		lerobot.Named(testLeRobotName): injectLeRobot,
		lerobot.Named(failLeRobotName): injectLeRobot2,
	}
	injectSvc, err := resource.NewAPIResourceCollection(lerobot.API, resourceMap)
	if err != nil {
		return nil, nil, nil, err
	}
	return lerobot.NewRPCServiceServer(injectSvc).(pb.LeRobotServiceServer), injectLeRobot, injectLeRobot2, nil
}

func TestStartRecording(t *testing.T) {
	server, workingService, failingService, err := newServer()
	test.That(t, err, test.ShouldBeNil)

	workingService.StartRecordingFunc = func(ctx context.Context, datasetName string, extra map[string]interface{}) (string, error) {
		return "session-123", nil
	}
	failingService.StartRecordingFunc = func(ctx context.Context, datasetName string, extra map[string]interface{}) (string, error) {
		return "", errTestFailed
	}

	resp, err := server.StartRecording(context.Background(), &pb.StartRecordingRequest{
		Name:        testLeRobotName,
		DatasetName: "test-dataset",
	})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, resp.SessionId, test.ShouldEqual, "session-123")

	_, err = server.StartRecording(context.Background(), &pb.StartRecordingRequest{
		Name:        failLeRobotName,
		DatasetName: "test-dataset",
	})
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())
}

func TestStopRecording(t *testing.T) {
	server, workingService, failingService, err := newServer()
	test.That(t, err, test.ShouldBeNil)

	workingService.StopRecordingFunc = func(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error) {
		return 123.5, nil
	}
	failingService.StopRecordingFunc = func(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error) {
		return 0, errTestFailed
	}

	resp, err := server.StopRecording(context.Background(), &pb.StopRecordingRequest{
		Name:      testLeRobotName,
		SessionId: "session-123",
	})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, resp.DurationS, test.ShouldEqual, float32(123.5))

	_, err = server.StopRecording(context.Background(), &pb.StopRecordingRequest{
		Name:      failLeRobotName,
		SessionId: "session-123",
	})
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())
}

func TestRecordEpisode(t *testing.T) {
	server, workingService, failingService, err := newServer()
	test.That(t, err, test.ShouldBeNil)

	workingService.RecordEpisodeFunc = func(ctx context.Context, req lerobot.RecordEpisodeRequest) (*lerobot.RecordEpisodeResponse, error) {
		return &lerobot.RecordEpisodeResponse{
			NumFrames:       100,
			ActualDurationS: 10.5,
			EpisodePath:     "/path/to/episode",
		}, nil
	}
	failingService.RecordEpisodeFunc = func(ctx context.Context, req lerobot.RecordEpisodeRequest) (*lerobot.RecordEpisodeResponse, error) {
		return nil, errTestFailed
	}

	resp, err := server.RecordEpisode(context.Background(), &pb.RecordEpisodeRequest{
		Name:         testLeRobotName,
		DatasetName:  "test-dataset",
		EpisodeIndex: 1,
		Fps:          30,
	})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, resp.NumFrames, test.ShouldEqual, int32(100))
	test.That(t, resp.ActualDurationS, test.ShouldEqual, float32(10.5))
	test.That(t, resp.EpisodePath, test.ShouldEqual, "/path/to/episode")

	_, err = server.RecordEpisode(context.Background(), &pb.RecordEpisodeRequest{
		Name:         failLeRobotName,
		DatasetName:  "test-dataset",
		EpisodeIndex: 1,
	})
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())
}

func TestReplayEpisode(t *testing.T) {
	server, workingService, failingService, err := newServer()
	test.That(t, err, test.ShouldBeNil)

	workingService.ReplayEpisodeFunc = func(ctx context.Context, req lerobot.ReplayEpisodeRequest) (*lerobot.ReplayEpisodeResponse, error) {
		return &lerobot.ReplayEpisodeResponse{
			NumFramesReplayed: 50,
			DurationS:         5.0,
		}, nil
	}
	failingService.ReplayEpisodeFunc = func(ctx context.Context, req lerobot.ReplayEpisodeRequest) (*lerobot.ReplayEpisodeResponse, error) {
		return nil, errTestFailed
	}

	resp, err := server.ReplayEpisode(context.Background(), &pb.ReplayEpisodeRequest{
		Name:         testLeRobotName,
		DatasetName:  "test-dataset",
		EpisodeIndex: 1,
		Fps:          30,
	})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, resp.NumFramesReplayed, test.ShouldEqual, int32(50))
	test.That(t, resp.DurationS, test.ShouldEqual, float32(5.0))

	_, err = server.ReplayEpisode(context.Background(), &pb.ReplayEpisodeRequest{
		Name:         failLeRobotName,
		DatasetName:  "test-dataset",
		EpisodeIndex: 1,
	})
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())
}

func TestStartTeleoperation(t *testing.T) {
	server, workingService, failingService, err := newServer()
	test.That(t, err, test.ShouldBeNil)

	workingService.StartTeleoperationFunc = func(ctx context.Context, req lerobot.StartTeleoperationRequest) (string, error) {
		return "teleop-session-456", nil
	}
	failingService.StartTeleoperationFunc = func(ctx context.Context, req lerobot.StartTeleoperationRequest) (string, error) {
		return "", errTestFailed
	}

	resp, err := server.StartTeleoperation(context.Background(), &pb.StartTeleoperationRequest{
		Name:             testLeRobotName,
		TeleopDeviceType: "keyboard",
		Fps:              30,
		DisplayCameras:   true,
	})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, resp.SessionId, test.ShouldEqual, "teleop-session-456")

	_, err = server.StartTeleoperation(context.Background(), &pb.StartTeleoperationRequest{
		Name:             failLeRobotName,
		TeleopDeviceType: "keyboard",
	})
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())
}

func TestStopTeleoperation(t *testing.T) {
	server, workingService, failingService, err := newServer()
	test.That(t, err, test.ShouldBeNil)

	workingService.StopTeleoperationFunc = func(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error) {
		return 60.0, nil
	}
	failingService.StopTeleoperationFunc = func(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error) {
		return 0, errTestFailed
	}

	resp, err := server.StopTeleoperation(context.Background(), &pb.StopTeleoperationRequest{
		Name:      testLeRobotName,
		SessionId: "teleop-session-456",
	})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, resp.DurationS, test.ShouldEqual, float32(60.0))

	_, err = server.StopTeleoperation(context.Background(), &pb.StopTeleoperationRequest{
		Name:      failLeRobotName,
		SessionId: "teleop-session-456",
	})
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())
}

func TestLoadPolicy(t *testing.T) {
	server, workingService, failingService, err := newServer()
	test.That(t, err, test.ShouldBeNil)

	workingService.LoadPolicyFunc = func(ctx context.Context, req lerobot.LoadPolicyRequest) (*lerobot.LoadPolicyResponse, error) {
		return &lerobot.LoadPolicyResponse{
			PolicyID:   "policy-789",
			PolicyType: "act",
		}, nil
	}
	failingService.LoadPolicyFunc = func(ctx context.Context, req lerobot.LoadPolicyRequest) (*lerobot.LoadPolicyResponse, error) {
		return nil, errTestFailed
	}

	resp, err := server.LoadPolicy(context.Background(), &pb.LoadPolicyRequest{
		Name:         testLeRobotName,
		PolicyRepoId: "huggingface/lerobot-act",
	})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, resp.PolicyId, test.ShouldEqual, "policy-789")
	test.That(t, resp.PolicyType, test.ShouldEqual, "act")

	_, err = server.LoadPolicy(context.Background(), &pb.LoadPolicyRequest{
		Name:         failLeRobotName,
		PolicyRepoId: "huggingface/lerobot-act",
	})
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())
}

func TestRunPolicyEpisode(t *testing.T) {
	server, workingService, failingService, err := newServer()
	test.That(t, err, test.ShouldBeNil)

	workingService.RunPolicyEpisodeFunc = func(ctx context.Context, req lerobot.RunPolicyEpisodeRequest) (*lerobot.RunPolicyEpisodeResponse, error) {
		return &lerobot.RunPolicyEpisodeResponse{
			NumSteps:    200,
			DurationS:   20.0,
			Success:     true,
			EpisodePath: "/path/to/policy/episode",
		}, nil
	}
	failingService.RunPolicyEpisodeFunc = func(ctx context.Context, req lerobot.RunPolicyEpisodeRequest) (*lerobot.RunPolicyEpisodeResponse, error) {
		return nil, errTestFailed
	}

	resp, err := server.RunPolicyEpisode(context.Background(), &pb.RunPolicyEpisodeRequest{
		Name:            testLeRobotName,
		PolicyId:        "policy-789",
		MaxSteps:        300,
		Fps:             30,
		RecordToDataset: true,
		DatasetName:     "eval-dataset",
		EpisodeIndex:    0,
	})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, resp.NumSteps, test.ShouldEqual, int32(200))
	test.That(t, resp.DurationS, test.ShouldEqual, float32(20.0))
	test.That(t, resp.Success, test.ShouldBeTrue)
	test.That(t, resp.EpisodePath, test.ShouldEqual, "/path/to/policy/episode")

	_, err = server.RunPolicyEpisode(context.Background(), &pb.RunPolicyEpisodeRequest{
		Name:     failLeRobotName,
		PolicyId: "policy-789",
	})
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())
}
