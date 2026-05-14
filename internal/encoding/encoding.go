package encoding

import (
	"github.com/cloudwego/base64x"
)

// EncodeBase64 codifica dados em base64 usando a implementação otimizada do cloudwego
func EncodeBase64(data []byte) string {
	return base64x.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodifica uma string base64 usando a implementação otimizada do cloudwego
func DecodeBase64(encoded string) ([]byte, error) {
	return base64x.StdEncoding.DecodeString(encoded)
}

// EncodeBase64URLSafe codifica dados em base64 URL-safe
func EncodeBase64URLSafe(data []byte) string {
	return base64x.URLEncoding.EncodeToString(data)
}

// DecodeBase64URLSafe decodifica uma string base64 URL-safe
func DecodeBase64URLSafe(encoded string) ([]byte, error) {
	return base64x.URLEncoding.DecodeString(encoded)
}
