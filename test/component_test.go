package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cloudposse/test-helpers/pkg/atmos"
	helper "github.com/cloudposse/test-helpers/pkg/atmos/component-helper"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

type ComponentSuite struct {
	helper.TestSuite
}

func (s *ComponentSuite) TestBasic() {
	const component = "ssm-parameter/basic"
	const stack = "default-test"
	const awsRegion = "us-east-2"

	randomID := strings.ToLower(random.UniqueId())
	path := fmt.Sprintf("/%s/testing", randomID)

	inputs := map[string]any{
		"params": map[string]any{
			path: map[string]any{
				"value":       randomID,
				"description": "This is a test.",
				"type":        "String",
			},
		},
	}

	defer s.DestroyAtmosComponent(s.T(), component, stack, &inputs)

	options, _ := s.DeployAtmosComponent(s.T(), component, stack, &inputs)
	assert.NotNil(s.T(), options)

	createdParam := aws.GetParameter(s.T(), awsRegion, path)
	assert.Equal(s.T(), randomID, createdParam)

	paramsOutput := atmos.Output(s.T(), options, "created_params")
	assert.Contains(s.T(), paramsOutput, path)

	s.DriftTest(component, stack, &inputs)
}

func (s *ComponentSuite) TestEnabledFlag() {
	const component = "ssm-parameter/disabled"
	const stack = "default-test"
	s.VerifyEnabledFlag(component, stack, nil)
}

func TestRunSuite(t *testing.T) {
	suite := new(ComponentSuite)
	helper.Run(t, suite)
}
