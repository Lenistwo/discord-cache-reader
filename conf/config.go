package conf

type Config struct {
	FileExtension        string `json:"file_extension"`
	OutputPath           string `json:"output_path"`
	WithModificationTime bool   `json:"with_modification_time"`
}
