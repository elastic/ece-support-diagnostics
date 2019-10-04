// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	goLintRepo     = "golang.org/x/lint/golint"
	goLicenserRepo = "github.com/elastic/go-licenser"
	buildDir       = "build"
	metaDir        = "_meta"
)

// Default set to build everything by default.
var Default = Build.All

// Build namespace used to build binaries.
type Build mg.Namespace

// Test namespace contains all the task for testing the projects.
type Test mg.Namespace

// Check namespace contains tasks related check the actual code quality.
type Check mg.Namespace

// Prepare tasks related to bootstrap the environment or get information about the environment.
type Prepare mg.Namespace

// Format automatically format the code.
type Format mg.Namespace

// Env returns information about the environment.
func (Prepare) Env() {
	RunGo("version")
	RunGo("env")
}

// InstallGoLicenser install go-licenser to check license of the files.
func (Prepare) InstallGoLicenser() error {
	return GoGet(goLicenserRepo)
}

// InstallGoLint for the code.
func (Prepare) InstallGoLint() error {
	return GoGet(goLintRepo)
}

// All build all the things for the current projects.
func (Build) All() {
	mg.Deps(Build.Binary)
}

// Binary build the ece-support-diagnostics artifact.
func (Build) BinaryLinux() error {
	mg.Deps(Prepare.Env)

	env := map[string]string{
		"GOOS":   "linux",
		"GOARCH": "amd64",
	}

	return RunGoEnv(env,
		"build",
		"-o", filepath.Join(buildDir, "linux", "ece-support-diagnostics"),
		"-ldflags", flags(),
		// "-i", "cmd/agent/agent.go",
	)

	// return sh.Run("mage", "-goos=linux", "-goarch=amd64",
	// "-compile", CreateDir(filepath.Join("build", "mage-linux-amd64")))
}

// Binary build the ece-support-diagnostics artifact.
func (Build) Binary() error {
	mg.Deps(Prepare.Env)
	return RunGo(
		"build",
		"-o", filepath.Join(buildDir, "ece-support-diagnostics"),
		"-ldflags", flags(),
		// "-i", "x-pack/cmd/agent/agent.go",
	)
}

// Clean up dev environment.
func (Build) Clean() {
	os.RemoveAll(buildDir)
}

// All run all the code checks.
func (Check) All() {
	mg.SerialDeps(Check.License, Check.GoLint)
}

// GoLint run the code through the linter.
func (Check) GoLint() error {
	mg.Deps(Prepare.InstallGoLint)
	packagesString, err := sh.Output("go", "list", "./...")
	if err != nil {
		return err
	}

	packages := strings.Split(packagesString, "\n")
	for _, pkg := range packages {
		if strings.Contains(pkg, "/vendor/") {
			continue
		}

		if e := sh.RunV("golint", "-set_exit_status", pkg); e != nil {
			err = multierror.Append(err, e)
		}
	}

	return err
}

// License makes sure that all the Golang files have the appropriate license header.
func (Check) License() error {
	mg.Deps(Prepare.InstallGoLicenser)
	// exclude copied files until we come up with a better option
	return sh.RunV("go-licenser", "-d", "-license", "ASL2", "-exclude", "vendor")
}

// All runs all the tests.
func (Test) All() {
	mg.SerialDeps(Test.Unit)
}

// Unit runs all the unit tests.
func (Test) Unit() error {
	mg.Deps(Prepare.Env)
	return RunGo("test", "-race", "-v", "-coverprofile", filepath.Join(buildDir, "coverage.out"), "./...")
}

// Coverage takes the coverages report from running all the tests and display the results in the browser.
func (Test) Coverage() error {
	mg.Deps(Prepare.Env)
	return RunGo("tool", "cover", "-html="+filepath.Join(buildDir, "coverage.out"))
}

// All format automatically all the codes.
func (Format) All() {
	mg.SerialDeps(Format.License)
}

// License applies the right license header.
func (Format) License() error {
	mg.Deps(Prepare.InstallGoLicenser)
	return sh.RunV("go-licenser", "-license", "ASL2", "-exclude", "vendor")
}

// // Package packages the Beat for distribution.
// // Use SNAPSHOT=true to build snapshots.
// // Use PLATFORMS to control the target platforms.
// // Use VERSION_QUALIFIER to control the version qualifier.
// func Package() {
// 	start := time.Now()
// 	defer func() { fmt.Println("package ran for", time.Since(start)) }()
// 	//mage.UseElasticBeatOSSPackaging()
// 	// mage.UseElasticBeatPackaging()

// 	// mg.Deps(Update)
// 	// mg.Deps(CrossBuild, CrossBuildGoDaemon)
// 	// mg.SerialDeps(mage.Package, TestPackages)
// }

// // TestPackages tests the generated packages (i.e. file modes, owners, groups).
// func TestPackages() error {
// 	return mage.TestPackages()
// }

// RunGo runs go command and output the feedback to the stdout and the stderr.
func RunGo(args ...string) error {
	return sh.RunV(mg.GoCmd(), args...)
}

// RunGo runs go command and output the feedback to the stdout and the stderr.
func RunGoEnv(env map[string]string, args ...string) error {
	return sh.RunWith(env, mg.GoCmd(), args...)
}

// GoGet fetch a remote dependencies.
func GoGet(link string) error {
	_, err := sh.Exec(map[string]string{"GO111MODULE": "on"}, os.Stdout, os.Stderr, "go", "get", link)
	return err
}

// Mkdir returns a function that create a directory.
func Mkdir(dir string) func() error {
	return func() error {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("failed to create directory: %v, error: %+v", dir, err)
		}
		return nil
	}
}

func commitID() string {
	commitID, err := sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "cannot retrieve hash"
	}
	return commitID
}

func flags() string {
	ts := time.Now().Format(time.RFC3339)
	commitID := commitID()

	return fmt.Sprintf(
		`-s -w -X "github.com/elastic/ece-support-diagnostics/pkg/release.buildTime=%s" -X "github.com/elastic/ece-support-diagnostics/pkg/release.commit=%s"`,
		ts,
		commitID,
	)
}
