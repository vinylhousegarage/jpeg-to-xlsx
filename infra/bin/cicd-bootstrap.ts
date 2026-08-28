#!/usr/bin/env node

import * as cdk from 'aws-cdk-lib';
import { CicdBootstrapStack } from '../lib/cicd-bootstrap-stack';

const app = new cdk.App();

new CicdBootstrapStack(
  app,
  'Staging-JpegToXlsx-CicdBootstrapStack',
  {
    env: {
      account: process.env.CDK_DEFAULT_ACCOUNT,
      region: 'ap-northeast-1',
    },
  },
);
