// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package ints

import (
	"os"
	"testing"

	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/testing/integration"
	"github.com/stretchr/testify/assert"
)

var step1Name interface{}

func TestHelmInputs(t *testing.T) {
	kubectx := os.Getenv("KUBERNETES_CONTEXT")

	if kubectx == "" {
		t.Skipf("Skipping test due to missing KUBERNETES_CONTEXT variable")
	}

	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir:          "step1",
		Dependencies: []string{"@pulumi/kubernetes"},
		Quick:        true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)

			//
			// Assert that the computed values made it through the `ChartOpts`.
			//

			serviceFound := false
			for _, r := range stackInfo.Deployment.Resources {
				name, _ := openapi.Pluck(r.Outputs, "metadata", "name")
				kind, _ := r.Outputs["kind"]
				if name != "nginx-config" {
					// Check that ConfigMap's computed creation time was set in the annotations of
					// Chart resources.
					cmcreation, _ := openapi.Pluck(r.Outputs, "metadata", "annotations",
						"cmcreation")
					assert.True(t, len(cmcreation.(string)) > 0,
						"ConfigMap creation time must be set in annotations")
				} else if name == "simple-nginx-nginx-lego" && kind == "Service" {
					// Check that computed .spec.type was propagated through Chart's values,
					// overriding the default.

					ty, _ := openapi.Pluck(r.Outputs, "spec", "type")
					assert.Equal(t, "ClusterIP", ty)
				}
			}

			assert.True(t, serviceFound, "There must be a service whose type was overridden")
		},
	})
}
