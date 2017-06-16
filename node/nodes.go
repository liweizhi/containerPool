package node



type(
	Info struct {
		HostName	string	`json:"host_name,omitempty"`
		CPU		string `json:"cpu,omitempty"`
		Memory		string `json:"memory,omitempty"`
		OS		string `json:"os,omitempty"`
		IP		[]string `json:"ip,omitempty"`


	}
)
