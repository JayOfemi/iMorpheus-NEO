package gui

import (
	"blockChainMorp/accounts"
	"blockChainMorp/encrypt/base58"
	"blockChainMorp/pow/proofofwork"
	"blockChainMorp/server"
	"blockChainMorp/transaction/transactions"
	"blockChainMorp/types"
	"blockChainMorp/utxo/utxo_set"
	"fmt"
	"github.com/visualfc/goqt/ui"
	"log"
	"strings"
	"blockChainMorp/blockchain"
)

var addressChecksumLen = accounts.AddressChecksumLen * 2
var mineOnFlag = false
var minerAddress = ""
var btcFlag = false

type MainWindow struct {
	*ui.QWidget
	addressesSel  *ui.QComboBox
	btnStart      *ui.QPushButton
	btnStop       *ui.QPushButton
	btnSend       *ui.QPushButton
	btnGetBalance *ui.QPushButton
	statusEdit    *ui.QPlainTextEdit
	imRadiobtn    *ui.QRadioButton
	btcRadiobtn   *ui.QRadioButton

	bc            *types.Blockchain
}

func NewMainWindow() (*MainWindow, error) {
	w := &MainWindow{}
	w.QWidget = ui.NewWidget()
	btcFlag = false

	w.setWindowCenterWithSize(640, 480)

	w.initMainWindowWidgets()

	w.setMainWindowLayout()

	w.listAddresses(btcFlag)

	w.setup()

	return w, nil
}

//init mainwindow widgets
func (w *MainWindow) initMainWindowWidgets() {
	w.addressesSel = ui.NewComboBox()

	w.btnStart = ui.NewPushButton()
	w.btnStart.SetText("Start Node")
	w.btnStart.OnClicked(w.BtnStartOnClick)
	w.btnStart.SetEnabled(true)

	w.btnStop = ui.NewPushButton()
	w.btnStop.SetText("Stop Node")
	w.btnStop.OnClicked(w.BtnStopOnClick)
	w.btnStop.SetEnabled(false)

	w.btnSend = ui.NewPushButton()
	w.btnSend.SetText("Send...")
	w.btnSend.OnClicked(w.BtnSendOnClick)
	//w.btnSend.SetEnabled(false)

	w.btnGetBalance = ui.NewPushButton()
	w.btnGetBalance.SetText("Get Balance")
	w.btnGetBalance.OnClicked(w.BtnGetBalanceOnClick)
	w.btnGetBalance.SetEnabled(false)

	w.statusEdit = ui.NewPlainTextEdit()
	w.statusEdit.SetReadOnly(true)

	w.imRadiobtn = ui.NewRadioButton()
	w.imRadiobtn.SetText("iMorpheus Coin")
	w.imRadiobtn.SetChecked(true)
	w.imRadiobtn.OnClicked(w.BtnIMOnClick)

	w.btcRadiobtn = ui.NewRadioButton()
	w.btcRadiobtn.SetText("Bit Coin")
	w.btcRadiobtn.SetChecked(false)
	w.btcRadiobtn.OnClicked(w.BtnBTCOnClick)
}

//set window position and size
func (win *MainWindow) setWindowCenterWithSize(w, h int) {
	app := ui.Application()
	desktopWidget := app.Desktop()
	screenRect := desktopWidget.ScreenGeometry()

	rect := ui.NewRect()
	width := int32(w)
	height := int32(h)
	x := int32(screenRect.Width()/2) - width/2
	y := int32(screenRect.Height()/2) - height/2

	rect.Adjust(x, y, x+width, y+height)
	win.SetGeometry(rect)
}

//set mainwindow layout
func (w *MainWindow) setMainWindowLayout() {
	coinTypeLable := ui.NewLabel()
	coinTypeLable.SetText("Use this coin type:")

	addrHintLable := ui.NewLabel()
	addrHintLable.SetText("Use this address:")

	opsHintLable := ui.NewLabel()
	opsHintLable.SetText("Operations:")

	statusHintLable := ui.NewLabel()
	statusHintLable.SetText("Status and informations:")

	hboxRadioBtns := ui.NewHBoxLayout()
	hboxRadioBtns.AddWidget(w.imRadiobtn)
	hboxRadioBtns.AddWidget(w.btcRadiobtn)

	hboxBtns := ui.NewHBoxLayout()
	hboxBtns.AddWidget(w.btnStart)
	hboxBtns.AddWidget(w.btnStop)
	hboxBtns.AddWidget(w.btnSend)
	hboxBtns.AddWidget(w.btnGetBalance)

	vbox := ui.NewVBoxLayout()
	vbox.AddWidget(coinTypeLable)
	vbox.AddLayout(hboxRadioBtns)
	vbox.AddWidget(addrHintLable)
	vbox.AddWidget(w.addressesSel)
	vbox.AddWidget(opsHintLable)
	vbox.AddLayout(hboxBtns)
	vbox.AddWidget(statusHintLable)
	vbox.AddWidget(w.statusEdit)

	w.QWidget.SetLayout(vbox)
	w.SetWindowTitle("BlockChain Node")
}

func (w *MainWindow) setup() {
	if !blockchain.DbFileOpenedFlag {
		w.bc = blockchain.NewBlockchain()
	}
}



func (w *MainWindow) listAddresses(btcFlag bool) {
	var accountsVar *accounts.Accounts
	var err error
	if btcFlag {
		accountsVar, err = accounts.NewBTCAccounts()
		if err != nil {
			log.Panic(err)
		}
	} else {
		accountsVar, err = accounts.NewAccounts()
		if err != nil {
			log.Panic(err)
		}
	}
	addresses := accountsVar.GetAddresses()

	//for _, address := range addresses {
	//
	//}

	w.addressesSel.InsertItems(0, addresses)
}

func (w *MainWindow) BtnStartOnClick() {
	minerAddress = w.addressesSel.CurrentText()
	fmt.Printf("Starting node\n")

	if (strings.Compare(minerAddress, "") == 0) || (len(minerAddress) <= addressChecksumLen) || (!accounts.ValidateAddress(minerAddress)) {
		log.Panic("Wrong miner address!")
	} else {
		fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
	}

	txVar := transactions.NewTrans()
	outVar := transactions.NewTXOutput()
	utxoVar := utxoset.NewUTXOSet()
	powVar := proofofwork.NewProofOfWork()

	w.bc.Db.Close()
	w.bc = nil
	blockchain.DbFileOpenedFlag = false
	go server.StartServer(txVar, outVar, powVar, utxoVar, minerAddress, w.bc, btcFlag)

	mineOnFlag = true

	//ui
	ui.Async(func() {
		w.statusEdit.AppendPlainText(fmt.Sprintf("Node started and mining is on. Address to receive rewards: %s\n", minerAddress))
	})
	w.btnStart.SetEnabled(false)
	w.btnStop.SetEnabled(true)
	w.btnSend.SetEnabled(true)
	w.btnGetBalance.SetEnabled(false)
}

func (w *MainWindow) BtnStopOnClick() {
	//ui
	ui.Async(func() {
		w.statusEdit.AppendPlainText(fmt.Sprintf("Stopping node..."))
	})
	w.btnSend.SetEnabled(false)
	w.btnGetBalance.SetEnabled(false)
	w.btnStop.SetEnabled(false)

	server.StopServer()
	if !blockchain.DbFileOpenedFlag {
		w.bc = blockchain.NewBlockchain()
	}

	mineOnFlag = false

	//ui
	ui.Async(func() {
		w.statusEdit.AppendPlainText(fmt.Sprintf("Node stopped and mining is off\n"))
	})
	w.btnStart.SetEnabled(true)
	w.btnSend.SetEnabled(true)
	w.btnGetBalance.SetEnabled(true)
}

func (w *MainWindow) BtnSendOnClick() {
	if mineOnFlag {
		server.StopServer()
	} else {
		w.bc.Db.Close()
		w.bc = nil
		blockchain.DbFileOpenedFlag = false
	}

	sw, err := NewSendWindow(w.addressesSel.CurrentText(), w, btcFlag)
	if err != nil {
		log.Fatalln(err)
	}
	sw.SetWindowModality(ui.Qt_ApplicationModal)
	sw.Show()
}

func (w *MainWindow) BtnGetBalanceOnClick() {
	address := w.addressesSel.CurrentText()
	if !accounts.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	UTXOSet := types.UTXOSet{w.bc}

	//bc := blockchain.NewBlockchain()
	//UTXOSet := types.UTXOSet{bc}
	//defer func() {
	//	if bc != nil {
	//		bc.Db.Close()
	//		blockchain.DbFileOpenedFlag = false
	//	}
	//}()


	utxos := utxoset.NewUTXOSet()
	base58coder := base58.NewBase58Coder()
	balance := 0
	pubKeyHash := base58coder.Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := utxos.FindUTXO(UTXOSet, pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	//ui
	str := fmt.Sprintf("Balance of '%s': %d\n", address, balance)
	ui.Async(func() {
		w.statusEdit.AppendPlainText(str)
	})
}

func (w *MainWindow) BtnIMOnClick() {
	btcFlag = false

	//ui
	w.addressesSel.Clear()
	w.listAddresses(btcFlag)
	w.addressesSel.Update()
}

func (w *MainWindow) BtnBTCOnClick() {
	btcFlag = true

	//ui
	w.addressesSel.Clear()
	w.listAddresses(btcFlag)
	w.addressesSel.Update()
}

func (w *MainWindow) UnableSendBtn() {
	w.btnSend.SetEnabled(false)
}

func (w *MainWindow) UnableGetBalanceBtn() {
	w.btnGetBalance.SetEnabled(false)
}

func (w *MainWindow) UnableStopBtn() {
	w.btnStop.SetEnabled(false)
}

func (w *MainWindow) SendDone() {
	if mineOnFlag {
		txVar := transactions.NewTrans()
		outVar := transactions.NewTXOutput()
		utxoVar := utxoset.NewUTXOSet()
		powVar := proofofwork.NewProofOfWork()

		var bc *types.Blockchain = nil
		go server.StartServer(txVar, outVar, powVar, utxoVar, minerAddress, bc, btcFlag)

		w.btnSend.SetEnabled(true)
		w.btnStop.SetEnabled(true)
	} else {
		w.bc = blockchain.NewBlockchain()
		w.btnSend.SetEnabled(true)
		w.btnGetBalance.SetEnabled(true)
	}
}
