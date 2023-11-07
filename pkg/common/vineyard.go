package common

// Runtime for Vineyard
const (
	VineyardRuntime = "vineyard"

	VineyardMountType = "vineyard-fuse"

	VineyardChart = VineyardRuntime

	VineyardRuntimeImageEnv = "Vineyard_RUNTIME_IMAGE_ENV"

	VineyardFuseImageEnv = "Vineyard_FUSE_IMAGE_ENV"

	DefaultVineyardRuntimeImage = "registry.cn-huhehaote.aliyuncs.com/Vineyard/Vineyard:2.3.0-SNAPSHOT-2c41226"

	DefaultVineyardFuseImage = "registry.cn-huhehaote.aliyuncs.com/Vineyard/Vineyard-fuse:2.3.0-SNAPSHOT-2c41226"
)

var (
	// Vineyard ufs root path
	VineyardMountPathFormat = RootDirPath + "%s"

	VineyardLocalStorageRootPath   = "/underFSStorage"
	VineyardLocalStoragePathFormat = VineyardLocalStorageRootPath + "/%s"
)
