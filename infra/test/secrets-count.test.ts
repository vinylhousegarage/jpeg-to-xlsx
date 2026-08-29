import * as cdk from 'aws-cdk-lib';
import { Template } from 'aws-cdk-lib/assertions';
import { InfraStack } from '../lib/infra-stack';

const testEnv = {
  APP_ENV: 'staging',
  BEDROCK_MODEL_ID: 'test-model-id',
  PROMPT_FILE_NAME: 'extractor.txt',
  SLACK_CLIENT_ID: 'test-slack-client-id',
  SLACK_REDIRECT_URI:
    'https://example.com/api/oauth/slack/callback',
};

const originalEnv = {
  APP_ENV: process.env.APP_ENV,
  BEDROCK_MODEL_ID: process.env.BEDROCK_MODEL_ID,
  PROMPT_FILE_NAME: process.env.PROMPT_FILE_NAME,
  SLACK_CLIENT_ID: process.env.SLACK_CLIENT_ID,
  SLACK_REDIRECT_URI: process.env.SLACK_REDIRECT_URI,
};

function restoreEnv(
  name: keyof typeof originalEnv,
): void {
  const value = originalEnv[name];

  if (value === undefined) {
    delete process.env[name];

    return;
  }

  process.env[name] = value;
}

describe('Secrets Manager resources', () => {
  let template: Template;

  beforeAll(() => {
    Object.assign(process.env, testEnv);

    const app = new cdk.App();

    const stack = new InfraStack(
      app,
      'TestSecretsCountStack',
      {
        env: {
          account: '123456789012',
          region: 'ap-northeast-1',
        },
      },
    );

    template = Template.fromStack(stack);
  });

  afterAll(() => {
    restoreEnv('APP_ENV');
    restoreEnv('BEDROCK_MODEL_ID');
    restoreEnv('PROMPT_FILE_NAME');
    restoreEnv('SLACK_CLIENT_ID');
    restoreEnv('SLACK_REDIRECT_URI');
  });

  test('creates two Secrets Manager secrets', () => {
    template.resourceCountIs(
      'AWS::SecretsManager::Secret',
      2,
    );
  });
});
