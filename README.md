# awslock

awslock is a small go-lang programm, that helps you to lock/unlock non-ephemeral instances on aws that are need to be secured by "DeleteOnTermination" and "DisableApiTermination".

Wether you lock or unlock (e.g. for deletion) all machines, the attributes "DisableApiTermination" on the instace and "DeleteOnTermination" on the instace root device will be set accordingly.

## Prerequesites

Setup and configure awscli to use profiles from ~/.aws/config and ~/.aws/credentials.

All instances, that are allowed to be stopped and started need to be tagged with:

```
Ephemeral=False
```

## Usage

```
awspause command [options]
```

### commands

  * lock - lock all non-ephemeral machines
  * unlock - unlock all non-ephemeral machines

### options

  * verbose,v - print return values and debugging information
  * profile,p - select the profile to use
