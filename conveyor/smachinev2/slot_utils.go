//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package smachine

func (s *Slot) activateSlot(worker FixedSlotWorker) {
	s.machine.updateSlotQueue(s, worker, activateSlot)
}

func (p SlotLink) activateSlot(worker FixedSlotWorker) {
	if p.IsValid() {
		p.s.activateSlot(worker)
	}
}

func (s *Slot) releaseDependency(worker FixedSlotWorker) {
	dep := s.dependency
	if dep == nil {
		return
	}
	s.dependency = nil
	dep.Release(func(link SlotLink) {
		s.machine._activateDependantByLink(link, worker)
	})
}
