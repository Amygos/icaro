/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Icaro project.
 *
 * Icaro is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Icaro is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Icaro.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package methods

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/icaro/sun/sun-api/database"
	"github.com/nethesis/icaro/sun/sun-api/models"
	"github.com/nethesis/icaro/sun/sun-api/utils"
)

func GetIntegrations(c *gin.Context) {
	var integrations []models.Integration

	db := database.Instance()
	db.Find(&integrations)

	if len(integrations) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No integrations found!"})
		return
	}

	c.JSON(http.StatusOK, integrations)
}

func GetHotspotIntegrations(c *gin.Context) {
	var integrations []models.HotspotIntegration
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Param("hotspot_id")

	// convert hotspot id to int
	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Hotspot id error", "error": err.Error()})
		return
	}

	db := database.Instance()
	db.Where("hotspot_id in (?)", utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).Find(&integrations)

	if len(integrations) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No integrations found!"})
		return
	}

	c.JSON(http.StatusOK, integrations)
}

func UpdateHotspotIntegrations(c *gin.Context) {
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Param("hotspot_id")
	integrationId := c.Param("integration_id")

	// convert hotspot id to int
	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Hotspot id error", "error": err.Error()})
		return
	}

	integrationIdInt, err := strconv.Atoi(integrationId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Hotspot id error", "error": err.Error()})
		return
	}

	// check hotspot ownership
	if utils.Contains(utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt), hotspotIdInt) {

		integration := models.HotspotIntegration{
			HotspotId:     hotspotIdInt,
			IntegrationId: integrationIdInt,
		}

		db := database.Instance()
		db.Save(&integration)

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This hotspot is not yours"})
	}
}

func DeleteHotspotIntegrations(c *gin.Context) {
	var integration models.HotspotIntegration
	accountId := c.MustGet("token").(models.AccessToken).AccountId

	hotspotId := c.Param("hotspot_id")
	integrationId := c.Param("integration_id")

	// convert hotspot id to int
	hotspotIdInt, err := strconv.Atoi(hotspotId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Hotspot id error", "error": err.Error()})
		return
	}

	db := database.Instance()
	db.Where("integration_id = ? AND hotspot_id in (?)", integrationId, utils.ExtractHotspotIds(accountId, (accountId == 1), hotspotIdInt)).First(&integration)

	if integration.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No integration found!"})
		return
	}

	db.Delete(&integration)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
