// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ZupIT/horusec-devkit/pkg/utils/logger"

	"github.com/ZupIT/horusec/internal/helpers/messages"
)

const (
	MinVersionGitAccept    = 2
	MinSubVersionGitAccept = 0o1
)

var (
	// ErrMinVersion is the error logged when the installed Git version is not the minimum supported.
	ErrMinVersion = fmt.Errorf("%v.%v", MinVersionGitAccept, MinSubVersionGitAccept)

	// ErrGitNotInstalled occurs when Git is not installed.
	ErrGitNotInstalled = errors.New("git not found. Please check and try again")

	// ErrGitLowerVersion occur when the installed Git version is not the minimum supported.
	ErrGitLowerVersion = errors.New("git version is lower of 2.01. Please check and try again")
)

type RequirementGit struct{}

func NewRequirementGit() *RequirementGit {
	return &RequirementGit{}
}

func (r *RequirementGit) ValidateGit() error {
	response, err := r.validateIfGitIsInstalled()
	if err != nil {
		return err
	}
	return r.validateIfGitIsSupported(response)
}

func (r *RequirementGit) validateIfGitIsInstalled() (string, error) {
	response, err := r.execGitVersion()
	if err != nil {
		return "", err
	}
	if !r.checkIfContainsGitVersion(response) {
		return "", ErrGitNotInstalled
	}
	return response, nil
}

func (r *RequirementGit) validateIfGitIsSupported(version string) error {
	err := r.validateIfGitIsRunningInMinVersion(version)
	if err != nil {
		logger.LogInfo(messages.MsgInfoHowToInstallGit)
		return err
	}
	return nil
}

func (r *RequirementGit) execGitVersion() (string, error) {
	responseBytes, err := exec.Command("git", "--version").CombinedOutput()
	if err != nil {
		logger.LogErrorWithLevel(messages.MsgErrorWhenCheckRequirementsGit, err)
		return "", err
	}
	return strings.ToLower(string(responseBytes)), nil
}

func (r *RequirementGit) validateIfGitIsRunningInMinVersion(response string) error {
	version, subversion, err := r.extractGitVersionFromString(response)
	if err != nil {
		return err
	}
	if version < MinVersionGitAccept {
		logger.LogErrorWithLevel(messages.MsgErrorWhenGitIsLowerVersion, ErrMinVersion)
		return ErrGitLowerVersion
	} else if version == MinVersionGitAccept && subversion < MinSubVersionGitAccept {
		logger.LogErrorWithLevel(messages.MsgErrorWhenGitIsLowerVersion, ErrMinVersion)
		return ErrGitLowerVersion
	}
	return nil
}

func (r *RequirementGit) extractGitVersionFromString(response string) (int, int, error) {
	responseSpited := strings.Split(strings.ToLower(response), "git version ")
	if len(responseSpited) < 1 || len(responseSpited) > 1 && len(responseSpited[1]) < 3 {
		return 0, 0, ErrGitNotInstalled
	}
	return r.getVersionAndSubVersion(responseSpited[1])
}

func (r *RequirementGit) checkIfContainsGitVersion(response string) bool {
	return strings.Contains(strings.ToLower(response), "git version ")
}

func (r *RequirementGit) getVersionAndSubVersion(fullVersion string) (int, int, error) {
	version, err := strconv.Atoi(fullVersion[0:1])
	if err != nil {
		return 0, 0, ErrGitNotInstalled
	}
	subversion, err := strconv.Atoi(fullVersion[2:4])
	if err != nil {
		return 0, 0, ErrGitNotInstalled
	}
	return version, subversion, nil
}
