import { Amplify } from 'aws-amplify';

function requireEnv(
  name: keyof ImportMetaEnv,
): string {
  const value = import.meta.env[name];

  if (!value || value.trim() === '') {
    throw new Error(`${name} is required`);
  }

  return value;
}

export function configureAmplify(): void {
  const domain = requireEnv(
    'VITE_COGNITO_DOMAIN',
  );
  const redirectSignIn = requireEnv(
    'VITE_COGNITO_REDIRECT_SIGN_IN',
  );
  const redirectSignOut = requireEnv(
    'VITE_COGNITO_REDIRECT_SIGN_OUT',
  );
  const userPoolClientId = requireEnv(
    'VITE_COGNITO_USER_POOL_CLIENT_ID',
  );
  const userPoolId = requireEnv(
    'VITE_COGNITO_USER_POOL_ID',
  );

  Amplify.configure({
    Auth: {
      Cognito: {
        userPoolId,
        userPoolClientId,
        loginWith: {
          oauth: {
            domain,
            scopes: [
              'openid',
              'email',
              'profile',
            ],
            redirectSignIn: [redirectSignIn],
            redirectSignOut: [redirectSignOut],
            responseType: 'code',
          },
        },
      },
    },
  });
}
