package profile

import (
	"fmt"
	"os"
	"runtime/pprof"
)

type Profile struct {
	CPUFile *os.File
	MemFile *os.File
}

func NewProfile(cpufileName string, memFileName string) (*Profile, error) {
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
		CPUFile: cpuFile,
		MemFile: memFile,
	}, nil
}

func (p *Profile) StartCPUProfile() error {
	err := pprof.StartCPUProfile(p.CPUFile)
	if err != nil {
		return fmt.Errorf("error starting CPU profile: %w", err)
	}

	return nil
}

func (p *Profile) StopCPUProfile() {
	pprof.StopCPUProfile()
}

func (p *Profile) WriteHeapProfile() error {
	err := pprof.WriteHeapProfile(p.MemFile)
	if err != nil {
		return fmt.Errorf("error writing head memory profile: %w", err)
	}

	return nil
}

func (p *Profile) Close() error {
	cerr := p.CPUFile.Close()
	merr := p.MemFile.Close()

	if cerr != nil {
		return fmt.Errorf("error closing CPU profile: %w", cerr)
	}

	if merr != nil {
		return fmt.Errorf("error closing memory profile: %w", merr)
	}

	return nil
}
