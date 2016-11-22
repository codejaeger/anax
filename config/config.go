package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type HorizonConfig struct {
	Edge         Config
	AgreementBot AGConfig
}

// This is the configuration options for Edge component flavor of Anax
type Config struct {
	WorkloadROStorage             string
	TorrentDir                    string
	APIListen                     string
	DBPath                        string
	GethURL                       string
	DockerEndpoint                string
	DefaultCPUSet                 string
	DefaultServiceRegistrationRAM int64
	StaticWebContent              string
	PublicKeyPath                 string
	CACertsPath                   string
	ExchangeURL                   string
	PolicyPath                    string
	ExchangeHeartbeat             int    // Seconds between heartbeats
	AgreementTimeoutS             uint64 // Number of seconds to wait before declaring agreement not finalized in blockchain
	DVPrefix                      string // When passing agreement ids into a workload container, add this prefix to the agreement id

	// these Ids could be provided in config or discovered after startup by the system
	BlockchainAccountId        string
	BlockchainDirectoryAddress string
}

// This is the configuration options for Agreement bot flavor of Anax
type AGConfig struct {
	TxLostDelayTolerationSeconds int
	AgreementWorkers             int
	DBPath                       string
	GethURL                      string
	AgreementTimeoutS            uint64 // Number of seconds to wait before declaring agreement not finalized in blockchain
	NoDataIntervalS              uint64 // default should be 15 mins == 15*60 == 900. Ignored if the policy has data verification disabled.
	ActiveAgreementsURL          string // This field is used when policy files indicate they want data verification but they dont specify a URL
	PolicyPath                   string // The directory where policy files are kept, default /etc/provider-tremor/policy/
	NewContractIntervalS         uint64 // default should be 1
	ProcessGovernanceIntervalS   uint64 // How long the gov sleeps before general gov checks (new payloads, interval payments, etc).
	IgnoreContractWithAttribs    string // A comma seperated list of contract attributes. If set, the contracts that contain one or more of the attributes will be ignored. The default is "ethereum_account".
	ExchangeURL                  string // The URL of the Horizon exchange. If not configured, the exchange will not be used.
	ExchangeHeartbeat            int    // Seconds between heartbeats to the exchange
	ExchangeId                   string // The id of the agbot, not the userid of the exchange user
	ExchangeToken                string // The agbot's authentication token
	DVPrefix                     string // When looking for agreement ids in the data verification API response, look for agreement ids with this prefix.
	ActiveDeviceTimeoutS         int    // The amount of time a device can go without heartbeating and still be considered active for the purposes of search

}

func Read(file string) (*HorizonConfig, error) {

	if _, err := os.Stat(file); err != nil {
		return nil, fmt.Errorf("Config file not found: %s. Error: %v", file, err)
	}

	// attempt to parse config file
	path, err := os.Open(filepath.Clean(file))
	if err != nil {
		return nil, fmt.Errorf("Unable to read config file: %s. Error: %v", file, err)
	} else {
		// instantiate empty which will be filled
		config := HorizonConfig{}

		err := json.NewDecoder(path).Decode(&config)
		if err != nil {
			return nil, fmt.Errorf("Unable to decode content of config file: %v", err)
		}

		// success at last!
		return &config, nil
	}
}