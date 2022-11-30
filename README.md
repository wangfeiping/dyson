# Dyson
Command wrapper capable of being invoked remotely or automatically executed

The purpose of this project is to provide decentralized automatic maintenance and management of account assets (e.g., Queries; Allows agencies to provide reverse management of their clients' private assets). Includes features like: command executes automatically; configurable prometheus-based data metrics monitoring; APIs for remote calls, etc.

## Key features:

* (Implemented) Program command executes automatically
* (Implemented) Configurable prometheus-based data metrics monitoring
* (Not yet implemented) APIs for remote calls
* (Not yet implemented) Password caching and automatically entering passwords if necessary

## In order to achieve the project objectives, the following issues need to be considered:

### Security
  * The password can be entered as needed each time the program is started, and the password is only cached in memory; (Password refers to the unlock password of the private key vault, through which the password is used to obtain the permission to use the account private key explicitly authorized in the configuration)
  * The private key vault is managed by the original program. The dyson program will not unlock the private key or obtain it by other means, and will only automatically enter the cached password when the program needs to sign and confirm the transaction;
  * Because the proxy call of the program is implemented through system calls, APIs for external calls need to be whitelisted for security considerations, and the default is not executable without configuration;

### Internal command execution
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

## Usage and configuration instructions


