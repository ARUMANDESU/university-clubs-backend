package domain

type File struct {
	Name  string `json:"file_name"`
	Bytes []byte `json:"file_bytes"`
	Size  int64  `json:"file_size"`
	Type  string `json:"file_type"`
}
