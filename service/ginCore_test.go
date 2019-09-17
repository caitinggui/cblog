package service

import (
	"reflect"
	"testing"
)

func TestPaginator(t *testing.T) {
	res := Paginator(1, 10, 200)
	tmpList := res["Pages"]
	if !reflect.DeepEqual(tmpList, []uint64{2, 3}) {
		t.Fatalf("test failed: %+v", res)
	}
	tmpList = Paginator(3, 10, 200)["Pages"]
	if !reflect.DeepEqual(tmpList, []uint64{2, 3, 4, 5}) {
		t.Fatalf("test failed: %+v", res)
	}
	tmpList = Paginator(5, 10, 200)["Pages"]
	if !reflect.DeepEqual(tmpList, []uint64{3, 4, 5, 6, 7}) {
		t.Fatalf("test failed: %+v", res)
	}
	tmpList = Paginator(18, 10, 200)["Pages"]
	if !reflect.DeepEqual(tmpList, []uint64{16, 17, 18, 19}) {
		t.Fatalf("test failed: %+v", res)
	}
	tmpList = Paginator(20, 10, 200)["Pages"]
	if !reflect.DeepEqual(tmpList, []uint64{18, 19}) {
		t.Fatalf("test failed: %+v", res)
	}
	tmpList = Paginator(1, 10, 10)["Pages"]
	if !reflect.DeepEqual(tmpList, []uint64{}) {
		t.Fatalf("test failed: %+v", res)
	}
}
