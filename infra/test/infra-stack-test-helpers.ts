import * as cdk from 'aws-cdk-lib';
import { Template } from 'aws-cdk-lib/assertions';
import { InfraStack } from '../lib/infra-stack';

export const testAppEnv = 'staging';
export const testApplicationUrl = 'https://example.com';
export const testAWSAccount = '123456789012';
export const testAWSRegion = 'ap-northeast-1';
export const testBedrockModelId = 'test-model-id';
export const testPromptFileName = 'extractor.txt';
export const testSlackClientId = 'test-slack-client-id';
export const testSlackRedirectUri = `${testApplicationUrl}/api/oauth/slack/callback`;
export const testGoogleClientId = 'test-google-client-id';
export const testEnvironment = {
  APP_ENV: testAppEnv,
  APPLICATION_URL: testApplicationUrl,
  BEDROCK_MODEL_ID: testBedrockModelId,
  PROMPT_FILE_NAME: testPromptFileName,
  SLACK_CLIENT_ID: testSlackClientId,
  SLACK_REDIRECT_URI: testSlackRedirectUri,
  GOOGLE_CLIENT_ID: testGoogleClientId,
} as const;

type TestEnvironmentName = keyof typeof testEnvironment;

const environmentNames = Object.keys(
  testEnvironment,
) as TestEnvironmentName[];

const originalEnvironment = environmentNames.reduce(
  (values, name) => {
    values[name] = process.env[name];

    return values;
  },
  {} as Record<
    TestEnvironmentName,
    string | undefined
  >,
);

export function setTestEnvironment(): void {
  for (const [name, value] of Object.entries(testEnvironment)) {
    process.env[name] = value;
  }
}

export function restoreTestEnvironment(): void {
  for (const name of environmentNames) {
    const value = originalEnvironment[name];

    if (value === undefined) {
      delete process.env[name];

      continue;
    }

    process.env[name] = value;
  }
}

export function createTestStack(id = 'TestStack'): cdk.Stack {
  const app = new cdk.App();

  return new cdk.Stack(app, id, {
    env: {
      account: testAWSAccount,
      region: testAWSRegion,
    },
  });
}

export function createInfraStackTemplate(): Template {
  const app = new cdk.App();

  const stack = new InfraStack(
    app,
    'TestInfraStack',
    {
      env: {
        account: testAWSAccount,
        region: testAWSRegion,
      },
    },
  );

  return Template.fromStack(stack);
}
