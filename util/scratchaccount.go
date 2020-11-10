// Copyright © 2020 Weald Technology Trading
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

package util

import (
	"context"
	"errors"

	"github.com/google/uuid"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

// ScratchAccount is an account that exists temporarily.
type ScratchAccount struct {
	id       uuid.UUID
	privKey  e2types.PrivateKey
	pubKey   e2types.PublicKey
	unlocked bool
}

// NewScratchAccount creates a new local account.
func NewScratchAccount(privKey []byte, pubKey []byte) (*ScratchAccount, error) {
	if len(privKey) > 0 {
		return newScratchAccountFromPrivKey(privKey)
	}
	return newScratchAccountFromPubKey(pubKey)
}

func newScratchAccountFromPrivKey(privKey []byte) (*ScratchAccount, error) {
	key, err := e2types.BLSPrivateKeyFromBytes(privKey)
	if err != nil {
		return nil, err
	}
	return &ScratchAccount{
		id:      uuid.New(),
		privKey: key,
		pubKey:  key.PublicKey(),
	}, nil
}

func newScratchAccountFromPubKey(pubKey []byte) (*ScratchAccount, error) {
	key, err := e2types.BLSPublicKeyFromBytes(pubKey)
	if err != nil {
		return nil, err
	}
	return &ScratchAccount{
		id:     uuid.New(),
		pubKey: key,
	}, nil
}

func (a *ScratchAccount) ID() uuid.UUID {
	return a.id
}

func (a *ScratchAccount) Name() string {
	return "scratch"
}

func (a *ScratchAccount) PublicKey() e2types.PublicKey {
	return a.pubKey
}

func (a *ScratchAccount) Path() string {
	return ""
}

func (a *ScratchAccount) Lock(ctx context.Context) error {
	a.unlocked = false
	return nil
}

func (a *ScratchAccount) Unlock(ctx context.Context, passphrase []byte) error {
	a.unlocked = true
	return nil
}

func (a *ScratchAccount) IsUnlocked(ctx context.Context) (bool, error) {
	return a.unlocked, nil
}

func (a *ScratchAccount) Sign(ctx context.Context, data []byte) (e2types.Signature, error) {
	if !a.unlocked {
		return nil, errors.New("locked")
	}
	if a.privKey == nil {
		return nil, errors.New("no private key")
	}
	return a.privKey.Sign(data), nil
}
