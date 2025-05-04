package main

import (
    "github.com/aws/constructs-go/constructs/v10"
    "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
    "github.com/cdk8s-team/cdk8s-plus-26-go/cdk8splus26/v2"
    "github.com/aws/jsii-runtime-go"
)

type MyChartProps struct {
    cdk8s.ChartProps
}

func NewMyChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
    chart := cdk8s.NewChart(scope, jsii.String(id), &props.ChartProps)

    appLabels := map[string]*string{
        "app": jsii.String("go-web-app"),
    }

    // Deployment
    deployment := cdk8splus26.NewDeployment(chart, jsii.String("go-web-app-deployment"), &cdk8splus26.DeploymentProps{
        Metadata: &cdk8s.ApiObjectMetadata{
            Name:   jsii.String("go-web-app"),
            Labels: &appLabels,
        },
        Replicas: jsii.Number(2),
        Containers: &[]*cdk8splus26.ContainerProps{
            {
                Image: jsii.String("your-dockerhub-username/go-web-app:latest"), // Replace with actual image
                Port:  jsii.Number(8080),
            },
        },
    })

    // Service
    svc := cdk8splus26.NewService(chart, jsii.String("go-web-app-service"), &cdk8splus26.ServiceProps{
        Metadata: &cdk8s.ApiObjectMetadata{
            Name: jsii.String("go-web-app"),
        },
        Type: cdk8splus26.ServiceType_LOAD_BALANCER,
        Ports: &[]*cdk8splus26.ServicePort{
            {
                Port:       jsii.Number(80),
                TargetPort: cdk8splus26.IntOrString_FromNumber(jsii.Number(8080)),
            },
        },
        Selector: &appLabels,
    })

    // Ingress
    cdk8splus26.NewIngress(chart, jsii.String("go-web-app-ingress"), &cdk8splus26.IngressProps{
        Metadata: &cdk8s.ApiObjectMetadata{
            Name: jsii.String("go-web-app"),
        },
        Rules: &[]*cdk8splus26.IngressRule{
            {
                Path:     jsii.String("/"),
                PathType: cdk8splus26.HttpIngressPathType_PREFIX,
                Backend:  cdk8splus26.IngressBackend_FromService(svc),
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
