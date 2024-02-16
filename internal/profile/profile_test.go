package profile_test

import (
	"context"
	"os"
	"testing"

	"github.com/mjdusa/github-fork-update/internal/profile"
	"github.com/stretchr/testify/assert"
)

func TestProfile(t *testing.T) {
	ctx := context.Background()

	cpuFile, _ := os.CreateTemp("", "cpu")
	memFile, _ := os.CreateTemp("", "mem")

	defer os.Remove(cpuFile.Name())
	defer os.Remove(memFile.Name())

	profile, err := profile.NewProfile(ctx, cpuFile.Name(), memFile.Name())
	if err != nil {
		t.Fatalf("Failed to create new profile: %v", err)
	}

	err = profile.StartCPUProfile()
	if err != nil {
		t.Errorf("Failed to start CPU profile: %v", err)
	}

	profile.StopCPUProfile()

	err = profile.WriteHeapProfile()
	if err != nil {
		t.Errorf("Failed to write heap profile: %v", err)
	}

	err = profile.Close()
	if err != nil {
		t.Errorf("Failed to close profile: %v", err)
	}
}

func TestProfile_bad_cpu(t *testing.T) {
	ctx := context.Background()

	memFile, _ := os.CreateTemp("", "mem")
	defer os.Remove(memFile.Name())

	profile, err := profile.NewProfile(ctx, "", memFile.Name())

	assert.Error(t, err)
	assert.Nil(t, profile)
}

func TestProfile_bad_mem(t *testing.T) {
	ctx := context.Background()

	cpuFile, _ := os.CreateTemp("", "cpu")
	defer os.Remove(cpuFile.Name())

	profile, err := profile.NewProfile(ctx, cpuFile.Name(), "")
	assert.Error(t, err)
	assert.Nil(t, profile)
}
