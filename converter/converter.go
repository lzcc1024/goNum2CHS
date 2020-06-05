package converter

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	lowerChinese map[string]string = map[string]string{"0": "零", "1": "一", "2": "二", "3": "三", "4": "四", "5": "五", "6": "六", "7": "七", "8": "八", "9": "九"}
	upperChinese map[string]string = map[string]string{"0": "零", "1": "壹", "2": "贰", "3": "叁", "4": "肆", "5": "伍", "6": "陆", "7": "柒", "8": "捌", "9": "玖"}
)

// 定义 Converter 结构
type Converter struct{}

// 创建 Converter 结构对象
func NewConverter() *Converter {
	return &Converter{}
}

// ----------------- 数字转字符 start -----------------

func (this *Converter) Num2Char(number interface{}) (chinese [2]string, err error) {
	var (
		lData string
		uData string
	)

	switch number.(type) {
	case string:
		if lData, err = this.n2CL(number.(string)); err == nil {
			uData, err = this.n2CU(number.(string))
		}
	case int:
		if lData, err = this.n2CL(strconv.Itoa(number.(int))); err == nil {
			uData, err = this.n2CU(strconv.Itoa(number.(int)))
		}
	case int64:
		if lData, err = this.n2CL(strconv.FormatInt(number.(int64), 10)); err == nil {
			uData, err = this.n2CU(strconv.FormatInt(number.(int64), 10))
		}
	default:
		err = errors.New("Not Supported Data Type")
	}

	if err != nil {
		return chinese, err
	}
	chinese[0] = lData
	chinese[1] = uData

	return
}

// 纯数字转中文字符 用 中文小写map
func (this *Converter) n2CL(data string) (chinese string, err error) {
	return this.numeric2Chinese(data, lowerChinese)
}

// 纯数字转中文字符 用 中文大写map
func (this *Converter) n2CU(data string) (chinese string, err error) {
	return this.numeric2Chinese(data, upperChinese)
}

// 纯数字转中文字符
// 负号 和 点 会忽略
func (this *Converter) numeric2Chinese(data string, chineseMap map[string]string) (chinese string, err error) {
	for _, v := range data {
		char := string(v)
		if char == "." || char == "-" {
			continue
		}

		if c, ok := chineseMap[char]; ok {
			chinese += c
		} else {
			return "", errors.New("Convert Failed")
		}
	}
	return
}

// 小数部分数字转中文小写
// 0.9999999999999999 最大长度数值
func (this *Converter) Decimal2ChineseL(number float64) (chinese string, err error) {
	if float64(1) <= number || number == 0 {
		return "", errors.New("only support decimal")
	}
	data := strconv.FormatFloat(number, 'f', -1, 64)

	if data[0:1] != "0" {
		return "", errors.New("only support decimal")
	}

	return this.n2CL(data[2:])
}

// 小数部分数字转中文大写
// 0.9999999999999999 最大长度数值
func (this *Converter) Decimal2ChineseU(number float64) (chinese string, err error) {
	if float64(1) <= number || number == 0 {
		return "", errors.New("only support decimal")
	}
	data := strconv.FormatFloat(number, 'f', -1, 64)

	if data[0:1] != "0" {
		return "", errors.New("only support decimal")
	}

	return this.n2CU(data[2:])
}

// ----------------- 数字转字符 end -----------------

// ----------------- 数字转rmb start -----------------
// https://www.jianshu.com/p/f6367c747798

func (this *Converter) Num2Rmb(number float64) (chinese string, err error) {
	strnum := strconv.FormatFloat(number*100, 'f', 0, 64)
	sliceUnit := []string{"仟", "佰", "拾", "亿", "仟", "佰", "拾", "万", "仟", "佰", "拾", "元", "角", "分"}

	s := sliceUnit[len(sliceUnit)-len(strnum):]

	upperDigitUnit := map[string]string{"0": "零", "1": "壹", "2": "贰", "3": "叁", "4": "肆", "5": "伍", "6": "陆", "7": "柒", "8": "捌", "9": "玖"}

	for k, v := range strnum[:] {
		chinese = chinese + upperDigitUnit[string(v)] + s[k]
	}

	reg, err := regexp.Compile(`零角零分$`)
	chinese = reg.ReplaceAllString(chinese, "整")

	reg, err = regexp.Compile(`零角`)
	chinese = reg.ReplaceAllString(chinese, "零")

	reg, err = regexp.Compile(`零分$`)
	chinese = reg.ReplaceAllString(chinese, "整")

	reg, err = regexp.Compile(`零[仟佰拾]`)
	chinese = reg.ReplaceAllString(chinese, "零")

	reg, err = regexp.Compile(`零{2,}`)
	chinese = reg.ReplaceAllString(chinese, "零")

	reg, err = regexp.Compile(`零亿`)
	chinese = reg.ReplaceAllString(chinese, "亿")

	reg, err = regexp.Compile(`零万`)
	chinese = reg.ReplaceAllString(chinese, "万")

	reg, err = regexp.Compile(`零*元`)
	chinese = reg.ReplaceAllString(chinese, "元")

	reg, err = regexp.Compile(`亿零{0, 3}万`)
	chinese = reg.ReplaceAllString(chinese, "^元")

	reg, err = regexp.Compile(`零元`)
	chinese = reg.ReplaceAllString(chinese, "零")

	return
}

func (this *Converter) Num2Cap(num float64) string {
	strnum := strconv.FormatFloat(num, 'f', 2, 64)
	capitalSlice := []string{"万", "亿", "兆"}
	index := 0
	result := ""
	sdivision := strings.Split(strnum, ".")
	sl := sdivision[0]
	if len(sdivision) > 1 {
		result = "." + sdivision[1]
	}
	// slength := len(sl)
	for len(sl) > 4 {
		result = capitalSlice[index] + sl[len(sl)-4:] + result
		index = index + 1
		sl = sl[0 : len(sl)-4]
	}
	result = sl + result
	result = strings.Replace(result, "万0000", "万", -1)
	result = strings.Replace(result, "亿0000", "亿", -1)
	result = strings.Replace(result, "兆0000", "兆", -1)
	result = strings.Replace(result, "亿万", "亿", -1)
	result = strings.Replace(result, "兆亿", "兆", -1)
	return result
}

// 人民币最大的单位是圆（元）
// 人民币（货币）单位是：元、角、分。
// 而“亿“指的是数目，也可以说是计数单位。不是人民币单位。
// 计数单位依次为
// 个、十、百、千、万、十万、百万、千万、亿、十亿、百亿、千亿 、兆、十兆、百兆、千兆、京、十京、百京、千京、垓、十垓、百垓、千垓、秭、十秭、百秭、千秭、穰、十穰、百穰、千穰、沟、十沟、百沟、千沟、涧、十涧、百涧、千涧、正、十正、百正、千正、载、十载、百载、千载、极、十极、百极、千极、恒河沙、十恒河沙、百恒河沙、千恒河沙、阿僧祗、十阿僧祗、百阿僧祗、千阿僧祗、那由他、十那由他、百那由他、千那由他、不可思议、十不可思议、百不可思议、千不可思议、 无量、十无量、百无量、千无量、大数、十大数、百大数、千大数
// 亦可以写作为:
// 万：10的四次方。 亿：10的八次方。 兆：10的十二次方。 京：10的十六次方。 垓：10的二十次方。 杼：10的二十四次方。 穰：10的二十八次方。 沟：10的三十二次方。 涧：10的三十六次方。 正：10的四十次方。 载：10的四十四次方。 极：10的四十八次方。 恒河沙：10的五十二次方。 阿僧祗：10的五十六次方。 那由他：10的六十次方。 不可思议：10的六十四次方。 无量：10的六十八次方。 大数：10的七十二次方

// ----------------- 数字转大写 end -----------------
