# awslock

awslock is a small go-lang programm, that helps you to lock/unlock non-ephemeral instances on aws that are need to be secured by "DeleteOnTermination" and "DisableApiTermination".

Wether you lock or unlock (e.g. for deletion) all machines, the attributes "DisableApiTermination" on the instace and "DeleteOnTermination" on the instace root device will be set accordingly.

## Prerequesites

Setup and configure awscli to use profiles from ~/.aws/config and ~/.aws/credentials.

All instances, that are allowed to be locked or unlocked need to be tagged with:

```
Ephemeral=False
```

## Usage

```
awslock command [options]
```

### commands

  * lock - lock all non-ephemeral machines
  * unlock - unlock all non-ephemeral machines

### options

  * verbose,v - print return values and debugging information
  * profile,p - select the profile to use
  * tags,t - add additional tags, to filter for (e.g. --tags="Name=fu,Environment=bar")
  * ids,i - add subset of instance IDs, to iterate over (e.g. --ids=i-1234*************,i-5678*************)

## Alternative

Alternatively you can lock/unlock instances by awscli as well, without the need of installing this little go helper.

### locking all non-ephemeral instances
```
aws ec2 describe-instances --filters Name=tag:Ephemeral,Values="False" | jq ".Reservations[].Instances[].InstanceId" -r | xargs -l aws ec2 modify-instance-attribute --disable-api-termination --instance-id
aws ec2 describe-instances --filters Name=tag:Ephemeral,Values="False" | jq ".Reservations[].Instances[].InstanceId" -r | xargs -l aws ec2 modify-instance-attribute --block-device-mappings DeviceName=/dev/sda1,Ebs={DeleteOnTermination=false} --instance-id
```

### unlocking all non-ephemeral instances
```
aws ec2 describe-instances --filters Name=tag:Ephemeral,Values="False" | jq ".Reservations[].Instances[].InstanceId" -r | xargs -l aws ec2 modify-instance-attribute --no-disable-api-termination --instance-id
aws ec2 describe-instances --filters Name=tag:Ephemeral,Values="False" | jq ".Reservations[].Instances[].InstanceId" -r | xargs -l aws ec2 modify-instance-attribute --block-device-mappings DeviceName=/dev/sda1,Ebs={DeleteOnTermination=true} --instance-id
```
