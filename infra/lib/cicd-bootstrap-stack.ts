import * as cdk from 'aws-cdk-lib';
import * as iam from 'aws-cdk-lib/aws-iam';
import { Construct } from 'constructs';

export class CicdBootstrapStack extends cdk.Stack {
  constructor(
    scope: Construct,
    id: string,
    props?: cdk.StackProps,
  ) {
    super(scope, id, props);

    // GitHub Actions OIDC Providerを参照
    const githubOidcProvider =
      iam.OpenIdConnectProvider.fromOpenIdConnectProviderArn(
        this,
        'GitHubOidcProvider',
        `arn:${this.partition}:iam::${this.account}:oidc-provider/token.actions.githubusercontent.com`,
      );

    // staging環境専用CI/CDロール
    const cicdRole = new iam.Role(this, 'JpegToXlsxCicdRole', {
      roleName: 'jpeg-to-xlsx-cicd-role',
      description:
        'Role assumed by GitHub Actions to deploy jpeg-to-xlsx staging resources',
      assumedBy: new iam.WebIdentityPrincipal(
        githubOidcProvider.openIdConnectProviderArn,
        {
          StringEquals: {
            'token.actions.githubusercontent.com:aud':
              'sts.amazonaws.com',
            'token.actions.githubusercontent.com:sub':
              'repo:vinylhousegarage@172001646/jpeg-to-xlsx@1347837190:environment:staging',
          },
        },
      ),
      maxSessionDuration: cdk.Duration.hours(1),
    });

    // GitHub ActionsからCDK標準bootstrapロールを引き受ける権限
    cicdRole.addToPolicy(
      new iam.PolicyStatement({
        actions: ['sts:AssumeRole'],
        resources: [
          `arn:${this.partition}:iam::${this.account}:role/cdk-hnb659fds-deploy-role-${this.account}-${this.region}`,
          `arn:${this.partition}:iam::${this.account}:role/cdk-hnb659fds-file-publishing-role-${this.account}-${this.region}`,
          `arn:${this.partition}:iam::${this.account}:role/cdk-hnb659fds-image-publishing-role-${this.account}-${this.region}`,
          `arn:${this.partition}:iam::${this.account}:role/cdk-hnb659fds-lookup-role-${this.account}-${this.region}`,
        ],
      }),
    );

    new cdk.CfnOutput(this, 'JpegToXlsxCicdRoleArn', {
      value: cicdRole.roleArn,
      description: 'IAM role ARN used by the jpeg-to-xlsx workflow',
    });
  }
}
