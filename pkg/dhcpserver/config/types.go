package config

type ConfigType struct {
	FilePath string `json:"-"`

	ServerConfiguration ServerConfigType
	StaticAssignments   []StaticAssignmentsConfigType
}

// Follows the libvirt naming
type StaticAssignmentsConfigType struct {
	MAC  string
	IP   string
	Name string
}

type ServerConfigType struct {
	IP            string
	StartFrom     string
	LeaseDuration int
	LeaseRange    int
	Options       map[string]string
}
