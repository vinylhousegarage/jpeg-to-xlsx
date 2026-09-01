import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as s3 from 'aws-cdk-lib/aws-s3';
import {
  Match,
  Template,
} from 'aws-cdk-lib/assertions';

import {
  createDeliveryResources,
  DeliveryResources,
} from '../../lib/constructs/delivery-resources';
import {
  createTestStack,
} from '../infra-stack-test-helpers';

describe('createDeliveryResources', () => {
  let template: Template;
  let resources: DeliveryResources;

  beforeAll(() => {
    const stack = createTestStack(
      'DeliveryResourcesTestStack',
    );

    const apiHandler =
      new lambda.Function(
        stack,
        'ApiHandler',
        {
          runtime:
            lambda.Runtime
              .NODEJS_24_X,
          handler:
            'index.handler',
          code:
            lambda.Code.fromInline(
              'exports.handler = async () => ({ statusCode: 200 });',
            ),
        },
      );

    const websiteBucket =
      new s3.Bucket(
        stack,
        'WebsiteBucket',
      );

    resources =
      createDeliveryResources(
        stack,
        {
          apiHandler,
          websiteBucket,
        },
      );

    template =
      Template.fromStack(stack);
  });

  test('creates HTTP API', () => {
    template.resourceCountIs(
      'AWS::ApiGatewayV2::Api',
      1,
    );

    template.hasResourceProperties(
      'AWS::ApiGatewayV2::Api',
      Match.objectLike({
        Name:
          'Jpeg To Xlsx HTTP API',
        ProtocolType:
          'HTTP',
      }),
    );
  });

  test('does not configure HTTP API CORS', () => {
    template.hasResourceProperties(
      'AWS::ApiGatewayV2::Api',
      Match.objectLike({
        ProtocolType:
          'HTTP',
        CorsConfiguration:
          Match.absent(),
      }),
    );
  });

  test('creates Cognito login route', () => {
    template.hasResourceProperties(
      'AWS::ApiGatewayV2::Route',
      Match.objectLike({
        RouteKey:
          'GET /api/auth/login',
        Target:
          Match.anyValue(),
      }),
    );
  });

  test('creates Cognito callback route', () => {
    template.hasResourceProperties(
      'AWS::ApiGatewayV2::Route',
      Match.objectLike({
        RouteKey:
          'GET /api/auth/callback',
        Target:
          Match.anyValue(),
      }),
    );
  });

  test('creates Slack login route', () => {
    template.hasResourceProperties(
      'AWS::ApiGatewayV2::Route',
      Match.objectLike({
        RouteKey:
          'GET /api/oauth/slack/login',
        Target:
          Match.anyValue(),
      }),
    );
  });

  test('creates Slack callback route', () => {
    template.hasResourceProperties(
      'AWS::ApiGatewayV2::Route',
      Match.objectLike({
        RouteKey:
          'GET /api/oauth/slack/callback',
        Target:
          Match.anyValue(),
      }),
    );
  });

  test('creates storage upload route', () => {
    template.hasResourceProperties(
      'AWS::ApiGatewayV2::Route',
      Match.objectLike({
        RouteKey:
          'POST /api/storage/upload',
        Target:
          Match.anyValue(),
      }),
    );
  });

  test('creates five API routes', () => {
    template.resourceCountIs(
      'AWS::ApiGatewayV2::Route',
      5,
    );
  });

  test('creates HTTP API default stage', () => {
    template.hasResourceProperties(
      'AWS::ApiGatewayV2::Stage',
      Match.objectLike({
        StageName:
          '$default',
        AutoDeploy:
          true,
      }),
    );
  });

  test('grants HTTP API permission to invoke API Lambda', () => {
    template.hasResourceProperties(
      'AWS::Lambda::Permission',
      Match.objectLike({
        Action:
          'lambda:InvokeFunction',
        Principal:
          'apigateway.amazonaws.com',
        FunctionName: {
          'Fn::GetAtt': [
            Match.stringLikeRegexp(
              '^ApiHandler',
            ),
            'Arn',
          ],
        },
      }),
    );
  });

  test('creates CloudFront distribution', () => {
    template.resourceCountIs(
      'AWS::CloudFront::Distribution',
      1,
    );

    template.hasResourceProperties(
      'AWS::CloudFront::Distribution',
      Match.objectLike({
        DistributionConfig:
          Match.objectLike({
            Enabled:
              true,
            DefaultRootObject:
              'index.html',
          }),
      }),
    );
  });

  test('redirects default behavior to HTTPS', () => {
    template.hasResourceProperties(
      'AWS::CloudFront::Distribution',
      Match.objectLike({
        DistributionConfig:
          Match.objectLike({
            DefaultCacheBehavior:
              Match.objectLike({
                ViewerProtocolPolicy:
                  'redirect-to-https',
              }),
          }),
      }),
    );
  });

  test('forwards API paths to API Gateway', () => {
    template.hasResourceProperties(
      'AWS::CloudFront::Distribution',
      Match.objectLike({
        DistributionConfig:
          Match.objectLike({
            CacheBehaviors:
              Match.arrayWith([
                Match.objectLike({
                  PathPattern:
                    '/api/*',
                  ViewerProtocolPolicy:
                    'redirect-to-https',
                  AllowedMethods:
                    Match.arrayWith([
                      'GET',
                      'HEAD',
                      'OPTIONS',
                      'PUT',
                      'PATCH',
                      'POST',
                      'DELETE',
                    ]),
                  Compress:
                    true,
                }),
              ]),
          }),
      }),
    );
  });

  test('creates origin access control', () => {
    template.resourceCountIs(
      'AWS::CloudFront::OriginAccessControl',
      1,
    );

    template.hasResourceProperties(
      'AWS::CloudFront::OriginAccessControl',
      Match.objectLike({
        OriginAccessControlConfig:
          Match.objectLike({
            OriginAccessControlOriginType:
              's3',
            SigningBehavior:
              'always',
            SigningProtocol:
              'sigv4',
          }),
      }),
    );
  });

  test('creates frontend bucket deployment', () => {
    template.hasResourceProperties(
      'Custom::CDKBucketDeployment',
      Match.objectLike({
        DestinationBucketName: {
          Ref:
            Match.stringLikeRegexp(
              '^WebsiteBucket',
            ),
        },
        DistributionId: {
          Ref:
            Match.stringLikeRegexp(
              '^WebsiteDistribution',
            ),
        },
        DistributionPaths: [
          '/*',
        ],
      }),
    );
  });

  test('returns HTTP API', () => {
    expect(
      resources.api,
    ).toBeDefined();
  });

  test('returns CloudFront distribution', () => {
    expect(
      resources.distribution,
    ).toBeDefined();
  });

  test('returns application URL', () => {
    expect(
      resources.applicationUrl,
    ).toContain(
      'https://',
    );
  });
});
