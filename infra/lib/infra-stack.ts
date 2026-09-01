import * as cdk from 'aws-cdk-lib';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as s3n from 'aws-cdk-lib/aws-s3-notifications';
import { Construct } from 'constructs';

import {
  createAuthResources,
} from './constructs/auth-resources';
import {
  createComputeResources,
} from './constructs/compute-resources';
import {
  createDeliveryResources,
} from './constructs/delivery-resources';
import {
  createSlackResources,
} from './constructs/slack-resources';
import {
  createStorageResources,
} from './constructs/storage-resources';

function requireEnv(
  name: string,
): string {
  const value = process.env[name];

  if (!value) {
    throw new Error(
      `${name} is required`,
    );
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

    const appEnv = requireEnv(
      'APP_ENV',
    );

    const bedrockModelId = requireEnv(
      'BEDROCK_MODEL_ID',
    );

    const promptFileName =
      process.env.PROMPT_FILE_NAME ||
      'extractor.txt';

    const slackClientId = requireEnv(
      'SLACK_CLIENT_ID',
    );

    const slackRedirectUri = requireEnv(
      'SLACK_REDIRECT_URI',
    );

    const googleClientId = requireEnv(
      'GOOGLE_CLIENT_ID',
    );

    // リソース削除設定

    const removalPolicy =
      cdk.RemovalPolicy.DESTROY;

    const autoDeleteObjects = true;

    // 1. Slackリソースの作成

    const slackResources =
      createSlackResources(
        this,
        {
          appEnv,
          removalPolicy,
        },
      );

    // 2. S3リソースの作成

    const storageResources =
      createStorageResources(
        this,
        {
          removalPolicy,
          autoDeleteObjects,
        },
      );

    // 3. Lambda関数の作成
    // （Goランタイム）

    const computeResources =
      createComputeResources(
        this,
        {
          appEnv,
          bedrockModelId,
          promptFileName,
          slackClientId,
          slackRedirectUri,
          inputBucket:
            storageResources
              .inputBucket,
          outputBucket:
            storageResources
              .outputBucket,
          slackSecret:
            slackResources
              .slackSecret,
          slackTokenTable:
            slackResources
              .slackTokenTable,
        },
      );

    // 4. IAM権限とイベントトリガー

    // API Handlerの権限

    storageResources.inputBucket
      .grantWrite(
        computeResources
          .apiHandler,
      );

    slackResources.slackTokenTable
      .grantReadWriteData(
        computeResources
          .apiHandler,
      );

    slackResources.slackSecret
      .grantRead(
        computeResources
          .apiHandler,
      );

    // Processor Handlerの権限

    storageResources.inputBucket
      .grantRead(
        computeResources
          .processorHandler,
      );

    storageResources.outputBucket
      .grantReadWrite(
        computeResources
          .processorHandler,
      );

    slackResources.slackTokenTable
      .grantReadData(
        computeResources
          .processorHandler,
      );

    // Processor Handlerに
    // Bedrockの実行権限を付与

    computeResources.processorHandler
      .addToRolePolicy(
        new iam.PolicyStatement({
          actions: [
            'bedrock:InvokeModel',
          ],
          resources: ['*'],
        }),
      );

    // Inputバケットへの画像保存時に
    // Processor Handlerを起動

    storageResources.inputBucket
      .addEventNotification(
        s3.EventType.OBJECT_CREATED,
        new s3n.LambdaDestination(
          computeResources
            .processorHandler,
        ),
      );

    // 5. 配信リソースの作成
    // （API Gateway・CloudFront・Frontend）

    const deliveryResources =
      createDeliveryResources(
        this,
        {
          apiHandler:
            computeResources
              .apiHandler,
          websiteBucket:
            storageResources
              .websiteBucket,
        },
      );

    // 6. Cognito認証リソースの作成

    const authResources =
      createAuthResources(
        this,
        {
          appEnv,
          applicationUrl:
            deliveryResources
              .applicationUrl,
          googleClientId,
          removalPolicy,
        },
      );

    // 7. Outputs

    new cdk.CfnOutput(
      this,
      'CloudFrontURL',
      {
        value:
          deliveryResources
            .applicationUrl,
      },
    );

    new cdk.CfnOutput(
      this,
      'SlackClientSecretArn',
      {
        value:
          slackResources
            .slackSecret
            .secretArn,
        description:
          'Secrets Manager ARN for the Slack client secret',
      },
    );

    new cdk.CfnOutput(
      this,
      'CognitoUserPoolId',
      {
        value:
          authResources
            .userPool
            .userPoolId,
        description:
          'Cognito user pool ID',
      },
    );

    new cdk.CfnOutput(
      this,
      'CognitoDomain',
      {
        value:
          authResources
            .userPoolDomain
            .baseUrl(),
        description:
          'Cognito managed login domain',
      },
    );

    new cdk.CfnOutput(
      this,
      'GoogleRedirectUri',
      {
        value:
          authResources
            .googleRedirectUri,
        description:
          'Redirect URI for the Google OAuth client',
      },
    );

    new cdk.CfnOutput(
      this,
      'GoogleOAuthClientSecretArn',
      {
        value:
          authResources
            .googleOAuthSecret
            .secretArn,
        description:
          'Secrets Manager ARN for the Google OAuth client secret',
      },
    );

    new cdk.CfnOutput(
      this,
      'CognitoUserPoolClientId',
      {
        value:
          authResources
            .userPoolClient
            .userPoolClientId,
        description:
          'Cognito user pool client ID',
      },
    );
  }
}
