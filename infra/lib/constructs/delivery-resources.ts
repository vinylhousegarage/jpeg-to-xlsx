import * as cdk from 'aws-cdk-lib';
import * as apigwv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {
  HttpLambdaIntegration,
} from 'aws-cdk-lib/aws-apigatewayv2-integrations';
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront';
import * as origins from 'aws-cdk-lib/aws-cloudfront-origins';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as s3deploy from 'aws-cdk-lib/aws-s3-deployment';
import { Construct } from 'constructs';

export type DeliveryResourcesProps = {
  apiHandler: lambda.IFunction;
  websiteBucket: s3.IBucket;
};

export type DeliveryResources = {
  api: apigwv2.HttpApi;
  distribution:
    cloudfront.Distribution;
  applicationUrl: string;
};

export const createDeliveryResources = (
  scope: Construct,
  props: DeliveryResourcesProps,
): DeliveryResources => {
  const stack = cdk.Stack.of(scope);

  const api = new apigwv2.HttpApi(
    scope,
    'JpegToXlsxHttpApi',
    {
      apiName:
        'Jpeg To Xlsx HTTP API',
      corsPreflight: {
        allowOrigins: ['*'],
        allowMethods: [
          apigwv2
            .CorsHttpMethod
            .ANY,
        ],
        allowHeaders: ['*'],
      },
    },
  );

  const apiIntegration =
    new HttpLambdaIntegration(
      'ApiIntegration',
      props.apiHandler,
    );

  api.addRoutes({
    path:
      '/api/storage/upload',
    methods: [
      apigwv2.HttpMethod.POST,
    ],
    integration:
      apiIntegration,
  });

  api.addRoutes({
    path:
      '/api/oauth/slack/login',
    methods: [
      apigwv2.HttpMethod.GET,
    ],
    integration:
      apiIntegration,
  });

  api.addRoutes({
    path:
      '/api/oauth/slack/callback',
    methods: [
      apigwv2.HttpMethod.GET,
    ],
    integration:
      apiIntegration,
  });

  const apiOrigin =
    new origins.HttpOrigin(
      `${api.apiId}.execute-api.${stack.region}.amazonaws.com`,
      {
        protocolPolicy:
          cloudfront
            .OriginProtocolPolicy
            .HTTPS_ONLY,
      },
    );

  const distribution =
    new cloudfront.Distribution(
      scope,
      'WebsiteDistribution',
      {
        defaultBehavior: {
          origin:
            origins
              .S3BucketOrigin
              .withOriginAccessControl(
                props.websiteBucket,
              ),
          viewerProtocolPolicy:
            cloudfront
              .ViewerProtocolPolicy
              .REDIRECT_TO_HTTPS,
        },
        additionalBehaviors: {
          '/api/*': {
            origin:
              apiOrigin,
            viewerProtocolPolicy:
              cloudfront
                .ViewerProtocolPolicy
                .REDIRECT_TO_HTTPS,
            allowedMethods:
              cloudfront
                .AllowedMethods
                .ALLOW_ALL,
            cachePolicy:
              cloudfront
                .CachePolicy
                .CACHING_DISABLED,
            originRequestPolicy:
              cloudfront
                .OriginRequestPolicy
                .ALL_VIEWER_EXCEPT_HOST_HEADER,
          },
        },
        defaultRootObject:
          'index.html',
      },
    );

  const applicationUrl =
    `https://${distribution.distributionDomainName}`;

  new s3deploy.BucketDeployment(
    scope,
    'DeployWebsite',
    {
      sources: [
        s3deploy.Source.asset(
          '../frontend/dist',
        ),
      ],
      destinationBucket:
        props.websiteBucket,
      distribution,
      distributionPaths: ['/*'],
    },
  );

  return {
    api,
    distribution,
    applicationUrl,
  };
};
