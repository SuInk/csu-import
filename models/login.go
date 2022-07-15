package models

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"csu-import/utils"
	"encoding/base64"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const JwcUnifiedUrl = "https://ca.csu.edu.cn/authserver/login?service=http%3A%2F%2Fcsujwc.its.csu.edu.cn%2Fsso.jsp"

type User struct {
	Id, Pwd string
}

// GetRandomString 随机字符串
func GetRandomString(n int) []byte {
	str := "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < n; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return result
}

// Padding 对明文进行填充
func Padding(plainText []byte, blockSize int) []byte {
	//计算要填充的长度
	n := blockSize - len(plainText)%blockSize
	//对原来的明文填充n个n
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

// AesCbcEncrypt AES加密（CBC模式）
func AesCbcEncrypt(plainText []byte, key []byte) string {
	//指定加密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//进行填充
	plainText = append(GetRandomString(64), plainText...)
	plainText = Padding(plainText, block.BlockSize())
	//指定初始向量vi,长度和block的块尺寸一致
	iv := GetRandomString(16)
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//加密连续数据库
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	//返回密文
	return base64.StdEncoding.EncodeToString(cipherText)
}

// Login 教务系统登录
func Login(user *User) (http.Client, error) {
	//获取cookie
	var client http.Client
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client.Jar = jar

	req, _ := http.NewRequest("GET", JwcUnifiedUrl, nil)
	response, err := client.Do(req)
	if err != nil {
		return client, err
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	encodePwd := AesCbcEncrypt([]byte(user.Pwd), []byte(doc.Find("#pwdEncryptSalt").AttrOr("value", "")))
	reqData := url.Values{
		"username":   {user.Id},
		"password":   {encodePwd},
		"captcha":    {"None"},
		"rememberMe": {"True"},
		"_eventId":   {"submit"},
		"cllt":       {"userNameLogin"},
		"dllt":       {"generalLogin"},
		"lt":         {"None"},
		"execution":  {doc.Find("#execution").AttrOr("value", "")},
	}
	if err != nil {
		return client, utils.ErrorServer
	}
	response1, err := client.Post(JwcUnifiedUrl, "application/x-www-form-urlencoded", strings.NewReader(reqData.Encode()))

	if err != nil {
		return client, utils.ErrorServer
	}
	body, _ := ioutil.ReadAll(response1.Body)
	defer response.Body.Close()
	//登陆成功
	if strings.Contains(string(body), "我的桌面") {
		return client, nil
	}
	//账号或密码错误
	return client, utils.ErrorIdPwd
}
