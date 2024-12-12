package codenc

import (
	"crypto/md5"
	"encoding/hex"
)

type MD5CodeEncoder struct{}

func (MD5CodeEncoder) Encode(code string) string {
	hash := md5.Sum([]byte(code))
	return hex.EncodeToString(hash[:])
}

