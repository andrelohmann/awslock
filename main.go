package main

import (
  "fmt"
  "os"
  flag "github.com/ogier/pflag"
  aws "github.com/aws/aws-sdk-go/aws"
  awserr "github.com/aws/aws-sdk-go/aws/awserr"
  awssession "github.com/aws/aws-sdk-go/aws/session"
  awsec2 "github.com/aws/aws-sdk-go/service/ec2"
)

// flags
var (
  profile   string
  verbose   bool
  session   *awssession.Session
  service   *awsec2.EC2
  filter    *awsec2.DescribeInstancesInput
  instances []*string
)

func main() {

  // if no parameters are set
  if len(os.Args) < 2 {
    printUsage()
  }

  switch os.Args[1] {
  case "lock":
      lockInstances()
    case "unlock":
      unlockInstances()
    default:
      printUsage()
  }
}

func init() {
  flag.StringVarP(&profile, "profile", "p", "default", "Select the profile to use")
  flag.BoolVarP(&verbose, "verbose", "v", false, "Print return value")
  flag.Parse()
  session = loadSession()
  service = loadService()
  // filter by Tag: Ephemeral=False
  filter = &awsec2.DescribeInstancesInput{
    Filters: []*awsec2.Filter{
      &awsec2.Filter{
        Name: aws.String("tag:Ephemeral"),
        Values: []*string{
          aws.String("False"),
        },
      },
      &awsec2.Filter{
        Name: aws.String("instance-state-name"),
        Values: []*string{
          aws.String("pending"),
          aws.String("running"),
          aws.String("shutting-down"),
          aws.String("stopping"),
          aws.String("stopped"),
        },
      },
    },
  }
  loadInstances()
}

func loadSession() *awssession.Session {
  // Force enable Shared Config support
  sess, err := awssession.NewSessionWithOptions(awssession.Options{
    SharedConfigState: awssession.SharedConfigEnable,
    Profile: profile,
  })

  if err != nil {
    fmt.Println("Error creating session ", err)
    os.Exit(1)
  }

  return sess
}

func loadService() *awsec2.EC2 {
  svc := awsec2.New(session)

  return svc
}

// Extract all instance IDs tagged with Ephemeral=False
func loadInstances() {
  result, err := service.DescribeInstances(filter)
  if err != nil {
    fmt.Println("Error on DescribeInstances", err)
    os.Exit(1)
  } else {
    if verbose {
      fmt.Println("Success on DescribeInstances", result)
    }

    for idx, res := range result.Reservations {
      if verbose {
        fmt.Println("  > Reservation Id", *res.ReservationId, " Num Instances: ", len(res.Instances))
      }

      for _, inst := range result.Reservations[idx].Instances {
        if verbose {
          fmt.Println("    - Instance ID: ", *inst.InstanceId)
        }
        instances = append(instances, inst.InstanceId)
      }
    }
  }
}

func unlockInstances() {
  if len(instances) > 0 {
    for _, instance := range instances {
      unlockInstance(instance)
    }

    fmt.Printf("Success: unlocked %v instance(s)\n", len(instances))
  } else {
    fmt.Println("No instances found")
  }
  os.Exit(0)
}

func unlockInstance(instanceId *string) {
  modifyApiTermination(instanceId, false)
  modifyDeleteOnTermination(instanceId, true)
}

func lockInstances() {
  if len(instances) > 0 {
    for _, instance := range instances {
      lockInstance(instance)
    }

    fmt.Printf("Success: locked %v instance(s)\n", len(instances))
  } else {
    fmt.Println("No instances found")
  }
  os.Exit(0)
}

func lockInstance(instanceId *string) {
  modifyApiTermination(instanceId, true)
  modifyDeleteOnTermination(instanceId, false)
}

func modifyApiTermination(instanceId *string, b bool) {

  input := &awsec2.ModifyInstanceAttributeInput{
    DisableApiTermination: &awsec2.AttributeBooleanValue{
      Value: aws.Bool(b),
    },
    InstanceId: instanceId,
    DryRun: aws.Bool(true),
  }

  result, err := service.ModifyInstanceAttribute(input)
  awsErr, ok := err.(awserr.Error)

  if ok && awsErr.Code() == "DryRunOperation" {
    input.DryRun = aws.Bool(false)
    result, err = service.ModifyInstanceAttribute(input)

    if err != nil {
      fmt.Println("Error", err)
      os.Exit(1)
    } else {
      if verbose {
        fmt.Println("Success", result)
      }
    }
  } else {
    fmt.Println("Error", err)
    os.Exit(1)
  }
}

func modifyDeleteOnTermination(instanceId *string, b bool) {

  input := &awsec2.ModifyInstanceAttributeInput{
    BlockDeviceMappings: []*awsec2.InstanceBlockDeviceMappingSpecification{
      &awsec2.InstanceBlockDeviceMappingSpecification{
        DeviceName: aws.String("/dev/sda1"),
        Ebs: &awsec2.EbsInstanceBlockDeviceSpecification{
          DeleteOnTermination: aws.Bool(b),
        },
      },
    },
    InstanceId: instanceId,
    DryRun: aws.Bool(true),
  }

  result, err := service.ModifyInstanceAttribute(input)
  awsErr, ok := err.(awserr.Error)

  if ok && awsErr.Code() == "DryRunOperation" {
    input.DryRun = aws.Bool(false)
    result, err = service.ModifyInstanceAttribute(input)

    if err != nil {
      fmt.Println("Error", err)
      os.Exit(1)
    } else {
      if verbose {
        fmt.Println("Success", result)
      }
    }
  } else {
    fmt.Println("Error", err)
    os.Exit(1)
  }
}

func printUsage() {
  fmt.Println("awslock helps you lock/unlock EC2 Instances on AWS, that need to be non-ephemeral.")
  fmt.Println("Tag all Instances with Name=Ephemeral,Values=\"False\".")
  fmt.Println("Wether you lock or unlock (e.g. for deletion) all machines, the attributes \"DisableApiTermination\" on the instace and \"DeleteOnTermination\" on the instace root device will be set accordingly.")
  fmt.Println("")
  fmt.Println("Usage:")
  fmt.Println("         awslock command [options]")
  fmt.Println("")
  fmt.Println("The commands are:")
  fmt.Println("         lock    lock all non-ephemeral machines")
  fmt.Println("         unlock  unlock all non-ephemeral machines")
  fmt.Println("")
  fmt.Println("Options:")
  flag.PrintDefaults()
  os.Exit(1)
}
