/*
   Copyright 2013 Am Laher

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package deb

import (
	"os"
	"path/filepath"
)

// BuildDevPackageFunc specifies a function which can build a DevPackage
type BuildBinaryArtifactFunc func(*BinaryArtifact, *BuildParams) error

// *BinaryPackage specifies functionality for building binary '.deb' packages.
// This encapsulates a Package plus information about platform-specific debs and executables
type BinaryPackage struct {
	*Package
	ExeDest         string
//	BinaryArtifacts []*BinaryArtifact //BinaryArtifact-specific builds
}

// Factory for BinaryPackage
func NewBinaryPackage(pkg *Package) *BinaryPackage {
	return &BinaryPackage{Package: pkg, ExeDest: "/usr/bin"}
}

func (pkg *BinaryPackage) GetArtifacts() (map[Architecture]*BinaryArtifact, error) {
	arches, err := pkg.GetArches()
	if err != nil {
		return nil, err
	}
	ret := map[Architecture]*BinaryArtifact{}
	for _, arch := range arches {
		archArtifact := pkg.GetBinaryArtifact(arch)
		ret[arch] = archArtifact
	}
	return ret, nil
}

// Builds debs for all arches.
func (pkg *BinaryPackage) Build(build *BuildParams, buildFunc BuildBinaryArtifactFunc) error {
	if buildFunc == nil {
		return ErrNoBuildFunc
	}
	archArtifacts, err := pkg.GetArtifacts()
	if err != nil {
		return err
	}
	err = build.Init()
	if err != nil {
		return err
	}

	for _, archArtifact := range archArtifacts {
		err = buildFunc(archArtifact, build)
		//even with an error, remove temp
		if build.IsRmtemp {
			os.RemoveAll(build.TmpDir)

		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Builds debs for a single architecture. Assumes default behaviours of a typical Go package.
// This allows for a limited amount of flexibility (e.g. the use of templates for metadata files).
// To get full flexibility, please use the more granular methods to return archives and add manual work
func (pkg *BinaryPackage) GetBinaryArtifact(arch Architecture) *BinaryArtifact {
	//platform := pkg.getBinaryArtifact(arch)
	/*	if platform == nil {
		platform = pkg.InitBinaryArtifact(arch)
	}*/
	platform := NewBinaryArtifact(pkg, arch)
	return platform
}

// Initialise and return the 'control.tar.gz' archive
func (pkg *BinaryPackage) InitControlArchive(build *BuildParams) (*TarGzWriter, error) {
	archiveFilename := filepath.Join(build.TmpDir, "control.tar.gz")
	tgzw, err := NewTarGzWriter(archiveFilename)
	if err != nil {
		return nil, err
	}
	return tgzw, err
}

// Initialise and return the 'data.tar.gz' archive
func (pkg *BinaryPackage) InitDataArchive(build *BuildParams) (*TarGzWriter, error) {
	archiveFilename := filepath.Join(build.TmpDir, "data.tar.gz")
	tgzw, err := NewTarGzWriter(archiveFilename)
	if err != nil {
		return nil, err
	}
	return tgzw, err
}

// Add executables from file system.
// Be careful to make sure these are the relevant executables for the correct architecture
func (pkg *BinaryPackage) AddExecutablesByFilepath(executablePaths []string, tgzw *TarGzWriter) error {
	return pkg.AddFilesByFilepath(pkg.ExeDest, executablePaths, tgzw)
}

func (pkg *BinaryPackage) AddFilesByFilepath(destinationDir string, executablePaths []string, tgzw *TarGzWriter) error {
	if executablePaths != nil {
		for _, executable := range executablePaths {

			exeName := filepath.Join(destinationDir, filepath.Base(executable))
			err := tgzw.AddFile(executable, exeName)
			if err != nil {
				tgzw.Close()
				return err
			}
		}
	}
	return nil
}
/*
//Initialise the generation parameters for a given architecture
func (pkg *BinaryPackage) InitBinaryArtifact(arch Architecture) *BinaryArtifact {
	if pkg.BinaryArtifacts == nil {
		pkg.BinaryArtifacts = []*BinaryArtifact{}
	}
	platform := NewBinaryArtifact(pkg, arch, targetFile)
	pkg.BinaryArtifacts = append(pkg.BinaryArtifacts, platform)
	return platform
}

// get the generation parameters for a given architecture
func (pkg *BinaryPackage) getBinaryArtifact(arch Architecture) *BinaryArtifact {
	if pkg.BinaryArtifacts == nil {
		pkg.BinaryArtifacts = []*BinaryArtifact{}
	}
	for _, platform := range pkg.BinaryArtifacts {
		if platform.Architecture == arch {
			return platform
		}
	}
	return nil
}
*/
