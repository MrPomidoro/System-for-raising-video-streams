package repositoryMock

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCommon := NewMockCommon(mockCtrl)
	res, err := mockCommon.Get(context.Background(), true)
	if err != nil {
		t.Error("unexpected error", err)
	}
	fmt.Println(res)
}
