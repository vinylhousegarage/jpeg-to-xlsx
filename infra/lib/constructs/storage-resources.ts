import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import { Construct } from 'constructs';

export type StorageResourcesProps = {
  removalPolicy: cdk.RemovalPolicy;
  autoDeleteObjects: boolean;
};

export type StorageResources = {
  inputBucket: s3.Bucket;

  outputBucket: s3.Bucket;

  websiteBucket: s3.Bucket;
};

export const createStorageResources = (
  scope: Construct,
  props: StorageResourcesProps,
): StorageResources => {
  // Inputバケット
  // （画像アップロード用：1日で自動削除）
  const inputBucket =
    new s3.Bucket(
      scope,
      'InputBucket',
      {
        removalPolicy:
          props.removalPolicy,
        autoDeleteObjects:
          props.autoDeleteObjects,
        lifecycleRules: [
          {
            expiration:
              cdk.Duration.days(1),
          },
        ],
        cors: [
          {
            allowedMethods: [
              s3.HttpMethods.PUT,
            ],
            allowedOrigins: ['*'],
            allowedHeaders: ['*'],
          },
        ],
      },
    );

  // Outputバケット
  // （生成したXLSX保存用：1日で自動削除）
  const outputBucket =
    new s3.Bucket(
      scope,
      'OutputBucket',
      {
        removalPolicy:
          props.removalPolicy,
        autoDeleteObjects:
          props.autoDeleteObjects,
        lifecycleRules: [
          {
            expiration:
              cdk.Duration.days(1),
          },
        ],
      },
    );

  // Websiteバケット
  const websiteBucket =
    new s3.Bucket(
      scope,
      'WebsiteBucket',
      {
        removalPolicy:
          props.removalPolicy,
        autoDeleteObjects:
          props.autoDeleteObjects,
        blockPublicAccess:
          s3.BlockPublicAccess
            .BLOCK_ALL,
      },
    );

  return {
    inputBucket,
    outputBucket,
    websiteBucket,
  };
};
