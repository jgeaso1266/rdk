package lerobot_test

import (
	"context"
	"net"
	"testing"

	"go.viam.com/test"
	"go.viam.com/utils/rpc"

	viamgrpc "go.viam.com/rdk/grpc"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/lerobot"
	"go.viam.com/rdk/testutils/inject"
)

func TestClient(t *testing.T) {
	logger := logging.NewTestLogger(t)
	listener, err := net.Listen("tcp", "localhost:0")
	test.That(t, err, test.ShouldBeNil)
	rpcServer, err := rpc.NewServer(logger, rpc.WithUnauthenticated())
	test.That(t, err, test.ShouldBeNil)

	workingLeRobot := &inject.LeRobotService{}
	failingLeRobot := &inject.LeRobotService{}

	// Setup working service functions
	workingLeRobot.StartRecordingFunc = func(ctx context.Context, datasetName string, extra map[string]interface{}) (string, error) {
		return "session-123", nil
	}
	workingLeRobot.StopRecordingFunc = func(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error) {
		return 123.5, nil
	}
	workingLeRobot.RecordEpisodeFunc = func(ctx context.Context, req lerobot.RecordEpisodeRequest) (*lerobot.RecordEpisodeResponse, error) {
		return &lerobot.RecordEpisodeResponse{
			NumFrames:       100,
			ActualDurationS: 10.5,
			EpisodePath:     "/path/to/episode",
		}, nil
	}
	workingLeRobot.ReplayEpisodeFunc = func(ctx context.Context, req lerobot.ReplayEpisodeRequest) (*lerobot.ReplayEpisodeResponse, error) {
		return &lerobot.ReplayEpisodeResponse{
			NumFramesReplayed: 50,
			DurationS:         5.0,
		}, nil
	}
	workingLeRobot.StartTeleoperationFunc = func(ctx context.Context, req lerobot.StartTeleoperationRequest) (string, error) {
		return "teleop-456", nil
	}
	workingLeRobot.StopTeleoperationFunc = func(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error) {
		return 60.0, nil
	}
	workingLeRobot.LoadPolicyFunc = func(ctx context.Context, req lerobot.LoadPolicyRequest) (*lerobot.LoadPolicyResponse, error) {
		return &lerobot.LoadPolicyResponse{
			PolicyID:   "policy-789",
			PolicyType: "act",
		}, nil
	}
	workingLeRobot.RunPolicyEpisodeFunc = func(ctx context.Context, req lerobot.RunPolicyEpisodeRequest) (*lerobot.RunPolicyEpisodeResponse, error) {
		return &lerobot.RunPolicyEpisodeResponse{
			NumSteps:    200,
			DurationS:   20.0,
			Success:     true,
			EpisodePath: "/path/to/policy/episode",
		}, nil
	}

	// Setup failing service functions
	failingLeRobot.StartRecordingFunc = func(ctx context.Context, datasetName string, extra map[string]interface{}) (string, error) {
		return "", errTestFailed
	}
	failingLeRobot.StopRecordingFunc = func(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error) {
		return 0, errTestFailed
	}
	failingLeRobot.RecordEpisodeFunc = func(ctx context.Context, req lerobot.RecordEpisodeRequest) (*lerobot.RecordEpisodeResponse, error) {
		return nil, errTestFailed
	}
	failingLeRobot.ReplayEpisodeFunc = func(ctx context.Context, req lerobot.ReplayEpisodeRequest) (*lerobot.ReplayEpisodeResponse, error) {
		return nil, errTestFailed
	}
	failingLeRobot.StartTeleoperationFunc = func(ctx context.Context, req lerobot.StartTeleoperationRequest) (string, error) {
		return "", errTestFailed
	}
	failingLeRobot.StopTeleoperationFunc = func(ctx context.Context, sessionID string, extra map[string]interface{}) (float32, error) {
		return 0, errTestFailed
	}
	failingLeRobot.LoadPolicyFunc = func(ctx context.Context, req lerobot.LoadPolicyRequest) (*lerobot.LoadPolicyResponse, error) {
		return nil, errTestFailed
	}
	failingLeRobot.RunPolicyEpisodeFunc = func(ctx context.Context, req lerobot.RunPolicyEpisodeRequest) (*lerobot.RunPolicyEpisodeResponse, error) {
		return nil, errTestFailed
	}

	resourceMap := map[resource.Name]lerobot.Service{
		lerobot.Named(testLeRobotName): workingLeRobot,
		lerobot.Named(failLeRobotName): failingLeRobot,
	}
	lerobotSvc, err := resource.NewAPIResourceCollection(lerobot.API, resourceMap)
	test.That(t, err, test.ShouldBeNil)
	resourceAPI, ok, err := resource.LookupAPIRegistration[lerobot.Service](lerobot.API)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, ok, test.ShouldBeTrue)
	test.That(t, resourceAPI.RegisterRPCService(context.Background(), rpcServer, lerobotSvc), test.ShouldBeNil)

	go rpcServer.Serve(listener)
	defer rpcServer.Stop()

	t.Run("Failing client due to canceled context", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err = viamgrpc.Dial(cancelCtx, listener.Addr().String(), logger)
		test.That(t, err, test.ShouldNotBeNil)
		test.That(t, err, test.ShouldBeError, context.Canceled)
	})

	t.Run("client tests for working lerobot", func(t *testing.T) {
		conn, err := viamgrpc.Dial(context.Background(), listener.Addr().String(), logger)
		test.That(t, err, test.ShouldBeNil)
		workingClient, err := lerobot.NewClientFromConn(context.Background(), conn, "", lerobot.Named(testLeRobotName), logger)
		test.That(t, err, test.ShouldBeNil)

		// Test StartRecording
		sessionID, err := workingClient.StartRecording(context.Background(), "test-dataset", nil)
		test.That(t, err, test.ShouldBeNil)
		test.That(t, sessionID, test.ShouldEqual, "session-123")

		// Test StopRecording
		duration, err := workingClient.StopRecording(context.Background(), "session-123", nil)
		test.That(t, err, test.ShouldBeNil)
		test.That(t, duration, test.ShouldEqual, float32(123.5))

		// Test RecordEpisode
		recordResp, err := workingClient.RecordEpisode(context.Background(), lerobot.RecordEpisodeRequest{
			DatasetName:  "test-dataset",
			EpisodeIndex: 1,
			Fps:          30,
		})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, recordResp.NumFrames, test.ShouldEqual, int32(100))
		test.That(t, recordResp.ActualDurationS, test.ShouldEqual, float32(10.5))
		test.That(t, recordResp.EpisodePath, test.ShouldEqual, "/path/to/episode")

		// Test ReplayEpisode
		replayResp, err := workingClient.ReplayEpisode(context.Background(), lerobot.ReplayEpisodeRequest{
			DatasetName:  "test-dataset",
			EpisodeIndex: 1,
			Fps:          30,
		})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, replayResp.NumFramesReplayed, test.ShouldEqual, int32(50))
		test.That(t, replayResp.DurationS, test.ShouldEqual, float32(5.0))

		// Test StartTeleoperation
		teleopSessionID, err := workingClient.StartTeleoperation(context.Background(), lerobot.StartTeleoperationRequest{
			TeleopDeviceType: "keyboard",
			Fps:              30,
			DisplayCameras:   true,
		})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, teleopSessionID, test.ShouldEqual, "teleop-456")

		// Test StopTeleoperation
		teleopDuration, err := workingClient.StopTeleoperation(context.Background(), "teleop-456", nil)
		test.That(t, err, test.ShouldBeNil)
		test.That(t, teleopDuration, test.ShouldEqual, float32(60.0))

		// Test LoadPolicy
		loadResp, err := workingClient.LoadPolicy(context.Background(), lerobot.LoadPolicyRequest{
			PolicyRepoID: "huggingface/lerobot-act",
		})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, loadResp.PolicyID, test.ShouldEqual, "policy-789")
		test.That(t, loadResp.PolicyType, test.ShouldEqual, "act")

		// Test RunPolicyEpisode
		runResp, err := workingClient.RunPolicyEpisode(context.Background(), lerobot.RunPolicyEpisodeRequest{
			PolicyID:        "policy-789",
			MaxSteps:        300,
			Fps:             30,
			RecordToDataset: true,
			DatasetName:     "eval-dataset",
			EpisodeIndex:    0,
		})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, runResp.NumSteps, test.ShouldEqual, int32(200))
		test.That(t, runResp.DurationS, test.ShouldEqual, float32(20.0))
		test.That(t, runResp.Success, test.ShouldBeTrue)
		test.That(t, runResp.EpisodePath, test.ShouldEqual, "/path/to/policy/episode")

		test.That(t, workingClient.Close(context.Background()), test.ShouldBeNil)
		test.That(t, conn.Close(), test.ShouldBeNil)
	})

	t.Run("client tests for failing lerobot", func(t *testing.T) {
		conn, err := viamgrpc.Dial(context.Background(), listener.Addr().String(), logger)
		test.That(t, err, test.ShouldBeNil)
		failingClient, err := lerobot.NewClientFromConn(context.Background(), conn, "", lerobot.Named(failLeRobotName), logger)
		test.That(t, err, test.ShouldBeNil)

		_, err = failingClient.StartRecording(context.Background(), "test-dataset", nil)
		test.That(t, err, test.ShouldNotBeNil)
		test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())

		_, err = failingClient.StopRecording(context.Background(), "session-123", nil)
		test.That(t, err, test.ShouldNotBeNil)
		test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())

		_, err = failingClient.RecordEpisode(context.Background(), lerobot.RecordEpisodeRequest{})
		test.That(t, err, test.ShouldNotBeNil)
		test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())

		_, err = failingClient.ReplayEpisode(context.Background(), lerobot.ReplayEpisodeRequest{})
		test.That(t, err, test.ShouldNotBeNil)
		test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())

		_, err = failingClient.StartTeleoperation(context.Background(), lerobot.StartTeleoperationRequest{})
		test.That(t, err, test.ShouldNotBeNil)
		test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())

		_, err = failingClient.StopTeleoperation(context.Background(), "teleop-456", nil)
		test.That(t, err, test.ShouldNotBeNil)
		test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())

		_, err = failingClient.LoadPolicy(context.Background(), lerobot.LoadPolicyRequest{})
		test.That(t, err, test.ShouldNotBeNil)
		test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())

		_, err = failingClient.RunPolicyEpisode(context.Background(), lerobot.RunPolicyEpisodeRequest{})
		test.That(t, err, test.ShouldNotBeNil)
		test.That(t, err.Error(), test.ShouldContainSubstring, errTestFailed.Error())

		test.That(t, failingClient.Close(context.Background()), test.ShouldBeNil)
		test.That(t, conn.Close(), test.ShouldBeNil)
	})
}
