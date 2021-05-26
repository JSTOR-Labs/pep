package globals

var BuildJobs map[string]int

func init() {
	BuildJobs = make(map[string]int)
}