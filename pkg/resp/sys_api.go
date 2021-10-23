package resp

type ApiResp struct {
	Base
	Method   string `json:"method"`
	Path     string `json:"path"`
	Category string `json:"category"`
	Desc     string `json:"desc"`
	Title    string `json:"title"`
}

type ApiGroupByCategoryResp struct {
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Children []ApiResp `json:"children"`
}

type ApiTreeWithAccessResp struct {
	List      []ApiGroupByCategoryResp `json:"list"`
	AccessIds []uint                   `json:"accessIds"`
}
