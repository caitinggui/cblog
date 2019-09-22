package service

import (
	"reflect"
	"testing"
)

func TestPaginator(t *testing.T) {
	res := Paginator(1, 10, 200)
	tmpList := res["Pages"]
	if !reflect.DeepEqual(tmpList, []int{2, 3}) {
		t.Fatalf("test failed: %+v", res)
	}
	tmpList = Paginator(3, 10, 200)["Pages"]
	if !reflect.DeepEqual(tmpList, []int{2, 3, 4, 5}) {
		t.Fatalf("test failed: %+v", res)
	}
	tmpList = Paginator(5, 10, 200)["Pages"]
	if !reflect.DeepEqual(tmpList, []int{3, 4, 5, 6, 7}) {
		t.Fatalf("test failed: %+v", res)
	}
	tmpList = Paginator(18, 10, 200)["Pages"]
	if !reflect.DeepEqual(tmpList, []int{16, 17, 18, 19}) {
		t.Fatalf("test failed: %+v", res)
	}
	tmpList = Paginator(20, 10, 200)["Pages"]
	if !reflect.DeepEqual(tmpList, []int{18, 19}) {
		t.Fatalf("test failed: %+v", res)
	}
	tmpList = Paginator(1, 10, 10)["Pages"]
	if !reflect.DeepEqual(tmpList, []int{}) {
		t.Fatalf("test failed: %+v", res)
	}
}
