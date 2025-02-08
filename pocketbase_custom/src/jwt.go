package main

import (
	"crypto/rand"
	// "crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/ed25519"
)

const (
	TokenExpiry = 900	//900秒 (15分)
)

var privKey ed25519.PrivateKey
var pubKey ed25519.PublicKey

type JwtPayload struct {
	UserID string	//ユーザーID
	Labels []string	//ラベルを入れ宇
}

func NowTime() int64 {
	return time.Now().Unix()
}

func init() {
	const privKeyFile = "./pb_data/private_key.pem"
	const pubKeyFile = "./pb_data/public_key.pem"

	// 鍵ペアを生成するか、既存の鍵を読み込む

	if _, err := os.Stat(privKeyFile); os.IsNotExist(err) {
		// 鍵が存在しない場合は生成
		var err error
		privKey, pubKey, err = generateAndSaveKeys(privKeyFile, pubKeyFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("鍵が正常に生成され、PEMファイルに保存されました。")
	} else {
		// 鍵が存在する場合は読み込む
		var err error
		privKey, pubKey, err = loadKeys(privKeyFile, pubKeyFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("既存の鍵を読み込みました。")
	}
}


func GenJwt(args JwtPayload) (string,error) {
	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, jwt.MapClaims{
		"UserID" : args.UserID,
		"Labels" : args.Labels,
		"exp" : NowTime() + TokenExpiry,	//有効期限設定
	})
	
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(privKey)
	
	// エラー処理
	if err != nil {
		return "",err
	}

	return tokenString,nil
}

// 鍵を生成し、PEMファイルに書き出す関数
func generateAndSaveKeys(privKeyFile, pubKeyFile string) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("鍵の生成に失敗しました: %w", err)
	}

	// 秘密鍵をPEM形式で書き出す
	privKeyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privKey.Seed(), // 秘密鍵のシードを使う
	}
	if err := writePEMFile(privKeyFile, privKeyBlock); err != nil {
		return nil, nil, err
	}

	// 公開鍵をPEM形式で書き出す
	pubKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKey,
	}
	if err := writePEMFile(pubKeyFile, pubKeyBlock); err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}

// PEMファイルに書き出す関数
func writePEMFile(filename string, block *pem.Block) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ファイルの作成に失敗しました: %w", err)
	}
	defer file.Close()

	if err := pem.Encode(file, block); err != nil {
		return fmt.Errorf("PEMエンコードに失敗しました: %w", err)
	}

	return nil
}

// PEMファイルから鍵を読み込む関数
func loadKeys(privKeyFile, pubKeyFile string) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	privKeyPEM, err := ioutil.ReadFile(privKeyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("秘密鍵ファイルの読み込みに失敗しました: %w", err)
	}

	privKeyBlock, _ := pem.Decode(privKeyPEM)
	if privKeyBlock == nil || privKeyBlock.Type != "PRIVATE KEY" {
		return nil, nil, errors.New("無効な秘密鍵形式")
	}

	pubKeyPEM, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("公開鍵ファイルの読み込みに失敗しました: %w", err)
	}

	pubKeyBlock, _ := pem.Decode(pubKeyPEM)
	if pubKeyBlock == nil || pubKeyBlock.Type != "PUBLIC KEY" {
		return nil, nil, errors.New("無効な公開鍵形式")
	}

	privKey := ed25519.NewKeyFromSeed(privKeyBlock.Bytes)
	pubKey := pubKeyBlock.Bytes

	return privKey, pubKey, nil
}