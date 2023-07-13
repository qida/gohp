package pinyin

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	tones = [][]rune{
		{'ā', 'ē', 'ī', 'ō', 'ū', 'ǖ', 'Ā', 'Ē', 'Ī', 'Ō', 'Ū', 'Ǖ'},
		{'á', 'é', 'í', 'ó', 'ú', 'ǘ', 'Á', 'É', 'Í', 'Ó', 'Ú', 'Ǘ'},
		{'ǎ', 'ě', 'ǐ', 'ǒ', 'ǔ', 'ǚ', 'Ǎ', 'Ě', 'Ǐ', 'Ǒ', 'Ǔ', 'Ǚ'},
		{'à', 'è', 'ì', 'ò', 'ù', 'ǜ', 'À', 'È', 'Ì', 'Ò', 'Ù', 'Ǜ'},
	}
	neutrals = []rune{'a', 'e', 'i', 'o', 'u', 'v', 'A', 'E', 'I', 'O', 'U', 'V'}
)

var (
	// 从带声调的声母到对应的英文字符的映射
	tonesMap map[rune]rune

	// 从汉字到声调的映射
	numericTonesMap map[rune]int

	// 从汉字到拼音的映射（带声调）
	pinyinMap map[rune]string

	initialized bool
)

type Mode int

const (
	WithoutTone        Mode = iota + 1 // 默认模式，例如：guo
	Tone                               // 带声调的拼音 例如：guó
	InitialsInCapitals                 // 首字母大写不带声调，例如：Guo
	Initials                           // 首字母大写不带声调，例如：G
)

type pinyin struct {
	origin string
	split  string
	mode   Mode
}

func init() {
	tonesMap = make(map[rune]rune)
	numericTonesMap = make(map[rune]int)
	pinyinMap = make(map[rune]string)
	for i, runes := range tones {
		for j, tone := range runes {
			tonesMap[tone] = neutrals[j]
			numericTonesMap[tone] = i + 1
		}
	}

	f, err := getFileContent()
	if err != nil {
		initialized = false
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		strs := strings.Split(scanner.Text(), "=>")
		if len(strs) < 2 {
			continue
		}
		i, err := strconv.ParseInt(strs[0], 16, 32)
		if err != nil {
			continue
		}
		pinyinMap[rune(i)] = strs[1]
	}
	initialized = true
}

func getFileContent() (io.ReadCloser, error) {
	// resp, err := http.Get("https://raw.githubusercontent.com/chain-zhang/pinyin/master/pinyin.txt")
	resp, err := http.Get("http://golang.cdn.sunqida.cn/golang/qida/pinyin/pinyin.txt")
	if err != nil {
		file, err1 := os.Open("./pinyin.txt")
		return file, err1

	}
	return resp.Body, err
}

func New(origin string) *pinyin {
	return &pinyin{
		origin: strings.ToLower(origin),
		split:  " ",
		mode:   WithoutTone,
	}
}

func (py *pinyin) Split(split string) *pinyin {
	py.split = split
	return py
}

func (py *pinyin) Mode(mode Mode) *pinyin {
	py.mode = mode
	return py
}

func (py *pinyin) Convert() (string, error) {
	if !initialized {
		return "", ErrInitialize
	}

	sr := []rune(py.origin)
	words := make([]string, 0)
	for _, s := range sr {
		word, err := getPinyin(s, py.mode)
		if err != nil {
			return "", err
		}
		if len(word) > 0 {
			words = append(words, word)
		}
	}
	return strings.Join(words, py.split), nil
}

func getPinyin(hanzi rune, mode Mode) (string, error) {
	if !initialized {
		return "", ErrInitialize
	}
	switch mode {
	case Tone:
		return getTone(hanzi), nil
	case InitialsInCapitals:
		return getInitialsInCapitals(hanzi), nil
	case Initials:
		return getInitials(hanzi), nil
	default:
		return getDefault(hanzi), nil
	}
}

func getTone(hanzi rune) string {
	if pinyinMap[hanzi] == "" {
		return string(hanzi)
	} else {
		return pinyinMap[hanzi]
	}
}

func getDefault(hanzi rune) string {
	tone := getTone(hanzi)

	if tone == "" {
		return tone
	}

	output := make([]rune, utf8.RuneCountInString(tone))

	count := 0
	for _, t := range tone {
		neutral, found := tonesMap[t]
		if found {
			output[count] = neutral
		} else {
			output[count] = t
		}
		count++
	}
	return string(output)
}
func getInitialsInCapitals(hanzi rune) string {
	def := getDefault(hanzi)
	if def == "" {
		return def
	}
	sr := []rune(def)
	if sr[0] > 32 {
		sr[0] = sr[0] - 32
	}
	return string(sr)
}

func getInitials(hanzi rune) string {
	def := getDefault(hanzi)
	if def == "" {
		return def
	}
	sr := []rune(def)
	if sr[0] > 32 {
		sr[0] = sr[0] - 32
	}
	return string(sr[0])
}

func GetPY(src string) (py string) {
	var err error
	src = strings.Trim(src, " ")
	if src != "" {
		py, err = New(src).Split("").Mode(Initials).Convert()
		if err != nil {
			py = src
			fmt.Printf("PinYin Error:%s\r\n", err)
		}
	}
	return
}

func GetPYF(src string) (py string) {
	src = GetPY(src)
	if len(src) < 1 {
		return ""
	}
	py = string(([]rune(src))[0])
	return
}
