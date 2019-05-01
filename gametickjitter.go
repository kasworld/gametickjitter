// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gametickjitter

import (
	"fmt"
	"sync"

	"github.com/kasworld/gametick"
	"github.com/kasworld/globalgametick"
)

type GameTickJitter struct {
	mutex       sync.RWMutex `webformhide:"" stringformhide:""`
	name        string
	startTime   gametick.GameTick
	lastActTime gametick.GameTick
	count       int64
	lastJitter  float64
	lastDur     gametick.GameTick
	avgDur      gametick.GameTick
}

func (jt *GameTickJitter) String() string {
	name := jt.name
	if jt.name == "" {
		name = "GameTickJitter"
	}
	return fmt.Sprintf(
		"%s[Count:%v Avg:%v Last[%v %4.2f%%]",
		name, jt.count, jt.GetAvg().ToTimeDuration(), jt.lastDur.ToTimeDuration(), jt.lastJitter)
}

func New(name string) *GameTickJitter {
	jt := &GameTickJitter{
		name:        name,
		startTime:   globalgametick.GetGameTick(),
		lastActTime: globalgametick.GetGameTick(),
		count:       0,
	}
	return jt
}

func (jt *GameTickJitter) GetAvg() gametick.GameTick {
	return jt.avgDur
	// if jt.count == 0 {
	// 	return 0
	// }
	// return (globalgametick.GetGameTick() - jt.startTime) / gametick.GameTick(jt.count)
}

func (jt *GameTickJitter) Act() float64 {
	jt.mutex.Lock()
	defer jt.mutex.Unlock()

	actTime := globalgametick.GetGameTick()
	return jt.actByValue(actTime)
}

func (jt *GameTickJitter) ActByValue(actTime gametick.GameTick) float64 {
	jt.mutex.Lock()
	defer jt.mutex.Unlock()
	return jt.actByValue(actTime)
}

func (jt *GameTickJitter) actByValue(actTime gametick.GameTick) float64 {
	if jt.count == 0 {
		jt.lastActTime = actTime
		jt.count++
		jt.avgDur = actTime - jt.startTime
		return jt.lastJitter
	}

	jt.count++
	thisDur := actTime - jt.lastActTime
	oldAvg := jt.GetAvg()
	jt.avgDur = (oldAvg + thisDur) / 2
	jt.lastDur = thisDur
	jt.lastJitter = float64(jt.lastDur-oldAvg) * 100 / float64(oldAvg)
	jt.lastActTime = actTime
	return jt.lastJitter
}

func (jt *GameTickJitter) GetLastJitter() float64 {
	jt.mutex.RLock()
	defer jt.mutex.RUnlock()
	return jt.lastJitter
}
func (jt *GameTickJitter) GetLastDur() gametick.GameTick {
	jt.mutex.RLock()
	defer jt.mutex.RUnlock()
	return jt.lastDur
}
