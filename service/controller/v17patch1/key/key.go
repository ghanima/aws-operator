package key

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/aws-operator/service/controller/v17patch1/templates/cloudconfig"
	"github.com/giantswarm/aws-operator/service/controller/v17patch1/templates/cloudformation/guest"
	"github.com/giantswarm/aws-operator/service/controller/v17patch1/templates/cloudformation/hostpost"
	"github.com/giantswarm/aws-operator/service/controller/v17patch1/templates/cloudformation/hostpre"
)

const (
	// CloudConfigVersion defines the version of k8scloudconfig in use.
	// It is used in the main stack output and S3 object paths.
	CloudConfigVersion = "v_3_6_3"

	// CloudProviderTagName is used to add Cloud Provider tags to AWS resources.
	CloudProviderTagName = "kubernetes.io/cluster/%s"

	// Cluster tag name for tagging all resources helping cost analysis in AWS.
	ClusterTagName = "giantswarm.io/cluster"

	// CloudProviderTagOwnedValue is used to indicate an AWS resource is owned
	// and managed by a cluster.
	CloudProviderTagOwnedValue = "owned"

	// EnableTerminationProtection is used to protect the CF stacks from deletion.
	EnableTerminationProtection = true

	// InstallationTagName is used for AWS resource tagging.
	InstallationTagName = "giantswarm.io/installation"

	// OrganizationTagName is used for AWS resource tagging.
	OrganizationTagName = "giantswarm.io/organization"

	// ProfileNameTemplate will be included in the IAM instance profile name.
	ProfileNameTemplate = "EC2-K8S-Role"
	// RoleNameTemplate will be included in the IAM role name.
	RoleNameTemplate = "EC2-K8S-Role"
	// PolicyNameTemplate will be included in the IAM policy name.
	PolicyNameTemplate = "EC2-K8S-Policy"
	// LogDeliveryURI is used for setting the correct ACL in the access log bucket
	LogDeliveryURI = "uri=http://acs.amazonaws.com/groups/s3/LogDelivery"

	InstanceIDAnnotation = "aws-operator.giantswarm.io/instance"

	chinaAWSCliContainerRegistry   = "docker://registry-intl.cn-shanghai.aliyuncs.com/giantswarm/awscli:latest"
	defaultAWSCliContainerRegistry = "quay.io/coreos/awscli:025a357f05242fdad6a81e8a6b520098aa65a600"
	defaultDockerVolumeSizeGB      = 100
)

const (
	DockerVolumeResourceNameKey   = "DockerVolumeResourceName"
	HostedZoneNameServers         = "HostedZoneNameServers"
	MasterImageIDKey              = "MasterImageID"
	MasterInstanceResourceNameKey = "MasterInstanceResourceName"
	MasterInstanceTypeKey         = "MasterInstanceType"
	MasterInstanceMonitoring      = "Monitoring"
	MasterCloudConfigVersionKey   = "MasterCloudConfigVersion"
	WorkerASGKey                  = "WorkerASGName"
	WorkerCountKey                = "WorkerCount"
	WorkerDockerVolumeSizeKey     = "WorkerDockerVolumeSizeGB"
	WorkerImageIDKey              = "WorkerImageID"
	WorkerInstanceMonitoring      = "Monitoring"
	WorkerInstanceTypeKey         = "WorkerInstanceType"
	WorkerCloudConfigVersionKey   = "WorkerCloudConfigVersion"
	VersionBundleVersionKey       = "VersionBundleVersion"
)

const (
	ClusterIDLabel = "giantswarm.io/cluster"

	AnnotationEtcdDomain        = "giantswarm.io/etcd-domain"
	AnnotationPrometheusCluster = "giantswarm.io/prometheus-cluster"

	LabelApp           = "app"
	LabelCluster       = "giantswarm.io/cluster"
	LabelCustomer      = "customer"
	LabelOrganization  = "giantswarm.io/organization"
	LabelVersionBundle = "giantswarm.io/version-bundle"

	LegacyLabelCluster = "cluster"
)

const (
	NodeDrainerLifecycleHookName = "NodeDrainer"
	WorkerASGRef                 = "workerAutoScalingGroup"
)

func ClusterAPIEndpoint(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.Cluster.Kubernetes.API.Domain
}

func AutoScalingGroupName(customObject v1alpha1.AWSConfig, groupName string) string {
	return fmt.Sprintf("%s-%s", ClusterID(customObject), groupName)
}

func AvailabilityZone(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.AWS.AZ
}

func AWSCliContainerRegistry(customObject v1alpha1.AWSConfig) string {
	if IsChinaRegion(customObject) {
		return chinaAWSCliContainerRegistry
	}
	return defaultAWSCliContainerRegistry
}

func BucketName(customObject v1alpha1.AWSConfig, accountID string) string {
	return fmt.Sprintf("%s-g8s-%s", accountID, ClusterID(customObject))
}

// BucketObjectName computes the S3 object path to the actual cloud config.
//
//     /version/3.4.0/cloudconfig/v_3_2_5/master
//     /version/3.4.0/cloudconfig/v_3_2_5/worker
//
func BucketObjectName(customObject v1alpha1.AWSConfig, role string) string {
	return fmt.Sprintf("version/%s/cloudconfig/%s/%s", VersionBundleVersion(customObject), CloudConfigVersion, role)
}

func CredentialName(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.AWS.CredentialSecret.Name
}

func CredentialNamespace(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.AWS.CredentialSecret.Namespace
}

func CloudConfigSmallTemplates() []string {
	return []string{
		cloudconfig.Small,
	}
}

func CloudFormationGuestTemplates() []string {
	return []string{
		guest.AutoScalingGroup,
		guest.IAMPolicies,
		guest.Instance,
		guest.InternetGateway,
		guest.LaunchConfiguration,
		guest.LoadBalancers,
		guest.Main,
		guest.NatGateway,
		guest.LifecycleHooks,
		guest.Outputs,
		guest.RecordSets,
		guest.RouteTables,
		guest.SecurityGroups,
		guest.Subnets,
		guest.VPC,
	}
}

func CloudFormationHostPostTemplates() []string {
	return []string{
		hostpost.Main,
		hostpost.RecordSets,
		hostpost.RouteTables,
	}
}

func CloudFormationHostPreTemplates() []string {
	return []string{
		hostpre.IAMRoles,
		hostpre.Main,
	}
}

func ClusterCloudProviderTag(customObject v1alpha1.AWSConfig) string {
	return fmt.Sprintf(CloudProviderTagName, ClusterID(customObject))
}

func ClusterCustomer(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.Cluster.Customer.ID
}

func ClusterEtcdDomain(customObject v1alpha1.AWSConfig) string {
	return fmt.Sprintf("%s:%d", customObject.Spec.Cluster.Etcd.Domain, customObject.Spec.Cluster.Etcd.Port)
}

func ClusterID(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.Cluster.ID
}

func ClusterNamespace(customObject v1alpha1.AWSConfig) string {
	return ClusterID(customObject)
}

// ClusterOrganization returns the org name from the custom object.
// It uses ClusterCustomer until this field is renamed in the custom object.
func ClusterOrganization(customObject v1alpha1.AWSConfig) string {
	return ClusterCustomer(customObject)
}

func ClusterTags(customObject v1alpha1.AWSConfig, installationName string) map[string]string {
	cloudProviderTag := ClusterCloudProviderTag(customObject)
	tags := map[string]string{
		cloudProviderTag:    CloudProviderTagOwnedValue,
		ClusterTagName:      ClusterID(customObject),
		InstallationTagName: installationName,
		OrganizationTagName: ClusterOrganization(customObject),
	}

	return tags
}

func ClusterVersion(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.Cluster.Version
}

func CustomerID(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.Cluster.Customer.ID
}

func DockerVolumeResourceName(customObject v1alpha1.AWSConfig) string {
	return getResourcenameWithTimeHash("DockerVolume", customObject)
}

func DockerVolumeName(customObject v1alpha1.AWSConfig) string {
	return fmt.Sprintf("%s-docker", ClusterID(customObject))
}

func EtcdVolumeName(customObject v1alpha1.AWSConfig) string {
	return fmt.Sprintf("%s-etcd", ClusterID(customObject))
}

func EC2ServiceDomain(customObject v1alpha1.AWSConfig) string {
	domain := "ec2.amazonaws.com"

	if IsChinaRegion(customObject) {
		domain += ".cn"
	}

	return domain
}

func BaseDomain(customObject v1alpha1.AWSConfig) string {
	// TODO remove other zones and make it a BaseDomain in the CR.
	// CloudFormation creates a separate HostedZone with the same name.
	// Probably the easiest way for now is to just allow single domain for
	// everything which we do now.
	return customObject.Spec.AWS.HostedZones.API.Name
}

func HostedZoneNameAPI(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.AWS.HostedZones.API.Name
}

func HostedZoneNameEtcd(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.AWS.HostedZones.Etcd.Name
}

func HostedZoneNameIngress(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.AWS.HostedZones.Ingress.Name
}

func IngressControllerInsecurePort(customObject v1alpha1.AWSConfig) int {
	return customObject.Spec.Cluster.Kubernetes.IngressController.InsecurePort
}

func IngressControllerSecurePort(customObject v1alpha1.AWSConfig) int {
	return customObject.Spec.Cluster.Kubernetes.IngressController.SecurePort
}

func InstanceProfileName(customObject v1alpha1.AWSConfig, profileType string) string {
	return fmt.Sprintf("%s-%s-%s", ClusterID(customObject), profileType, ProfileNameTemplate)
}

func IsChinaRegion(customObject v1alpha1.AWSConfig) bool {
	return strings.HasPrefix(Region(customObject), "cn-")
}

func IsDeleted(customObject v1alpha1.AWSConfig) bool {
	return customObject.GetDeletionTimestamp() != nil
}

func KubernetesAPISecurePort(customObject v1alpha1.AWSConfig) int {
	return customObject.Spec.Cluster.Kubernetes.API.SecurePort
}

// LoadBalancerName produces a unique name for the load balancer.
// It takes the domain name, extracts the first subdomain, and combines it with the cluster name.
func LoadBalancerName(domainName string, cluster v1alpha1.AWSConfig) (string, error) {
	if ClusterID(cluster) == "" {
		return "", microerror.Maskf(missingCloudConfigKeyError, "spec.cluster.id")
	}

	componentName, err := componentName(domainName)
	if err != nil {
		return "", microerror.Maskf(malformedCloudConfigKeyError, "spec.cluster.id")
	}

	lbName := fmt.Sprintf("%s-%s", ClusterID(cluster), componentName)

	return lbName, nil
}

func MainGuestStackName(customObject v1alpha1.AWSConfig) string {
	clusterID := ClusterID(customObject)

	return fmt.Sprintf("cluster-%s-guest-main", clusterID)
}

func MainHostPreStackName(customObject v1alpha1.AWSConfig) string {
	clusterID := ClusterID(customObject)

	return fmt.Sprintf("cluster-%s-host-setup", clusterID)
}

func MainHostPostStackName(customObject v1alpha1.AWSConfig) string {
	clusterID := ClusterID(customObject)

	return fmt.Sprintf("cluster-%s-host-main", clusterID)
}

func MasterCount(customObject v1alpha1.AWSConfig) int {
	return len(customObject.Spec.AWS.Masters)
}

func MasterImageID(customObject v1alpha1.AWSConfig) string {
	var imageID string

	if len(customObject.Spec.AWS.Masters) > 0 {
		imageID = customObject.Spec.AWS.Masters[0].ImageID
	}

	return imageID
}

func MasterInstanceResourceName(customObject v1alpha1.AWSConfig) string {
	return getResourcenameWithTimeHash("MasterInstance", customObject)
}

func MasterInstanceName(customObject v1alpha1.AWSConfig) string {
	clusterID := ClusterID(customObject)

	return fmt.Sprintf("%s-master", clusterID)
}

func MasterInstanceType(customObject v1alpha1.AWSConfig) string {
	var instanceType string

	if len(customObject.Spec.AWS.Masters) > 0 {
		instanceType = customObject.Spec.AWS.Masters[0].InstanceType
	}

	return instanceType
}

func MasterRoleARN(customObject v1alpha1.AWSConfig, accountID string) string {
	return baseRoleARN(customObject, accountID, "master")
}

func PeerAccessRoleName(customObject v1alpha1.AWSConfig) string {
	return fmt.Sprintf("%s-vpc-peer-access", ClusterID(customObject))
}

func PeerID(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.AWS.VPC.PeerID
}

func PolicyName(customObject v1alpha1.AWSConfig, profileType string) string {
	return fmt.Sprintf("%s-%s-%s", ClusterID(customObject), profileType, PolicyNameTemplate)
}

func PrivateSubnetCIDR(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.AWS.VPC.PrivateSubnetCIDR
}

func CIDR(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.AWS.VPC.CIDR
}

func Region(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.AWS.Region
}

func RegionARN(customObject v1alpha1.AWSConfig) string {
	regionARN := "aws"

	if IsChinaRegion(customObject) {
		regionARN += "-cn"
	}

	return regionARN
}

func RoleName(customObject v1alpha1.AWSConfig, profileType string) string {
	return fmt.Sprintf("%s-%s-%s", ClusterID(customObject), profileType, RoleNameTemplate)
}

func RouteTableName(customObject v1alpha1.AWSConfig, suffix string) string {
	return fmt.Sprintf("%s-%s", ClusterID(customObject), suffix)
}

func S3ServiceDomain(customObject v1alpha1.AWSConfig) string {
	s3Domain := fmt.Sprintf("s3.%s.amazonaws.com", Region(customObject))

	if IsChinaRegion(customObject) {
		s3Domain += ".cn"
	}

	return s3Domain
}

func SecurityGroupName(customObject v1alpha1.AWSConfig, groupName string) string {
	return fmt.Sprintf("%s-%s", ClusterID(customObject), groupName)
}

func SmallCloudConfigPath(customObject v1alpha1.AWSConfig, accountID string, role string) string {
	return fmt.Sprintf("%s/%s", BucketName(customObject, accountID), BucketObjectName(customObject, role))
}

func SmallCloudConfigS3HTTPURL(customObject v1alpha1.AWSConfig, accountID string, role string) string {
	return fmt.Sprintf("https://%s/%s", S3ServiceDomain(customObject), SmallCloudConfigPath(customObject, accountID, role))
}

func SmallCloudConfigS3URL(customObject v1alpha1.AWSConfig, accountID string, role string) string {
	return fmt.Sprintf("s3://%s", SmallCloudConfigPath(customObject, accountID, role))
}

func SubnetName(customObject v1alpha1.AWSConfig, suffix string) string {
	return fmt.Sprintf("%s-%s", ClusterID(customObject), suffix)
}

func TargetLogBucketName(customObject v1alpha1.AWSConfig) string {
	return fmt.Sprintf("%s-g8s-access-logs", ClusterID(customObject))
}

func ToClusterEndpoint(v interface{}) (string, error) {
	customObject, err := ToCustomObject(v)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return ClusterAPIEndpoint(customObject), nil
}

func ToClusterID(v interface{}) (string, error) {
	customObject, err := ToCustomObject(v)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return ClusterID(customObject), nil
}

func ToClusterStatus(v interface{}) (v1alpha1.StatusCluster, error) {
	customObject, err := ToCustomObject(v)
	if err != nil {
		return v1alpha1.StatusCluster{}, microerror.Mask(err)
	}

	return customObject.Status.Cluster, nil
}

func ToCustomObject(v interface{}) (v1alpha1.AWSConfig, error) {
	if v == nil {
		return v1alpha1.AWSConfig{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &v1alpha1.AWSConfig{}, v)
	}

	customObjectPointer, ok := v.(*v1alpha1.AWSConfig)
	if !ok {
		return v1alpha1.AWSConfig{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &v1alpha1.AWSConfig{}, v)
	}
	customObject := *customObjectPointer

	customObject = *customObject.DeepCopy()

	return customObject, nil
}

func ToNodeCount(v interface{}) (int, error) {
	customObject, err := ToCustomObject(v)
	if err != nil {
		return 0, microerror.Mask(err)
	}

	nodeCount := MasterCount(customObject) + WorkerCount(customObject)

	return nodeCount, nil
}

func ToVersionBundleVersion(v interface{}) (string, error) {
	customObject, err := ToCustomObject(v)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return VersionBundleVersion(customObject), nil
}

// VersionBundleVersion returns the version contained in the Version Bundle.
func VersionBundleVersion(customObject v1alpha1.AWSConfig) string {
	return customObject.Spec.VersionBundle.Version
}

func WorkerCount(customObject v1alpha1.AWSConfig) int {
	return len(customObject.Spec.AWS.Workers)
}

// WorkerDockerVolumeSizeGB returns size of a docker volume configured for
// worker nodes. If there are no workers in custom object, 0 is returned as
// size.
func WorkerDockerVolumeSizeGB(customObject v1alpha1.AWSConfig) int {
	if len(customObject.Spec.AWS.Workers) <= 0 {
		return defaultDockerVolumeSizeGB
	}

	if customObject.Spec.AWS.Workers[0].DockerVolumeSizeGB <= 0 {
		return defaultDockerVolumeSizeGB
	}

	return customObject.Spec.AWS.Workers[0].DockerVolumeSizeGB
}

func WorkerImageID(customObject v1alpha1.AWSConfig) string {
	var imageID string

	if len(customObject.Spec.AWS.Workers) > 0 {
		imageID = customObject.Spec.AWS.Workers[0].ImageID
	}

	return imageID
}

func WorkerInstanceType(customObject v1alpha1.AWSConfig) string {
	var instanceType string

	if len(customObject.Spec.AWS.Workers) > 0 {
		instanceType = customObject.Spec.AWS.Workers[0].InstanceType

	}

	return instanceType
}

func WorkerRoleARN(customObject v1alpha1.AWSConfig, accountID string) string {
	return baseRoleARN(customObject, accountID, "worker")
}

func baseRoleARN(customObject v1alpha1.AWSConfig, accountID string, kind string) string {
	clusterID := ClusterID(customObject)
	partition := RegionARN(customObject)

	return fmt.Sprintf("arn:%s:iam::%s:role/%s-%s-%s", partition, accountID, clusterID, kind, RoleNameTemplate)
}

// componentName returns the first component of a domain name.
// e.g. apiserver.example.customer.cloud.com -> apiserver
func componentName(domainName string) (string, error) {
	splits := strings.SplitN(domainName, ".", 2)

	if len(splits) != 2 {
		return "", microerror.Mask(malformedCloudConfigKeyError)
	}

	return splits[0], nil
}

// ImageID returns the EC2 AMI for the configured region.
func ImageID(customObject v1alpha1.AWSConfig) (string, error) {
	region := Region(customObject)

	/*
		Container Linux AMIs for each active AWS region.

		NOTE 1: AMIs should always be for HVM virtualisation and not PV.
		NOTE 2: You also need to update the tests.

		service/controller/v17patch1/key/key_test.go
		service/controller/v17patch1/adapter/adapter_test.go
		service/controller/v17patch1/resource/cloudformation/main_stack_test.go

		Current Release: CoreOS Container Linux stable 1745.4.0 (HVM)
		AMI IDs copied from https://stable.release.core-os.net/amd64-usr/1745.4.0/coreos_production_ami_hvm.txt.
	*/
	imageIDs := map[string]string{
		"ap-northeast-1": "ami-21d9265e",
		"ap-northeast-2": "ami-efe94281",
		"ap-south-1":     "ami-0799b568",
		"ap-southeast-1": "ami-73b28f0f",
		"ap-southeast-2": "ami-8fbf6ced",
		"ca-central-1":   "ami-fb39b99f",
		"cn-north-1":     "ami-2a875e47",
		"cn-northwest-1": "ami-00b8ac62",
		"eu-central-1":   "ami-32042fd9",
		"eu-west-1":      "ami-82645dfb",
		"eu-west-3":      "ami-8d6cddf0",
		"sa-east-1":      "ami-c32d73af",
		"us-east-1":      "ami-93d3b4ec",
		"us-east-2":      "ami-e5cdf180",
		"us-gov-west-1":  "ami-2960f648",
		"us-west-1":      "ami-5e63783e",
		"us-west-2":      "ami-574f362f",
	}

	imageID, ok := imageIDs[region]
	if !ok {
		return "", microerror.Maskf(invalidConfigError, "no image id for region '%s'", region)
	}

	return imageID, nil
}

// getResourcenameWithTimeHash returns the string compared from specific prefix,
// time hash and cluster ID.
func getResourcenameWithTimeHash(prefix string, customObject v1alpha1.AWSConfig) string {
	clusterID := strings.Replace(ClusterID(customObject), "-", "", -1)

	h := sha1.New()
	h.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
	timeHash := fmt.Sprintf("%x", h.Sum(nil))[0:5]

	upperTimeHash := strings.ToUpper(timeHash)
	upperClusterID := strings.ToUpper(clusterID)

	return fmt.Sprintf("%s%s%s", prefix, upperClusterID, upperTimeHash)
}
