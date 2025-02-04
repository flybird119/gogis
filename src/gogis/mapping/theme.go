package mapping

import (
	"gogis/base"
	"gogis/data"
	"gogis/draw"
)

// 专题类型定义
type ThemeType string

const (
	ThemeSimple ThemeType = "SimpleTheme" // 普通图层的单一风格
	ThemeGrid   ThemeType = "GridTheme"   // 格网聚合图
	ThemeUnique ThemeType = "UniqueTheme" // 单值专题图
	ThemeRange  ThemeType = "RangeTheme"  // 范围分段专题图

	// Theme      ThemeType = ""
)

type NewThemeFunc func() Theme

var gNewThemes map[ThemeType]NewThemeFunc

// 支持用户自定义专题图类型
func RegisterTheme(themetype ThemeType, newfunc NewThemeFunc) {
	if gNewThemes == nil {
		gNewThemes = make(map[ThemeType]NewThemeFunc)
	}
	gNewThemes[themetype] = newfunc
}

func NewTheme(themeType ThemeType) Theme {
	newfunc, ok := gNewThemes[themeType]
	if ok {
		return newfunc()
	}
	return nil
}

// func NewTheme(themeType ThemeType) Theme {
// 	switch themeType {
// 	case ThemeSimple:
// 		return new(SimpleTheme)
// 	case ThemeGrid:
// 		return new(GridTheme)
// 	case ThemeUnique:
// 		return new(UniqueTheme)
// 	case ThemeRange:
// 		return new(RangeTheme)
// 	}
// 	return nil
// }

type Theme interface {
	Draw(canvas draw.Canvas, feaItr data.FeatureIterator, prjc *base.PrjConvert) int64
	MakeDefault(feaset data.Featureset) // 设置默认值，New出来的时候调用
	WhenOpenning()
	GetType() ThemeType
}
