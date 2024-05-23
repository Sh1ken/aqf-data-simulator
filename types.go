package main

type Column struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Position int    `json:"position"`
}

type DateConfig struct {
	Format string `json:"format"`
	Gmt    string `json:"gmt"`
}

type File struct {
	Name       string     `json:"name"`
	Separator  string     `json:"separator"`
	Interval   int        `json:"intervalMinutes"`
	DateConfig DateConfig `json:"dateConfig"`
	Columns    []Column   `json:"columns"`
}

type Config struct {
	FolderData string `json:"folderData"`
	FolderTemp string `json:"folderTemp"`
	Files      []File `json:"files"`
}
