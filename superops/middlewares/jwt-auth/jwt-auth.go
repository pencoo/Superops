package jwt_auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"superops/libs/config"
	"superops/libs/rcache"
	"superops/modules"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var ExpireTime int64
var TL *time.Location

func init() {
	TL, _ = time.LoadLocation(config.Timezone)
}

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Access-Token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    40000,
				"data":    nil,
				"message": "请求未携带token，无权限访问",
			})
			c.Abort()
			return
		}
		// parseToken 解析token包含的信息
		claims, err := ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusOK, gin.H{
					"code":    40000,
					"data":    nil,
					"message": "token已过期",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code":    40000,
				"data":    nil,
				"message": fmt.Sprint(err),
			})
			c.Abort()
			return
		}
		//用户提交的token为apitoken
		if claims.ExpiresAt == 0 {
			//当redis不可用时从mysql中查询用户的API Token是否合法
			var ut modules.Userapitoken
			ut.Username, ut.Token = claims.Name, token
			if !ut.UsersApiTokenIsActive() {
				c.JSON(http.StatusOK, gin.H{
					"code":    40000,
					"data":    nil,
					"message": fmt.Sprint("ApiToken无效"),
				})
				c.Abort()
				return
			}
			//当redis可用时使用下面的方法从redis中确认用户的API Token
			claims.Name = rcache.Userinfo_api_token + claims.Name
		} else {
			claims.Name = rcache.Userinfo_token + claims.Name
		}
		//redis可用时判断用户token是否合法
		v, err := rcache.RedisCacheStringGet(claims.Name)
		if v != token || err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    40000,
				"data":    nil,
				"message": fmt.Sprint("token无效"),
			})
			c.Abort()
			return
		}
		if strings.Contains(claims.Name, rcache.Userinfo_api_token) {
			s := strings.Split(claims.Name, rcache.Userinfo_api_token)
			claims.Name = s[1]
		}
		if strings.Contains(claims.Name, rcache.Userinfo_token) {
			s := strings.Split(claims.Name, rcache.Userinfo_token)
			claims.Name = s[1]
		}
		c.Set("Id", claims.Id)
		c.Set("Name", claims.Name)
	}
}

// 一些常量
var (
	TokenExpired     error = errors.New("Token is expired")
	TokenNotValidYet error = errors.New("Token not active yet")
	TokenMalformed   error = errors.New("That's not even a token")
	TokenInvalid     error = errors.New("Couldn't handle this token:")
)

// 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

// CreateToken 生成一个token
func CreateToken(claims CustomClaims) (string, error) {
	if ExpireTime == 0 || ExpireTime < time.Now().In(TL).Unix() {
		ExpireTime = WeekEnd()
	}
	claims.ExpiresAt = ExpireTime
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.TokenSecret))
}

func WeekEnd() int64 {
	//todayLast := time.Now().In(TL).Format("2006-01-02") + " 23:59:59"
	//todayLastTime, _ := time.ParseInLocation("2006-01-02 15:04:05", todayLast, TL)
	todayLast := time.Now().Format("2006-01-02") + " 23:59:59"
	todayLastTime, _ := time.ParseInLocation("2006-01-02 15:04:05", todayLast, time.Local)
	var d string
	if w := fmt.Sprint(time.Now().Weekday()); w == "Monday" {
		d = fmt.Sprint(6*24) + "h"
	} else if w == "Tuesday" {
		d = fmt.Sprint(5*24) + "h"
	} else if w == "Wednesday" {
		d = fmt.Sprint(4*24) + "h"
	} else if w == "Thursday" {
		d = fmt.Sprint(3*24) + "h"
	} else if w == "Friday" {
		d = fmt.Sprint(2*24) + "h"
	} else if w == "Saturday" {
		d = fmt.Sprint(1*24) + "h"
	} else {
		d = fmt.Sprint(7*24) + "h"
	}
	h, _ := time.ParseDuration(d)
	return todayLastTime.Add(h).Unix()
}

func CreateTokenForFuture(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.TokenSecret))
}

// 解析Tokne
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.TokenSecret), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// 更新token
func RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.TokenSecret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return CreateToken(*claims)
	}
	return "", TokenInvalid
}
