package main
 
import (
	"cdk8s"
	cdk8splus "github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26"
	constructs "github.com/aws/constructs-go/constructs/v10"
)

type MyChartProps struct {
	cdk8s.ChartProps
}

func NewMyChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, &id, &props.ChartProps)

	appLabels := map[string]*string{
		"app": cdk8s.String("go-web-app"),
	}

	// Deployment
	cdk8splus.NewDeployment(chart, cdk8s.String("go-web-app-deployment"), &cdk8splus.DeploymentProps{
		Metadata: &cdk8s.ApiObjectMetadata{
			Name:   cdk8s.String("go-web-app"),
			Labels: &appLabels,
		},
		Replicas: cdk8s.Float64(2),
		Containers: &[]*cdk8splus.ContainerProps{
			{
				Image: cdk8s.String("your-dockerhub-username/go-web-app:latest"), // Replace with actual image
				Port:  cdk8s.Float64(8080),
			},
		},
	})

	// Service
	cdk8splus.NewService(chart, cdk8s.String("go-web-app-service"), &cdk8splus.ServiceProps{
		Metadata: &cdk8s.ApiObjectMetadata{
			Name: cdk8s.String("go-web-app"),
		},
		Type: cdk8splus.ServiceType_LOAD_BALANCER,
		Ports: &[]*cdk8splus.ServicePort{
			{
				Port:       cdk8s.Float64(80),
				TargetPort: cdk8splus.IntOrString_FromNumber(cdk8s.Float64(8080)),
			},
		},
		Selector: &appLabels,
	})

	// Ingress (optional)
	cdk8splus.NewIngress(chart, cdk8s.String("go-web-app-ingress"), &cdk8splus.IngressProps{
		Metadata: &cdk8s.ApiObjectMetadata{
			Name: cdk8s.String("go-web-app"),
		},
		Rules: &[]*cdk8splus.IngressRule{
			{
				Path:     cdk8s.String("/"),
				PathType: cdk8splus.HttpIngressPathType_PREFIX,
				Backend:  cdk8splus.IngressBackend_FromService(cdk8splus.Service_FromServiceName(chart, cdk8s.String("go-web-app-service-ref"), cdk8s.String("go-web-app"), cdk8splus.ServiceType_LOAD_BALANCER)),
			},
		},
	})

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewMyChart(app, "go-web-app", &MyChartProps{})
	app.Synth()
}
