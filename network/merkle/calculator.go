/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package merkle

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/pkg/errors"
)

type Calculator interface {
	GetPulseProof(context.Context, *PulseEntry) ([]byte, *PulseProof, error)
	GetGlobuleProof(context.Context, *NodeEntry) (*GlobuleProof, error)
	GetCloudProof(context.Context) (*CloudProof, error)
}

type calculator struct {
	Ledger      core.Ledger      `inject:""`
	NodeNetwork core.NodeNetwork `inject:""`
	Certificate core.Certificate `inject:""`
}

func NewCalculator() Calculator {
	return &calculator{}
}

func (c *calculator) getPulseHash(ctx context.Context, entry *PulseEntry) []byte {
	return pulseHash(entry.Pulse)
}

func (c *calculator) getGlobuleHash(ctx context.Context, entry *NodeEntry) ([]byte, error) {
	globuleHash := make([]byte, 0) // TODO: calculate tree
	return globuleHash, nil
}

func (c *calculator) getCloudHash(ctx context.Context) ([]byte, error) {
	cloudHash := make([]byte, 0) // TODO: calculate tree
	return cloudHash, nil
}

func (c *calculator) getStateHash(role core.NodeRole) ([]byte, error) {
	// TODO: do something with role
	return c.Ledger.GetArtifactManager().State()
}

func (c *calculator) GetPulseProof(ctx context.Context, entry *PulseEntry) ([]byte, *PulseProof, error) {
	role := c.NodeNetwork.GetOrigin().Roles()[0] // TODO: remove after switch to single role model
	stateHash, err := c.getStateHash(role)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ GetPulseProof ] Could't get node stateHash")
	}

	pulseHash := c.getPulseHash(ctx, entry)
	nodeInfoHash := nodeInfoHash(pulseHash, stateHash)

	signature, err := ecdsa.Sign(nodeInfoHash, c.Certificate.GetEcdsaPrivateKey())
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ GetPulseProof ] Could't sign node info hash")
	}

	return pulseHash, &PulseProof{
		StateHash: stateHash,
		Signature: signature,
	}, nil
}

func (c *calculator) GetGlobuleProof(ctx context.Context, entry *NodeEntry) (*GlobuleProof, error) {
	globuleHash, err := c.getGlobuleHash(ctx, entry)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetGlobuleProof ] Could't get globule hash")
	}

	signature, err := ecdsa.Sign(globuleHash, c.Certificate.GetEcdsaPrivateKey())
	if err != nil {
		return nil, errors.Wrap(err, "[ GetGlobuleProof ] Could't sign globule hash")
	}

	return &GlobuleProof{
		Signature: signature,
	}, nil
}

func (c *calculator) GetCloudProof(ctx context.Context) (*CloudProof, error) {
	cloudHash, err := c.getCloudHash(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetCloudProof ] Could't get cloud hash")
	}

	signature, err := ecdsa.Sign(cloudHash, c.Certificate.GetEcdsaPrivateKey())
	if err != nil {
		return nil, errors.Wrap(err, "[ GetCloudProof ] Could't sign cloud hash")
	}

	return &CloudProof{
		Signature: signature,
	}, nil
}
