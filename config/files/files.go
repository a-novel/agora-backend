package configfiles

import (
	_ "embed"
	"encoding/json"
)

var (
	//go:embed development.yml
	DevCFG json.RawMessage
	//go:embed production.yml
	ProdCFG json.RawMessage
	//go:embed generic.yml
	GenericCFG json.RawMessage
)
