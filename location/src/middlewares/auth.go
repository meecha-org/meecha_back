package middlewares

import (
	"crypto/ed25519"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"location/utils"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// PEMファイルから公開鍵を読み込む関数
func loadPublicKey(pubKeyFile string) (ed25519.PublicKey, error) {
	// 公開鍵の読み込み
	pubKeyPEM, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("公開鍵ファイルの読み込みに失敗しました: %w", err)
	}

	pubKeyBlock, _ := pem.Decode(pubKeyPEM)
	if pubKeyBlock == nil || pubKeyBlock.Type != "PUBLIC KEY" {
		return nil, errors.New("無効な公開鍵形式")
	}

	pubKey := pubKeyBlock.Bytes
	return pubKey, nil
}

var (
	pubKey = ed25519.PublicKey{}
)

func init() {
	// 鍵を読み込む
	loadeed,err := loadPublicKey("./authPublic/public_key.pem")

	// エラー処理
	if err != nil {
		log.Println("failed to load key")
		log.Fatalln(err)
	}
	
	// グローバル変数に割り当てる
	pubKey = loadeed
}

type JwtPayload struct {
	UserID string	//ユーザーID
	Labels []string	//ラベルを入れ宇
}

func ValidToken(tokenString string) (JwtPayload,error) {
	// トークンのパース
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return pubKey, nil
	})

	// エラー処理
	if err != nil {
		return JwtPayload{},err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// 成功したとき
		return JwtPayload{
			UserID: claims["UserID"].(string),
			Labels: decodeLabel(claims["Labels"].([]interface{})),
		},nil
	} else {
		utils.Println(err)
	}

	return JwtPayload{},err
}

func decodeLabel(labels []interface{}) ([]string) {
	rlabels := []string{}

	for _, val := range labels {
		rlabels = append(rlabels, val.(string))
	}

	return rlabels
}

func PocketAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Authorize Header 取得
			token := ctx.Request().Header.Get("Authorization")

			// 文字列検証
			if !strings.HasPrefix(token,"Bearer ") {
				// 始まっていない場合
				return ctx.NoContent(http.StatusUnauthorized)
			}

			// ７文字目以降を取り出す			
			splitToken := token[7:]

			// トークンを検証する
			parseData,err := ValidToken(splitToken)

			// エラー
			if err != nil {
				utils.Println(err)
				return ctx.NoContent(http.StatusUnauthorized)
			}

			// データ設定
			ctx.Set("UserID",parseData.UserID)
			ctx.Set("Labels",parseData.Labels)

			return next(ctx)
		}
	}
}