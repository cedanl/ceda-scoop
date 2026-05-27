package internal

// UpdateCheckDoneMsg carries the result of the async GitHub version check.
type UpdateCheckDoneMsg struct {
	Latest string
	Err    error
}

// UpdateDownloadDoneMsg carries the result of a binary download.
type UpdateDownloadDoneMsg struct {
	Path string // path to new binary (Unix: current exe replaced; Windows: *-update.exe)
	Err  error
}
