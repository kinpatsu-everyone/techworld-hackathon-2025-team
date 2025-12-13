package config

var Environment string

const (
	EnvLocal       = "local"
	EnvDevelopment = "development"
	EnvTest        = "test"
	EnvProduction  = "production"
)

func IsLocal() bool {
	return ENV == EnvLocal
}

// IsDebugMode はデバッグモードかどうかを返します。
// Production環境ではない場合にデバッグモードとみなします。
func IsDebugMode() bool {
	return ENV == EnvLocal || ENV == EnvDevelopment || ENV == EnvTest
}

func IsCloud() bool {
	return ENV == EnvDevelopment || ENV == EnvProduction
}

func IsDevelopment() bool {
	return ENV == EnvDevelopment
}

func IsTest() bool {
	return ENV == EnvTest
}

func IsProduction() bool {
	return ENV == EnvProduction
}
