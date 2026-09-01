import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import {
  Match,
  Template,
} from 'aws-cdk-lib/assertions';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager';

import {
  createComputeResources,
} from '../../lib/constructs/compute-resources';
import {
  createTestStack,
  testAppEnv,
  testApplicationUrl,
  testBedrockModelId,
  testPromptFileName,
  testSlackClientId,
  testSlackRedirectUri,
} from '../infra-stack-test-helpers';

const testCognitoClientId =
  'test-cognito-client-id';

const testCognitoClientSecretArn =
  'arn:aws:secretsmanager:' +
  'ap-northeast-1:123456789012:' +
  'secret:test-cognito-client-secret';

const testCognitoIssuer =
  'https://cognito-idp.ap-northeast-1.' +
  'amazonaws.com/ap-northeast-1_test';

const testCognitoAuthorizationEndpoint =
  'https://test.auth.ap-northeast-1.' +
  'amazoncognito.com/oauth2/authorize';

const testCognitoTokenEndpoint =
  'https://test.auth.ap-northeast-1.' +
  'amazoncognito.com/oauth2/token';

const testCognitoRedirectUri =
  `${testApplicationUrl}/api/auth/callback`;

const testOAuthStateTableName =
  'test-cognito-oauth-states';

const testSessionTableName =
  'test-auth-sessions';

describe('createComputeResources', () => {
  let template: Template;

  beforeAll(() => {
    const stack = createTestStack(
      'ComputeResourcesTestStack',
    );

    const inputBucket = new s3.Bucket(
      stack,
      'InputBucket',
    );

    const outputBucket = new s3.Bucket(
      stack,
      'OutputBucket',
    );

    const slackSecret =
      new secretsmanager.Secret(
        stack,
        'SlackClientSecret',
      );

    const slackTokenTable =
      new dynamodb.Table(
        stack,
        'SlackTokenTable',
        {
          partitionKey: {
            name:
              'id',
            type:
              dynamodb.AttributeType
                .STRING,
          },
          billingMode:
            dynamodb.BillingMode
              .PAY_PER_REQUEST,
        },
      );

    createComputeResources(
      stack,
      {
        appEnv:
          testAppEnv,

        bedrockModelId:
          testBedrockModelId,

        promptFileName:
          testPromptFileName,

        slackClientId:
          testSlackClientId,

        slackRedirectUri:
          testSlackRedirectUri,

        inputBucket,

        outputBucket,

        slackSecret,

        slackTokenTable,

        cognitoClientId:
          testCognitoClientId,

        cognitoClientSecretArn:
          testCognitoClientSecretArn,

        cognitoIssuer:
          testCognitoIssuer,

        cognitoAuthorizationEndpoint:
          testCognitoAuthorizationEndpoint,

        cognitoTokenEndpoint:
          testCognitoTokenEndpoint,

        cognitoRedirectUri:
          testCognitoRedirectUri,

        postLoginRedirectUrl:
          testApplicationUrl,

        oauthStateTableName:
          testOAuthStateTableName,

        sessionTableName:
          testSessionTableName,
      },
    );

    template =
      Template.fromStack(stack);
  });

  test('creates API and processor Lambda functions', () => {
    template.resourceCountIs(
      'AWS::Lambda::Function',
      2,
    );
  });

  test('creates API Lambda function', () => {
    template.hasResourceProperties(
      'AWS::Lambda::Function',
      Match.objectLike({
        Runtime:
          'provided.al2023',
        Handler:
          'bootstrap',
        Architectures: [
          'arm64',
        ],
        Timeout:
          15,
        Environment: {
          Variables:
            Match.objectLike({
              APP_ENV:
                testAppEnv,

              INPUT_BUCKET_NAME:
                Match.anyValue(),

              SLACK_CLIENT_ID:
                testSlackClientId,

              SLACK_CLIENT_SECRET_ARN:
                Match.anyValue(),

              SLACK_REDIRECT_URI:
                testSlackRedirectUri,

              SLACK_TOKEN_TABLE_NAME:
                Match.anyValue(),

              COGNITO_CLIENT_ID:
                testCognitoClientId,

              COGNITO_CLIENT_SECRET_ARN:
                testCognitoClientSecretArn,

              COGNITO_ISSUER:
                testCognitoIssuer,

              COGNITO_AUTHORIZATION_ENDPOINT:
                testCognitoAuthorizationEndpoint,

              COGNITO_TOKEN_ENDPOINT:
                testCognitoTokenEndpoint,

              COGNITO_REDIRECT_URI:
                testCognitoRedirectUri,

              AUTH_REDIRECT_URL:
                testApplicationUrl,

              COGNITO_OAUTH_STATE_TABLE_NAME:
                testOAuthStateTableName,

              AUTH_SESSION_TABLE_NAME:
                testSessionTableName,
            }),
        },
      }),
    );
  });

  test('passes input bucket name to API Lambda', () => {
    template.hasResourceProperties(
      'AWS::Lambda::Function',
      Match.objectLike({
        Environment: {
          Variables:
            Match.objectLike({
              INPUT_BUCKET_NAME: {
                Ref:
                  Match.stringLikeRegexp(
                    '^InputBucket',
                  ),
              },
            }),
        },
      }),
    );
  });

  test('passes Slack secret ARN to API Lambda', () => {
    template.hasResourceProperties(
      'AWS::Lambda::Function',
      Match.objectLike({
        Environment: {
          Variables:
            Match.objectLike({
              SLACK_CLIENT_SECRET_ARN: {
                Ref:
                  Match.stringLikeRegexp(
                    '^SlackClientSecret',
                  ),
              },
            }),
        },
      }),
    );
  });

  test('passes Slack token table name to API Lambda', () => {
    template.hasResourceProperties(
      'AWS::Lambda::Function',
      Match.objectLike({
        Environment: {
          Variables:
            Match.objectLike({
              SLACK_TOKEN_TABLE_NAME: {
                Ref:
                  Match.stringLikeRegexp(
                    '^SlackTokenTable',
                  ),
              },
              COGNITO_CLIENT_ID:
                testCognitoClientId,
            }),
        },
      }),
    );
  });

  test('creates processor Lambda function', () => {
    template.hasResourceProperties(
      'AWS::Lambda::Function',
      Match.objectLike({
        Runtime:
          'provided.al2023',
        Handler:
          'bootstrap',
        Architectures: [
          'arm64',
        ],
        Timeout:
          30,
        Environment: {
          Variables:
            Match.objectLike({
              APP_ENV:
                testAppEnv,

              INPUT_BUCKET_NAME:
                Match.anyValue(),

              OUTPUT_BUCKET_NAME:
                Match.anyValue(),

              BEDROCK_MODEL_ID:
                testBedrockModelId,

              PROMPT_FILE_NAME:
                testPromptFileName,

              SLACK_TOKEN_TABLE_NAME:
                Match.anyValue(),

              COGNITO_CLIENT_ID:
                Match.absent(),

              COGNITO_CLIENT_SECRET_ARN:
                Match.absent(),

              AUTH_SESSION_TABLE_NAME:
                Match.absent(),
            }),
        },
      }),
    );
  });

  test('passes input bucket name to processor Lambda', () => {
    template.hasResourceProperties(
      'AWS::Lambda::Function',
      Match.objectLike({
        Timeout:
          30,
        Environment: {
          Variables:
            Match.objectLike({
              INPUT_BUCKET_NAME: {
                Ref:
                  Match.stringLikeRegexp(
                    '^InputBucket',
                  ),
              },
            }),
        },
      }),
    );
  });

  test('passes output bucket name to processor Lambda', () => {
    template.hasResourceProperties(
      'AWS::Lambda::Function',
      Match.objectLike({
        Environment: {
          Variables:
            Match.objectLike({
              OUTPUT_BUCKET_NAME: {
                Ref:
                  Match.stringLikeRegexp(
                    '^OutputBucket',
                  ),
              },
            }),
        },
      }),
    );
  });

  test('passes Slack token table name to processor Lambda', () => {
    template.hasResourceProperties(
      'AWS::Lambda::Function',
      Match.objectLike({
        Timeout:
          30,
        Environment: {
          Variables:
            Match.objectLike({
              SLACK_TOKEN_TABLE_NAME: {
                Ref:
                  Match.stringLikeRegexp(
                    '^SlackTokenTable',
                  ),
              },
            }),
        },
      }),
    );
  });
});
