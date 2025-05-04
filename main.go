package main

import (
	"fmt"
	"log"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/pkg/apis/meta/v1"
)

type MyApp struct {
	cdk8s.App
}

func main() {
	app := MyApp{}
	stack := cdk8s.NewStack(&app, "my-stack", nil)

	// Creating a Kubernetes deployment resource using cdk8s plus go
	deployment := cdk8splus26.NewDeployment(stack, "MyDeployment", &cdk8splus26.DeploymentProps{
		Metadata: &v1.ObjectMeta{
			Name: "example-deployment",
		},
		Spec: &cdk8splus26.DeploymentSpec{
			Replicas: 2,
			Template: &cdk8splus26.PodTemplateSpec{
				Spec: &cdk8splus26.PodSpec{
					Containers: &[]*cdk8splus26.Container{
						{
							Name:  "example-container",
							Image: "nginx",
						},
					},
				},
			},
		},
	})

	// Synthesize the Kubernetes YAML into the output directory
	app.Synth()
	fmt.Println("CDK8s application synthesized successfully!")
}
