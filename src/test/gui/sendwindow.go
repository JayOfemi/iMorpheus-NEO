package gui

import (
	"github.com/visualfc/goqt/ui"
	"blockChainMorp/accounts"
	"strings"
	"fmt"
	"blockChainMorp/client"
	"strconv"
	"log"
	"blockChainMorp/transaction/transactions"
	"blockChainMorp/utxo/utxo_set"
	"blockChainMorp/pow/proofofwork"
	"blockChainMorp/types"
	"blockChainMorp/server"
)

type SendWindow struct {
	*ui.QWidget
	btnSend *ui.QPushButton
	btnCancel *ui.QPushButton
	fromEdit *ui.QPlainTextEdit
	toEdit *ui.QPlainTextEdit
	amountEdit *ui.QPlainTextEdit

	hintLable *ui.QLabel

	from string
	useBtcFlag bool

	mw *MainWindow
}

func NewSendWindow(from string, w *MainWindow, BtcFlag bool) (*SendWindow, error) {
	sw := &SendWindow{}
	sw.QWidget = ui.NewWidget()

	sw.from = from
	sw.useBtcFlag = BtcFlag

	sw.mw = w

	sw.setWindowCenterWithSize(400, 240)

	sw.initSendWindowWidgets()

	sw.setSendWindowLayout()

	return sw, nil
}

//init sendwindow widgets
func (sw *SendWindow) initSendWindowWidgets() {
	sw.btnSend = ui.NewPushButton()
	sw.btnSend.SetText("Send Coin")
	sw.btnSend.OnClicked(sw.btnSendOnClick)

	sw.btnCancel = ui.NewPushButton()
	sw.btnCancel.SetText("Cancel")
	sw.btnCancel.OnClicked(sw.btnCancelOnClick)

	sw.fromEdit = ui.NewPlainTextEdit()
	sw.fromEdit.SetReadOnly(true)
	sw.fromEdit.AppendPlainText(sw.from)
	sw.fromEdit.SetFixedHeight(28)

	sw.hintLable = ui.NewLabel()
	sw.hintLable.SetText(" ")

	sw.toEdit = ui.NewPlainTextEdit()
	sw.toEdit.SetFixedHeight(28)

	sw.amountEdit = ui.NewPlainTextEdit()
	sw.amountEdit.SetFixedHeight(28)
}

//set window position and size
func (sw *SendWindow) setWindowCenterWithSize(w, h int) {
	app := ui.Application()
	desktopWidget := app.Desktop()
	screenRect := desktopWidget.ScreenGeometry()

	rect := ui.NewRect()
	width := int32(w)
	height := int32(h)
	x := int32(screenRect.Width() / 2) - width / 2
	y := int32(screenRect.Height() / 2) - height / 2

	rect.Adjust(x, y, x + width, y + height)
	sw.SetGeometry(rect)
}

//set sendwindow layout
func (sw *SendWindow) setSendWindowLayout() {
	addrFromHintLable := ui.NewLabel()
	addrFromHintLable.SetText("From:")

	addrToHintLable := ui.NewLabel()
	addrToHintLable.SetText("To:")

	amountHintLable := ui.NewLabel()
	amountHintLable.SetText("Amount:")

	hbox := ui.NewHBoxLayout()
	hbox.AddWidget(sw.btnSend)
	hbox.AddWidget(sw.btnCancel)

	vbox := ui.NewVBoxLayout()
	vbox.AddWidget(addrFromHintLable)
	vbox.AddWidget(sw.fromEdit)
	vbox.AddWidget(addrToHintLable)
	vbox.AddWidget(sw.toEdit)
	vbox.AddWidget(amountHintLable)
	vbox.AddWidget(sw.amountEdit)
	vbox.AddWidget(sw.hintLable)
	vbox.AddLayout(hbox)

	sw.QWidget.SetLayout(vbox)
	sw.SetWindowTitle("Send...")
}

func (sw *SendWindow) btnSendOnClick() {
	addrFrom := sw.fromEdit.ToPlainText()
	addrTo := sw.toEdit.ToPlainText()
	amountStr := sw.amountEdit.ToPlainText()
	amount := 0
	var err error

	if (strings.Compare(addrTo, "") == 0) || (len(addrTo) <= addressChecksumLen) || (!accounts.ValidateAddress(addrTo)) {
		ui.Async(func() {
			sw.hintLable.SetText("ERROR: Recipient address is not valid")
		})

		return
	} else {
		amount, err = strconv.Atoi(amountStr)
		if err != nil {
			log.Panic(err)
		}

		sw.mw.UnableSendBtn()
		sw.mw.UnableGetBalanceBtn()
		if mineOnFlag {
			sw.mw.UnableStopBtn()
		}
		go client.SendAndDo(addrFrom, addrTo, amount, mineOnFlag, sw.useBtcFlag, sw.newBlockChain)
	}

	ui.Async(func() {
		sw.mw.statusEdit.AppendPlainText(fmt.Sprintf("Sent %d coins from %s to %s\n", amount, addrFrom, addrTo))
		if mineOnFlag {
			sw.mw.statusEdit.AppendPlainText("Mining a new block...\n")
		}
	})

	sw.Close()
}

func (sw *SendWindow) btnCancelOnClick() {
	if mineOnFlag {
		txVar := transactions.NewTrans()
		outVar := transactions.NewTXOutput()
		utxoVar := utxoset.NewUTXOSet()
		powVar := proofofwork.NewProofOfWork()

		var bc *types.Blockchain = nil
		go server.StartServer(txVar, outVar, powVar, utxoVar, minerAddress, bc, btcFlag)

		sw.Close()
	} else {
		sw.newBlockChain()
		sw.Close()
	}
}

func (sw *SendWindow) newBlockChain() {
		sw.mw.SendDone()
}