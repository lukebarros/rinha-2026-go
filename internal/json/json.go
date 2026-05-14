package json

import (
	"github.com/bytedance/sonic"
)

// Unmarshal desserializa JSON usando sonic
func Unmarshal(data []byte, v any) error {
	return sonic.Unmarshal(data, v)
}

// Marshal serializa v para JSON usando sonic
func Marshal(v any) ([]byte, error) {
	return sonic.Marshal(v)
}

// MarshalToString serializa v para string JSON usando sonic
func MarshalToString(v any) (string, error) {
	return sonic.MarshalString(v)
}

// UnmarshalFromString desserializa JSON string usando sonic
func UnmarshalFromString(data string, v any) error {
	return sonic.UnmarshalString(data, v)
}
