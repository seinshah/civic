package generator

type CV struct {
	resumeConfigPath string
}

func NewCV(configPath string) *CV {
	return &CV{
		resumeConfigPath: configPath,
	}
}

func (*CV) Generate() error {
	return nil
}
