package common

//  app新闻读取请求 将http get请求转换成该结构
type NewsListRequest struct {
	LocationId uint32 `json:"location_id"`  // 地区码
	UserId     uint32 `json:"user_id"`      // 用户id
	LastNewsId uint32 `json:"last_news_id"` // 上一条新闻id
	NewsCnt    uint32 `json:"news_cnt"`     // 获取条数
	CookieStr  string `json:"cookie_str"`
}

//  app请求返回值
type NewsListResponse struct {
	Ret  int32    `json:"ret"`
	Data NewsList `json:"data"`
}

//  新闻列表
type NewsList struct {
	News []NewsItemDesc `json:"list"`
}

const (
	kDispTypeNorm    = 1 // 正常新闻
	kDispTypeMutiPic = 2 // 组图
)

//  一个新闻的描述
type NewsItemDesc struct {
	Id       uint32      `json:"id"`
	Type     uint32      `json:"disp_type"`
	Title    int8        `json:"title"`
	Abstract string      `json:"abstract"`
	ViewUrl  string      `json:"view_url"`
	Images   []ImageInfo `json:"images"`
}

// 图片列表

type ImageInfo struct {
	Id      uint32 `json:"id"`
	ImgUrl  string `json:"img_url"`
	ImgDesc string `json:"img_desc"`
}

//   索引中一组新闻 ，一个地区为一组
type NewsCluster struct {
	news []NewsMetaInfo
}

//  一条新闻 元信息
type NewsMetaInfo struct {
	Id      uint32      `json:"id"`
	Type    uint32      `json:"disp_type"`
	Title   int8        `json:"title"`
	ViewUrl string      `json:"view_url"`
	Images  []ImageInfo `json:"images"`
}
