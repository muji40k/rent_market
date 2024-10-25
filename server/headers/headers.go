package headers

import "github.com/gin-contrib/cors"

const API_KEY string = "X-API-KEY"
const API_RENEW string = "X-API-RENEW"

func CorsFiller(config *cors.Config) {
	config.AddAllowHeaders(API_KEY, API_RENEW)
}

