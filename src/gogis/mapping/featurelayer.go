package mapping

import (
	"encoding/json"
	"fmt"
	"gogis/base"
	"gogis/data"
	"gogis/draw"
	"gogis/geometry"

	"github.com/tidwall/mvt"
)

// 要素图层类
type FeatureLayer struct {
	LayerType LayerType
	Name      string          `json:"LayerName"` // 图层名
	feaset    data.Featureset // 数据来源
	Params    data.ConnParams `json:"ConnParams"` // 存储和打开地图文档时用的数据连接信息
	Type      ThemeType       `json:"ThemeType"`
	theme     Theme           // 专题风格
	Object    interface{}     `json:"Theme"` // 好一招狸猫换太子
}

func newFeatureLayer(feaset data.Featureset, theme Theme) *FeatureLayer {
	layer := new(FeatureLayer)
	layer.LayerType = LayerFeature
	layer.Setting(feaset)
	// 默认图层名 等于 数据集名
	layer.Name = layer.Params["name"].(string)
	if theme == nil {
		layer.theme = new(SimpleTheme)
		layer.theme.MakeDefault(feaset)
	} else {
		layer.theme = theme
		layer.Name += "_" + string(theme.GetType())
	}
	layer.theme.MakeDefault(feaset)
	layer.Type = layer.theme.GetType()

	return layer
}

func (this *FeatureLayer) UnmarshalJSON(data []byte) error {
	type cloneType FeatureLayer
	rawMsg := json.RawMessage{}
	this.Object = &rawMsg
	json.Unmarshal(data, (*cloneType)(this))

	this.theme = NewTheme(this.Type)
	json.Unmarshal(rawMsg, this.theme)
	return nil
}

// new出来的时候，做统一设置
func (this *FeatureLayer) Setting(feaset data.Featureset) bool {
	this.feaset = feaset
	store := feaset.GetStore()
	if store != nil {
		this.Params = store.GetConnParams()
		this.Params["name"] = feaset.GetName()
		return true
	}
	return false
}

func (this *FeatureLayer) GetBounds() base.Rect2D { // base.Bounds
	return this.feaset.GetBounds()
}

func (this *FeatureLayer) GetProjection() *base.ProjInfo { // 得到投影坐标系，没有返回nil
	return this.feaset.GetProjection()
}

func (this *FeatureLayer) GetName() string {
	return this.Name
}

func (this *FeatureLayer) GetType() LayerType {
	return LayerFeature
}

func (this *FeatureLayer) GetConnParams() data.ConnParams {
	return this.Params
}

func (this *FeatureLayer) Close() {
	this.feaset.Close()
}

// 地图 Save时，内部存储调整为相对路径
func (this *FeatureLayer) WhenSaving(mappath string) {
	filename, ok := this.Params["filename"]
	if ok {
		storename := filename.(string)
		if len(storename) > 0 {
			this.Params["filename"] = base.GetRelativePath(mappath, storename)
		}
	}
	// 保证专题图类型的存储和读取
	this.Object = this.theme
}

// 地图Open时调用，加载实际数据，准备绘制
func (this *FeatureLayer) WhenOpening(mappath string) {
	store := data.NewDatastore(data.StoreType(this.Params["type"].(string)))
	if store != nil {
		// 如果有文件路径，则需要恢复为绝对路径
		filename, ok := this.Params["filename"]
		if ok {
			storename := filename.(string)
			if len(storename) > 0 {
				this.Params["filename"] = base.GetAbsolutePath(mappath, storename)
			}
		}
		ok, _ = store.Open(this.Params)
		if ok {
			this.feaset, _ = store.GetFeasetByName(this.Params["name"].(string))
			this.feaset.Open()
			// 缓存到内存中
			cache := this.Params["cache"]
			if cache != nil && cache.(bool) {
				this.feaset = data.Cache(this.feaset, []string{})
			}
		}
	}
	if this.theme != nil {
		this.theme.WhenOpenning()
	}
}

func (this *FeatureLayer) Select(obj interface{}) (geos []geometry.Geometry) {
	// todo 暂时只支持拉框选择
	bbox, ok := obj.(base.Rect2D)
	if ok {
		// tr := base.NewTimeRecorder()
		var def data.QueryDef
		// def.SpatialMode = data.Intersects
		def.SpatialObj = bbox
		feait := this.feaset.Query(&def)
		// tr.Output("layer query bounds")
		geos = make([]geometry.Geometry, 0, feait.Count())
		objCount := 1000
		forCount := feait.BeforeNext(objCount)

		for i := 0; i < forCount; i++ {
			// todo 批量处理
			if feas, ok := feait.BatchNext(i); ok {
				// temps := make([]geometry.Geometry, len(feas))
				// var wg *sync.WaitGroup = new(sync.WaitGroup)
				for _, v := range feas {
					polygon, ok := v.Geo.(*geometry.GeoPolygon)
					// if ok {
					// 	wg.Add(1)
					// 	go polygonIsIntersectRect(polygon, bbox, temps, j, wg)
					// }
					// if ok && polygon.IsIntersectsRect(bbox) {
					if ok {
						var geoBbox geometry.GeoPolygon
						geoBbox.Make(bbox)
						var rel geometry.GeoRelation
						rel.A = &geoBbox
						rel.B = polygon
						if rel.IsIntersects() {
							geos = append(geos, polygon)
						}

					} else {
						fmt.Println("id:", polygon.GetID())
					}
				}
				// wg.Wait()
				// for _, v := range temps {
				// 	if v != nil {
				// 		geos = append(geos, v)
				// 	}
				// }
			}
		}
		// tr.Output("layer fetch data")
	}
	return
}

func (this *FeatureLayer) Draw(canvas draw.Canvas, proj *base.ProjInfo) (objCount int64) {
	feaPrj := this.feaset.GetProjection()
	prjc := base.NewPrjConvert(proj, feaPrj)
	bbox := canvas.GetBounds()
	// 查询数据的bbox，要反过来先做投影转化；这样才能查出实际数据来
	if prjc != nil {
		bbox.Min = prjc.DoPnt(bbox.Min)
		bbox.Max = prjc.DoPnt(bbox.Max)
	}
	// tr := base.NewTimeRecorder()
	var def data.QueryDef
	def.SpatialObj = bbox
	def.Fields = []string{} // 不要字段
	feait := this.feaset.Query(&def)
	// tr.Output("query")
	// fmt.Println("count:", feait.Count())

	if this.theme != nil {
		prjc := base.NewPrjConvert(feaPrj, proj)
		objCount = this.theme.Draw(canvas, feait, prjc)
	}
	feait.Close()
	return
}

func (this *FeatureLayer) OutputMvt(mvtLayer *mvt.Layer, canvas draw.Canvas) (count int64) {
	// feait := this.feaset.QueryByBounds(canvas.Params.GetBounds())
	bbox := canvas.GetBounds()
	var def data.QueryDef
	def.SpatialObj = &bbox
	feait := this.feaset.Query(&def)
	goCount := feait.BeforeNext(1000)
	for i := 0; i < goCount; i++ {
		feas, ok := feait.BatchNext(i)
		if ok {
			for _, v := range feas {
				if v.Geo != nil {
					mvtGeotype := geometry.GeoType2Mvt(v.Geo.Type())
					mvtFea := mvtLayer.AddFeature(mvtGeotype)
					geometry.Geo2Mvt(v.Geo, mvtFea, canvas)
					count++
				}
			}
		}
	}
	feait.Close()
	return
}
