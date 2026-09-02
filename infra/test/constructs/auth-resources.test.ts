import * as cdk from 'aws-cdk-lib';
import {
  Match,
  Template,
} from 'aws-cdk-lib/assertions';

import {
  AuthResources,
  createAuthResources,
} from '../../lib/constructs/auth-resources';
import {
  createTestStack,
  testAppEnv,
  testApplicationUrl,
  testGoogleClientId,
} from '../infra-stack-test-helpers';

describe('createAuthResources', () => {
  let template: Template;
  let resources: AuthResources;

  beforeAll(() => {
    const stack = createTestStack(
      'AuthResourcesTestStack',
    );

    resources = createAuthResources(
      stack,
      {
        appEnv: testAppEnv,
        applicationUrl: testApplicationUrl,
        googleClientId: testGoogleClientId,
        removalPolicy: cdk.RemovalPolicy.DESTROY,
      },
    );

    template = Template.fromStack(stack);
  });

  test('creates Cognito user pool', () => {
    template.resourceCountIs(
      'AWS::Cognito::UserPool',
      1,
    );

    template.hasResourceProperties(
      'AWS::Cognito::UserPool',
      Match.objectLike({
        UserPoolName: 'jpeg-to-xlsx-staging-users',
        AdminCreateUserConfig: Match.objectLike({ AllowAdminCreateUserOnly: true }),
        UsernameConfiguration: { CaseSensitive: false },
      }),
    );
  });

  test('creates Cognito user pool domain', () => {
    template.resourceCountIs(
      'AWS::Cognito::UserPoolDomain',
      1,
    );

    template.hasResourceProperties(
      'AWS::Cognito::UserPoolDomain',
      Match.objectLike({
        Domain: 'jpeg-to-xlsx-staging-123456789012',
        UserPoolId: { Ref: Match.stringLikeRegexp('^UserPool') },
      }),
    );
  });

  test('creates Google OAuth client secret', () => {
    template.hasResourceProperties(
      'AWS::SecretsManager::Secret',
      Match.objectLike({
        Name: 'jpeg-to-xlsx/staging/google-oauth-client',
        Description: 'Google OAuth client secret for jpeg-to-xlsx staging',
      }),
    );
  });

  test('creates Google identity provider', () => {
    template.resourceCountIs(
      'AWS::Cognito::UserPoolIdentityProvider',
      1,
    );

    template.hasResourceProperties(
      'AWS::Cognito::UserPoolIdentityProvider',
      Match.objectLike({
        ProviderName: 'Google',
        ProviderType: 'Google',
        ProviderDetails: Match.objectLike({
          client_id: testGoogleClientId,
          client_secret: Match.anyValue(),
          authorize_scopes: 'openid email profile',
        }),
        AttributeMapping: Match.objectLike({
          email: 'email',
          given_name: 'given_name',
          family_name: 'family_name',
        }),
        UserPoolId: {
          Ref: Match.stringLikeRegexp('^UserPool'),
        },
      }),
    );
  });

  test('creates Cognito user pool client for BFF', () => {
    template.resourceCountIs(
      'AWS::Cognito::UserPoolClient',
      1,
    );

    template.hasResourceProperties(
      'AWS::Cognito::UserPoolClient',
      Match.objectLike({
        ClientName: 'jpeg-to-xlsx-staging-client',
        GenerateSecret: true,
        PreventUserExistenceErrors: 'ENABLED',
        SupportedIdentityProviders: ['Google'],
        AllowedOAuthFlows: ['code'],
        AllowedOAuthFlowsUserPoolClient: true,
        AllowedOAuthScopes: Match.arrayWith([
          'openid',
          'email',
          'profile',
        ]),
        CallbackURLs: [`${testApplicationUrl}/api/auth/callback`],
        LogoutURLs: [testApplicationUrl],
        UserPoolId: { Ref: Match.stringLikeRegexp('^UserPool') },
      }),
    );
  });

  test('creates Cognito client secret', () => {
    template.hasResourceProperties(
      'AWS::SecretsManager::Secret',
      Match.objectLike({
        Name: 'jpeg-to-xlsx/staging/cognito-client-secret',
        Description: 'Cognito client secret for jpeg-to-xlsx staging',
        SecretString: Match.anyValue(),
      }),
    );
  });

  test('creates two authentication secrets', () => {
    template.resourceCountIs(
      'AWS::SecretsManager::Secret',
      2,
    );
  });

  test('creates Cognito OAuth state table', () => {
    template.hasResourceProperties(
      'AWS::DynamoDB::Table',
      Match.objectLike({
        TableName: 'jpeg-to-xlsx-staging-cognito-oauth-states',
        BillingMode: 'PAY_PER_REQUEST',
        AttributeDefinitions: [
          {
            AttributeName: 'state',
            AttributeType: 'S',
          },
        ],
        KeySchema: [
          {
            AttributeName: 'state',
            KeyType: 'HASH',
          },
        ],
        TimeToLiveSpecification: {
          AttributeName: 'expires_at',
          Enabled: true,
        },
      }),
    );
  });

  test('creates authentication session table', () => {
    template.hasResourceProperties(
      'AWS::DynamoDB::Table',
      Match.objectLike({
        TableName: 'jpeg-to-xlsx-staging-auth-sessions',
        BillingMode: 'PAY_PER_REQUEST',
        AttributeDefinitions: [
          {
            AttributeName: 'id_hash',
            AttributeType: 'S',
          },
        ],
        KeySchema: [
          {
            AttributeName: 'id_hash',
            KeyType: 'HASH',
          },
        ],
        TimeToLiveSpecification: {
          AttributeName: 'expires_at',
          Enabled: true,
        },
      }),
    );
  });

  test('creates two authentication tables', () => {
    template.resourceCountIs(
      'AWS::DynamoDB::Table',
      2,
    );
  });

  test('returns Cognito redirect URI', () => {
    expect(resources.redirectUri).toBe(
      `${testApplicationUrl}/api/auth/callback`,
    );
  });

  test('returns Cognito authorization endpoint', () => {
    expect(
      resources.authorizationEndpoint,
    ).toContain('/oauth2/authorize');
  });

  test('returns Cognito token endpoint', () => {
    expect(
      resources.tokenEndpoint,
    ).toContain('/oauth2/token');
  });

  test('returns Google OAuth redirect URI', () => {
    expect(
      resources.googleRedirectUri,
    ).toContain('/oauth2/idpresponse');
  });

  test('returns Cognito issuer', () => {
    expect(resources.issuer).toMatch(
      /^https:\/\/cognito-idp\.ap-northeast-1\..+\/.+$/,
    );
  });
});
