package core

import "xfschainbrowser/chainsync"

func chainSyncCore() {
	syncServe := chainsync.NewSyncService()
	syncServe.Start()
}
