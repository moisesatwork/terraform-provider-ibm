// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package dnsservices

import (
	"fmt"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	pdnsSecondaryZoneID           = "secondary_zone_id"
	pdnsSecondaryZoneZone         = "zone"
	pdnsSecondaryZoneTransferFrom = "transfer_from"
	pdnsSecondaryZoneCreatedOn    = "created_on"
	pdnsSecondaryZoneModifiedOn   = "modified_on"
	pdnsSecondaryZoneEnabled      = "enabled"
	pdnsSecondaryZoneDescription  = "description"
)

func ResourceIBMPrivateDNSSecondaryZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMPrivateDNSSecondaryZoneCreate,
		Read:   resourceIBMPrivateDNSSecondaryZoneRead,
		Update: resourceIBMPrivateDNSSecondaryZoneUpdate,
		Delete: resourceIBMPrivateDNSSecondaryZoneDelete,
		Exists: resourceIBMPrivateDNSSecondaryZoneExists,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			pdnsSecondaryZoneID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Secondary Zone ID",
			},

			pdnsSecondaryZoneCreatedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Secondary Zone Creation date",
			},

			pdnsSecondaryZoneModifiedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Secondary Zone Modification date",
			},
		},
	}
}

func resourceIBMPrivateDNSSecondaryZoneCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	instanceID := d.Get(pdnsInstanceID).(string)
	zoneID := d.Get(pdnsZoneID).(string)
	resolverID := d.Get(pdnsResolverID).(string)
	zone := d.Get(pdnsSecondaryZoneZone).(string)
	transferFrom := d.Get(pdnsSecondaryZoneTransferFrom).(string)

	createSecondaryZoneOptions := sess.NewCreateSecondaryZoneOptions(instanceID, resolverID)

	createSecondaryZoneOptions.SetZone(zone)
	createSecondaryZoneOptions.SetTransferFrom(
		[]string{transferFrom},
	)

	mk := "private_dns_secondary_zone_" + instanceID + zoneID
	conns.IbmMutexKV.Lock(mk)
	defer conns.IbmMutexKV.Unlock(mk)

	response, detail, err := sess.CreateSecondaryZone(createSecondaryZoneOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating pdns secondary zone:%s\n%s", err, detail)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", instanceID, zoneID, *response.ID))
	return resourceIBMPrivateDNSSecondaryZoneRead(d, meta)
}

func resourceIBMPrivateDNSSecondaryZoneRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	idSet := strings.Split(d.Id(), "/")
	if len(idSet) < 3 {
		return fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of InstanceID/zoneID/permittedNetworkID", d.Id())
	}
	instanceID := idSet[0]
	resolverID := idSet[1]
	secondaryZoneID := idSet[2]
	getSecondaryZoneOptions := sess.NewGetSecondaryZoneOptions(instanceID, resolverID, secondaryZoneID)
	response, detail, err := sess.GetSecondaryZone(getSecondaryZoneOptions)

	if err != nil {
		return fmt.Errorf("[ERROR] Error reading pdns permitted network:%s\n%s", err, detail)
	}

	d.Set(pdnsInstanceID, idSet[0])
	d.Set(pdnsZoneID, idSet[1])
	d.Set(pdnsSecondaryZoneID, response.ID)
	d.Set(pdnsSecondaryZoneCreatedOn, response.CreatedOn)
	d.Set(pdnsSecondaryZoneModifiedOn, response.ModifiedOn)
	d.Set(pdnsSecondaryZoneEnabled, response.Enabled)

	return nil
}

func resourceIBMPrivateDNSSecondaryZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	idSet := strings.Split(d.Id(), "/")
	if len(idSet) < 3 {
		return fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of InstanceID/zoneID/permittedNetworkID", d.Id())
	}
	instanceID := idSet[0]
	resolverID := idSet[1]
	secondaryZoneID := idSet[2]

	// Check DNS zone is present
	getZoneOptions := sess.NewGetSecondaryZoneOptions(instanceID, resolverID, secondaryZoneID)
	_, response, err := sess.GetSecondaryZone(getZoneOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error fetching secondary zone:%s\n%s", err, response)
	}

	// Update DNS zone if attributes has any change
	if d.HasChange(pdnsSecondaryZoneTransferFrom) ||
		d.HasChange(pdnsSecondaryZoneDescription) ||
		d.HasChange(pdnsSecondaryZoneEnabled) {
		updateSecondaryZoneOptions := sess.NewUpdateSecondaryZoneOptions(instanceID, resolverID, secondaryZoneID)
		transferFrom := d.Get(pdnsSecondaryZoneTransferFrom).(string)
		description := d.Get(pdnsSecondaryZoneDescription).(string)
		enabled := d.Get(pdnsSecondaryZoneEnabled).(bool)
		updateSecondaryZoneOptions.SetTransferFrom(
			[]string{
				transferFrom,
			},
		)
		updateSecondaryZoneOptions.SetDescription(description)
		updateSecondaryZoneOptions.SetEnabled(enabled)

		_, detail, err := sess.UpdateSecondaryZone(updateSecondaryZoneOptions)

		if err != nil {
			return fmt.Errorf("[ERROR] Error updating pdns zone:%s\n%s", err, detail)
		}
	}

	return resourceIBMPrivateDNSSecondaryZoneRead(d, meta)
}

func resourceIBMPrivateDNSSecondaryZoneDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIBMPrivateDNSSecondaryZoneExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return false, nil
}
