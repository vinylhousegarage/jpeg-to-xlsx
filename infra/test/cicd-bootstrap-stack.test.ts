import * as cdk from 'aws-cdk-lib';
import { Match, Template } from 'aws-cdk-lib/assertions';
import { CicdBootstrapStack } from '../lib/cicd-bootstrap-stack';

const account = '123456789012';
const region = 'ap-northeast-1';

function createStack() {
  const app = new cdk.App();

  const stack = new CicdBootstrapStack(
    app,
    'TestCicdBootstrapStack',
    {
      env: {
        account,
        region,
      },
    },
  );

  return {
    stack,
    template: Template.fromStack(stack),
  };
}

describe('CicdBootstrapStack', () => {
  test('creates CI/CD role with GitHub OIDC trust policy', () => {
    const { stack, template } = createStack();

    const expectedProviderArn = stack.resolve(
      `arn:${stack.partition}:iam::${stack.account}:` +
        'oidc-provider/token.actions.githubusercontent.com',
    );

    template.resourceCountIs('AWS::IAM::Role', 1);

    template.hasResourceProperties(
      'AWS::IAM::Role',
      Match.objectLike({
        RoleName: 'jpeg-to-xlsx-cicd-role',
        Description:
          'Role assumed by GitHub Actions to deploy jpeg-to-xlsx staging resources',
        MaxSessionDuration: 3600,
        AssumeRolePolicyDocument: {
          Statement: Match.arrayWith([
            Match.objectLike({
              Action: 'sts:AssumeRoleWithWebIdentity',
              Effect: 'Allow',
              Principal: {
                Federated: expectedProviderArn,
              },
              Condition: {
                StringEquals: {
                  'token.actions.githubusercontent.com:aud':
                    'sts.amazonaws.com',
                  'token.actions.githubusercontent.com:sub':
                    'repo:vinylhousegarage/jpeg-to-xlsx:' +
                    'environment:staging',
                },
              },
            }),
          ]),
        },
      }),
    );
  });

  test('does not create another GitHub OIDC provider', () => {
    const { template } = createStack();

    template.resourceCountIs('AWS::IAM::OIDCProvider', 0);
  });

  test('allows assuming only CDK bootstrap roles', () => {
    const { stack, template } = createStack();

    const bootstrapRoleNames = [
      'deploy-role',
      'file-publishing-role',
      'image-publishing-role',
      'lookup-role',
    ];

    const expectedRoleArns = bootstrapRoleNames.map((roleName) =>
      stack.resolve(
        `arn:${stack.partition}:iam::${stack.account}:role/` +
          `cdk-hnb659fds-${roleName}-${stack.account}-${stack.region}`,
      ),
    );

    template.resourceCountIs('AWS::IAM::Policy', 1);

    template.hasResourceProperties(
      'AWS::IAM::Policy',
      Match.objectLike({
        PolicyDocument: {
          Statement: Match.arrayWith([
            Match.objectLike({
              Action: 'sts:AssumeRole',
              Effect: 'Allow',
              Resource: expectedRoleArns,
            }),
          ]),
        },
      }),
    );
  });

  test('outputs CI/CD role ARN', () => {
    const { template } = createStack();

    template.hasOutput(
      'JpegToXlsxCicdRoleArn',
      Match.objectLike({
        Description:
          'IAM role ARN used by the jpeg-to-xlsx workflow',
        Value: Match.anyValue(),
      }),
    );
  });
});
