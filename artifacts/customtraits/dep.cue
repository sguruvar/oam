"dep-svc": {
	annotations: {}
	attributes: workload: definition: {
		apiVersion: "apps/v1"
		kind:       "Deployment"
	}
	description: "My component."
	labels: {}
	type: "component"
}

template: {
	output: {
		apiVersion: "apps/v1"
		kind:       "Deployment"
		metadata: name: "oam"
		spec: {
			replicas: 1
			selector: matchLabels: "app": context.name
			template: {
				metadata: labels: "app":  context.name
				spec: containers: [{
					image: "nginx"
					name:  "cn"
					ports: [{
						containerPort: 80
						name:          "http"
						protocol:      "TCP"
					}]
				}]
			}
		}
	}
	outputs: "service": {
		apiVersion: "v1"
		kind:       "Service"
		metadata: name:  context.name
		spec: {
			ports: [{
				name:       "http"
				port:       80
				protocol:   "TCP"
				targetPort: 8080
			}]
			selector: app:  context.name
			type: "ClusterIP"
		}
	}
	parameter: {}
}

