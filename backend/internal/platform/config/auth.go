package config

func loadBFFAuthConfig() (
	BFFAuthConfig,
	error,
) {
	cognitoClientID, err := loadRequiredEnv(
		"COGNITO_CLIENT_ID",
	)
	if err != nil {
		return BFFAuthConfig{}, err
	}

	cognitoClientSecretARN, err :=
		loadRequiredEnv(
			"COGNITO_CLIENT_SECRET_ARN",
		)
	if err != nil {
		return BFFAuthConfig{}, err
	}

	cognitoIssuer, err := loadRequiredEnv(
		"COGNITO_ISSUER",
	)
	if err != nil {
		return BFFAuthConfig{}, err
	}

	cognitoAuthorizationEndpoint, err :=
		loadRequiredEnv(
			"COGNITO_AUTHORIZATION_ENDPOINT",
		)
	if err != nil {
		return BFFAuthConfig{}, err
	}

	cognitoTokenEndpoint, err :=
		loadRequiredEnv(
			"COGNITO_TOKEN_ENDPOINT",
		)
	if err != nil {
		return BFFAuthConfig{}, err
	}

	cognitoRedirectURI, err :=
		loadRequiredEnv(
			"COGNITO_REDIRECT_URI",
		)
	if err != nil {
		return BFFAuthConfig{}, err
	}

	postLoginRedirectURL, err :=
		loadRequiredEnv(
			"AUTH_REDIRECT_URL",
		)
	if err != nil {
		return BFFAuthConfig{}, err
	}

	oauthStateTableName, err :=
		loadRequiredEnv(
			"COGNITO_OAUTH_STATE_TABLE_NAME",
		)
	if err != nil {
		return BFFAuthConfig{}, err
	}

	sessionTableName, err := loadRequiredEnv(
		"AUTH_SESSION_TABLE_NAME",
	)
	if err != nil {
		return BFFAuthConfig{}, err
	}

	return BFFAuthConfig{
		CognitoClientID:              cognitoClientID,
		CognitoClientSecretARN:       cognitoClientSecretARN,
		CognitoIssuer:                cognitoIssuer,
		CognitoAuthorizationEndpoint: cognitoAuthorizationEndpoint,
		CognitoTokenEndpoint:         cognitoTokenEndpoint,
		CognitoRedirectURI:           cognitoRedirectURI,
		PostLoginRedirectURL:         postLoginRedirectURL,
		OAuthStateTableName:          oauthStateTableName,
		OAuthStateTTL:                defaultOAuthStateTTL,
		SessionTableName:             sessionTableName,
		SessionLifetime:              defaultAuthSessionLifetime,
	}, nil
}
