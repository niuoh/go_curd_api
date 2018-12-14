package  function

import (
  "bytes"
  "crypto/cipher"
  "crypto/aes"
  "encoding/base64"
  "strconv"
  "strings"
  . "../config"
  "time"

)

// sid|time|version|randstr
var rand_str_length=16

func CreateToken(id int) string{
  var str string
  token_time:=int(time.Now().Unix()+int64(Config["token_time"].(int)))
  str=strconv.Itoa(id)+"|"+strconv.Itoa(token_time)+"|"+Config["version"].(string)
  str+="|"+RandString(rand_str_length)
  token:=AesEncrypt(str)
  return token
}

func CheckToken(token string) int{
  token=AesDecrypt(token)
  if(token==""){
    return 0
  }
  info:=strings.Split(token,"|")

  if(len(info)!=4){
    return 0
  }
  token_time,_:=strconv.Atoi(info[1])

  if(token_time<int(time.Now().Unix())){
    return 0
  }
  if(info[2]!=Config["version"]){
    return 0
  }
  if(len(info[3])!=rand_str_length){
    return 0
  }

  sid,_:=strconv.Atoi(info[0])
  return sid
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext) % blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

func AesEncrypt(origData0 string) (string) {
    key := []byte("0000000000000001")
    block, _ := aes.NewCipher(key)
    origData:=[]byte(origData0)
    blockSize := block.BlockSize()
    origData = PKCS7Padding(origData, blockSize)
    blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
    crypted := make([]byte, len(origData))
    blockMode.CryptBlocks(crypted, origData)

    result:=base64.StdEncoding.EncodeToString(crypted)
    return result
}

func AesDecrypt(str string) (string) {

    str2,err:=base64.StdEncoding.DecodeString(str)

    if(err!=nil){
      return ""
    }
    str=string(str2)
    crypted:=[]byte(str)
    key := []byte("0000000000000001")
    block, err := aes.NewCipher(key)
    if err != nil {
        return ""
    }
    blockSize := block.BlockSize()
    blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
    origData := make([]byte, len(crypted))
    blockMode.CryptBlocks(origData, crypted)
    origData = PKCS7UnPadding(origData)
    return string(origData)
}
