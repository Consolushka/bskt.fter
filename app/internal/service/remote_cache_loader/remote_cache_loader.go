package remote_cache_loader

import (
	"IMP/app/internal/ports"
	"IMP/app/pkg/logger"
	"encoding/json"
	"fmt"
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
				"error":    err,
			})
		} else {
			logger.Info("Resource was successfully downloaded", map[string]interface{}{
				"resource": resource.LocalFileName(),
			})
		}
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
		return fmt.Errorf("mkdir %s returned error: %w", localCacheStorageDir, err)
	}

	data, err := resource.Load()
	if err != nil {
		return fmt.Errorf("load resource %v returned error: %w", resource, err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data %v returned error: %w", data, err)
	}

	err = os.WriteFile(localCacheStorageDir+resource.LocalFileName(), jsonData, 0644)
	if err != nil {
		return fmt.Errorf("write file to %s with data %v returned error: %w", localCacheStorageDir+resource.LocalFileName(), jsonData, err)
	}
	return nil
}

func getFromLocalStorage[T any](filePath string) (T, error) {
	var result T

	data, err := os.ReadFile(filePath)
	if err != nil {
		return result, fmt.Errorf("read file %s: %w", filePath, err)
	}

	if err = json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("unmarshal cache file %s: %w", filePath, err)
	}

	return result, nil
}
