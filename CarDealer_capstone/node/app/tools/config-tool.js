'use strict'

var fsx = require('fs-extra');
var Client = require('fabric-client');
var base_config_path = "../"

var ConfigTool = class {

    constructor() {

    }

    cleanUpConfigCache(orgName) {
        let client = Client.loadFromConfig(base_config_path + orgName + '.yaml');
        let client_config = client.getClientConfig();

        let store_path = client_config.credentialStore.path;
        fsx.removeSync(store_path);

        let crypto_path = client_config.credentialStore.cryptoStore.path;
        fsx.removeSync(crypto_path);
    }

    initClient() {
        var client = Client.loadFromConfig(base_config_path + 'network.yaml');
        return client.initCredentialStores().then((nothing) => {
            return Promise.resolve(client);
        });
    }
}

module.exports = ConfigTool;
