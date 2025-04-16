package config

type ProcessesSettings struct {
	Common              CommonProcessSetting `yaml:"common" json:"common"`
	CustomFilterSetting CustomFilterSetting  `yaml:"customFilterSetting" json:"customFilterSetting"`
}

type CommonProcessSetting struct {
	Size int `yaml:"size" json:"size"`
}

type CustomFilterSetting struct {
	Common   CommonProcessSetting `yaml:"common" json:"common"`
	MinValue int                  `yaml:"minValue" json:"minValue"`
}
