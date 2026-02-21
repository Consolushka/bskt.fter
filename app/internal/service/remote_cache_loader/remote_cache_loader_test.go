package remote_cache_loader

import (
	"IMP/app/internal/adapters/cached_remote_resource"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func withTempCwd(t *testing.T) {
	t.Helper()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("failed to chdir temp dir: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(prev)
	})
}

func TestLoadLocalFile_DownloadAndReuseCache(t *testing.T) {
	withTempCwd(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resource := cached_remote_resource.NewMockCachedRemoteResource(ctrl)
	resource.EXPECT().LocalFileName().Return("resource.json").AnyTimes()
	resource.EXPECT().GetLifeTime().Return(time.Hour).AnyTimes()
	resource.EXPECT().Load().Return([]int{1, 2, 3}, nil).Times(1)

	got1, err := LoadLocalFile[[]int](resource)
	if err != nil {
		t.Fatalf("unexpected error on first load: %v", err)
	}
	if len(got1) != 3 {
		t.Fatalf("unexpected first result length: %d", len(got1))
	}

	got2, err := LoadLocalFile[[]int](resource)
	if err != nil {
		t.Fatalf("unexpected error on second load: %v", err)
	}
	if len(got2) != 3 {
		t.Fatalf("unexpected second result length: %d", len(got2))
	}
}

func TestLoadLocalFile_ReloadsWhenExpired(t *testing.T) {
	withTempCwd(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resource := cached_remote_resource.NewMockCachedRemoteResource(ctrl)
	resource.EXPECT().LocalFileName().Return("expired.json").AnyTimes()
	resource.EXPECT().GetLifeTime().Return(-time.Second).AnyTimes()
	resource.EXPECT().Load().Return(map[string]int{"v": 1}, nil).Times(2)

	_, err := LoadLocalFile[map[string]int](resource)
	if err != nil {
		t.Fatalf("unexpected error on first load: %v", err)
	}
	_, err = LoadLocalFile[map[string]int](resource)
	if err != nil {
		t.Fatalf("unexpected error on second load: %v", err)
	}
}

func TestLoadLocalFile_ReturnsErrorWhenRemoteLoadFails(t *testing.T) {
	withTempCwd(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resource := cached_remote_resource.NewMockCachedRemoteResource(ctrl)
	resource.EXPECT().LocalFileName().Return("failed.json").AnyTimes()
	resource.EXPECT().GetLifeTime().Return(time.Hour).AnyTimes()
	resource.EXPECT().Load().Return(nil, errors.New("remote error")).Times(1)

	_, err := LoadLocalFile[map[string]int](resource)
	if err == nil {
		t.Fatal("expected error when remote load fails and cache file is absent")
	}
}

func TestGetFromLocalStorage_ReturnsErrorOnCorruptedCache(t *testing.T) {
	withTempCwd(t)

	cacheFile := filepath.Join(localCacheStorageDir, "broken.json")
	if err := os.MkdirAll(localCacheStorageDir, 0755); err != nil {
		t.Fatalf("failed to create cache dir: %v", err)
	}
	if err := os.WriteFile(cacheFile, []byte("not-json"), 0644); err != nil {
		t.Fatalf("failed to create corrupted cache file: %v", err)
	}

	_, err := getFromLocalStorage[map[string]int](cacheFile)
	if err == nil {
		t.Fatal("expected unmarshal error for corrupted cache file")
	}
}

func TestSaveToLocalCacheStorage_ReturnsMarshalError(t *testing.T) {
	withTempCwd(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resource := cached_remote_resource.NewMockCachedRemoteResource(ctrl)
	resource.EXPECT().Load().Return(make(chan int), nil).Times(1)

	err := saveToLocalCacheStorage(resource)
	if err == nil {
		t.Fatal("expected marshal error for unsupported json type")
	}
}

func TestSaveToLocalCacheStorage_ReturnsMkdirAllError(t *testing.T) {
	withTempCwd(t)

	// Make "tmp" a file to force MkdirAll("tmp/cache") to fail with ENOTDIR.
	if err := os.WriteFile("tmp", []byte("blocking file"), 0644); err != nil {
		t.Fatalf("failed to create blocking file: %v", err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resource := cached_remote_resource.NewMockCachedRemoteResource(ctrl)

	err := saveToLocalCacheStorage(resource)
	if err == nil {
		t.Fatal("expected mkdir error when cache parent path is a file")
	}
}

func TestIsNeedToBeDownload(t *testing.T) {
	withTempCwd(t)

	if !isNeedToBeDownload("missing.json", time.Hour) {
		t.Fatal("expected missing file to require download")
	}

	cacheFile := filepath.Join(localCacheStorageDir, "fresh.json")
	if err := os.MkdirAll(localCacheStorageDir, 0755); err != nil {
		t.Fatalf("failed to create cache dir: %v", err)
	}
	if err := os.WriteFile(cacheFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to write fresh file: %v", err)
	}

	if isNeedToBeDownload(cacheFile, time.Hour) {
		t.Fatal("expected fresh file to skip download")
	}
}
