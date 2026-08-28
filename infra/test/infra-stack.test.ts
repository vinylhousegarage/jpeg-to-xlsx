import * as cdk from 'aws-cdk-lib';
import { Match, Template } from 'aws-cdk-lib/assertions';
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

describe('InfraStack', () => {
  let template: Template;

  beforeAll(() => {
    Object.assign(process.env, testEnv);

    const app = new cdk.App();

    const stack = new InfraStack(
      app,
      'TestInfraStack',
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

  test('creates Slack client secret', () => {
    template.resourceCountIs(
      'AWS::SecretsManager::Secret',
      1,
    );

    template.hasResourceProperties(
      'AWS::SecretsManager::Secret',
      Match.objectLike({
        Name:
          'jpeg-to-xlsx/staging/slack-client-secret',
        Description:
          'Slack client secret for jpeg-to-xlsx staging',
      }),
    );
  });

  test('passes Slack client secret ARN to API Lambda', () => {
    template.hasResourceProperties(
      'AWS::Lambda::Function',
      Match.objectLike({
        Environment: {
          Variables: Match.objectLike({
            SLACK_CLIENT_SECRET_ARN: {
              Ref: Match.stringLikeRegexp(
                '^SlackClientSecret',
              ),
            },
          }),
        },
      }),
    );
  });

  test('grants API Lambda permission to read secret', () => {
    template.hasResourceProperties(
      'AWS::IAM::Policy',
      Match.objectLike({
        PolicyDocument: {
          Statement: Match.arrayWith([
            Match.objectLike({
              Action: Match.arrayWith([
                'secretsmanager:GetSecretValue',
                'secretsmanager:DescribeSecret',
              ]),
              Effect: 'Allow',
              Resource: {
                Ref: Match.stringLikeRegexp(
                  '^SlackClientSecret',
                ),
              },
            }),
          ]),
        },
      }),
    );
  });

  test('outputs Slack client secret ARN', () => {
    template.hasOutput(
      'SlackClientSecretArn',
      Match.objectLike({
        Description:
          'Secrets Manager ARN for the Slack client secret',
        Value: {
          Ref: Match.stringLikeRegexp(
            '^SlackClientSecret',
          ),
        },
      }),
    );
  });
});
