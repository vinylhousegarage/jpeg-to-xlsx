import {
  Match,
  Template,
} from 'aws-cdk-lib/assertions';

import {
  createInfraStackTemplate,
  restoreTestEnvironment,
  setTestEnvironment,
  testAWSAccount,
} from './infra-stack-test-helpers';

describe('InfraStack', () => {
  let template: Template;

  beforeAll(() => {
    setTestEnvironment();

    template =
      createInfraStackTemplate();
  });

  afterAll(() => {
    restoreTestEnvironment();
  });

  test('grants API Lambda permission to read Slack secret', () => {
    template.hasResourceProperties(
      'AWS::IAM::Policy',
      Match.objectLike({
        PolicyDocument: {
          Statement:
            Match.arrayWith([
              Match.objectLike({
                Action:
                  Match.arrayWith([
                    'secretsmanager:GetSecretValue',
                    'secretsmanager:DescribeSecret',
                  ]),
                Effect:
                  'Allow',
                Resource: {
                  Ref:
                    Match.stringLikeRegexp(
                      '^SlackClientSecret',
                    ),
                },
              }),
            ]),
        },
      }),
    );
  });

  test('grants API Lambda permission to read Cognito client secret', () => {
    template.hasResourceProperties(
      'AWS::IAM::Policy',
      Match.objectLike({
        PolicyDocument: {
          Statement:
            Match.arrayWith([
              Match.objectLike({
                Action:
                  Match.arrayWith([
                    'secretsmanager:GetSecretValue',
                    'secretsmanager:DescribeSecret',
                  ]),
                Effect:
                  'Allow',
                Resource: {
                  Ref:
                    Match.stringLikeRegexp(
                      '^CognitoClientSecret',
                    ),
                },
              }),
            ]),
        },
      }),
    );
  });

  test('grants API Lambda read and write access to DynamoDB', () => {
    template.hasResourceProperties(
      'AWS::IAM::Policy',
      Match.objectLike({
        Roles:
          Match.arrayWith([
            {
              Ref:
                Match.stringLikeRegexp(
                  '^ApiHandlerServiceRole',
                ),
            },
          ]),
        PolicyDocument: {
          Statement:
            Match.arrayWith([
              Match.objectLike({
                Action:
                  Match.arrayWith([
                    'dynamodb:GetItem',
                    'dynamodb:PutItem',
                    'dynamodb:UpdateItem',
                    'dynamodb:DeleteItem',
                  ]),
                Effect:
                  'Allow',
              }),
            ]),
        },
      }),
    );
  });

  test('grants processor Lambda read access to Slack token table', () => {
    template.hasResourceProperties(
      'AWS::IAM::Policy',
      Match.objectLike({
        Roles:
          Match.arrayWith([
            {
              Ref:
                Match.stringLikeRegexp(
                  '^ProcessorHandlerServiceRole',
                ),
            },
          ]),
        PolicyDocument: {
          Statement:
            Match.arrayWith([
              Match.objectLike({
                Action:
                  Match.arrayWith([
                    'dynamodb:GetItem',
                  ]),
                Effect:
                  'Allow',
              }),
            ]),
        },
      }),
    );
  });

  test('grants processor Lambda permission to invoke Bedrock model', () => {
    template.hasResourceProperties(
      'AWS::IAM::Policy',
      Match.objectLike({
        Roles:
          Match.arrayWith([
            {
              Ref:
                Match.stringLikeRegexp(
                  '^ProcessorHandlerServiceRole',
                ),
            },
          ]),
        PolicyDocument: {
          Statement:
            Match.arrayWith([
              Match.objectLike({
                Action:
                  'bedrock:InvokeModel',
                Effect:
                  'Allow',
                Resource:
                  '*',
              }),
            ]),
        },
      }),
    );
  });

  test('configures input bucket to invoke processor Lambda', () => {
    template.hasResourceProperties(
      'Custom::S3BucketNotifications',
      Match.objectLike({
        BucketName: {
          Ref:
            Match.stringLikeRegexp(
              '^InputBucket',
            ),
        },
        NotificationConfiguration:
          Match.objectLike({
            LambdaFunctionConfigurations:
              Match.arrayWith([
                Match.objectLike({
                  Events: [
                    's3:ObjectCreated:*',
                  ],
                  LambdaFunctionArn: {
                    'Fn::GetAtt': [
                      Match.stringLikeRegexp(
                        '^ProcessorHandler',
                      ),
                      'Arn',
                    ],
                  },
                }),
              ]),
          }),
      }),
    );
  });

  test('grants S3 permission to invoke processor Lambda', () => {
    template.hasResourceProperties(
      'AWS::Lambda::Permission',
      Match.objectLike({
        Action:
          'lambda:InvokeFunction',
        Principal:
          's3.amazonaws.com',
        FunctionName: {
          'Fn::GetAtt': [
            Match.stringLikeRegexp(
              '^ProcessorHandler',
            ),
            'Arn',
          ],
        },
        SourceAccount:
          testAWSAccount,
        SourceArn: {
          'Fn::GetAtt': [
            Match.stringLikeRegexp(
              '^InputBucket',
            ),
            'Arn',
          ],
        },
      }),
    );
  });

  test('outputs CloudFront URL', () => {
    template.hasOutput(
      'CloudFrontURL',
      Match.objectLike({
        Value:
          Match.anyValue(),
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
          Ref:
            Match.stringLikeRegexp(
              '^SlackClientSecret',
            ),
        },
      }),
    );
  });

  test('outputs Cognito user pool ID', () => {
    template.hasOutput(
      'CognitoUserPoolId',
      Match.objectLike({
        Description:
          'Cognito user pool ID',
        Value: {
          Ref:
            Match.stringLikeRegexp(
              '^UserPool',
            ),
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
          Ref:
            Match.stringLikeRegexp(
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
        Value:
          Match.anyValue(),
      }),
    );
  });

  test('outputs Google OAuth redirect URI', () => {
    template.hasOutput(
      'GoogleRedirectUri',
      Match.objectLike({
        Description:
          'Redirect URI for the Google OAuth client',
        Value:
          Match.anyValue(),
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
          Ref:
            Match.stringLikeRegexp(
              '^GoogleOAuthClientSecret',
            ),
        },
      }),
    );
  });
});
