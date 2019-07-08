package cryptology

import (
	"crypto/md5"
	"fmt"
	"io"
)

/**
 * 使用Sum进行md5加密
 */
func Md5Sum(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

/**
 * 使用New进行md5加密
 */
func Md5New(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	return fmt.Sprintf("%x", w.Sum(nil))
}
