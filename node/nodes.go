package node



type(
	info struct {
		hostname	string	`json:"host_name,omitempty"`
		os		string `json:"os,omitempty"`
		ip		string `json:"ip,omitempty"`
		memory		string `json:"memory,omitempty"`
		cpu		string `json:"cpu,omitempty"`
	}
)
