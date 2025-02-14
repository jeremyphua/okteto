// Copyright 2023 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"context"
	"testing"

	"github.com/okteto/okteto/pkg/model"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeployNotRemovingEnvFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	_, err := fs.Create(".env")
	require.NoError(t, err)
	opts := &Options{
		Manifest: &model.Manifest{
			Deploy: &model.DeployInfo{},
		},
	}
	localDeployer := localDeployer{
		ConfigMapHandler: &fakeCmapHandler{},
		Fs:               fs,
	}
	err = localDeployer.runDeploySection(context.Background(), opts)
	assert.NoError(t, err)
	_, err = fs.Stat(".env")
	require.NoError(t, err)

}
