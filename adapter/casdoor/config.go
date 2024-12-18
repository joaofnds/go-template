package casdoor

type Config struct {
	Endpoint         string `mapstructure:"endpoint"`
	ClientID         string `mapstructure:"client_id"`
	ClientSecret     string `mapstructure:"client_secret"`
	Certificate      string `mapstructure:"certificate"`
	OrganizationName string `mapstructure:"organization_name"`
	ApplicationName  string `mapstructure:"application_name"`
}
