/*
 * Copyright 2021 Skyscanner Limited.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package git

import (
	"strings"
	"testing"

	"github.com/skyscanner/turbolift/cmd/flags"
	"github.com/skyscanner/turbolift/internal/executor"
	"github.com/stretchr/testify/assert"
)

func TestItReturnsErrorOnFailure(t *testing.T) {
	fakeExecutor := executor.NewAlwaysFailsFakeExecutor()
	execInstance = fakeExecutor

	_, err := runAndCaptureOutput()
	assert.Error(t, err)

	fakeExecutor.AssertCalledWith(t, [][]string{
		{"work/org/repo1", "git", "checkout", "-b", "some_branch"},
	})
}

func TestItReturnsNilErrorOnSuccess(t *testing.T) {
	fakeExecutor := executor.NewAlwaysSucceedsFakeExecutor()
	execInstance = fakeExecutor

	_, err := runAndCaptureOutput()
	assert.NoError(t, err)

	fakeExecutor.AssertCalledWith(t, [][]string{
		{"work/org/repo1", "git", "checkout", "-b", "some_branch"},
	})
}

func TestItReturnsNilErrorOnSuccessWithDryRun(t *testing.T) {
	flags.DryRun = true
	t.Cleanup(func() {
		flags.DryRun = false
	})

	output, err := runAndCaptureOutput()
	assert.NoError(t, err)
	assert.Equal(t, "Dry-run mode: git [checkout -b some_branch]. Working dir: work/org/repo1", output)
}

func runAndCaptureOutput() (string, error) {
	sb := strings.Builder{}
	err := NewRealGit().Checkout(&sb, "work/org/repo1", "some_branch")

	if err != nil {
		return sb.String(), err
	}
	return sb.String(), nil
}
