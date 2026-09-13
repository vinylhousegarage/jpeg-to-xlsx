import * as cdk from 'aws-cdk-lib';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager';
import { Construct } from 'constructs';

export type SlackResourcesProps = {
  appEnv: string;
  removalPolicy: cdk.RemovalPolicy;
};

export type SlackResources = {
  slackSecret: secretsmanager.Secret;
  slackTokenTable: dynamodb.Table;
};

export const createSlackResources = (
  scope: Construct,
  props: SlackResourcesProps,
): SlackResources => {
  // Slack Client Secret保存用Secret
  const slackSecret = new secretsmanager.Secret(
    scope,
    'SlackClientSecret',
    {
      secretName: `jpeg-to-xlsx/${props.appEnv}/slack-client-secret`,
      description: `Slack client secret for jpeg-to-xlsx ${props.appEnv}`,
      removalPolicy: props.removalPolicy,
    },
  );

  // CognitoユーザーごとのSlack OAuthトークン保存用テーブル
  const slackTokenTable = new dynamodb.Table(
    scope,
    'SlackTokenTable',
    {
      partitionKey: {
        name: 'cognito_sub',
        type: dynamodb.AttributeType.STRING,
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: props.removalPolicy,
    },
  );

  return {
    slackSecret,
    slackTokenTable,
  };
};
