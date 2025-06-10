package notification

import (
	"fmt"
	"github.com/gen2brain/beeep"
)

const title = "Reflo"

type BeeepNotifier struct{}

func NewBeeepNotifier() *BeeepNotifier {
	return &BeeepNotifier{}
}

func (b *BeeepNotifier) NotifyFocusComplete() error {
	fmt.Print("\a")
	return beeep.Alert(title, "フォーカスタイムが終了しました！", "")
}

func (b *BeeepNotifier) NotifyBreakComplete() error {
	fmt.Print("\a")
	return beeep.Alert(title, "休憩時間が終了しました！", "")
}
