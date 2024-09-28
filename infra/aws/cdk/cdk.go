package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	vpcCIDR                    = "10.0.0.0/16"
	publicSubnetName           = "PublicSubnet"
	privateSubnetName          = "PrivateSubnet"
	clusterName                = "DmytroEKSCluster"
	eksClusterSGName           = "DmytroEKSClusterSG"
	iamUserArn                 = "arn:aws:iam::406477933661:user/dmytro"
	nodeInstanceType           = "t3.medium"
	eksAddonVPCVersion         = "v1.18.3-eksbuild.3"
	eksAddonCoreDNSVersion     = "v1.11.3-eksbuild.1"
	eksAddonKubeProxyVersion   = "v1.30.3-eksbuild.5"
	eksAddonPodIdentityVersion = "v1.3.2-eksbuild.2"
	accountID                  = "406477933661"
	region                     = "us-west-1"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	vpc := awsec2.NewVpc(stack, jsii.String("DmytroVpc"), &awsec2.VpcProps{
		IpAddresses: awsec2.IpAddresses_Cidr(jsii.String(vpcCIDR)),
		MaxAzs:      jsii.Number(2),
		NatGateways: jsii.Number(1),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				CidrMask:            jsii.Number(24),
				Name:                jsii.String(publicSubnetName),
				SubnetType:          awsec2.SubnetType_PUBLIC,
				MapPublicIpOnLaunch: jsii.Bool(true),
			},
			{
				CidrMask:   jsii.Number(24),
				Name:       jsii.String(privateSubnetName),
				SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
			},
		},
	})

	sg := awsec2.NewSecurityGroup(stack, jsii.String(eksClusterSGName), &awsec2.SecurityGroupProps{
		Vpc:               vpc,
		SecurityGroupName: jsii.String(eksClusterSGName),
		AllowAllOutbound:  jsii.Bool(true),
		Description:       jsii.String("Allow access to EKS Cluster"),
	})

	sg.AddIngressRule(awsec2.Peer_Ipv4(jsii.String("0.0.0.0/0")), awsec2.Port_Tcp(jsii.Number(22)), jsii.String("Allow SSH access"), jsii.Bool(false))
	sg.AddIngressRule(awsec2.Peer_Ipv4(jsii.String("0.0.0.0/0")), awsec2.Port_Tcp(jsii.Number(80)), jsii.String("Allow HTTP access"), jsii.Bool(false))
	sg.AddIngressRule(awsec2.Peer_Ipv4(jsii.String("0.0.0.0/0")), awsec2.Port_Tcp(jsii.Number(443)), jsii.String("Allow HTTPS access"), jsii.Bool(false))

	eksRole := awsiam.NewRole(stack, jsii.String("DmytroEKSRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("eks.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSClusterPolicy")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSServicePolicy")),
		},
	})

	eksCluster := awseks.NewCluster(stack, jsii.String(clusterName), &awseks.ClusterProps{
		Version:         awseks.KubernetesVersion_V1_30(),
		Vpc:             vpc,
		SecurityGroup:   sg,
		ClusterName:     jsii.String(clusterName),
		DefaultCapacity: jsii.Number(0),
		Role:            eksRole,
	})

	eksCluster.AwsAuth().AddUserMapping(awsiam.User_FromUserArn(stack, jsii.String("dmytro"), jsii.String(iamUserArn)), &awseks.AwsAuthMapping{
		Username: jsii.String("dmytro"),
		Groups: &[]*string{
			jsii.String("system:masters"),
		},
	})

	nodeRole := awsiam.NewRole(stack, jsii.String("DmytroEKSNodeRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSWorkerNodePolicy")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2ContainerRegistryReadOnly")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2FullAccess")),
		},
	})

	awseks.NewNodegroup(stack, jsii.String("DmytroEKSNodeGroup"), &awseks.NodegroupProps{
		Cluster:       eksCluster,
		NodeRole:      nodeRole,
		InstanceTypes: &[]awsec2.InstanceType{awsec2.NewInstanceType(jsii.String(nodeInstanceType))},
		MinSize:       jsii.Number(1),
		MaxSize:       jsii.Number(10),
		DesiredSize:   jsii.Number(2),
		AmiType:       awseks.NodegroupAmiType_AL2_X86_64,
	})

	awseks.NewCfnAddon(stack, jsii.String("VPCCNIAddon"), &awseks.CfnAddonProps{
		ClusterName:      eksCluster.ClusterName(),
		AddonName:        jsii.String("vpc-cni"),
		AddonVersion:     jsii.String(eksAddonVPCVersion),
		ResolveConflicts: jsii.String("OVERWRITE"),
	})

	awseks.NewCfnAddon(stack, jsii.String("CoreDNSAddon"), &awseks.CfnAddonProps{
		ClusterName:      eksCluster.ClusterName(),
		AddonName:        jsii.String("coredns"),
		AddonVersion:     jsii.String(eksAddonCoreDNSVersion),
		ResolveConflicts: jsii.String("OVERWRITE"),
	})

	awseks.NewCfnAddon(stack, jsii.String("KubeProxyAddon"), &awseks.CfnAddonProps{
		ClusterName:      eksCluster.ClusterName(),
		AddonName:        jsii.String("kube-proxy"),
		AddonVersion:     jsii.String(eksAddonKubeProxyVersion),
		ResolveConflicts: jsii.String("OVERWRITE"),
	})

	awseks.NewCfnAddon(stack, jsii.String("PodIdentityAddon"), &awseks.CfnAddonProps{
		ClusterName:      eksCluster.ClusterName(),
		AddonName:        jsii.String("eks-pod-identity-agent"),
		AddonVersion:     jsii.String(eksAddonPodIdentityVersion),
		ResolveConflicts: jsii.String("OVERWRITE"),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkStack(app, "CdkStack", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(accountID),
		Region:  jsii.String(region),
	}
}
