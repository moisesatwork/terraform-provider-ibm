package dnsservices

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMPrivateDNSSecondaryZones() *schema.Resource {
	return &schema.Resource{}
}

func dataSourceIBMPrivateDNSSecondaryZonesRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	DnszoneID := d.Get(pdnsZoneID).(string)
	listDNSSecondaryZoneOptions := sess.NewListSecondaryZonesOptions(instanceID, DnszoneID)
	// availableSecondaryZones, detail, err := sess.ListSecondaryZones(listDNSSecondaryZoneOptions)
	_, detail, err := sess.ListSecondaryZones(listDNSSecondaryZoneOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error reading list of pdns resource records:%s\n%s", err, detail)
	}
	// secondaryZones := make([]map[string]interface{}, 0)
	// for _, instance := range secondaryZones.SecondaryZones {
	// }
	return nil
}

// dataSourceIBMPrivateDNSSecondaryZonesID returns a reasonable ID for dns zones list.
func dataSourceIBMPrivateDNSSecondaryZonesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
