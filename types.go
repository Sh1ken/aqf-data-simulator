package main

type Column struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Position int    `json:"position"`
}

type DateConfig struct {
	Format   string `json:"format"`
	Gmt      string `json:"gmt"`
	Interval int    `json:"intervalMinutes"`
}

type File struct {
	Name       string     `json:"name"`
	Separator  string     `json:"separator"`
	DateConfig DateConfig `json:"dateConfig"`
	Columns    []Column   `json:"columns"`
}

type Config struct {
	DataFolder string `json:"dataFolder"`
	TempFolder string `json:"tempFolder"`
	Files      []File `json:"files"`
}
