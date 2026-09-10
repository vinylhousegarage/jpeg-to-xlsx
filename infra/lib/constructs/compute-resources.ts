import * as cdk from 'aws-cdk-lib';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager';
import { Construct } from 'constructs';

export type ComputeResourcesProps = {
  appEnv: string;
  bedrockModelId: string;
  cognitoAuthorizationEndpoint: string;
  cognitoClientId: string;
  cognitoClientSecretArn: string;
  cognitoIssuer: string;
  cognitoLogoutEndpoint: string;
  cognitoRedirectUri: string;
  cognitoTokenEndpoint: string;
  inputBucket: s3.IBucket;
  oauthStateTableName: string;
  outputBucket: s3.IBucket;
  postLoginRedirectUrl: string;
  postLogoutRedirectUrl: string;
  promptFileName: string;
  sessionTableName: string;
  slackClientId: string;
  slackRedirectUri: string;
  slackSecret: secretsmanager.ISecret;
  slackTokenTable: dynamodb.ITable;
};

export type ComputeResources = {
  apiHandler: lambda.Function;
  processorHandler: lambda.Function;
};

export const createComputeResources = (
  scope: Construct,
  props: ComputeResourcesProps,
): ComputeResources => {
  // API Handler

  const apiHandler = new lambda.Function(
    scope,
    'ApiHandler',
    {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      architecture: lambda.Architecture.ARM_64,
      code: lambda.Code.fromAsset('../backend/bin/api'),
      timeout: cdk.Duration.seconds(15),
      environment: {
        APP_ENV: props.appEnv,
        AUTH_LOGOUT_REDIRECT_URL: props.postLogoutRedirectUrl,
        AUTH_REDIRECT_URL: props.postLoginRedirectUrl,
        AUTH_SESSION_TABLE_NAME: props.sessionTableName,
        COGNITO_AUTHORIZATION_ENDPOINT: props.cognitoAuthorizationEndpoint,
        COGNITO_CLIENT_ID: props.cognitoClientId,
        COGNITO_CLIENT_SECRET_ARN: props.cognitoClientSecretArn,
        COGNITO_ISSUER: props.cognitoIssuer,
        COGNITO_LOGOUT_ENDPOINT: props.cognitoLogoutEndpoint,
        COGNITO_OAUTH_STATE_TABLE_NAME: props.oauthStateTableName,
        COGNITO_REDIRECT_URI: props.cognitoRedirectUri,
        COGNITO_TOKEN_ENDPOINT: props.cognitoTokenEndpoint,
        INPUT_BUCKET_NAME: props.inputBucket.bucketName,
        SLACK_CLIENT_ID: props.slackClientId,
        SLACK_CLIENT_SECRET_ARN: props.slackSecret.secretArn,
        SLACK_REDIRECT_URI: props.slackRedirectUri,
        SLACK_TOKEN_TABLE_NAME: props.slackTokenTable.tableName,
      },
    },
  );

  // Processor Handler（Bedrock解析・XLSX生成・Slack通知）

  const processorHandler = new lambda.Function(
    scope,
    'ProcessorHandler',
    {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      architecture: lambda.Architecture.ARM_64,
      code: lambda.Code.fromAsset('../backend/bin/processor'),
      timeout: cdk.Duration.seconds(30),
      environment: {
        APP_ENV: props.appEnv,
        BEDROCK_MODEL_ID: props.bedrockModelId,
        INPUT_BUCKET_NAME: props.inputBucket.bucketName,
        OUTPUT_BUCKET_NAME: props.outputBucket.bucketName,
        PROMPT_FILE_NAME: props.promptFileName,
        SLACK_TOKEN_TABLE_NAME: props.slackTokenTable.tableName,
      },
    },
  );

  return {
    apiHandler,
    processorHandler,
  };
};
