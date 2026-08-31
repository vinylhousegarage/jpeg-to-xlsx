package config

import "testing"

func TestLoadAWSConfig(
	t *testing.T,
) {
	tests := []struct {
		name               string
		region             string
		lambdaFunctionName string
		wantRegion         string
		wantIsLambda       bool
	}{
		{
			name:               "default values",
			region:             "",
			lambdaFunctionName: "",
			wantRegion:         defaultAWSRegion,
			wantIsLambda:       false,
		},
		{
			name:               "custom values",
			region:             "us-east-1",
			lambdaFunctionName: "my-lambda-function",
			wantRegion:         "us-east-1",
			wantIsLambda:       true,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Setenv(
					"AWS_REGION",
					test.region,
				)
				t.Setenv(
					"AWS_LAMBDA_FUNCTION_NAME",
					test.lambdaFunctionName,
				)

				config := loadAWSConfig()

				if config.Region !=
					test.wantRegion {
					t.Errorf(
						"Region = %q, want %q",
						config.Region,
						test.wantRegion,
					)
				}

				if config.IsLambda !=
					test.wantIsLambda {
					t.Errorf(
						"IsLambda = %v, want %v",
						config.IsLambda,
						test.wantIsLambda,
					)
				}
			},
		)
	}
}
