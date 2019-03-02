package uniqueid

import (
	"testing"
)

// 测试生成id的性能
func BenchmarkNextId(b *testing.B) {
	sf := NewUniqueId(1, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sf.NextId()
	}
}

func TestNextIdOnece(t *testing.T) {
	t.Log("测试id参数是否正确")
	var WorkerId = uint16(1)
	var ReserveId = uint8(1)
	sf := NewUniqueId(1, 1)
	uid, _ := sf.NextId()
	result := Prase(uid)
	if result["workerId"] != uint64(WorkerId) {
		t.Fatal("workerId error, wanted: ", WorkerId, " received: ", result["workerId"])

	}
	if result["reserveId"] != uint64(ReserveId) {
		t.Fatal("reserveId error, wanted: ", ReserveId, " received: ", result["reserveId"])

	}

}

func TestNextIdIfUnique(t *testing.T) {
	t.Log("测试id是否重复")
	var WorkerId = uint16(1)
	var ReserveId = uint8(1)
	sf := NewUniqueId(WorkerId, ReserveId)
	result := map[uint64]bool{}
	for i := 0; i < 5000000; i++ {
		uid, _ := sf.NextId()
		if ok, _ := result[uid]; ok {
			t.Fatal("uid有重复值：", uid, Prase(uid))

		}
		result[uid] = true

	}

}
