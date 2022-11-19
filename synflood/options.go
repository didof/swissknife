package synflood

type options struct {
	Verbose bool
}

type SynFloodOptions struct {
	Port                      int
	PayloadLength             int
	FloodDurationMilliseconds int

	options
}

var synFloodOptions = &SynFloodOptions{}

func GetSynFloodOptions() *SynFloodOptions {
	return synFloodOptions
}
