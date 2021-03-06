// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sign

import (
	"crypto"

	"github.com/insolar/insolar/insolar"
)

type AlgorithmProvider interface {
	DataSigner(crypto.PrivateKey, insolar.Hasher) insolar.Signer
	DigestSigner(crypto.PrivateKey) insolar.Signer
	DataVerifier(crypto.PublicKey, insolar.Hasher) insolar.Verifier
	DigestVerifier(crypto.PublicKey) insolar.Verifier
}
