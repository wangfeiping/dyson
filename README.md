# Dyson
Command wrapper capable of being invoked remotely or automatically executed

The purpose of this project is to provide decentralized automatic maintenance and management of account assets (e.g., Queries; Allows agencies to provide reverse management of their clients' private assets). Includes features like: command executes automatically; configurable prometheus-based data metrics monitoring; APIs for remote calls, etc.

## Key features:

* (Implemented) Program command executes automatically
* (Implemented) Configurable prometheus-based data metrics monitoring
* (Not yet implemented) APIs for remote calls
* (Not yet implemented) Password caching and automatically entering passwords if necessary

## Issues need to be considered:

### Security
  * The password can be entered as needed each time the program is started, and the password is only cached in memory; (Password refers to the unlock password of the private key vault, through which the password is used to obtain the permission to use the account private key explicitly authorized in the configuration)
  * The private key vault is managed by the original program. The dyson program will not unlock the private key or obtain it by other means, and will only automatically enter the cached password when the program needs to sign and confirm the transaction;
  * Because the proxy call of the program is implemented through system calls, APIs for external calls need to be whitelisted for security considerations, and the default is not executable without configuration;

### Command execution
  * Configurable;
  * With the help of the key store function of the original program, multi-account operations can be managed and executed;

### Monitor
  * Count the number of execution warnings and exceptions;
  * Functional availability, health;
  * Execution statistics of internal commands and remote calls;
  * Parse the returned data of configuration commands and generate monitoring metrics;

### Logs
  * Record the command executed each time;
  * Record execution warnings and exceptions;

### APIs for remote calls
  * Only commands authorized by configuration in the whitelist can be executed;

## Instructions

### Build

```bash
# Clone
$ git clone https://github.com/wangfeiping/dyson.git
$ cd ./dyson/

# Look for the tag
$ git tag
v0.0.2

# Check out the source code
$ git checkout v0.0.2

# Build
$ sh ./build.sh
$ ./build/dyson version
```

### Config

```plain
executor:
# Query the block height
- command: '/path_to/gaiad q block'
  parser:
# The returned data is parsed by JSON Path and the result is cached in height.
  - '$.block.header.height'
# Count the number of validators in staking.
- command: '/path_to/gaiad q staking validators -o json | jq ''.validators | length'''
  parser:
# The result is cached in validators_staking.
  - 'validators_staking=$'
  exporter:
# Generate monitoring metric data.
# metric_name:metric_description{"label_name":"label_value"} ${cached_data}
  - 'validators_staking:Validators in staking.{"chain":"GOC","height":"${height}"} ${validators_staking}'
# Query proposals, parse and generate monitoring metrics
- command: '/path_to/gaiad q gov proposals --output json --count-total --limit 10 --status voting_period'
  parser:
  - '$.proposals[0].voting_start_time'
  - '$.proposals[0].voting_end_time'
  - '$.proposals[0].proposal_id'
  exporter:
  - 'proposal:The proposal in voting period.{"chain":"GOC","start":"${voting_start_time}","end":"${voting_end_time}"} ${proposal_id}'
```

### Run

```bash
# Show help info
$ dyson -h
$ dyson start -h

# Use the ./config.yml to execute every 60 seconds and listen on port 25559
$ dyson start -c ./config.yml -d 60 -l :25559
```
