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
  GOOGLE_CLIENT_ID: 'test-google-client-id',
};

const originalEnv = {
  APP_ENV: process.env.APP_ENV,
  BEDROCK_MODEL_ID: process.env.BEDROCK_MODEL_ID,
  PROMPT_FILE_NAME: process.env.PROMPT_FILE_NAME,
  SLACK_CLIENT_ID: process.env.SLACK_CLIENT_ID,
  SLACK_REDIRECT_URI: process.env.SLACK_REDIRECT_URI,
  GOOGLE_CLIENT_ID: process.env.GOOGLE_CLIENT_ID,
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
    restoreEnv('GOOGLE_CLIENT_ID');
  });

  test('creates Cognito user pool', () => {
    template.resourceCountIs(
      'AWS::Cognito::UserPool',
      1,
    );

    template.hasResourceProperties(
      'AWS::Cognito::UserPool',
      Match.objectLike({
        UserPoolName:
          'jpeg-to-xlsx-staging-users',
        AdminCreateUserConfig: Match.objectLike({
          AllowAdminCreateUserOnly: true,
        }),
        UsernameConfiguration: {
          CaseSensitive: false,
        },
      }),
    );
  });

  test('creates Cognito user pool domain', () => {
    template.resourceCountIs(
      'AWS::Cognito::UserPoolDomain',
      1,
    );

    template.hasResourceProperties(
      'AWS::Cognito::UserPoolDomain',
      Match.objectLike({
        Domain:
          'jpeg-to-xlsx-staging-123456789012',
        UserPoolId: {
          Ref: Match.stringLikeRegexp('^UserPool'),
        },
      }),
    );
  });

  test('creates Google identity provider', () => {
    template.resourceCountIs(
      'AWS::Cognito::UserPoolIdentityProvider',
      1,
    );

    template.hasResourceProperties(
      'AWS::Cognito::UserPoolIdentityProvider',
      Match.objectLike({
        ProviderName: 'Google',
        ProviderType: 'Google',
        ProviderDetails: Match.objectLike({
          client_id: 'test-google-client-id',
          client_secret: Match.anyValue(),
          authorize_scopes:
            'openid email profile',
        }),
        AttributeMapping: Match.objectLike({
          email: 'email',
          given_name: 'given_name',
          family_name: 'family_name',
        }),
        UserPoolId: {
          Ref: Match.stringLikeRegexp('^UserPool'),
        },
      }),
    );
  });

  test('creates Cognito user pool client', () => {
    template.resourceCountIs(
      'AWS::Cognito::UserPoolClient',
      1,
    );

    template.hasResourceProperties(
      'AWS::Cognito::UserPoolClient',
      Match.objectLike({
        ClientName:
          'jpeg-to-xlsx-staging-client',
        GenerateSecret: false,
        PreventUserExistenceErrors: 'ENABLED',
        SupportedIdentityProviders: ['Google'],
        AllowedOAuthFlows: ['code'],
        AllowedOAuthFlowsUserPoolClient: true,
        AllowedOAuthScopes: Match.arrayWith([
          'openid',
          'email',
          'profile',
        ]),
        CallbackURLs: Match.anyValue(),
        LogoutURLs: Match.anyValue(),
        UserPoolId: {
          Ref: Match.stringLikeRegexp('^UserPool'),
        },
      }),
    );
  });

  test('outputs Cognito user pool ID', () => {
    template.hasOutput(
      'CognitoUserPoolId',
      Match.objectLike({
        Description: 'Cognito user pool ID',
        Value: {
          Ref: Match.stringLikeRegexp('^UserPool'),
        },
      }),
    );
  });

  test('outputs Cognito user pool client ID', () => {
    template.hasOutput(
      'CognitoUserPoolClientId',
      Match.objectLike({
        Description:
          'Cognito user pool client ID',
        Value: {
          Ref: Match.stringLikeRegexp(
            '^UserPoolUserPoolClient',
          ),
        },
      }),
    );
  });

  test('outputs Cognito domain', () => {
    template.hasOutput(
      'CognitoDomain',
      Match.objectLike({
        Description:
          'Cognito managed login domain',
        Value: Match.anyValue(),
      }),
    );
  });

  test('outputs Google OAuth redirect URI', () => {
    template.hasOutput(
      'GoogleRedirectUri',
      Match.objectLike({
        Description:
          'Redirect URI for the Google OAuth client',
        Value: Match.anyValue(),
      }),
    );
  });

  test('creates Slack client secret', () => {
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

  test('creates Google OAuth client secret', () => {
    template.hasResourceProperties(
      'AWS::SecretsManager::Secret',
      Match.objectLike({
        Name:
          'jpeg-to-xlsx/staging/google-oauth-client',
        Description:
          'Google OAuth client secret for jpeg-to-xlsx staging',
      }),
    );
  });

  test('outputs Google OAuth client secret ARN', () => {
    template.hasOutput(
      'GoogleOAuthClientSecretArn',
      Match.objectLike({
        Description:
          'Secrets Manager ARN for the Google OAuth client secret',
        Value: {
          Ref: Match.stringLikeRegexp(
            '^GoogleOAuthClientSecret',
          ),
        },
      }),
    );
  });
});
