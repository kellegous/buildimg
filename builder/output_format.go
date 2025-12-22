package builder

import "fmt"

type OutputFormat string

const (
	OutputFormatJSON  OutputFormat = "rawjson"
	OutputFormatPlain OutputFormat = "plain"
	OutputFormatAuto  OutputFormat = "auto"
	OutputFormatNone  OutputFormat = "none"
	OutputFormatTTY   OutputFormat = "tty"
	OutputFormatQuiet OutputFormat = "quiet"
)

func (o OutputFormat) IsValid() bool {
	switch o {
	case OutputFormatJSON, OutputFormatPlain, OutputFormatAuto, OutputFormatNone, OutputFormatTTY, OutputFormatQuiet:
		return true
	default:
		return false
	}
}

func (o *OutputFormat) Set(v string) error {
	*o = OutputFormat(v)
	if !o.IsValid() {
		return fmt.Errorf("invalid output format: %s", v)
	}
	return nil
}

func (o *OutputFormat) String() string {
	return string(*o)
}

func (o *OutputFormat) Type() string {
	return "output format"
}
