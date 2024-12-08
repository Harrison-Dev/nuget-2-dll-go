package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenGUID 產生簡單的假GUID (32 hex chars)
func GenGUID() string {
	const hexChars = "0123456789abcdef"
	b := make([]byte, 32)
	for i := 0; i < 32; i++ {
		b[i] = hexChars[rand.Intn(len(hexChars))]
	}
	return string(b)
}

// GenerateMeta 為一個資產產生簡單的meta檔案內容
func GenerateMeta(guid string) []byte {
	return []byte(fmt.Sprintf(`fileFormatVersion: 2
guid: %s
timeCreated: %d
licenseType: Free
DefaultImporter:
  externalObjects: {}
  userData: 
  assetBundleName: 
  assetBundleVariant: 
`, guid, time.Now().Unix()))
}
