package testdata

import "embed"

var (
	//go:embed passwd api-spec.yaml
	SpecFS embed.FS

	//go:embed api-spec.yaml
	APISpecFile []byte

	//go:embed template-single.html
	TemplateSingleFile []byte

	//go:embed template-multiple.html
	TemplateMultipleFile []byte
)

const (
	PasswdFilePath  = "passwd"
	APISpecFilePath = "api-spec.yaml"
)
