package data

import (
	"gogis/base"
	"strconv"
	"testing"
)

func TestConvert(t *testing.T) {
	tr := base.NewTimeRecorder()

	title := "JBNTBHTB" // JBNTBHTB chinapnt_84

	fromParams := NewConnParams()
	fromParams["filename"] = "c:/temp/" + title + ".shp"
	fromParams["type"] = string(StoreShapeMemory)

	toParams := NewConnParams()
	toParams["address"] = "localhost:2181"
	toParams["type"] = string(StoreHbase)

	var cvt Converter
	cvt.Convert(fromParams, title, toParams)

	tr.Output("convert")
}

func TestDelete(t *testing.T) {
	tr := base.NewTimeRecorder()

	title := "JBNTBHTB" // JBNTBHTB

	params := NewConnParams()
	params["address"] = "localhost:2181"
	params["type"] = string(StoreHbase)
	var store HbaseStore
	store.Open(params)
	store.DeleteFeaset(title)

	tr.Output("DeleteFeaset")
}

func TestOpenShape(t *testing.T) {
	var shpStore = new(ShapeStore)
	params := NewConnParams()
	params["filename"] = "./testdata/chinapnt.shp"
	ok, err := shpStore.Open(params)
	if !ok || err != nil {
		t.Errorf(err.Error())
	}
	shp, _ := shpStore.GetFeasetByNum(0)
	if shp.GetCount() != 2391 {
		t.Errorf("对象数量不对")
	}
}

func TestOpenSpatialite(t *testing.T) {
	var store = new(SqliteStore)
	params := NewConnParams()
	params["filename"] = "c:/temp/DLTB.sqlite"
	ok, err := store.Open(params)
	if !ok || err != nil {
		t.Errorf(err.Error())
	}
	feaset, _ := store.GetFeasetByNum(0)
	feaset.Open()
	count := feaset.GetCount()
	if count != 2391 {
		t.Errorf("对象数量不对:" + strconv.Itoa(int(count)))
	}
}
