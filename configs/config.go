package configs

var Conf = Config{}

type Config struct {
	Env           string `mapstructure:"ENV"`
	Port          int    `mapstructure:"PORT"`
	DB            string `mapstructure:"DB"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	MigrateTables bool   `mapstructure:"MIGRATE_TABLES"`
	AuthKey       string `mapstructure:"AUTH_KEY"`
}
