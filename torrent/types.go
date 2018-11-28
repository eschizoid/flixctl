package torrent

type Result struct {
	FileURL    string
	Magnet     string
	DescURL    string
	Name       string
	Size       string
	Quality    string
	Seeders    int
	Leechers   int
	UploadDate string
	Source     string
	FilePath   string
}

type Search struct {
	In              string
	Out             []Result
	SourcesToLookup []string
}
