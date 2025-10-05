package remote_cache_loader

import (
	"IMP/app/internal/ports"
	"IMP/app/pkg/logger"
	"encoding/json"
	"os"
	"time"
)

const localCacheStorageDir = "tmp/cache/"

func LoadLocalFile[T any](resource ports.CachedRemoteResource) (T, error) {
	if isNeedToBeDownload(localCacheStorageDir+resource.LocalFileName(), resource.GetLifeTime()) {
		logger.Info("Need to download resource from remote", map[string]interface{}{
			"resource": resource.LocalFileName(),
		})
		err := saveToLocalCacheStorage(resource)
		if err != nil {
			logger.Error("Could not save resource to local cache", map[string]interface{}{
				"resource": resource.LocalFileName(),
			})
		}
		logger.Info("Resource was successfully downloaded", map[string]interface{}{
			"resource": resource.LocalFileName(),
		})
	}

	return getFromLocalStorage[T](localCacheStorageDir + resource.LocalFileName())
}

func isNeedToBeDownload(filePath string, exp time.Duration) bool {
	fileInfo, err := os.Stat(filePath)

	if err != nil {
		return true
	}

	return time.Now().After(fileInfo.ModTime().Add(exp))
}

func saveToLocalCacheStorage(resource ports.CachedRemoteResource) error {
	if err := os.MkdirAll(localCacheStorageDir, 0755); err != nil {
		return err
	}

	data, err := resource.Load()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(localCacheStorageDir+resource.LocalFileName(), jsonData, 0644)
}

func getFromLocalStorage[T any](filePath string) (T, error) {
	var result T

	data, err := os.ReadFile(filePath)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)

	return result, nil
}
