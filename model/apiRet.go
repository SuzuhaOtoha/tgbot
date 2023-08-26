package model

type LoliconApiRet struct {
	Error string `json:"error"`
	Data  []Data `json:"data"`
}
type Urls struct {
	Original string `json:"original"`
}
type Data struct {
	Pid        int      `json:"pid"`
	P          int      `json:"p"`
	UID        int      `json:"uid"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	R18        bool     `json:"r18"`
	Width      int      `json:"width"`
	Height     int      `json:"height"`
	Tags       []string `json:"tags"`
	Ext        string   `json:"ext"`
	AiType     int      `json:"aiType"`
	UploadDate int64    `json:"uploadDate"`
	Urls       Urls     `json:"urls"`
}
