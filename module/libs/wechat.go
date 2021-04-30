package libs

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"iris_project_foundation/common/api_error"
	"iris_project_foundation/config"
	"net/http"
	"net/url"
)

type SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int64  `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func WechatAuthorize(code string) (openID, secretKey string, err error) {
	var (
		baseUrl *url.URL
		params  *url.Values
		resp    *http.Response
		result  SessionResponse
	)

	if baseUrl, err = url.Parse(config.GConfig.Wechat.Code2SessionAPI); err != nil {
		return
	}

	params = &url.Values{}
	params.Set("appid", config.GConfig.Wechat.WechatAppID)
	params.Set("secret", config.GConfig.Wechat.WechatSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")

	baseUrl.RawQuery = params.Encode()
	urlPath := baseUrl.String()

	if resp, err = http.Get(urlPath); err != nil {
		return "", "", api_error.WechatAuthorizeError
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", &api_error.BaseAPIError{ErrorCode: 10010, Message: err.Error()}
	}

	return result.OpenID, result.SessionKey, nil
}

func WechatDecryptUserInfo(encrypted []byte, sessionKey, iv string) interface{} {
	decrypted := AESDecryptCBC(encrypted, []byte(sessionKey), []byte(iv))
	return decrypted
}

func AESDecryptCBC(encrypted []byte, key, iv []byte) (decrypted []byte) {
	var (
		block     cipher.Block
		blockMode cipher.BlockMode
	)
	block, _ = aes.NewCipher(key)
	blockMode = cipher.NewCBCDecrypter(block, iv)
	decrypted = make([]byte, len(encrypted))
	blockMode.CryptBlocks(decrypted, encrypted)
	return pkcs5UnPadding(decrypted)
}

func pkcs5UnPadding(origData []byte) []byte {
	var (
		length    int
		unpadding int
	)
	length = len(origData)
	unpadding = int(origData[length-1])
	return origData[:(length - unpadding)]
}
