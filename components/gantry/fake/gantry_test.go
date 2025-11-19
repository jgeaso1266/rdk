package fake

import (
	"context"
	"testing"

	"go.viam.com/test"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

func TestGeometries(t *testing.T) {
	logger := logging.NewTestLogger(t)
	fakecfg := resource.Config{Name: "fake_gantry"}
	fake := NewGantry(fakecfg.ResourceName(), logger)

	geoms, err := fake.Geometries(context.Background(), nil)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, len(geoms), test.ShouldEqual, 2)
}

func TestKinematics(t *testing.T) {
	logger := logging.NewTestLogger(t)
	fakecfg := resource.Config{Name: "fake_gantry"}
	fake := NewGantry(fakecfg.ResourceName(), logger)

	model, err := fake.Kinematics(context.Background())
	test.That(t, err, test.ShouldBeNil)
	test.That(t, model.Name(), test.ShouldEqual, "test_gantry")

	currInput, err := fake.CurrentInputs(context.Background())
	test.That(t, err, test.ShouldBeNil)
	test.That(t, currInput, test.ShouldResemble, []float64{120})

	pose, err := model.Transform([]float64{250})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, pose.Point().X, test.ShouldAlmostEqual, 250)
	test.That(t, pose.Point().Y, test.ShouldAlmostEqual, 0)
	test.That(t, pose.Point().Z, test.ShouldAlmostEqual, 0)
}
