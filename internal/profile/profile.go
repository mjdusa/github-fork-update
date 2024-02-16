package profile

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"
)

type Profile struct {
	ctx     context.Context
	cpuFile *os.File
	memFile *os.File
}

func NewProfile(ctx context.Context, cpufileName string, memFileName string) (*Profile, error) {
	cpuFile, cerr := os.Create(cpufileName)
	if cerr != nil {
		return nil, fmt.Errorf("error creating CPU profile: %w", cerr)
	}

	memFile, merr := os.Create(memFileName)
	if merr != nil {
		cpuFile.Close()
		return nil, fmt.Errorf("error creating memory profile: %w", merr)
	}

	return &Profile{
		ctx:     ctx,
		cpuFile: cpuFile,
		memFile: memFile,
	}, nil
}

func (p *Profile) StartCPUProfile() error {
	err := pprof.StartCPUProfile(p.cpuFile)
	if err != nil {
		return fmt.Errorf("error starting CPU profile: %w", err)
	}

	return nil
}

func (p *Profile) StopCPUProfile() {
	pprof.StopCPUProfile()
}

func (p *Profile) WriteHeapProfile() error {
	err := pprof.WriteHeapProfile(p.memFile)
	if err != nil {
		return fmt.Errorf("error writing head memory profile: %w", err)
	}

	return nil
}

func (p *Profile) Close() error {
	var cerr error
	var merr error

	if p.cpuFile != nil {
		cerr = p.cpuFile.Close()
	}

	if p.memFile != nil {
		merr = p.memFile.Close()
	}

	if cerr != nil {
		return fmt.Errorf("error closing CPU profile: %w", cerr)
	}

	if merr != nil {
		return fmt.Errorf("error closing memory profile: %w", merr)
	}

	return nil
}
