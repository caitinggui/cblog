package V

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	EmptyIntId uint64 = 0  // 用来做数据插入的空主键
	EmptyStrId string = "" // 用来做数据插入的空主键

	MaxPageSize     uint = 1000 // 分页的每页最大条数
	DefaultPageSize uint = 10   // 分页的每页默认条数

	DefaultExpiration time.Duration = 8 * time.Hour      // cache过期时间
	CleanupInterval   time.Duration = time.Hour          // 定期清理cache的间隔时间
	NoExpiration      time.Duration = cache.NoExpiration // cache永不过期

	CurrentUser string = "CurrentUser" // 保存当前用户的key

	PraseIpTimeout time.Duration = 3 * time.Second // 解析ip地区的超时时间

	// cache keys
	VisitorSum         = "VisitorSum"
	OtherArticleDomain = "otherArticle.domain"

	// 附件所在目录
	AttachmentDirectory = "media/blog/attachments/"
	AttachmentSeparator = "/"
)
