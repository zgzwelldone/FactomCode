// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package restapi

import (
	//	"fmt"
	"github.com/FactomProject/btcd/wire"
	"time"
)

// BlockTimer is set to sent End-Of-Minute messages to processor
type BlockTimer struct {
	nextDBlockHeight uint64
	inCtlMsgQueue    chan wire.FtmInternalMsg //incoming message queue for factom control messages
}

// Send End-Of-Minute messages to processor for the current open directory block
func (bt *BlockTimer) StartBlockTimer() {
	//wait till the end of minute
	//the first minute section might be bigger than others. To be improved.
/*	t := time.Now()
	time.Sleep(time.Duration((60 - t.Second()) * 1000000000))
*/
	roundTime := time.Now().Round(time.Minute)
	minutesPassed := roundTime.Minute() - (roundTime.Minute()/10)*10

	for minutesPassed < 10 {

		// Sleep till the end of minute
		t0 := time.Now()
		t0_round := t0.Round(time.Minute)
		if t0.Before(t0_round) {
			time.Sleep(time.Duration((60 + t0.Second()) * 1000000000))
		} else {
			time.Sleep(time.Duration((60 - t0.Second()) * 1000000000))
		}

		eomMsg := &wire.MsgInt_EOM{
			EOM_Type:         wire.END_MINUTE_1 + byte(minutesPassed),
			NextDBlockHeight: bt.nextDBlockHeight,
		}

		//send the end-of-minute message to processor
		bt.inCtlMsgQueue <- eomMsg
		
		minutesPassed++
	}

}
