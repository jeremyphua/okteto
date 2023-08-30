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

package context

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/okteto/okteto/cmd/utils"
	"github.com/okteto/okteto/pkg/analytics"
	oktetoErrors "github.com/okteto/okteto/pkg/errors"
	oktetoLog "github.com/okteto/okteto/pkg/log"
	"github.com/okteto/okteto/pkg/okteto"
	"github.com/spf13/cobra"
)

// DeleteCMD removes a cluster from okteto context
func DeleteCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Args:  utils.MinimumNArgsAccepted(1, "https://okteto.com/docs/reference/cli/#delete"),
		Short: "Delete one or more contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			for idx, arg := range args {
				args[idx] = okteto.AddSchema(arg)
				args[idx] = strings.TrimSuffix(arg, "/")
			}
			errs := Delete(args)
			var errLen int
			var totalContextsDeleted int
			if errs != nil {
				if merr, ok := errs.(*multierror.Error); ok {
					errLen = len(merr.Errors)
				}
			}
			totalContextsDeleted = len(args) - errLen
			if len(args) == errLen {
				analytics.TrackContextDelete(totalContextsDeleted, false)
			}
			analytics.TrackContextDelete(totalContextsDeleted, true)
			return errs
		},
	}
	return cmd
}

func Delete(okCtxs []string) error {
	ctxStore := okteto.ContextStore()
	var errs error
	for _, okCtx := range okCtxs {
		if okCtx == ctxStore.CurrentContext {
			ctxStore.CurrentContext = ""
		}

		if _, ok := ctxStore.Contexts[okCtx]; ok {
			delete(ctxStore.Contexts, okCtx)
			if err := okteto.NewContextConfigWriter().Write(); err != nil {
				return err
			}
			oktetoLog.Success("'%s' deleted successfully", okCtx)
		} else {
			validOptions := make([]string, 0)
			for k, v := range ctxStore.Contexts {
				if v.IsOkteto {
					validOptions = append(validOptions, k)
				}
			}
			errs = multierror.Append(errs, oktetoErrors.UserError{
				E:    fmt.Errorf("'%s' context doesn't exist. Valid options are: [%s]", okCtx, strings.Join(validOptions, ", ")),
				Hint: fmt.Sprintf("To delete a Kubernetes context run 'kubectl config delete-context %s'", okCtx),
			})
		}
	}

	return errs
}
