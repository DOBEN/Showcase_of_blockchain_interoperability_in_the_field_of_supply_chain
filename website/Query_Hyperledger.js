module.exports =async function Query_Hyperledger(_function_name, _key) {

    'use strict';

    const { Wallets, Gateway } = require('fabric-network');
    const fs = require('fs');
    const path = require('path');

    const channelid_1 = "mychannel";
    const smartContractID_1 = "ttcc";

    //var Org=process.argv.slice(2);
    var Org = 1;

    try {
        // Parse the connection profile. This would be the path to the file downloaded
        // from the IBM Blockchain Platform operational console.
        const ccpPath = path.resolve(__dirname, '..', 'fabric-network', 'first-network', 'connection-org' + Org + '.json');
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

        // Configure a wallet. This wallet must already be primed with an identity that
        // the application can use to interact with the peer node.
        const walletPath = path.resolve(__dirname, 'wallet/Org' + Org + 'MSP');
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        // Create a new gateway, and connect to the gateway peer node(s). The identity
        // specified must already exist in the specified wallet.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'appUser', discovery: { enabled: true, asLocalhost: true } });

        // Get the network channel that the smart contract is deployed to.
        const network = await gateway.getNetwork(channelid_1);

        // Get the smart contract from the network channel.
        const contract = network.getContract(smartContractID_1);

        const response = await contract.evaluateTransaction(_function_name, _key);
     
        await gateway.disconnect();
        return response.toString();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}