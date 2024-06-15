package fshare

// custom constants
const (
	FsAPIURL = "https://api.fshare.vn"
	TimFsURL = "https://timfshare.com"
)

// TokenResponse model
type TokenResponse struct {
	Code      int64  `json:"code"`
	Msg       string `json:"msg"`
	SessionID string `json:"session_id"`
	Token     string `json:"token"`
}

// FileResponse model
type FileResponse struct {
	Crc32         string `json:"crc32"`
	Created       string `json:"created"`
	Deleted       string `json:"deleted"`
	Description   string `json:"description"`
	Directlink    string `json:"directlink"`
	Downloadcount string `json:"downloadcount"`
	FileType      string `json:"file_type"`
	FolderPath    string `json:"folder_path"`
	HashIndex     string `json:"hash_index"`
	ID            string `json:"id"`
	Linkcode      string `json:"linkcode"`
	Mimetype      string `json:"mimetype"`
	Modified      string `json:"modified"`
	Name          string `json:"name"`
	OwnerID       string `json:"owner_id"`
	Path          string `json:"path"`
	Pid           string `json:"pid"`
	Public        string `json:"public"`
	Pwd           string `json:"pwd"`
	Realname      string `json:"realname"`
	Secure        string `json:"secure"`
	Shared        string `json:"shared"`
	Size          string `json:"size"`
	StorageID     string `json:"storage_id"`
}

// FolderResponse model
type FolderResponse struct {
	Data []*FileResponse
}

// FileDownloadResponse model
type FileDownloadResponse struct {
	Location string `json:"location"`
}

// QueryResponse model
type QueryResponse struct {
	Code    int64                  `json:"code"`
	Message string                 `json:"message"`
	Data    []*QueryDetailResponse `json:"data"`
}

// QueryDetailResponse model
type QueryDetailResponse struct {
	FileType int64  `json:"file_type"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	URL      string `json:"url"`
}
