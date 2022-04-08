package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	vpc_model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getVpcAddressGroupResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcVpcV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Huaweicloud VPC client: %s", err)
	}

	request := &vpc_model.ShowAddressGroupRequest{
		AddressGroupId: state.Primary.ID,
	}

	return client.ShowAddressGroup(request)
}

func TestAccVpcAddressGroup_basic(t *testing.T) {
	var group vpc_model.ShowAddressGroupResponse

	rName := acceptance.RandomAccResourceName()
	rNameUpdate := rName + "_updated"
	resourceName := "huaweicloud_vpc_address_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getVpcAddressGroupResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testVpcAdressGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "addresses.#", "2"),
				),
			},
			{
				Config: testVpcAdressGroup_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by acc test"),
					resource.TestCheckResourceAttr(resourceName, "addresses.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testVpcAdressGroup_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_address_group" "test" {
  name        = "%s"
  description = "created by acc test"
  addresses   = [
    "192.168.3.2",
    "192.168.3.20-192.168.3.100"
  ]
}
`, rName)
}

func testVpcAdressGroup_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_address_group" "test" {
  name        = "%s"
  description = "updated by acc test"
  addresses   = [
    "192.168.5.0/24",
    "192.168.3.2",
    "192.168.3.20-192.168.3.100"
  ]
}
`, rName)
}
