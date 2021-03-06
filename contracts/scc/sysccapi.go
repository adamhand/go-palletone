/*
	This file is part of go-palletone.
	go-palletone is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.
	go-palletone is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.
	You should have received a copy of the GNU General Public License
	along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.
*/
/*
 * Copyright IBM Corp. All Rights Reserved.
 * @author PalletOne core developers <dev@pallet.one>
 * @date 2018
 */

package scc

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/palletone/go-palletone/core/vmContractPub/flogging"
	"github.com/palletone/go-palletone/core/vmContractPub/util"
	"github.com/palletone/go-palletone/contracts/shim"
	"github.com/palletone/go-palletone/core/vmContractPub/ccprovider"
	"github.com/palletone/go-palletone/vm/inproccontroller"

	"github.com/spf13/viper"

	pb "github.com/palletone/go-palletone/core/vmContractPub/protos/peer"
)

var sysccLogger = flogging.MustGetLogger("sccapi")

// SystemChaincode defines the metadata needed to initialize system chaincode
// when the comes up. SystemChaincodes are installed by adding an
// entry in importsysccs.go
type SystemChaincode struct {
	//Unique name of the system chaincode
	Name string

	//Path to the system chaincode; currently not used
	Path string

	//InitArgs initialization arguments to startup the system chaincode
	InitArgs [][]byte

	// Chaincode is the actual chaincode object
	Chaincode shim.Chaincode

	// InvokableExternal keeps track of whether
	// this system chaincode can be invoked
	// through a proposal sent to this peer
	InvokableExternal bool

	// InvokableCC2CC keeps track of whether
	// this system chaincode can be invoked
	// by way of a chaincode-to-chaincode
	// invocation
	InvokableCC2CC bool

	// Enabled a convenient switch to enable/disable system chaincode without
	// having to remove entry from importsysccs.go
	Enabled bool
}

// registerSysCC registers the given system chaincode with the peer
func registerSysCC(syscc *SystemChaincode) (bool, error) {
	if !syscc.Enabled || !isWhitelisted(syscc) {
		sysccLogger.Info(fmt.Sprintf("system chaincode (%s,%s,%t) disabled", syscc.Name, syscc.Path, syscc.Enabled))
		return false, nil
	}

	err := inproccontroller.Register(syscc.Path, syscc.Chaincode)
	if err != nil {
		//if the type is registered, the instance may not be... keep going
		if _, ok := err.(inproccontroller.SysCCRegisteredErr); !ok {
			errStr := fmt.Sprintf("could not register (%s,%v): %s", syscc.Path, syscc, err)
			sysccLogger.Error(errStr)
			return false, fmt.Errorf(errStr)
		}
	}

	sysccLogger.Infof("system chaincode %s(%s) registered", syscc.Name, syscc.Path)
	return true, err
}

// deploySysCC deploys the given system chaincode on a chain
func deploySysCC(chainID string, syscc *SystemChaincode) error {
	if !syscc.Enabled || !isWhitelisted(syscc) {
		sysccLogger.Info(fmt.Sprintf("system chaincode (%s,%s) disabled", syscc.Name, syscc.Path))
		return nil
	}

	var err error

	ccprov := ccprovider.GetChaincodeProvider()
	txid := util.GenerateUUID()
	ctxt := context.Background()
	//glh
	/*
	if chainID != "" {
		lgr := peer.GetLedger(chainID)
		if lgr == nil {
			panic(fmt.Sprintf("syschain %s start up failure - unexpected nil ledger for channel %s", syscc.Name, chainID))
		}

		//init can do GetState (and other Get's) even if Puts cannot be
		//be handled. Need ledger for this
		ctxt2, txsim, err := ccprov.GetContext(lgr, txid)
		if err != nil {
			return err
		}

		ctxt = ctxt2

		defer txsim.Done()
	}
*/
	chaincodeID := &pb.ChaincodeID{Path: syscc.Path, Name: syscc.Name}
	spec := &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: syscc.InitArgs}}

	// First build and get the deployment spec
	chaincodeDeploymentSpec, err := buildSysCC(ctxt, spec)

	if err != nil {
		sysccLogger.Error(fmt.Sprintf("Error deploying chaincode spec: %v\n\n error: %s", spec, err))
		return err
	}
	sysccLogger.Info("buildSysCC chaincodeDeploymentSpec =%v", chaincodeDeploymentSpec)
	version := util.GetSysCCVersion()
	cccid := ccprov.GetCCContext(chainID, chaincodeDeploymentSpec.ChaincodeSpec.ChaincodeId.Name, version, txid, true, nil, nil)

	_, _, err = ccprov.ExecuteWithErrorFilter(ctxt, cccid, chaincodeDeploymentSpec)

	if err != nil {
		sysccLogger.Errorf("ExecuteWithErrorFilter with syscc.Name[%s] chainId[%s] err !!", syscc.Name, chainID)
	}
	sysccLogger.Infof("system chaincode %s/%s(%s) deployed", syscc.Name, chainID, syscc.Path)

	return err
}

// DeDeploySysCC stops the system chaincode and deregisters it from inproccontroller
func DeDeploySysCC(chainID string, syscc *SystemChaincode) error {
	chaincodeID := &pb.ChaincodeID{Path: syscc.Path, Name: syscc.Name}
	spec := &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: syscc.InitArgs}}

	ctx := context.Background()
	// First build and get the deployment spec
	chaincodeDeploymentSpec, err := buildSysCC(ctx, spec)

	if err != nil {
		sysccLogger.Error(fmt.Sprintf("Error deploying chaincode spec: %v\n\n error: %s", spec, err))
		return err
	}

	ccprov := ccprovider.GetChaincodeProvider()

	version := util.GetSysCCVersion()

	cccid := ccprov.GetCCContext(chainID, syscc.Name, version, "", true, nil, nil)

	err = ccprov.Stop(ctx, cccid, chaincodeDeploymentSpec)

	return err
}

// buildLocal builds a given chaincode code
func buildSysCC(context context.Context, spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
	var codePackageBytes []byte
	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ExecEnv: pb.ChaincodeDeploymentSpec_SYSTEM, ChaincodeSpec: spec, CodePackage: codePackageBytes}
	return chaincodeDeploymentSpec, nil
}

func isWhitelisted(syscc *SystemChaincode) bool {
	chaincodes := viper.GetStringMapString("chaincode.system")
	val, ok := chaincodes[syscc.Name]
	enabled := val == "enable" || val == "true" || val == "yes"
	return ok && enabled
}
//RegisterSysCCs is the hook for system chaincodes where system chaincodes are registered
//note the chaincode must still be deployed and launched like a user chaincode will be
func RegisterSysCCs() {
	for _, sysCC := range systemChaincodes {
		sysccLogger.Infof("<%v>", sysCC)
		registerSysCC(sysCC)
	}
}

