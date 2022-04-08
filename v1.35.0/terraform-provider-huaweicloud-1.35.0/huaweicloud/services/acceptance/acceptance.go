package acceptance

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

var (
	HW_REGION_NAME                = os.Getenv("HW_REGION_NAME")
	HW_CUSTOM_REGION_NAME         = os.Getenv("HW_CUSTOM_REGION_NAME")
	HW_AVAILABILITY_ZONE          = os.Getenv("HW_AVAILABILITY_ZONE")
	HW_ACCESS_KEY                 = os.Getenv("HW_ACCESS_KEY")
	HW_SECRET_KEY                 = os.Getenv("HW_SECRET_KEY")
	HW_USER_ID                    = os.Getenv("HW_USER_ID")
	HW_PROJECT_ID                 = os.Getenv("HW_PROJECT_ID")
	HW_DOMAIN_ID                  = os.Getenv("HW_DOMAIN_ID")
	HW_DOMAIN_NAME                = os.Getenv("HW_DOMAIN_NAME")
	HW_ENTERPRISE_PROJECT_ID_TEST = os.Getenv("HW_ENTERPRISE_PROJECT_ID_TEST")

	HW_FLAVOR_ID             = os.Getenv("HW_FLAVOR_ID")
	HW_FLAVOR_NAME           = os.Getenv("HW_FLAVOR_NAME")
	HW_IMAGE_ID              = os.Getenv("HW_IMAGE_ID")
	HW_IMAGE_NAME            = os.Getenv("HW_IMAGE_NAME")
	HW_VPC_ID                = os.Getenv("HW_VPC_ID")
	HW_NETWORK_ID            = os.Getenv("HW_NETWORK_ID")
	HW_SUBNET_ID             = os.Getenv("HW_SUBNET_ID")
	HW_ENTERPRISE_PROJECT_ID = os.Getenv("HW_ENTERPRISE_PROJECT_ID")
	HW_MAPREDUCE_CUSTOM      = os.Getenv("HW_MAPREDUCE_CUSTOM")
	HW_ADMIN                 = os.Getenv("HW_ADMIN")
	HW_OBS_BUCKET_NAME       = os.Getenv("HW_OBS_BUCKET_NAME")

	HW_DEPRECATED_ENVIRONMENT = os.Getenv("HW_DEPRECATED_ENVIRONMENT")

	HW_WAF_ENABLE_FLAG = os.Getenv("HW_WAF_ENABLE_FLAG")

	HW_DEST_REGION         = os.Getenv("HW_DEST_REGION")
	HW_DEST_PROJECT_ID     = os.Getenv("HW_DEST_PROJECT_ID")
	HW_CHARGING_MODE       = os.Getenv("HW_CHARGING_MODE")
	HW_SWR_SHARING_ACCOUNT = os.Getenv("HW_SWR_SHARING_ACCOUNT")

	HW_CERTIFICATE_KEY_PATH         = os.Getenv("HW_CERTIFICATE_KEY_PATH")
	HW_CERTIFICATE_CHAIN_PATH       = os.Getenv("HW_CERTIFICATE_CHAIN_PATH")
	HW_CERTIFICATE_PRIVATE_KEY_PATH = os.Getenv("HW_CERTIFICATE_PRIVATE_KEY_PATH")
	HW_CERTIFICATE_SERVICE          = os.Getenv("HW_CERTIFICATE_SERVICE")
	HW_CERTIFICATE_PROJECT          = os.Getenv("HW_CERTIFICATE_PROJECT")
	HW_CERTIFICATE_PROJECT_UPDATED  = os.Getenv("HW_CERTIFICATE_PROJECT_UPDATED")
	HW_DMS_ENVIRONMENT              = os.Getenv("HW_DMS_ENVIRONMENT")

	HW_DLI_FLINK_JAR_OBS_PATH = os.Getenv("HW_DLI_FLINK_JAR_OBS_PATH")
)

// TestAccProviders is a static map containing only the main provider instance.
//
// Deprecated: Terraform Plugin SDK version 2 uses TestCase.ProviderFactories
// but supports this value in TestCase.Providers for backwards compatibility.
// In the future Providers: TestAccProviders will be changed to
// ProviderFactories: TestAccProviderFactories
var TestAccProviders map[string]*schema.Provider

// TestAccProviderFactories is a static map containing only the main provider instance
var TestAccProviderFactories map[string]func() (*schema.Provider, error)

// TestAccProvider is the "main" provider instance
var TestAccProvider *schema.Provider

func init() {
	TestAccProvider = huaweicloud.Provider()

	TestAccProviders = map[string]*schema.Provider{
		"huaweicloud": TestAccProvider,
	}

	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"huaweicloud": func() (*schema.Provider, error) {
			return TestAccProvider, nil
		},
	}
}

// ServiceFunc the HuaweiCloud resource query functions.
type ServiceFunc func(*config.Config, *terraform.ResourceState) (interface{}, error)

// resourceCheck resource check object, only used in the package.
type resourceCheck struct {
	resourceName    string
	resourceObject  interface{}
	getResourceFunc ServiceFunc
	resourceType    string
}

const (
	resourceTypeCode   = "resource"
	dataSourceTypeCode = "dataSource"

	checkAttrRegexpStr = `^\$\{([^\}]+)\}$`
)

/*
InitDataSourceCheck build a 'resourceCheck' object. Only used to check datasource attributes.
  Parameters:
    resourceName:    The resource name is used to check in the terraform.State.e.g. : huaweicloud_waf_domain.domain_1.
  Return:
    *resourceCheck: resourceCheck object
*/
func InitDataSourceCheck(sourceName string) *resourceCheck {
	return &resourceCheck{
		resourceName: sourceName,
		resourceType: dataSourceTypeCode,
	}
}

/*
InitResourceCheck build a 'resourceCheck' object. The common test methods are provided in 'resourceCheck'.
  Parameters:
    resourceName:    The resource name is used to check in the terraform.State.e.g. : huaweicloud_waf_domain.domain_1.
    resourceObject:  Resource object, used to check whether the resource exists in HuaweiCloud.
    getResourceFunc: The function used to get the resource object.
  Return:
    *resourceCheck: resourceCheck object
*/
func InitResourceCheck(resourceName string, resourceObject interface{}, getResourceFunc ServiceFunc) *resourceCheck {
	return &resourceCheck{
		resourceName:    resourceName,
		resourceObject:  resourceObject,
		getResourceFunc: getResourceFunc,
		resourceType:    resourceTypeCode,
	}
}

func parseVariableToName(varStr string) (string, string, error) {
	var resName, keyName string
	// Check the format of the variable.
	match, _ := regexp.MatchString(checkAttrRegexpStr, varStr)
	if !match {
		return resName, keyName, fmtp.Errorf("The type of 'variable' is error, "+
			"expected ${resourceType.name.field} got %s", varStr)
	}

	reg, err := regexp.Compile(checkAttrRegexpStr)
	if err != nil {
		return resName, keyName, fmtp.Errorf("The acceptance function is wrong.")
	}
	mArr := reg.FindStringSubmatch(varStr)
	if len(mArr) != 2 {
		return resName, keyName, fmtp.Errorf("The type of 'variable' is error, "+
			"expected ${resourceType.name.field} got %s", varStr)
	}

	// Get resName and keyName from variable.
	strs := strings.Split(mArr[1], ".")
	for i, s := range strs {
		if strings.Contains(s, "huaweicloud_") {
			resName = strings.Join(strs[0:i+2], ".")
			keyName = strings.Join(strs[i+2:], ".")
			break
		}
	}
	return resName, keyName, nil
}

/*
TestCheckResourceAttrWithVariable validates the variable in state for the given name/key combination.
  Parameters:
    resourceName: The resource name is used to check in the terraform.State.
    key:          The field name of the resource.
    variable:     The variable name of the value to be checked.

    variable such like ${huaweicloud_waf_certificate.certificate_1.id}
    or ${data.huaweicloud_waf_policies.policies_2.policies.0.id}
*/
func TestCheckResourceAttrWithVariable(resourceName, key, varStr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resName, keyName, err := parseVariableToName(varStr)
		if err != nil {
			return err
		}

		if strings.EqualFold(resourceName, resName) {
			return fmtp.Errorf("Meaningless verification. " +
				"The referenced resource cannot be the current resource.")
		}

		// Get the value based on resName and keyName from the state.
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmtp.Errorf("Can't find %s in state : %v.", resName, ok)
		}
		value := rs.Primary.Attributes[keyName]

		return resource.TestCheckResourceAttr(resourceName, key, value)(s)
	}
}

// CheckResourceDestroy check whether resources destroyed in HuaweiCloud.
func (rc *resourceCheck) CheckResourceDestroy() resource.TestCheckFunc {
	if strings.Compare(rc.resourceType, dataSourceTypeCode) == 0 {
		fmtp.Errorf("Error, you built a resourceCheck with 'InitDataSourceCheck', " +
			"it cannot run CheckResourceDestroy().")
		return nil
	}
	return func(s *terraform.State) error {
		strs := strings.Split(rc.resourceName, ".")
		var resourceType string
		for _, str := range strs {
			if strings.Contains(str, "huaweicloud_") {
				resourceType = strings.Trim(str, " ")
				break
			}
		}

		if resourceType == "" {
			return fmtp.Errorf("The format of the resource name is invalid, please check your configuration.")
		}

		conf := TestAccProvider.Meta().(*config.Config)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			if rc.getResourceFunc != nil {
				if _, err := rc.getResourceFunc(conf, rs); err == nil {
					return fmtp.Errorf("failed to destroy resource. The resource of %s : %s still exists.",
						resourceType, rs.Primary.ID)
				}
			} else {
				return fmtp.Errorf("The 'getResourceFunc' is nil, please set it during initialization.")
			}
		}
		return nil
	}
}

func (rc *resourceCheck) checkResourceExists(s *terraform.State) error {
	rs, ok := s.RootModule().Resources[rc.resourceName]
	if !ok {
		return fmtp.Errorf("Can not found the resource or data source in state: %s", rc.resourceName)
	}
	if rs.Primary.ID == "" {
		return fmtp.Errorf("No id set for the resource or data source: %s", rc.resourceName)
	}
	if strings.EqualFold(rc.resourceType, dataSourceTypeCode) {
		return nil
	}

	if rc.getResourceFunc != nil {
		conf := TestAccProvider.Meta().(*config.Config)
		r, err := rc.getResourceFunc(conf, rs)
		if err != nil {
			return fmtp.Errorf("checking resource %s %s exists error: %s ",
				rc.resourceName, rs.Primary.ID, err)
		}
		if rc.resourceObject != nil {
			b, err := json.Marshal(r)
			if err != nil {
				return fmtp.Errorf("marshaling resource %s %s error: %s ",
					rc.resourceName, rs.Primary.ID, err)
			}
			json.Unmarshal(b, rc.resourceObject)
		} else {
			logp.Printf("[WARN] The 'resourceObject' is nil, please set it during initialization.")
		}
	} else {
		return fmtp.Errorf("The 'getResourceFunc' is nil, please set it.")
	}

	return nil
}

// CheckResourceExists check whether resources exist in HuaweiCloud.
func (rc *resourceCheck) CheckResourceExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return rc.checkResourceExists(s)
	}
}

/*
CheckMultiResourcesExists checks whether multiple resources created by count in HuaweiCloud are both existed.
  Parameters:
    expCount: the expected number of resources that will be created.
*/
func (rc *resourceCheck) CheckMultiResourcesExists(expCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var err error
		for i := 0; i < expCount; i++ {
			rcCopy := *rc
			rcCopy.resourceName = fmt.Sprintf("%s.%d", rcCopy.resourceName, i)
			err = rcCopy.checkResourceExists(s)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func preCheckRequiredEnvVars(t *testing.T) {
	if HW_REGION_NAME == "" {
		t.Fatal("HW_REGION_NAME must be set for acceptance tests")
	}
}

//lintignore:AT003
func TestAccPreCheck(t *testing.T) {
	// Do not run the test if this is a deprecated testing environment.
	if HW_DEPRECATED_ENVIRONMENT != "" {
		t.Skip("This environment only runs deprecated tests")
	}

	preCheckRequiredEnvVars(t)
}

//lintignore:AT003
func TestAccPrecheckCustomRegion(t *testing.T) {
	if HW_CUSTOM_REGION_NAME == "" {
		t.Skip("HW_CUSTOM_REGION_NAME must be set for acceptance tests")
	}
}

//lintignore:AT003
func TestAccPreCheckDeprecated(t *testing.T) {
	if HW_DEPRECATED_ENVIRONMENT == "" {
		t.Skip("This environment does not support deprecated tests")
	}

	preCheckRequiredEnvVars(t)
}

//lintignore:AT003
func TestAccPreCheckEpsID(t *testing.T) {
	if HW_ENTERPRISE_PROJECT_ID == "" {
		t.Skip("This environment does not support Enterprise Project ID tests")
	}
}

//lintignore:AT003
func TestAccPreCheckBms(t *testing.T) {
	if HW_USER_ID == "" {
		t.Skip("HW_USER_ID must be set for BMS acceptance tests")
	}
}

//lintignore:AT003
func TestAccPreCheckMrsCustom(t *testing.T) {
	if HW_MAPREDUCE_CUSTOM == "" {
		t.Skip("HW_MAPREDUCE_CUSTOM must be set for acceptance tests:custom type cluster of map reduce")
	}
}

func RandomAccResourceName() string {
	return fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
}

func RandomAccResourceNameWithDash() string {
	return fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
}

func RandomCidr() string {
	return fmt.Sprintf("172.16.%d.0/24", acctest.RandIntRange(0, 255))
}

func RandomCidrAndGatewayIp() (string, string) {
	seed := acctest.RandIntRange(0, 255)
	return fmt.Sprintf("172.16.%d.0/24", seed), fmt.Sprintf("172.16.%d.1", seed)
}

//lintignore:AT003
func TestAccPrecheckWafInstance(t *testing.T) {
	if HW_WAF_ENABLE_FLAG == "" {
		t.Skip("Jump the WAF acceptance tests.")
	}
}

//lintignore:AT003
func TestAccPreCheckAdminOnly(t *testing.T) {
	if HW_ADMIN == "" {
		t.Skip("Skipping test because it requires the admin privileges")
	}
}

//lintignore:AT003
func TestAccPreCheckReplication(t *testing.T) {
	if HW_DEST_REGION == "" || HW_DEST_PROJECT_ID == "" {
		t.Skip("Jump the replication policy acceptance tests.")
	}
}

//lintignore:AT003
func TestAccPreCheckProject(t *testing.T) {
	if HW_ENTERPRISE_PROJECT_ID_TEST != "" {
		t.Skip("This environment does not support project tests")
	}
}

//lintignore:AT003
func TestAccPreCheckOBS(t *testing.T) {
	if HW_ACCESS_KEY == "" || HW_SECRET_KEY == "" {
		t.Skip("HW_ACCESS_KEY and HW_SECRET_KEY must be set for OBS acceptance tests")
	}
}

//lintignore:AT003
func TestAccPreCheckOBSBucket(t *testing.T) {
	if HW_OBS_BUCKET_NAME == "" {
		t.Skip("HW_OBS_BUCKET_NAME must be set for OBS object acceptance tests")
	}
}

//lintignore:AT003
func TestAccPreCheckChargingMode(t *testing.T) {
	if HW_CHARGING_MODE != "prePaid" {
		t.Skip("This environment does not support prepaid tests")
	}
}

//lintignore:AT003
func TestAccPreCheckScm(t *testing.T) {
	if HW_CERTIFICATE_KEY_PATH == "" || HW_CERTIFICATE_CHAIN_PATH == "" ||
		HW_CERTIFICATE_PRIVATE_KEY_PATH == "" || HW_CERTIFICATE_SERVICE == "" ||
		HW_CERTIFICATE_PROJECT == "" || HW_CERTIFICATE_PROJECT_UPDATED == "" {
		t.Skip("HW_CERTIFICATE_KEY_PATH, HW_CERTIFICATE_CHAIN_PATH, HW_CERTIFICATE_PRIVATE_KEY_PATH, " +
			"HW_CERTIFICATE_SERVICE, HW_CERTIFICATE_PROJECT and HW_CERTIFICATE_TARGET_UPDATED " +
			"can not be empty for SCM certificate tests")
	}
}

//lintignore:AT003
func TestAccPreCheckSWRDomian(t *testing.T) {
	if HW_SWR_SHARING_ACCOUNT == "" {
		t.Skip("HW_SWR_SHARING_ACCOUNT must be set for swr domian tests, " +
			"the value of HW_SWR_SHARING_ACCOUNT should be another IAM user name")
	}
}

//lintignore:AT003
func TestAccPreCheckDms(t *testing.T) {
	if HW_DMS_ENVIRONMENT == "" {
		t.Skip("This environment does not support DMS tests")
	}
}

//lintignore:AT003
func TestAccPreCheckDliJarPath(t *testing.T) {
	if HW_DLI_FLINK_JAR_OBS_PATH == "" {
		t.Skip("HW_DLI_FLINK_JAR_OBS_PATH must be set for DLI Flink Jar job acceptance tests.")
	}
}
