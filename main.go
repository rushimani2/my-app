package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
)

// MyChartProps defines the properties of the chart
type MyChartProps struct {
	cdk8s.ChartProps
}

// NewMyChart creates a new chart
func NewMyChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, &id, &props.ChartProps)

	appLabels := map[string]*string{
		"app": cdk8s.String("go-web-app"),
	}

	// Deployment
	cdk8splus26.NewDeployment(chart, cdk8s.String("go-web-app-deployment"), &cdk8splus26.DeploymentProps{
		Metadata: &cdk8s.ApiObjectMetadata{
			Name:   cdk8s.String("go-web-app"),
			Labels: &appLabels,
		},
		Replicas: cdk8s.Float64(2),
		Containers: &[]*cdk8splus26.ContainerProps{
			{
				Image: cdk8s.String("your-dockerhub-username/go-web-app:latest"), // Replace with actual image
				Port:  cdk8s.Float64(8080),
			},
		},
	})

	// Service
	cdk8splus26.NewService(chart, cdk8s.String("go-web-app-service"), &cdk8splus26.ServiceProps{
		Metadata: &cdk8s.ApiObjectMetadata{
			Name: cdk8s.String("go-web-app"),
		},
		Type: cdk8splus26.ServiceType_LOAD_BALANCER,
		Ports: &[]*cdk8splus26.ServicePort{
			{
				Port:       cdk8s.Float64(80),
				TargetPort: cdk8splus26.IntOrString_FromNumber(cdk8s.Float64(8080)),
			},
		},
		Selector: &appLabels,
	})

	// Ingress (optional)
	cdk8splus26.NewIngress(chart, cdk8s.String("go-web-app-ingress"), &cdk8splus26.IngressProps{
		Metadata: &cdk8s.ApiObjectMetadata{
			Name: cdk8s.String("go-web-app"),
		},
		Rules: &[]*cdk8splus26.IngressRule{
			{
				Path:     cdk8s.String("/"),
				PathType: cdk8splus26.HttpIngressPathType_PREFIX,
				Backend:  cdk8splus26.IngressBackend_FromService(cdk8splus26.Service_FromServiceName(chart, cdk8s.String("go-web-app-service-ref"), cdk8s.String("go-web-app"), cdk8splus26.ServiceType_LOAD_BALANCER)),
			},
		},
	})

	return chart
}

// Main function that triggers the chart synthesis
func main() {
	app := cdk8s.NewApp(nil)
	NewMyChart(app, "go-web-app", &MyChartProps{})
	app.Synth()
}
