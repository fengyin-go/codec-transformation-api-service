package model

import (
	"strings"
	"time"
)

const (
	AlgorithmBase64     = "base64"
	AlgorithmURL        = "url"
	AlgorithmHTML       = "html"
	AlgorithmUnicode    = "unicode"
	AlgorithmHex        = "hex"
	AlgorithmBase62     = "base62"
	AlgorithmMorse      = "morse"
	AlgorithmMD5        = "md5"
	AlgorithmSHA1       = "sha1"
	AlgorithmSHA256     = "sha256"
	AlgorithmHMACSHA256 = "hmac_sha256"
	AlgorithmJWT        = "jwt"
)

const (
	CategoryEncoding = "encoding"
	CategoryHash     = "hash"
	CategoryEscape   = "escape"
)

const (
	ModeEncode = "encode"
	ModeDecode = "decode"
)

var ValidAlgorithms = map[string]bool{
	AlgorithmBase64:     true,
	AlgorithmURL:        true,
	AlgorithmHTML:       true,
	AlgorithmUnicode:    true,
	AlgorithmHex:        true,
	AlgorithmBase62:     true,
	AlgorithmMorse:      true,
	AlgorithmMD5:        true,
	AlgorithmSHA1:       true,
	AlgorithmSHA256:     true,
	AlgorithmHMACSHA256: true,
	AlgorithmJWT:        true,
}

var ValidCategories = map[string]bool{
	CategoryEncoding: true,
	CategoryHash:     true,
	CategoryEscape:   true,
}

var AlgorithmReversible = map[string]bool{
	AlgorithmBase64:     true,
	AlgorithmURL:        true,
	AlgorithmHTML:       true,
	AlgorithmUnicode:    true,
	AlgorithmHex:        true,
	AlgorithmBase62:     true,
	AlgorithmMorse:      true,
	AlgorithmMD5:        false,
	AlgorithmSHA1:       false,
	AlgorithmSHA256:     false,
	AlgorithmHMACSHA256: false,
	AlgorithmJWT:        false,
}

var AlgorithmCategoryMap = map[string]string{
	AlgorithmBase64:     CategoryEncoding,
	AlgorithmURL:        CategoryEscape,
	AlgorithmHTML:       CategoryEscape,
	AlgorithmUnicode:    CategoryEscape,
	AlgorithmHex:        CategoryEncoding,
	AlgorithmBase62:     CategoryEncoding,
	AlgorithmMorse:      CategoryEncoding,
	AlgorithmMD5:        CategoryHash,
	AlgorithmSHA1:       CategoryHash,
	AlgorithmSHA256:     CategoryHash,
	AlgorithmHMACSHA256: CategoryHash,
	AlgorithmJWT:        CategoryHash,
}

type Algorithm struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Reversible  bool      `json:"reversible"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (a *Algorithm) Validate() error {
	a.Name = strings.TrimSpace(a.Name)
	a.Category = strings.TrimSpace(a.Category)
	a.Description = strings.TrimSpace(a.Description)
	if a.Name == "" {
		return NewValidationError("name", "算法名称不能为空")
	}
	if !ValidAlgorithms[a.Name] {
		return NewValidationError("name", "算法名称不合法")
	}
	if a.Category == "" {
		a.Category = AlgorithmCategoryMap[a.Name]
	}
	if !ValidCategories[a.Category] {
		return NewValidationError("category", "分类不合法")
	}
	return nil
}

type AlgorithmFilter struct {
	Category   string
	Reversible *bool
	Keyword    string
}

func (f AlgorithmFilter) Match(a *Algorithm) bool {
	if f.Category != "" && a.Category != f.Category {
		return false
	}
	if f.Reversible != nil && a.Reversible != *f.Reversible {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(a.Name), k) &&
			!strings.Contains(strings.ToLower(a.Description), k) {
			return false
		}
	}
	return true
}
