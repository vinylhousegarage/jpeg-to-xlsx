import * as cdk from 'aws-cdk-lib';
import {
  Match,
  Template,
} from 'aws-cdk-lib/assertions';

import {
  createStorageResources,
  StorageResources,
} from '../../lib/constructs/storage-resources';
import {
  createTestStack,
} from '../infra-stack-test-helpers';

describe('createStorageResources', () => {
  let template: Template;
  let resources: StorageResources;

  beforeAll(() => {
    const stack = createTestStack(
      'StorageResourcesTestStack',
    );

    resources =
      createStorageResources(
        stack,
        {
          removalPolicy:
            cdk.RemovalPolicy.DESTROY,
          autoDeleteObjects:
            true,
        },
      );

    template =
      Template.fromStack(stack);
  });

  test('creates three S3 buckets', () => {
    template.resourceCountIs(
      'AWS::S3::Bucket',
      3,
    );
  });

  test('creates input bucket with one-day expiration', () => {
    template.hasResourceProperties(
      'AWS::S3::Bucket',
      Match.objectLike({
        LifecycleConfiguration:
          Match.objectLike({
            Rules:
              Match.arrayWith([
                Match.objectLike({
                  ExpirationInDays:
                    1,
                  Status:
                    'Enabled',
                }),
              ]),
          }),
        CorsConfiguration:
          Match.anyValue(),
      }),
    );
  });

  test('allows PUT requests to input bucket', () => {
    template.hasResourceProperties(
      'AWS::S3::Bucket',
      Match.objectLike({
        CorsConfiguration:
          Match.objectLike({
            CorsRules:
              Match.arrayWith([
                Match.objectLike({
                  AllowedHeaders: [
                    '*',
                  ],
                  AllowedMethods: [
                    'PUT',
                  ],
                  AllowedOrigins: [
                    '*',
                  ],
                }),
              ]),
          }),
      }),
    );
  });

  test('creates output bucket with one-day expiration', () => {
    template.hasResourceProperties(
      'AWS::S3::Bucket',
      Match.objectLike({
        LifecycleConfiguration:
          Match.objectLike({
            Rules:
              Match.arrayWith([
                Match.objectLike({
                  ExpirationInDays:
                    1,
                  Status:
                    'Enabled',
                }),
              ]),
          }),
        CorsConfiguration:
          Match.absent(),
        PublicAccessBlockConfiguration:
          Match.absent(),
      }),
    );
  });

  test('blocks public access to website bucket', () => {
    template.hasResourceProperties(
      'AWS::S3::Bucket',
      Match.objectLike({
        PublicAccessBlockConfiguration: {
          BlockPublicAcls:
            true,
          BlockPublicPolicy:
            true,
          IgnorePublicAcls:
            true,
          RestrictPublicBuckets:
            true,
        },
      }),
    );
  });

  test('does not configure lifecycle expiration on website bucket', () => {
    template.hasResourceProperties(
      'AWS::S3::Bucket',
      Match.objectLike({
        PublicAccessBlockConfiguration:
          Match.anyValue(),
        LifecycleConfiguration:
          Match.absent(),
      }),
    );
  });

  test('configures automatic object deletion for all buckets', () => {
    template.resourceCountIs(
      'Custom::S3AutoDeleteObjects',
      3,
    );
  });

  test('applies delete policy to all buckets', () => {
    const buckets =
      template.findResources(
        'AWS::S3::Bucket',
      );

    expect(
      Object.values(buckets),
    ).toHaveLength(3);

    for (
      const bucket of Object.values(
        buckets,
      )
    ) {
      expect(
        bucket.DeletionPolicy,
      ).toBe(
        'Delete',
      );

      expect(
        bucket.UpdateReplacePolicy,
      ).toBe(
        'Delete',
      );
    }
  });

  test('returns input bucket', () => {
    expect(
      resources.inputBucket,
    ).toBeDefined();
  });

  test('returns output bucket', () => {
    expect(
      resources.outputBucket,
    ).toBeDefined();
  });

  test('returns website bucket', () => {
    expect(
      resources.websiteBucket,
    ).toBeDefined();
  });
});
