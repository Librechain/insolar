/*
 *
 *  Copyright  2019. Insolar Technologies GmbH
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package migrationadmin

import (
	"fmt"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/migrationshard"
	"github.com/pkg/errors"
	"strings"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// MigrationAdmin manage and change status for  migration daemon.
type MigrationAdmin struct {
	foundation.BaseContract
	MigrationDaemons     foundation.StableMap
	MigrationAdminMember insolar.Reference
	MigrationAddressShards [insolar.GenesisAmountMigrationAddressShards]insolar.Reference
	VestingParams        *VestingParams
}

type VestingParams struct {
	Lokup       int64 `json:"lokupInPulses"`
	Vesting     int64 `json:"vestingInPulses"`
	VestingStep int64 `json:"vestingStepInPulses"`
}

const (
	StatusActive     = "active"
	StatusInactivate = "inactive"
)

func (mA *MigrationAdmin) MigrationAdminCall(params map[string]interface{}, nameMethod string, caller insolar.Reference) (interface{}, error) {

	switch nameMethod {
	case "addAddresses":
		return mA.addMigrationAddressesCall(params, caller)

	case "activateDaemon":
		return mA.activateDaemonCall(params, caller)

	case "deactivateDaemon":
		return mA.deactivateDaemonCall(params, caller)

	case "checkDaemon":
		return mA.checkDaemonCall(params, caller)
	}
	return nil, fmt.Errorf("unknown method: migration.'%s'", nameMethod)
}

func (mA *MigrationAdmin) activateDaemonCall(params map[string]interface{}, memberRef insolar.Reference) (interface{}, error) {
	migrationDaemon, ok := params["reference"].(string)
	if !ok && len(migrationDaemon) == 0 {
		return nil, fmt.Errorf("incorect input: failed to get 'reference' param")
	}
	return nil, mA.ActivateDaemon(strings.TrimSpace(migrationDaemon), memberRef)
}

func (mA *MigrationAdmin) deactivateDaemonCall(params map[string]interface{}, memberRef insolar.Reference) (interface{}, error) {
	migrationDaemon, ok := params["reference"].(string)

	if !ok && len(migrationDaemon) == 0 {
		return nil, fmt.Errorf("incorect input: failed to get 'reference' param")
	}
	return nil, mA.DeactivateDaemon(strings.TrimSpace(migrationDaemon), memberRef)
}

func (mA *MigrationAdmin) addMigrationAddressesCall(params map[string]interface{}, memberRef insolar.Reference) (interface{}, error) {
	migrationAddresses, ok := params["migrationAddresses"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'migrationAddresses' param")
	}

	if memberRef != mA.MigrationAdminMember {
		return nil, fmt.Errorf("only migration daemon admin can call this method")
	}

	migrationAddressesStr := make([]string, len(migrationAddresses))

	for i, ba := range migrationAddresses {
		migrationAddress, ok := ba.(string)
		if !ok {
			return nil, fmt.Errorf("failed to 'migrationAddresses' param")
		}
		migrationAddressesStr[i] = migrationAddress
	}
	err := mA.addMigrationAddresses(migrationAddressesStr)
	if err != nil {
		return nil, fmt.Errorf("failed to add migration address: %s", err.Error())
	}

	return nil, nil
}

type CheckDaemonResponse struct {
	Status string `json:"status"`
}

func (mA *MigrationAdmin) checkDaemonCall(params map[string]interface{}, caller insolar.Reference) (interface{}, error) {
	migrationDaemon, ok := params["reference"].(string)

	_, err := mA.CheckDaemon(strings.TrimSpace(caller.String()))
	if caller != mA.MigrationAdminMember && err != nil {
		return nil, fmt.Errorf(" permission denied to information about migration daemons: %s", err.Error())
	}

	if !ok && len(migrationDaemon) == 0 {
		return nil, fmt.Errorf("incorect input: failed to get 'reference' param")
	}

	result, err := mA.CheckDaemon(strings.TrimSpace(migrationDaemon))
	if err != nil {
		return nil, fmt.Errorf(" check status migration daemon failed: %s", err.Error())
	}
	if result {
		return CheckDaemonResponse{Status: StatusActive}, nil
	}
	return CheckDaemonResponse{Status: StatusInactivate}, nil
}

// Return stable map migration daemon.
// ins:immutable
func (mA *MigrationAdmin) GetAllMigrationDaemon() (foundation.StableMap, error) {
	sizeMap := len(mA.MigrationDaemons)
	if sizeMap != insolar.GenesisAmountMigrationDaemonMembers {
		return foundation.StableMap{}, fmt.Errorf(" MigrationAdmin contains the wrong amount migration daemon %d", sizeMap)
	}
	return mA.MigrationDaemons, nil
}

// Activate migration daemon.
func (mA *MigrationAdmin) ActivateDaemon(daemonMember string, caller insolar.Reference) error {
	if caller != mA.MigrationAdminMember {
		return fmt.Errorf(" only migration admin can activate migration demons ")
	}
	switch mA.MigrationDaemons[daemonMember] {
	case StatusActive:
		return fmt.Errorf(" daemon member already activated - %s", daemonMember)
	case StatusInactivate:
		mA.MigrationDaemons[daemonMember] = StatusActive
		return nil
	default:
		return fmt.Errorf(" this referense is not daemon member ")
	}
}

// Deactivate migration daemon.
func (mA *MigrationAdmin) DeactivateDaemon(daemonMember string, caller insolar.Reference) error {
	if caller != mA.MigrationAdminMember {
		return fmt.Errorf(" only migration admin can deactivate migration demons ")
	}

	switch mA.MigrationDaemons[daemonMember] {
	case StatusActive:
		mA.MigrationDaemons[daemonMember] = StatusInactivate
		return nil
	case StatusInactivate:
		return fmt.Errorf(" daemon member already deactivated - %s", daemonMember)
	default:
		return fmt.Errorf(" this referense is not daemon member ")
	}
}

// Check this member is migration daemon or mot.
// ins:immutable
func (mA *MigrationAdmin) CheckDaemon(daemonMember string) (bool, error) {
	switch mA.MigrationDaemons[daemonMember] {
	case StatusActive:
		return true, nil
	case StatusInactivate:
		return false, nil
	default:
		return false, fmt.Errorf(" this reference is not daemon member %s", daemonMember)
	}
}

// Return only active daemons.
// ins:immutable
func (mA *MigrationAdmin) GetActiveDaemons() ([]string, error) {
	var activeDaemons []string
	for daemonsRef, status := range mA.MigrationDaemons {
		if status == StatusActive {
			activeDaemons = append(activeDaemons, daemonsRef)
		}
	}
	return activeDaemons, nil
}

func (mA *MigrationAdmin) GetDepositParameters() (*VestingParams, error) {
	return mA.VestingParams, nil
}

// GetMemberByMigrationAddress gets member reference by burn address.
// ins:immutable
func (mA *MigrationAdmin) GetMemberByMigrationAddress(migrationAddress string) (*insolar.Reference, error) {
	trimmedMigrationAddress := foundation.TrimAddress(migrationAddress)
	i := foundation.GetShardIndex(trimmedMigrationAddress, insolar.GenesisAmountMigrationAddressShards)
	if i >= len(mA.MigrationAddressShards) {
		return nil, fmt.Errorf("incorect shard index")
	}
	s := migrationshard.GetObject(mA.MigrationAddressShards[i])
	refStr, err := s.GetRef(trimmedMigrationAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get reference in shard")
	}
	ref, err := insolar.NewReferenceFromBase58(refStr)
	if err != nil {
		return nil, errors.Wrap(err, "bad member reference for this migration address")
	}

	return ref, nil
}

// AddMigrationAddresses adds migration addresses to list.
// ins:immutable
func (mA *MigrationAdmin) addMigrationAddresses(migrationAddresses []string) error {
	newMA := [insolar.GenesisAmountMigrationAddressShards][]string{}
	for _, ma := range migrationAddresses {
		trimmedMigrationAddress := foundation.TrimAddress(ma)
		i := foundation.GetShardIndex(trimmedMigrationAddress, insolar.GenesisAmountMigrationAddressShards)
		if i >= len(newMA) {
			return fmt.Errorf("incorect migration shard index")
		}
		newMA[i] = append(newMA[i], trimmedMigrationAddress)
	}

	for i, ma := range newMA {
		if len(ma) == 0 {
			continue
		}
		s := migrationshard.GetObject(mA.MigrationAddressShards[i])
		err := s.AddFreeMigrationAddresses(ma)
		if err != nil {
			return errors.New("failed to add migration addresses to shard")
		}
	}

	return nil
}

// AddMigrationAddress adds migration address to list.
// ins:immutable
func (mA *MigrationAdmin) addMigrationAddress(migrationAddress string) error {
	trimmedMigrationAddress := foundation.TrimAddress(migrationAddress)
	i := foundation.GetShardIndex(trimmedMigrationAddress, insolar.GenesisAmountMigrationAddressShards)
	if i >= len(mA.MigrationAddressShards) {
		return fmt.Errorf("incorect migration shard index")
	}
	s := migrationshard.GetObject(mA.MigrationAddressShards[i])
	err := s.AddFreeMigrationAddresses([]string{trimmedMigrationAddress})
	if err != nil {
		return errors.New("failed to add migration address to shard")
	}

	return nil
}

// ins:immutable
func (mA *MigrationAdmin) GetFreeMigrationAddress(publicKey string) (string, error) {
	trimmedPublicKey := foundation.TrimPublicKey(publicKey)
	shardIndex := foundation.GetShardIndex(trimmedPublicKey, insolar.GenesisAmountPublicKeyShards)
	if shardIndex >= len(mA.MigrationAddressShards) {
		return "", fmt.Errorf("incorect migration address shard index")
	}

	for i := shardIndex; i < len(mA.MigrationAddressShards); i++ {
		mas := migrationshard.GetObject(mA.MigrationAddressShards[i])
		ma, err := mas.GetFreeMigrationAddress()

		if err == nil {
			return ma, nil
		}

		if err != nil {
			if !strings.Contains(err.Error(), "no more migration address left") {
				return "", errors.Wrap(err, "failed to set reference in migration address shard")
			}
		}
	}

	for i := 0; i < shardIndex; i++ {
		mas := migrationshard.GetObject(mA.MigrationAddressShards[i])
		ma, err := mas.GetFreeMigrationAddress()

		if err == nil {
			return ma, nil
		}

		if err != nil {
			if !strings.Contains(err.Error(), "no more migration address left") {
				return "", errors.Wrap(err, "failed to set reference in migration address shard")
			}
		}
	}

	return "", errors.New("no more migration addresses left in any shard")
}

// AddNewMemberToMaps adds new member to MigrationAddressMap.
// ins:immutable
func (mA *MigrationAdmin) AddNewMigrationAddressToMaps(migrationAddress string, memberRef insolar.Reference) error {
	trimmedMigrationAddress := foundation.TrimAddress(migrationAddress)
	shardIndex := foundation.GetShardIndex(trimmedMigrationAddress, insolar.GenesisAmountPublicKeyShards)
	if shardIndex >= len(mA.MigrationAddressShards) {
		return fmt.Errorf("incorect migration address shard index")
	}
	mas := migrationshard.GetObject(mA.MigrationAddressShards[shardIndex])
	err := mas.SetRef(migrationAddress, memberRef.String())
	if err != nil {
		return errors.Wrap(err, "failed to set reference in migration address shard")
	}

	return nil
}
