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
 * author: Matteo Valentini <matteo.valentini@nethesis.it>
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

func ReportsCurrentSessions(c *gin.Context) {
	var active_sessions models.ActiveSessions

	accountId := c.MustGet("token").(models.AccessToken).AccountId
	hotspotId := c.Query("hotspot")

	db := database.Instance()

	hotspotIdInt, err := strconv.Atoi(hotspotId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "hostspot id must be a number!"})
	} else {

		if hotspotIdInt == 0 || len(utils.ExtractHotspotIds(accountId, accountId == 1, hotspotIdInt)) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No hostspot found!"})
		} else {
			db.Model(&models.Session{}).Where("hostspot_id = ?", hotspotIdInt).Count(&active_sessions.Count)
			c.JSON(http.StatusOK, active_sessions)

		}
	}

	return
}
