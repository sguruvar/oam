"ingress-trait": {
	description: "My ingress route trait."
	type: "trait"
	attributes: {
        appliesToWorkloads: ["dep-svc"]
    }
}
template: {
	output: {
		apiVersion: "networking.k8s.io/v1"
		kind:       "Ingress"
		metadata: {
			annotations: {
				"alb.ingress.kubernetes.io/scheme":      "internet-facing"
				"alb.ingress.kubernetes.io/target-type": "ip"
			}
			name:      context.name
			namespace:  context.namespace
		}
		spec: {
			ingressClassName: "aws-alb"
			rules: [{
				http: paths: [{
					backend: service: {
						name:  "oamt1"
						port: number: 80
					}
					path:     "/"
					pathType: "Prefix"
				}]
			}]
		}
	}
	outputs: {}
	parameter: {}
}

