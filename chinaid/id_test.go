/*
 * @Description:
 * @Version: 2.0
 * @Autor: qida
 * @Date: 2019-10-23 15:31:21
 */
package chinaid

import (
	"testing"
	"time"
)

// TestIDNum 正确性测试
func TestIDNum(t *testing.T) {
	testID := "341302199006041233"
	detail, err := IDCard(testID).Decode()
	if err != nil {
		t.Error("Valid id parse failed")
	}
	if detail.Province != "安徽省" {
		t.Error("Province check failed")
	}
	if detail.City != "宿州市" {
		t.Error("City check failed")
	}
	if detail.District != "埇桥区" {
		t.Error("District check failed")
	}
	birth, err := time.Parse("20060102", "19900604")
	if err != nil {
		return
	}
	if detail.Birthday != birth {
		t.Error("Birthday check failed")
	}
	if detail.Sex != Male {
		t.Error("Sex check failed")
	}
}
