package md5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"rent_service/internal/domain/models"
	"time"
)

func Hash(user *models.User) string {
	base := fmt.Sprintf(
		"%v%v%v%v%v",
		user.Id, user.Name, user.Email, user.Password, time.Now(),
	)
	hash := md5.Sum([]byte(base))

	return hex.EncodeToString(hash[:])
}

