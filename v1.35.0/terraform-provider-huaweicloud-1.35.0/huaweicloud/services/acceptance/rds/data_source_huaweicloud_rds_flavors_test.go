package rds

import (
	"regexp"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRdsFlavorDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_rds_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsFlavorDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "flavors.#", regexp.MustCompile("\\d+")),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.memory"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.instance_mode", "ha"),
				),
			},
		},
	})
}

var testAccRdsFlavorDataSource_basic = `
data "huaweicloud_rds_flavors" "test" {
  db_type       = "PostgreSQL"
  db_version    = "12"
  instance_mode = "ha"
}
`

func TestAccRdsFlavorDataSource_all(t *testing.T) {
	dataSourceName := "data.huaweicloud_rds_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsFlavorDataSource_all,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "flavors.#", regexp.MustCompile("\\d+")),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.name"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.vcpus", "16"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.memory", "32"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.instance_mode", "replica"),
				),
			},
		},
	})
}

var testAccRdsFlavorDataSource_all = `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_rds_flavors" "test" {
  db_type           = "MySQL"
  db_version        = "8.0"
  instance_mode     = "replica"
  vcpus             = 16
  memory            = 32
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}
`
