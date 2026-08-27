import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront';
import * as origins from 'aws-cdk-lib/aws-cloudfront-origins';
import * as s3deploy from 'aws-cdk-lib/aws-s3-deployment';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigwv2 from 'aws-cdk-lib/aws-apigatewayv2';
import { HttpLambdaIntegration } from 'aws-cdk-lib/aws-apigatewayv2-integrations';
import * as s3n from 'aws-cdk-lib/aws-s3-notifications';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager';

function requireEnv(name: string): string {
  const value = process.env[name];

  if (!value) {
    throw new Error(`${name} is required`);
  }

  return value;
}

export class InfraStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const appEnv = requireEnv('APP_ENV');
    const bedrockModelId = requireEnv('BEDROCK_MODEL_ID');
    const promptFileName = process.env.PROMPT_FILE_NAME || 'extractor.txt';
    const slackClientId = requireEnv('SLACK_CLIENT_ID');
    const slackClientSecretArn = requireEnv('SLACK_CLIENT_SECRET_ARN');
    const slackRedirectUri = requireEnv('SLACK_REDIRECT_URI');

    const slackSecret = secretsmanager.Secret.fromSecretCompleteArn(
      this,
      'SlackClientSecret',
      slackClientSecretArn,
    );

    // 1. S3 バケットの作成

    // S3 バケット削除設定
    const removalPolicy = cdk.RemovalPolicy.DESTROY;
    const autoDeleteObjects = true;

    // Inputバケット（画像アップロード用：1日で自動削除）
    const inputBucket = new s3.Bucket(this, 'InputBucket', {
      removalPolicy,
      autoDeleteObjects,
      lifecycleRules: [{ expiration: cdk.Duration.days(1) }],
      cors: [{
        allowedMethods: [s3.HttpMethods.PUT],
        allowedOrigins: ['*'],
        allowedHeaders: ['*'],
      }],
    });

    // Outputバケット（生成したJSON保存用：1日で自動削除）
    const outputBucket = new s3.Bucket(this, 'OutputBucket', {
      removalPolicy,
      autoDeleteObjects,
      lifecycleRules: [{ expiration: cdk.Duration.days(1) }],
    });

    // Websiteバケット（フロントエンドの静的ファイルホスティング用）
    const websiteBucket = new s3.Bucket(this, 'WebsiteBucket', {
      removalPolicy,
      autoDeleteObjects,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL, 
    });

    // Slack OAuthトークン保存用テーブル
    const slackTokenTable = new dynamodb.Table(this, 'SlackTokenTable', {
      partitionKey: { name: 'id', type: dynamodb.AttributeType.STRING },
      removalPolicy,
    });

    // 2. Lambda 関数の作成（Goランタイム）

    // API Handler（HTTP API）
    const apiHandler = new lambda.Function(this, 'ApiHandler', {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      architecture: lambda.Architecture.ARM_64,
      code: lambda.Code.fromAsset('../backend/bin/api'),
      timeout: cdk.Duration.seconds(15),
      environment: {
        APP_ENV: appEnv,
        INPUT_BUCKET_NAME: inputBucket.bucketName,
        SLACK_CLIENT_ID: slackClientId,
        SLACK_CLIENT_SECRET_ARN: slackSecret.secretArn,
        SLACK_REDIRECT_URI: slackRedirectUri,
        SLACK_TOKEN_TABLE_NAME: slackTokenTable.tableName,
      },
    });

    // Processor Handler（Bedrock解析・JSON生成・Slack通知）
    const processorHandler = new lambda.Function(
      this,
      'ProcessorHandler',
      {
        runtime: lambda.Runtime.PROVIDED_AL2023,
        handler: 'bootstrap',
        architecture: lambda.Architecture.ARM_64,
        code: lambda.Code.fromAsset('../backend/bin/processor'),
        timeout: cdk.Duration.seconds(30),
        environment: {
          APP_ENV: appEnv,
          INPUT_BUCKET_NAME: inputBucket.bucketName,
          OUTPUT_BUCKET_NAME: outputBucket.bucketName,
          BEDROCK_MODEL_ID: bedrockModelId,
          PROMPT_FILE_NAME: promptFileName,
          SLACK_TOKEN_TABLE_NAME: slackTokenTable.tableName,
        },
      },
    );

    // 3. 権限（IAM）と トリガー（Event）の設定

    // API Handlerには、Inputバケットへの「書き込み」権限と Slack Token Tableへの「読み・書き」と Slack Secret の「読み取り」権限を付与
    inputBucket.grantWrite(apiHandler);
    slackTokenTable.grantReadWriteData(apiHandler);
    slackSecret.grantRead(apiHandler);

    // ProcessorHandlerには、Inputから「読み取り」権限、Outputへ「読み書き」権限、Slack Token Tableへの「読み取り」権限を付与
    inputBucket.grantRead(processorHandler);
    outputBucket.grantReadWrite(processorHandler);
    slackTokenTable.grantReadData(processorHandler);

    // ProcessorHandlerにBedrockの実行権限を付与
    processorHandler.addToRolePolicy(new iam.PolicyStatement({
      actions: ['bedrock:InvokeModel'],
      resources: ['*'],
    }));

    // Inputバケットに画像が入ったら ProcessorHandler を自動起動
    inputBucket.addEventNotification(
      s3.EventType.OBJECT_CREATED,
      new s3n.LambdaDestination(processorHandler)
    );

    // 4. API Gateway の構築 (HTTP API)

    // HTTP API
    const api = new apigwv2.HttpApi(this, 'JpegToJsonHttpApi', {
      apiName: 'Jpeg To Json HTTP API',
      corsPreflight: {
        allowOrigins: ['*'],
        allowMethods: [apigwv2.CorsHttpMethod.ANY],
        allowHeaders: ['*'],
      },
    });

    // API Gateway の各ルートを API Handler に接続
    const apiIntegration = new HttpLambdaIntegration(
      'ApiIntegration',
      apiHandler,
    );

    api.addRoutes({
      path: '/api/storage/upload',
      methods: [apigwv2.HttpMethod.POST],
      integration: apiIntegration,
    });

    api.addRoutes({
      path: '/api/oauth/slack/login',
      methods: [apigwv2.HttpMethod.GET],
      integration: apiIntegration,
    });

    api.addRoutes({
      path: '/api/oauth/slack/callback',
      methods: [apigwv2.HttpMethod.GET],
      integration: apiIntegration,
    });

    // 5. CloudFront の作成

    // CloudFront
    const apiOrigin = new origins.HttpOrigin(
      `${api.apiId}.execute-api.${this.region}.amazonaws.com`,
      {
        protocolPolicy: cloudfront.OriginProtocolPolicy.HTTPS_ONLY,
      },
    );

    const distribution = new cloudfront.Distribution(this, 'WebsiteDistribution', {
      defaultBehavior: {
        origin: origins.S3BucketOrigin.withOriginAccessControl(websiteBucket),
        viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
      },

      additionalBehaviors: {
        '/api/*': {
          origin: apiOrigin,
          viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,

          allowedMethods: cloudfront.AllowedMethods.ALLOW_ALL,

          cachePolicy: cloudfront.CachePolicy.CACHING_DISABLED,

          originRequestPolicy:
            cloudfront.OriginRequestPolicy.ALL_VIEWER_EXCEPT_HOST_HEADER,
        },
      },

      defaultRootObject: 'index.html',
    });

    // デプロイ時は CloudFront のキャッシュを最新に更新
    new s3deploy.BucketDeployment(this, 'DeployWebsite', {
      sources: [s3deploy.Source.asset('../frontend/dist')],
      destinationBucket: websiteBucket,
      distribution: distribution,
      distributionPaths: ['/*'],
    });

    // 6. ログでURLを出力
    new cdk.CfnOutput(this, 'CloudFrontURL', { value: `https://${distribution.distributionDomainName}`});
  }
}
