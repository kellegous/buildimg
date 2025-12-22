package builder

type OutputFormat string

const (
	OutputFormatJSON  OutputFormat = "rawjson"
	OutputFormatPlain OutputFormat = "plain"
	OutputFormatAuto  OutputFormat = "auto"
	OutputFormatNone  OutputFormat = "none"
	OutputFormatTTY   OutputFormat = "tty"
	OutputFormatQuiet OutputFormat = "quiet"
)
