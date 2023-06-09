# OAM/Kubevela Command line Tool

This is a CLI for Kubevela to generate YAML files

As of now it supports only WebService and one component. 

## Usage for OAMCTL

If you want to use the tool, download the binary `oamctl` from this repo and you can try like below
`oamctl create deploy --appName="backend" --imagePath="github.com/sg09/oam:backendv1" --port=8080 --replicas=1`

And it will emit the YAML file to the output. Also stores as file as below
`apiversion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: backend
spec:
  components:
  - name: backend
    type: webservice
    properties:
      image: public.ecr.aws/m8u2y2m7/oam:backendv1
      ports:
      - port: 8080
        expose: true
    traits:
    - type: scaler
      properties:
        replicas: 1`
        
## For the CRD
1. Deploy crd definition
```bash
kubectl apply -f artifacts/crd/crd-definition.yaml
```
2. Create roles and rolebinding
```bash
kubectl create role oam-trait-role --verb=get,list,create,update,delete --resource=oamtraits,deployments,services,ingresses,horizontalpodautoscalers
kubectl create rolebinding oam-trait-rolebinding --role=oam-trait-role --serviceaccount=default:default
kubectl create clusterrole oam-trait-role-cluster --verb=get,list,create,update,delete --resource=oamtraits,deployments,services,ingresses,horizontalpodautoscalers
kubectl create clusterrolebinding oam-trait-rolebinding-cluster --clusterrole=oam-trait-role-cluster --serviceaccount=default:default

```
3. Deploy CRD 
```bash
kubectl apply -f artifacts/crd/crddeploy.yaml
```


## Build from Source

You can checkout the code and start building as below

`go build .`
        
