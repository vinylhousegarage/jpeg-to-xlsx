package config

func loadAPIStorageConfig() (
	StorageConfig,
	error,
) {
	inputBucketName, err :=
		loadInputBucketName()
	if err != nil {
		return StorageConfig{}, err
	}

	return StorageConfig{
		InputBucketName: inputBucketName,
	}, nil
}

func loadProcessorStorageConfig() (
	StorageConfig,
	error,
) {
	inputBucketName, err :=
		loadInputBucketName()
	if err != nil {
		return StorageConfig{}, err
	}

	outputBucketName, err := loadRequiredEnv(
		"OUTPUT_BUCKET_NAME",
	)
	if err != nil {
		return StorageConfig{}, err
	}

	return StorageConfig{
		InputBucketName:  inputBucketName,
		OutputBucketName: outputBucketName,
	}, nil
}

func loadInputBucketName() (
	string,
	error,
) {
	return loadRequiredEnv(
		"INPUT_BUCKET_NAME",
	)
}
