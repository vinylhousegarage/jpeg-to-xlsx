import * as cdk from 'aws-cdk-lib';
import {
  Match,
  Template,
} from 'aws-cdk-lib/assertions';

import {
  createSlackResources,
  SlackResources,
} from '../../lib/constructs/slack-resources';
import {
  createTestStack,
  testAppEnv,
} from '../infra-stack-test-helpers';

describe('createSlackResources', () => {
  let template: Template;
  let resources: SlackResources;

  beforeAll(() => {
    const stack = createTestStack('SlackResourcesTestStack');

    resources = createSlackResources(
      stack,
      {
        appEnv: testAppEnv,
        removalPolicy: cdk.RemovalPolicy.DESTROY,
      },
    );

    template = Template.fromStack(stack);
  });

  test('creates Slack client secret', () => {
    template.resourceCountIs(
      'AWS::SecretsManager::Secret',
      1,
    );

    template.hasResourceProperties(
      'AWS::SecretsManager::Secret',
      Match.objectLike({
        Name: 'jpeg-to-xlsx/staging/slack-client-secret',
        Description: 'Slack client secret for jpeg-to-xlsx staging',
      }),
    );
  });

  test('creates Slack token table', () => {
    template.resourceCountIs(
      'AWS::DynamoDB::Table',
      1,
    );

    template.hasResourceProperties(
      'AWS::DynamoDB::Table',
      Match.objectLike({
        AttributeDefinitions: [
          {
            AttributeName: 'id',
            AttributeType: 'S',
          },
        ],
        KeySchema: [
          {
            AttributeName: 'id',
            KeyType: 'HASH',
          },
        ],
        ProvisionedThroughput: Match.anyValue(),
      }),
    );
  });

  test('does not configure TTL on Slack token table', () => {
    template.hasResourceProperties(
      'AWS::DynamoDB::Table',
      Match.objectLike({
        TimeToLiveSpecification: Match.absent(),
      }),
    );
  });

  test('applies delete policy to Slack client secret', () => {
    template.hasResource(
      'AWS::SecretsManager::Secret',
      Match.objectLike({
        DeletionPolicy: 'Delete',
        UpdateReplacePolicy: 'Delete',
      }),
    );
  });

  test('applies delete policy to Slack token table', () => {
    template.hasResource(
      'AWS::DynamoDB::Table',
      Match.objectLike({
        DeletionPolicy: 'Delete',
        UpdateReplacePolicy: 'Delete',
      }),
    );
  });

  test('returns Slack client secret', () => {
    expect(resources.slackSecret).toBeDefined();
  });

  test('returns Slack token table', () => {
    expect(resources.slackTokenTable).toBeDefined();
  });
});
