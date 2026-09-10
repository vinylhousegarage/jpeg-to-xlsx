import * as cdk from 'aws-cdk-lib';
import * as cognito from 'aws-cdk-lib/aws-cognito';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
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
  userPoolDomain: cognito.UserPoolDomain;
  userPoolClient: cognito.UserPoolClient;
  googleOAuthSecret: secretsmanager.Secret;
  cognitoClientSecret: secretsmanager.Secret;
  oauthStateTable: dynamodb.Table;
  sessionTable: dynamodb.Table;
  googleRedirectUri: string;
  issuer: string;
  authorizationEndpoint: string;
  tokenEndpoint: string;
  logoutEndpoint: string;
  redirectUri: string;
};

export const createAuthResources = (
  scope: Construct,
  props: AuthResourcesProps,
): AuthResources => {
  const stack = cdk.Stack.of(scope);

  // Googleアカウント連携用User Pool

  const userPool = new cognito.UserPool(
    scope,
    'UserPool',
    {
      userPoolName:
        `jpeg-to-xlsx-${props.appEnv}-users`,
      selfSignUpEnabled: false,
      signInCaseSensitive: false,
      removalPolicy:
        props.removalPolicy,
    },
  );

  // Cognito Managed Login用Domain

  const cognitoDomainPrefix =
    `jpeg-to-xlsx-${props.appEnv}-${stack.account}`;

  const userPoolDomain = userPool.addDomain(
    'UserPoolDomain',
    { cognitoDomain: { domainPrefix: cognitoDomainPrefix } },
  );

  // Google OAuth Client Secret保存用Secret

  const googleOAuthSecret = new secretsmanager.Secret(
    scope,
    'GoogleOAuthClientSecret',
    {
      secretName: `jpeg-to-xlsx/${props.appEnv}/google-oauth-client`,
      description: `Google OAuth client secret for jpeg-to-xlsx ${props.appEnv}`,
      removalPolicy: props.removalPolicy,
    },
  );

  // CognitoとGoogle OAuthの連携

  const googleProvider = new cognito.UserPoolIdentityProviderGoogle(
    scope,
    'GoogleIdentityProvider',
    {
      userPool,
      clientId: props.googleClientId,
      clientSecretValue: googleOAuthSecret.secretValueFromJson('client_secret'),
      scopes: [
        'openid',
        'email',
        'profile',
      ],
      attributeMapping: {
        email: cognito.ProviderAttribute.GOOGLE_EMAIL,
        givenName: cognito.ProviderAttribute.GOOGLE_GIVEN_NAME,
        familyName: cognito.ProviderAttribute.GOOGLE_FAMILY_NAME,
      },
    },
  );

  // Cognito認証後にBFFへ戻るURI

  const redirectUri = `${props.applicationUrl}/api/auth/callback`;

  // Cognito User Pool Client

  const userPoolClient = userPool.addClient(
    'UserPoolClient',
    {
      userPoolClientName: `jpeg-to-xlsx-${props.appEnv}-client`,
      generateSecret: true,
      preventUserExistenceErrors: true,
      supportedIdentityProviders: [
        cognito.UserPoolClientIdentityProvider.GOOGLE,
      ],
      oAuth: {
        flows: { authorizationCodeGrant: true },
        scopes: [
          cognito.OAuthScope.OPENID,
          cognito.OAuthScope.EMAIL,
          cognito.OAuthScope.PROFILE,
        ],
        callbackUrls: [redirectUri],
        logoutUrls: [props.applicationUrl],
      },
    },
  );

  // Google IdP作成後にUser Pool Clientを作成

  userPoolClient.node.addDependency(googleProvider);

  // Cognito App Client Secret保存用Secret

  const cognitoClientSecret = new secretsmanager.Secret(
    scope,
    'CognitoClientSecret',
    {
      secretName: `jpeg-to-xlsx/${props.appEnv}/cognito-client-secret`,
      description: `Cognito client secret for jpeg-to-xlsx ${props.appEnv}`,
      secretObjectValue: {
        client_secret: userPoolClient.userPoolClientSecret,
      },
      removalPolicy: props.removalPolicy,
    },
  );

  // Authorization Code Flowで使用するstate・nonce・PKCE verifier保存用テーブル

  const oauthStateTable = new dynamodb.Table(
    scope,
    'CognitoOAuthStateTable',
    {
      tableName: `jpeg-to-xlsx-${props.appEnv}-cognito-oauth-states`,
      partitionKey: {
        name: 'state',
        type: dynamodb.AttributeType.STRING,
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      timeToLiveAttribute: 'expires_at',
      removalPolicy: props.removalPolicy,
    },
  );

  // BFF認証セッション保存用テーブル

  const sessionTable = new dynamodb.Table(
    scope,
    'AuthSessionTable',
    {
      tableName: `jpeg-to-xlsx-${props.appEnv}-auth-sessions`,
      partitionKey: {
        name: 'id_hash',
        type: dynamodb.AttributeType.STRING,
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      timeToLiveAttribute: 'expires_at',
      removalPolicy: props.removalPolicy,
    },
  );

  // Cognito関連URI

  const googleRedirectUri = `${userPoolDomain.baseUrl()}/oauth2/idpresponse`;
  const issuer = `https://cognito-idp.${stack.region}.${stack.urlSuffix}/${userPool.userPoolId}`;
  const authorizationEndpoint = `${userPoolDomain.baseUrl()}/oauth2/authorize`;
  const tokenEndpoint = `${userPoolDomain.baseUrl()}/oauth2/token`;
  const logoutEndpoint = `${userPoolDomain.baseUrl()}/logout`;

  return {
    userPool,
    userPoolDomain,
    userPoolClient,
    googleOAuthSecret,
    cognitoClientSecret,
    oauthStateTable,
    sessionTable,
    googleRedirectUri,
    issuer,
    authorizationEndpoint,
    tokenEndpoint,
    logoutEndpoint,
    redirectUri,
  };
};
