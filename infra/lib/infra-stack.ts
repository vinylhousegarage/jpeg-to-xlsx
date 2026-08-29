import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront';
import * as origins from 'aws-cdk-lib/aws-cloudfront-origins';
import * as s3deploy from 'aws-cdk-lib/aws-s3-deployment';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigwv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {
  HttpLambdaIntegration,
} from 'aws-cdk-lib/aws-apigatewayv2-integrations';
import * as s3n from 'aws-cdk-lib/aws-s3-notifications';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager';
import * as cognito from 'aws-cdk-lib/aws-cognito';

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
    const bedrockModelId = requireEnv('BEDROCK_MODEL_ID');
    const promptFileName =
      process.env.PROMPT_FILE_NAME || 'extractor.txt';
    const slackClientId = requireEnv('SLACK_CLIENT_ID');
    const slackRedirectUri = requireEnv(
      'SLACK_REDIRECT_URI',
    );

    // リソース削除設定
    const removalPolicy = cdk.RemovalPolicy.DESTROY;
    const autoDeleteObjects = true;

    // 1. Cognitoの作成

    // Googleアカウント連携用User Pool
    const userPool = new cognito.UserPool(
      this,
      'UserPool',
      {
        userPoolName: `jpeg-to-xlsx-${appEnv}-users`,
        selfSignUpEnabled: false,
        signInCaseSensitive: false,
        removalPolicy,
      },
    );

    // Cognito Managed Login用Domain
    const cognitoDomainPrefix =
      `jpeg-to-xlsx-${appEnv}-${this.account}`;

    const userPoolDomain = userPool.addDomain(
      'UserPoolDomain',
      {
        cognitoDomain: {
          domainPrefix: cognitoDomainPrefix,
        },
      },
    );

    // 2. Slackリソースの作成

    // Slack Client Secret保存用Secret
    const slackSecret = new secretsmanager.Secret(
      this,
      'SlackClientSecret',
      {
        secretName:
          `jpeg-to-xlsx/${appEnv}/slack-client-secret`,
        description:
          `Slack client secret for jpeg-to-xlsx ${appEnv}`,
        removalPolicy,
      },
    );

    // Slack OAuthトークン保存用テーブル
    const slackTokenTable = new dynamodb.Table(
      this,
      'SlackTokenTable',
      {
        partitionKey: {
          name: 'id',
          type: dynamodb.AttributeType.STRING,
        },
        removalPolicy,
      },
    );

    // 3. Googleリソースの作成

    // Google OAuth Client Secret保存用Secret
    const googleOAuthSecret = new secretsmanager.Secret(
      this,
      'GoogleOAuthClientSecret',
      {
        secretName:
          `jpeg-to-xlsx/${appEnv}/google-oauth-client`,
        description:
          `Google OAuth client secret for jpeg-to-xlsx ${appEnv}`,
        removalPolicy,
      },
    );

    // 4. S3バケットの作成

    // Inputバケット（画像アップロード用：1日で自動削除）
    const inputBucket = new s3.Bucket(
      this,
      'InputBucket',
      {
        removalPolicy,
        autoDeleteObjects,
        lifecycleRules: [
          {
            expiration: cdk.Duration.days(1),
          },
        ],
        cors: [
          {
            allowedMethods: [s3.HttpMethods.PUT],
            allowedOrigins: ['*'],
            allowedHeaders: ['*'],
          },
        ],
      },
    );

    // Outputバケット（生成したXLSX保存用：1日で自動削除）
    const outputBucket = new s3.Bucket(
      this,
      'OutputBucket',
      {
        removalPolicy,
        autoDeleteObjects,
        lifecycleRules: [
          {
            expiration: cdk.Duration.days(1),
          },
        ],
      },
    );

    // Websiteバケット
    const websiteBucket = new s3.Bucket(
      this,
      'WebsiteBucket',
      {
        removalPolicy,
        autoDeleteObjects,
        blockPublicAccess:
          s3.BlockPublicAccess.BLOCK_ALL,
      },
    );

    // 5. Lambda関数の作成（Goランタイム）

    // API Handler（HTTP API）
    const apiHandler = new lambda.Function(
      this,
      'ApiHandler',
      {
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
          SLACK_TOKEN_TABLE_NAME:
            slackTokenTable.tableName,
        },
      },
    );

    // Processor Handler
    // （Bedrock解析・XLSX生成・Slack通知）
    const processorHandler = new lambda.Function(
      this,
      'ProcessorHandler',
      {
        runtime: lambda.Runtime.PROVIDED_AL2023,
        handler: 'bootstrap',
        architecture: lambda.Architecture.ARM_64,
        code: lambda.Code.fromAsset(
          '../backend/bin/processor',
        ),
        timeout: cdk.Duration.seconds(30),
        environment: {
          APP_ENV: appEnv,
          INPUT_BUCKET_NAME: inputBucket.bucketName,
          OUTPUT_BUCKET_NAME: outputBucket.bucketName,
          BEDROCK_MODEL_ID: bedrockModelId,
          PROMPT_FILE_NAME: promptFileName,
          SLACK_TOKEN_TABLE_NAME:
            slackTokenTable.tableName,
        },
      },
    );

    // 6. IAM権限とイベントトリガーの設定

    // API Handlerの権限
    inputBucket.grantWrite(apiHandler);
    slackTokenTable.grantReadWriteData(apiHandler);
    slackSecret.grantRead(apiHandler);

    // Processor Handlerの権限
    inputBucket.grantRead(processorHandler);
    outputBucket.grantReadWrite(processorHandler);
    slackTokenTable.grantReadData(processorHandler);

    // Processor HandlerにBedrockの実行権限を付与
    processorHandler.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['bedrock:InvokeModel'],
        resources: ['*'],
      }),
    );

    // Inputバケットへの画像保存時にProcessor Handlerを起動
    inputBucket.addEventNotification(
      s3.EventType.OBJECT_CREATED,
      new s3n.LambdaDestination(processorHandler),
    );

    // 7. API Gatewayの構築（HTTP API）

    const api = new apigwv2.HttpApi(
      this,
      'JpegToXlsxHttpApi',
      {
        apiName: 'Jpeg To Xlsx HTTP API',
        corsPreflight: {
          allowOrigins: ['*'],
          allowMethods: [
            apigwv2.CorsHttpMethod.ANY,
          ],
          allowHeaders: ['*'],
        },
      },
    );

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

    // 8. CloudFrontの作成

    const apiOrigin = new origins.HttpOrigin(
      `${api.apiId}.execute-api.${this.region}.amazonaws.com`,
      {
        protocolPolicy:
          cloudfront.OriginProtocolPolicy.HTTPS_ONLY,
      },
    );

    const distribution = new cloudfront.Distribution(
      this,
      'WebsiteDistribution',
      {
        defaultBehavior: {
          origin:
            origins.S3BucketOrigin.withOriginAccessControl(
              websiteBucket,
            ),
          viewerProtocolPolicy:
            cloudfront.ViewerProtocolPolicy
              .REDIRECT_TO_HTTPS,
        },
        additionalBehaviors: {
          '/api/*': {
            origin: apiOrigin,
            viewerProtocolPolicy:
              cloudfront.ViewerProtocolPolicy
                .REDIRECT_TO_HTTPS,
            allowedMethods:
              cloudfront.AllowedMethods.ALLOW_ALL,
            cachePolicy:
              cloudfront.CachePolicy.CACHING_DISABLED,
            originRequestPolicy:
              cloudfront.OriginRequestPolicy
                .ALL_VIEWER_EXCEPT_HOST_HEADER,
          },
        },
        defaultRootObject: 'index.html',
      },
    );

    new s3deploy.BucketDeployment(
      this,
      'DeployWebsite',
      {
        sources: [
          s3deploy.Source.asset('../frontend/dist'),
        ],
        destinationBucket: websiteBucket,
        distribution,
        distributionPaths: ['/*'],
      },
    );

    // 9. Outputs

    new cdk.CfnOutput(this, 'CloudFrontURL', {
      value:
        `https://${distribution.distributionDomainName}`,
    });

    new cdk.CfnOutput(this, 'SlackClientSecretArn', {
      value: slackSecret.secretArn,
      description:
        'Secrets Manager ARN for the Slack client secret',
    });

    new cdk.CfnOutput(this, 'CognitoUserPoolId', {
      value: userPool.userPoolId,
      description: 'Cognito user pool ID',
    });

    new cdk.CfnOutput(this, 'CognitoDomain', {
      value: userPoolDomain.baseUrl(),
      description: 'Cognito managed login domain',
    });

    new cdk.CfnOutput(this, 'GoogleRedirectUri', {
      value:
        `${userPoolDomain.baseUrl()}/oauth2/idpresponse`,
      description:
        'Redirect URI for the Google OAuth client',
    });

    new cdk.CfnOutput(
      this,
      'GoogleOAuthClientSecretArn',
      {
        value: googleOAuthSecret.secretArn,
        description:
          'Secrets Manager ARN for the Google OAuth client secret',
      },
    );
  }
}
