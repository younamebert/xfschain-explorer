package conf

var (
	SyncRquesetURL        string = "http://localhost:9012" // JSON RPC  网关地址
	SyncMaxHashesFetch    int    = 100                     // 请求区块哈希列表最大跨度
	SyncBlockCacheSize    int    = 20                      // 区块缓存数量
	SyncMaxBlocksFetch    int    = 20                      // 区块数据请求并发量
	SyncTimeoutBlockFetch string = "2s"                    // 区块请求数据超时时间
	SyncForceSyncPeriod   string = "10s"                   //强制同步周期
	GeneisisHashPrevBlock string = "0x0000000000000000000000000000000000000000000000000000000000000000"
)
