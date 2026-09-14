import * as cdk from 'aws-cdk-lib';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as s3n from 'aws-cdk-lib/aws-s3-notifications';
import { Construct } from 'constructs';

import { createAuthResources } from './constructs/auth-resources';
import { createComputeResources } from './constructs/compute-resources';
import { createDeliveryResources } from './constructs/delivery-resources';
import { createSlackResources } from './constructs/slack-resources';
import { createStorageResources } from './constructs/storage-resources';

function requireEnv(name: string): string {
  const value = process.env[name];

  if (!value) {
    throw new Error(`${name} is required`);
  }

  return value;
}

export class InfraStack extends cdk.Stack {
  constructor(
    scope: Construct,
    id: string,
    props?: cdk.StackProps,
  ) {
    super(scope, id, props);

    const appEnv = requireEnv('APP_ENV');
    const applicationUrl = requireEnv('APPLICATION_URL');
    const bedrockModelId = requireEnv('BEDROCK_MODEL_ID');
    const promptFileName =
      process.env.PROMPT_FILE_NAME || 'extractor.txt';
    const slackClientId = requireEnv('SLACK_CLIENT_ID');
    const slackRedirectUri = requireEnv('SLACK_REDIRECT_URI');
    const googleClientId = requireEnv('GOOGLE_CLIENT_ID');

    // リソース削除設定

    const removalPolicy = cdk.RemovalPolicy.DESTROY;
    const autoDeleteObjects = true;

    // 1. Slackリソースの作成

    const slackResources = createSlackResources(this, {
      appEnv,
      removalPolicy,
    });

    // 2. S3リソースの作成

    const storageResources = createStorageResources(this, {
      removalPolicy,
      autoDeleteObjects,
    });

    // 3. Cognito認証リソースの作成

    const authResources = createAuthResources(this, {
      appEnv,
      applicationUrl,
      googleClientId,
      removalPolicy,
    });

    // 4. Lambda関数の作成（Goランタイム）

    const computeResources = createComputeResources(this, {
      appEnv,
      bedrockModelId,
      cognitoAuthorizationEndpoint:
        authResources.authorizationEndpoint,
      cognitoClientId:
        authResources.userPoolClient.userPoolClientId,
      cognitoClientSecretArn:
        authResources.cognitoClientSecret.secretArn,
      cognitoIssuer: authResources.issuer,
      cognitoLogoutEndpoint: authResources.logoutEndpoint,
      cognitoRedirectUri: authResources.redirectUri,
      cognitoTokenEndpoint: authResources.tokenEndpoint,
      inputBucket: storageResources.inputBucket,
      oauthStateTableName:
        authResources.oauthStateTable.tableName,
      outputBucket: storageResources.outputBucket,
      postLoginRedirectUrl: applicationUrl,
      postLogoutRedirectUrl: applicationUrl,
      promptFileName,
      sessionTableName: authResources.sessionTable.tableName,
      slackClientId,
      slackRedirectUri,
      slackSecret: slackResources.slackSecret,
      slackTokenTable: slackResources.slackTokenTable,
    });

    const { apiHandler, processorHandler } = computeResources;

    // 5. IAM権限とイベントトリガー

    // API HandlerのS3権限

    storageResources.inputBucket.grantWrite(apiHandler);

    // API HandlerのSlack権限

    slackResources.slackTokenTable.grantReadWriteData(
      apiHandler,
    );
    slackResources.slackSecret.grantRead(apiHandler);

    // API Handlerの認証権限

    authResources.cognitoClientSecret.grantRead(apiHandler);
    authResources.oauthStateTable.grantReadWriteData(
      apiHandler,
    );
    authResources.sessionTable.grantReadWriteData(apiHandler);

    // Processor Handlerの権限

    storageResources.inputBucket.grantRead(processorHandler);
    storageResources.outputBucket.grantReadWrite(
      processorHandler,
    );
    slackResources.slackTokenTable.grantReadData(
      processorHandler,
    );

    // Processor HandlerにBedrockの実行権限を付与

    processorHandler.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['bedrock:InvokeModel'],
        resources: ['*'],
      }),
    );

    // Inputバケットへの画像保存時にProcessor Handlerを起動

    storageResources.inputBucket.addEventNotification(
      s3.EventType.OBJECT_CREATED,
      new s3n.LambdaDestination(processorHandler),
    );

    // 6. 配信リソースの作成
    // （API Gateway・CloudFront・Frontend）

    createDeliveryResources(this, {
      apiHandler,
      websiteBucket: storageResources.websiteBucket,
    });
  }
}
