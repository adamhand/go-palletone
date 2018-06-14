// Copyright 2016 The go-palletone Authors
// This file is part of the go-palletone library.
//
// The go-palletone library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-palletone library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-palletone library. If not, see <http://www.gnu.org/licenses/>.

package panclient

import "github.com/palletone/go-palletone"

// Verify that Client implements the palletone interfaces.
var (
	_ = palletone.ChainReader(&Client{})
	_ = palletone.TransactionReader(&Client{})
	_ = palletone.ChainStateReader(&Client{})
	_ = palletone.ChainSyncReader(&Client{})
	_ = palletone.ContractCaller(&Client{})
	_ = palletone.GasEstimator(&Client{})
	_ = palletone.GasPricer(&Client{})
	_ = palletone.LogFilterer(&Client{})
	_ = palletone.PendingStateReader(&Client{})
	// _ = palletone.PendingStateEventer(&Client{})
	_ = palletone.PendingContractCaller(&Client{})
)
