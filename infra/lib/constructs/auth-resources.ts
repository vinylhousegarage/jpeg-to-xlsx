import * as cdk from 'aws-cdk-lib';
import * as cognito from 'aws-cdk-lib/aws-cognito';
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager';
import { Construct } from 'constructs';

export type AuthResourcesProps = {
  appEnv: string;
  applicationUrl: string;
  googleClientId: string;
  removalPolicy: cdk.RemovalPolicy;
};

export type AuthResources = {
  userPool: cognito.UserPool;

  userPoolDomain:
    cognito.UserPoolDomain;

  userPoolClient:
    cognito.UserPoolClient;

  googleOAuthSecret:
    secretsmanager.Secret;

  googleRedirectUri: string;
};

export const createAuthResources = (
  scope: Construct,
  props: AuthResourcesProps,
): AuthResources => {
  // Googleアカウント連携用User Pool
  const userPool = new cognito.UserPool(
    scope,
    'UserPool',
    {
      userPoolName:
        `jpeg-to-xlsx-${props.appEnv}-users`,
      selfSignUpEnabled: false,
      signInCaseSensitive: false,
      removalPolicy: props.removalPolicy,
    },
  );

  // Cognito Managed Login用Domain
  const stack = cdk.Stack.of(scope);

  const cognitoDomainPrefix =
    `jpeg-to-xlsx-${props.appEnv}-${stack.account}`;

  const userPoolDomain =
    userPool.addDomain(
      'UserPoolDomain',
      {
        cognitoDomain: {
          domainPrefix:
            cognitoDomainPrefix,
        },
      },
    );

  // Google OAuth Client Secret保存用Secret
  const googleOAuthSecret =
    new secretsmanager.Secret(
      scope,
      'GoogleOAuthClientSecret',
      {
        secretName:
          `jpeg-to-xlsx/${props.appEnv}/google-oauth-client`,
        description:
          `Google OAuth client secret for jpeg-to-xlsx ${props.appEnv}`,
        removalPolicy:
          props.removalPolicy,
      },
    );

  // CognitoとGoogle OAuthの連携
  const googleProvider =
    new cognito.UserPoolIdentityProviderGoogle(
      scope,
      'GoogleIdentityProvider',
      {
        userPool,
        clientId:
          props.googleClientId,
        clientSecretValue:
          googleOAuthSecret
            .secretValueFromJson(
              'client_secret',
            ),
        scopes: [
          'openid',
          'email',
          'profile',
        ],
        attributeMapping: {
          email:
            cognito.ProviderAttribute
              .GOOGLE_EMAIL,
          givenName:
            cognito.ProviderAttribute
              .GOOGLE_GIVEN_NAME,
          familyName:
            cognito.ProviderAttribute
              .GOOGLE_FAMILY_NAME,
        },
      },
    );

  // Cognito User Pool Client
  const userPoolClient =
    userPool.addClient(
      'UserPoolClient',
      {
        userPoolClientName:
          `jpeg-to-xlsx-${props.appEnv}-client`,
        generateSecret: false,
        preventUserExistenceErrors: true,
        supportedIdentityProviders: [
          cognito
            .UserPoolClientIdentityProvider
            .GOOGLE,
        ],
        oAuth: {
          flows: {
            authorizationCodeGrant: true,
          },
          scopes: [
            cognito.OAuthScope.OPENID,
            cognito.OAuthScope.EMAIL,
            cognito.OAuthScope.PROFILE,
          ],
          callbackUrls: [
            props.applicationUrl,
          ],
          logoutUrls: [
            props.applicationUrl,
          ],
        },
      },
    );

  // Google IdP作成後にUser Pool Clientを作成
  userPoolClient.node.addDependency(
    googleProvider,
  );

  const googleRedirectUri =
    `${userPoolDomain.baseUrl()}/oauth2/idpresponse`;

  return {
    userPool,
    userPoolDomain,
    userPoolClient,
    googleOAuthSecret,
    googleRedirectUri,
  };
};
