package crypto

import (
	//"fmt"
	"testing"

	"github.com/bmizerany/assert"
)

var (
	aesKey         = "astaxie12798akljzmknm.ahkjkljl;k"
	input          = "test skcloud crypto"
	md5Output      = "37d5bf6f759c49527e78f8db7c051276"
	base64Output   = "dGVzdCBza2Nsb3VkIGNyeXB0bw=="
	aesOutput      = "XGFffDlOwueH4DCVD+5Um/S/OA=="
	sha1Output     = "82d2ba7d706239e6f3bfcf0b752cf173f91b5d9c"
	hmacKey        = "hash key"
	hmacSha1Output = "aa38b69e64725e3687c549120645fdbbf7281923"
)

func TestMD5(t *testing.T) {
	out := GetMD5(input)
	assert.Equal(t, md5Output, out)
	assert.Equal(t, true, CheckMD5(input, out))
	assert.Equal(t, false, CheckMD5(input, input))
}

func TestSha1(t *testing.T) {
	out := GetSha1(input)
	assert.Equal(t, sha1Output, out)

	out2 := GetHmacSha1(input, hmacKey)
	assert.Equal(t, hmacSha1Output, out2)
}

func TestBase64(t *testing.T) {
	out := Base64Encode(input)
	assert.Equal(t, base64Output, out)

	in, e := Base64Decode(base64Output)
	assert.Equal(t, nil, e)
	assert.Equal(t, input, in)
}

func TestGetRandomKey(t *testing.T) {
	out1 := GetRandomKey()
	out2 := GetRandomKey()

	assert.Equal(t, 40, len(out1))
	assert.Equal(t, 40, len(out2))
	assert.NotEqual(t, out1, out2)

}

func TestAES(t *testing.T) {
	out1, err1 := AESEncode(input, aesKey)
	assert.Equal(t, nil, err1)
	assert.Equal(t, aesOutput, out1)

	out2, err2 := AESDecode(aesOutput, aesKey)
	assert.Equal(t, nil, err2)
	assert.Equal(t, input, out2)
}
