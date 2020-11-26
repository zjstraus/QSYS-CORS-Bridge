/*
Copyright (C) 2020  Zach Strauss

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, version 3.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// Design is the top-level structure returned by a Core
type Design struct {
	Snapshots        []Snapshot        `xml:"Snapshots>Snapshot",json:"snapshots"`
	ExternalControls []ExternalControl `xml:"ExternalControls>Control",json:"externalControls"`
	DesignName       string            `xml:"DesignName,attr",json:"designName"`
	CompileGUID      string            `xml:"CompileGUID,attr",json:"compileGUID"`
}

// Snapshot holds some state for a snapshot, untested
type Snapshot struct {
	Name     string `xml:"Name,attr",json:"name"`
	Count    string `xml:"Count,attr",json:"count"`
	CodeName string `xml:"CodeName,attr",json:"codeName"`
}

// ExternalControl structs encode the internal name, external name, and metadata for a Control
type ExternalControl struct {
	Id             string `xml:"Id,attr",json:"id"`
	ControlId      string `xml:"ControlId,attr",json:"controlId"`
	ControlName    string `xml:"ControlName,attr",json:"controlName"`
	ComponentId    string `xml:"ComponentId,attr",json:"componentId"`
	ComponentName  string `xml:"ComponentName,attr",json:"componentName"`
	ComponentLabel string `xml:"ComponentLabel,attr",json:"componentLabel"`
	MappingName    string
	Type           string `xml:"Type,attr",json:"type"`
	Mode           string `xml:"Mode,attr",json:"mode"`
	MinimumValue   string `xml:"MinimumValue,attr",json:"minimumValue"`
	MaximumValue   string `xml:"MaximumValue,attr",json:"maximumValue"`
}

// GetDesign requests the ExternalControls.xml file from the core and pre-processes the control names
func GetDesign(address string) (*Design, error) {
	resp, httperr := http.Get("https://" + address + "/designs/current_design/ExternalControls.xml")
	if httperr != nil {
		return nil, httperr
	}
	defer resp.Body.Close()

	data, readerr := ioutil.ReadAll(resp.Body)
	if readerr != nil {
		return nil, readerr
	}

	var parsed Design

	unmarshalerr := xml.Unmarshal(data, &parsed)
	if unmarshalerr != nil {
		return nil, unmarshalerr
	}

	// You have to compute these to attach Named Control listeners to internal controls anyway, so might as well do
	// it right away
	for i, control := range parsed.ExternalControls {
		parsed.ExternalControls[i].MappingName = control.ComponentId + "_" + control.ControlId
	}

	return &parsed, nil
}
