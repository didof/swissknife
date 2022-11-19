package synflood

type Options struct {
	Verbose bool
}

type SynFloodOptions struct {
	Port          int
	PayloadLength int

	Options
}

var synFloodOptions = &SynFloodOptions{}

func GetSynFloodOptions() *SynFloodOptions {
	return synFloodOptions
}
