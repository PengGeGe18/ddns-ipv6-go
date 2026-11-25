//go:build !windows

package cert

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

// applyOwner 以 root 运行时把证书授权给指定用户, 便于非 root 的 Nginx 读取私钥
func applyOwner(path, owner string) error {
	owner = strings.TrimSpace(owner)
	if owner == "" || os.Geteuid() != 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	out, err := runCommand(ctx, "chown", owner, path)
	if err != nil {
		return fmt.Errorf("chown %s 失败: %s: %s", owner, err, truncate(out, 200))
	}
	return nil
}
