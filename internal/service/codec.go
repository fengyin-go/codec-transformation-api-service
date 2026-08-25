package service

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"codec/internal/model"
	"codec/pkg/idgen"
)

var morseEncodeMap = map[rune]string{
	'A': ".-", 'B': "-...", 'C': "-.-.", 'D': "-..", 'E': ".",
	'F': "..-.", 'G': "--.", 'H': "....", 'I': "..", 'J': ".---",
	'K': "-.-", 'L': ".-..", 'M': "--", 'N': "-.", 'O': "---",
	'P': ".--.", 'Q': "--.-", 'R': ".-.", 'S': "...", 'T': "-",
	'U': "..-", 'V': "...-", 'W': ".--", 'X': "-..-", 'Y': "-.--",
	'Z': "--..",
	'1': ".----", '2': "..---", '3': "...--", '4': "....-", '5': ".....",
	'6': "-....", '7': "--...", '8': "---..", '9': "----.", '0': "-----",
	'.': ".-.-.-", ',': "--..--", '?': "..--..", '!': "-.-.--", '/': "-..-.",
	'@': ".--.-.", '(': "-.--.", ')': "-.--.-", '&': ".-...", ':': "---...",
	';': "-.-.-.", '=': "-...-", '+': ".-.-.", '-': "-....-", '_': "..--.-",
	'"': ".-..-.", '$': "...-..-", ' ': "/",
}

var morseDecodeMap func() map[string]rune

func init() {
	morseDecodeMap = func() map[string]rune {
		m := make(map[string]rune, len(morseEncodeMap))
		for r, s := range morseEncodeMap {
			m[s] = r
		}
		return m
	}
}

// base62Encode encodes bytes to base62 string.
func base62Encode(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	// treat as big integer base 256 -> base 62
	var n uint64
	for _, b := range data {
		n = n*256 + uint64(b)
	}
	if n == 0 {
		return "0"
	}
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var buf []byte
	for n > 0 {
		buf = append(buf, chars[n%62])
		n /= 62
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

// base62Decode decodes base62 string to bytes (best effort, may overflow for large inputs).
func base62Decode(s string) ([]byte, error) {
	if s == "" {
		return []byte{}, nil
	}
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	charToVal := make(map[byte]uint64, len(chars))
	for i := 0; i < len(chars); i++ {
		charToVal[chars[i]] = uint64(i)
	}
	var n uint64
	for i := 0; i < len(s); i++ {
		v, ok := charToVal[s[i]]
		if !ok {
			return nil, model.NewValidationError("input", "base62 包含非法字符")
		}
		n = n*62 + v
	}
	if n == 0 {
		return []byte{0}, nil
	}
	var buf []byte
	for n > 0 {
		buf = append(buf, byte(n%256))
		n /= 256
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return buf, nil
}

// Encode 执行编码/转换，返回输出结果。key 仅用于 HMAC/JWT。
func (s *Service) Encode(algorithm, input, key string) (string, error) {
	if input == "" {
		return "", model.NewValidationError("input", "输入不能为空")
	}
	alg := strings.TrimSpace(algorithm)
	if !model.ValidAlgorithms[alg] {
		return "", model.NewValidationError("algorithm", "不支持的算法")
	}

	var output string
	var err error

	switch alg {
	case model.AlgorithmBase64:
		output = base64.StdEncoding.EncodeToString([]byte(input))
	case model.AlgorithmURL:
		output = url.QueryEscape(input)
	case model.AlgorithmHTML:
		output = html.EscapeString(input)
	case model.AlgorithmUnicode:
		output = unicodeEscape(input)
	case model.AlgorithmHex:
		output = hex.EncodeToString([]byte(input))
	case model.AlgorithmBase62:
		output = base62Encode([]byte(input))
	case model.AlgorithmMorse:
		output = morseEncode(input)
	case model.AlgorithmMD5:
		sum := md5.Sum([]byte(input))
		output = hex.EncodeToString(sum[:])
	case model.AlgorithmSHA1:
		sum := sha1.Sum([]byte(input))
		output = hex.EncodeToString(sum[:])
	case model.AlgorithmSHA256:
		sum := sha256.Sum256([]byte(input))
		output = hex.EncodeToString(sum[:])
	case model.AlgorithmHMACSHA256:
		if key == "" {
			return "", model.NewValidationError("key", "HMAC-SHA256 需要 key")
		}
		mac := hmac.New(sha256.New, []byte(key))
		mac.Write([]byte(input))
		output = hex.EncodeToString(mac.Sum(nil))
	case model.AlgorithmJWT:
		if key == "" {
			return "", model.NewValidationError("key", "JWT 需要 secret key")
		}
		output, err = jwtEncode(input, key)
	default:
		return "", model.NewValidationError("algorithm", "未实现的算法")
	}

	if err != nil {
		return "", err
	}

	_ = s.recordOperation(alg, model.ModeEncode, input, output)
	return output, nil
}

// Decode 执行解码/反转。hash 类不可逆返回 ValidationError。
func (s *Service) Decode(algorithm, input, key string) (string, error) {
	if input == "" {
		return "", model.NewValidationError("input", "输入不能为空")
	}
	alg := strings.TrimSpace(algorithm)
	if !model.ValidAlgorithms[alg] {
		return "", model.NewValidationError("algorithm", "不支持的算法")
	}

	var output string
	var err error

	switch alg {
	case model.AlgorithmBase64:
		b, e := base64.StdEncoding.DecodeString(input)
		if e != nil {
			return "", model.NewValidationError("input", "Base64 解码失败: "+e.Error())
		}
		output = string(b)
	case model.AlgorithmURL:
		output, err = url.QueryUnescape(input)
		if err != nil {
			return "", model.NewValidationError("input", "URL 解码失败: "+err.Error())
		}
	case model.AlgorithmHTML:
		output = html.UnescapeString(input)
	case model.AlgorithmUnicode:
		output, err = unicodeUnescape(input)
		if err != nil {
			return "", err
		}
	case model.AlgorithmHex:
		b, e := hex.DecodeString(input)
		if e != nil {
			return "", model.NewValidationError("input", "Hex 解码失败: "+e.Error())
		}
		output = string(b)
	case model.AlgorithmBase62:
		b, e := base62Decode(input)
		if e != nil {
			return "", e
		}
		output = string(b)
	case model.AlgorithmMorse:
		output, err = morseDecode(input)
		if err != nil {
			return "", err
		}
	case model.AlgorithmMD5, model.AlgorithmSHA1, model.AlgorithmSHA256, model.AlgorithmHMACSHA256:
		return "", model.NewValidationError("algorithm", "哈希算法不可逆")
	case model.AlgorithmJWT:
		if key == "" {
			return "", model.NewValidationError("key", "JWT 校验需要 secret key")
		}
		output, err = jwtDecode(input, key)
		if err != nil {
			return "", err
		}
	default:
		return "", model.NewValidationError("algorithm", "未实现的算法")
	}

	_ = s.recordOperation(alg, model.ModeDecode, input, output)
	return output, nil
}

// BatchConvert 批量转换：对多个输入执行同一算法。
func (s *Service) BatchConvert(algorithm, mode string, inputs []string, key string) ([]string, error) {
	if len(inputs) == 0 {
		return nil, model.NewValidationError("inputs", "输入列表不能为空")
	}
	mode = strings.TrimSpace(mode)
	if mode != model.ModeEncode && mode != model.ModeDecode {
		return nil, model.NewValidationError("mode", "模式必须为 encode 或 decode")
	}
	results := make([]string, 0, len(inputs))
	for _, in := range inputs {
		var out string
		var err error
		if mode == model.ModeEncode {
			out, err = s.Encode(algorithm, in, key)
		} else {
			out, err = s.Decode(algorithm, in, key)
		}
		if err != nil {
			return nil, err
		}
		results = append(results, out)
	}
	return results, nil
}

func (s *Service) recordOperation(alg, mode, input, output string) error {
	inHash := sha256.Sum256([]byte(input))
	inHashStr := hex.EncodeToString(inHash[:])[:16]
	preview := output
	if len(preview) > 64 {
		preview = preview[:64] + "..."
	}
	op := &model.Operation{
		ID:            idgen.Hex(),
		Algorithm:     alg,
		Mode:          mode,
		InputHash:     inHashStr,
		OutputPreview: preview,
		CreatedAt:     time.Now(),
	}
	return s.store.CreateOperation(op)
}

func unicodeEscape(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r < 0x80 && r != '\\' {
			b.WriteRune(r)
		} else {
			b.WriteString(fmt.Sprintf("\\u%04x", r))
		}
	}
	return b.String()
}

func unicodeUnescape(s string) (string, error) {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) && s[i+1] == 'u' && i+5 < len(s) {
			hexPart := s[i+2 : i+6]
			r, err := strconv.ParseInt(hexPart, 16, 32)
			if err != nil {
				return "", model.NewValidationError("input", "Unicode 转义格式非法")
			}
			b.WriteRune(rune(r))
			i += 5
		} else {
			b.WriteByte(s[i])
		}
	}
	return b.String(), nil
}

func morseEncode(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 {
			b.WriteByte(' ')
		}
		v, ok := morseEncodeMap[unicode.ToUpper(r)]
		if ok {
			b.WriteString(v)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func morseDecode(s string) (string, error) {
	m := morseDecodeMap()
	parts := strings.Split(s, " ")
	var b strings.Builder
	for _, p := range parts {
		if p == "/" {
			b.WriteByte(' ')
			continue
		}
		r, ok := m[p]
		if ok {
			b.WriteRune(r)
		} else {
			return "", model.NewValidationError("input", "Morse 码包含无法识别的序列")
		}
	}
	return b.String(), nil
}

// jwtEncode 生成 JWT HS256 token（header.payload.signature）。
func jwtEncode(payloadJSON, secret string) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerBytes, _ := json.Marshal(header)
	if payloadJSON == "" {
		payloadJSON = "{}"
	}
	b64Header := base64.RawURLEncoding.EncodeToString(headerBytes)
	b64Payload := base64.RawURLEncoding.EncodeToString([]byte(payloadJSON))
	signingInput := b64Header + "." + b64Payload
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return signingInput + "." + sig, nil
}

// jwtDecode 解析并校验 JWT HS256 token，返回 payload 字符串。
func jwtDecode(token, secret string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", model.NewValidationError("input", "JWT 格式非法")
	}
	signingInput := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expectedSig), []byte(parts[2])) {
		return "", model.NewValidationError("input", "JWT 签名验证失败")
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", model.NewValidationError("input", "JWT payload 解码失败")
	}
	return string(payloadBytes), nil
}

// CodecResult 包含单次转换结果及可选错误。
type CodecResult struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

// SafeEncode 与 Encode 相同，但返回 error 时不 panic，封装成 CodecResult。
func (s *Service) SafeEncode(algorithm, input, key string) CodecResult {
	out, err := s.Encode(algorithm, input, key)
	if err != nil {
		return CodecResult{Input: input, Error: err.Error()}
	}
	return CodecResult{Input: input, Output: out}
}

// SafeDecode 与 Decode 相同，但返回 error 时不 panic，封装成 CodecResult。
func (s *Service) SafeDecode(algorithm, input, key string) CodecResult {
	out, err := s.Decode(algorithm, input, key)
	if err != nil {
		return CodecResult{Input: input, Error: err.Error()}
	}
	return CodecResult{Input: input, Output: out}
}

// BatchConvertResult 批量转换，每个输入独立处理并收集结果与错误。
func (s *Service) BatchConvertResult(algorithm, mode string, inputs []string, key string) ([]CodecResult, error) {
	if len(inputs) == 0 {
		return nil, model.NewValidationError("inputs", "输入列表不能为空")
	}
	mode = strings.TrimSpace(mode)
	if mode != model.ModeEncode && mode != model.ModeDecode {
		return nil, model.NewValidationError("mode", "模式必须为 encode 或 decode")
	}
	results := make([]CodecResult, len(inputs))
	for i, in := range inputs {
		if mode == model.ModeEncode {
			results[i] = s.SafeEncode(algorithm, in, key)
		} else {
			results[i] = s.SafeDecode(algorithm, in, key)
		}
	}
	return results, nil
}

// CompareHash 比对两个字符串的 MD5/SHA1/SHA256/HMAC 值是否相等（用于校验）。
func (s *Service) CompareHash(algorithm, a, b, key string) (bool, error) {
	var ha, hb string
	var err error
	switch algorithm {
	case model.AlgorithmMD5:
		ha = fmt.Sprintf("%x", md5.Sum([]byte(a)))
		hb = fmt.Sprintf("%x", md5.Sum([]byte(b)))
	case model.AlgorithmSHA1:
		ha = fmt.Sprintf("%x", sha1.Sum([]byte(a)))
		hb = fmt.Sprintf("%x", sha1.Sum([]byte(b)))
	case model.AlgorithmSHA256:
		ha = fmt.Sprintf("%x", sha256.Sum256([]byte(a)))
		hb = fmt.Sprintf("%x", sha256.Sum256([]byte(b)))
	case model.AlgorithmHMACSHA256:
		if key == "" {
			return false, model.NewValidationError("key", "HMAC-SHA256 需要 key")
		}
		ma := hmac.New(sha256.New, []byte(key))
		ma.Write([]byte(a))
		ha = hex.EncodeToString(ma.Sum(nil))
		mb := hmac.New(sha256.New, []byte(key))
		mb.Write([]byte(b))
		hb = hex.EncodeToString(mb.Sum(nil))
	default:
		return false, model.NewValidationError("algorithm", "仅支持 md5/sha1/sha256/hmac_sha256 比对")
	}
	return ha == hb, err
}

// GenerateRandomBytes 生成指定长度的随机字节并以 hex 返回（工具占位，未实际实现随机源）。
func (s *Service) GenerateRandomBytes(n int) (string, error) {
	if n <= 0 || n > 4096 {
		return "", model.NewValidationError("n", "长度必须在 1-4096 之间")
	}
	return "", errors.New("未实现")
}

// IsReversible 判断算法是否可逆。
func (s *Service) IsReversible(algorithm string) bool {
	return model.AlgorithmReversible[strings.TrimSpace(algorithm)]
}

// SupportedAlgorithms 返回所有支持的算法名称列表（排序后）。
func (s *Service) SupportedAlgorithms() []string {
	list := make([]string, 0, len(model.ValidAlgorithms))
	for k := range model.ValidAlgorithms {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

// stringRunes 辅助：统计字符串的 rune 数量。
func stringRunes(s string) int {
	return utf8.RuneCountInString(s)
}
