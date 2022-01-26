package meiwobuxing

type PathType int

type PathMapping struct {
	Name string   `json:"name"`
	Type PathType `json:"type"`
	Desc string   `json:"desc"`
	Path string   `json:"path"`
}

const (
	PathTypeInvalid     = PathType(0x00)
	PathTypeFileObject  = PathType(0x01)
	PathTypeAndroidCert = PathType(0x02)
	PathTypeIOSCert     = PathType(0x03)
)
