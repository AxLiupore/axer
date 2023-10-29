package container

const (
	RUNNING       = "running"
	STOP          = "stopped"
	Exit          = "exited"
	InfoLoc       = "/var/run/axer/"
	InfoLocFormat = InfoLoc + "%s/"
	ConfigName    = "config.json"
	IDLength      = 32
	LogFile       = "container.log"
)

// Info container info
type Info struct {
	Pid         string `json:"pid"`        // PID of the container's init process on the host machine
	Id          string `json:"id"`         // Container ID
	Name        string `json:"name"`       // Container name
	Command     string `json:"command"`    // Command used to run the init process inside the container
	CreatedTime string `json:"createTime"` // Creation time of the container
	Status      string `json:"status"`     // Status of the container
}
