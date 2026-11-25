//go:build windows

package cert

// applyOwner Windows 不支持文件归属调整, 直接跳过
func applyOwner(path, owner string) error {
	return nil
}
