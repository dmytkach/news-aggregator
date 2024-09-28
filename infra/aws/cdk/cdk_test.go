package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
	"testing"
)

func TestCdkStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkStack(app, "MyStack", nil)

	// THEN
	template := assertions.Template_FromStack(stack, nil)

	template.HasResourceProperties(jsii.String("AWS::EC2::VPC"), map[string]interface{}{
		"CidrBlock": jsii.String(vpcCIDR),
	})

	template.HasResourceProperties(jsii.String("AWS::EC2::SecurityGroup"), map[string]interface{}{
		"GroupDescription": jsii.String("Allow access to EKS Cluster"),
	})

	template.HasResourceProperties(jsii.String("AWS::EKS::Nodegroup"), map[string]interface{}{
		"ScalingConfig": map[string]interface{}{
			"MinSize":     jsii.Number(1),
			"MaxSize":     jsii.Number(10),
			"DesiredSize": jsii.Number(2),
		},
	})

	template.HasResourceProperties(jsii.String("AWS::IAM::Role"), map[string]interface{}{
		"AssumeRolePolicyDocument": map[string]interface{}{
			"Statement": []interface{}{
				map[string]interface{}{
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"Service": jsii.String("eks.amazonaws.com"),
					},
					"Action": jsii.String("sts:AssumeRole"),
				},
			},
		},
	})

	template.HasResourceProperties(jsii.String("AWS::IAM::Role"), map[string]interface{}{
		"AssumeRolePolicyDocument": map[string]interface{}{
			"Statement": []interface{}{
				map[string]interface{}{
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"Service": jsii.String("ec2.amazonaws.com"),
					},
					"Action": jsii.String("sts:AssumeRole"),
				},
			},
		},
	})

	template.HasResourceProperties(jsii.String("AWS::EKS::Addon"), map[string]interface{}{
		"AddonName":    jsii.String("vpc-cni"),
		"AddonVersion": jsii.String(eksAddonVPCVersion),
	})

	template.HasResourceProperties(jsii.String("AWS::EKS::Addon"), map[string]interface{}{
		"AddonName":    jsii.String("coredns"),
		"AddonVersion": jsii.String(eksAddonCoreDNSVersion),
	})

	template.HasResourceProperties(jsii.String("AWS::EKS::Addon"), map[string]interface{}{
		"AddonName":    jsii.String("kube-proxy"),
		"AddonVersion": jsii.String(eksAddonKubeProxyVersion),
	})

	template.HasResourceProperties(jsii.String("AWS::EKS::Addon"), map[string]interface{}{
		"AddonName":    jsii.String("eks-pod-identity-agent"),
		"AddonVersion": jsii.String(eksAddonPodIdentityVersion),
	})

}
