/*
 *
 * Module:    BIG Modelling Bus
 * Package:   Languages/Conceptual Domain Modelling, Version 1
 *
 * This package implements the Conceptual Domain Modelling language, version 1, for the BIG Modelling Bus.
 *
 * Creator: Henderik A. Proper (e.proper@acm.org), TU Wien, Austria
 *
 * Version of: XX.11.2025
 *
 */

package cdm_v1

import (
	"encoding/json"

	"github.com/erikproper/big-modelling-bus.go.v1/connect"
	"github.com/erikproper/big-modelling-bus.go.v1/generics"
)

/*
 * Defining key constants
 */

const (
	ModelJSONVersion = "cdm-1.0-1.0"
)

/*
 * Defining the CDM model structure, including the JSON structure
 */

type (
	TRelationReading struct {
		InvolvementTypes []string `json:"involvement types"` // The involvement types used in the relation type readings
		ReadingElements  []string `json:"reading elements"`  // The strings used in relation type reading
	}

	TCDMModel struct {
		ModelName                  string                                 `json:"model name"` // The name of the model
		ModellingBusArtefactPoster connect.TModellingBusArtefactConnector `json:"-"`          // The Modelling Bus Artefact Poster used to post the model
		//		TypeIDCount                int                                    `json:"-"`          // The counter for type IDs
		InstanceIDCount int `json:"-"` // The counter for instance IDs

		// For types
		TypeName map[string]string `json:"type names"` // The names of the types, by their IDs

		// For concrete individual types
		ConcreteIndividualTypes map[string]bool `json:"concrete individual types"` // The concrete individual types

		// For quality types
		QualityTypes        map[string]bool   `json:"quality types"`            // The quality types
		DomainOfQualityType map[string]string `json:"domains of quality types"` // The domain of each quality type

		// For involvement types
		InvolvementTypes              map[string]bool   `json:"involvement types"`                   // The involvement types
		BaseTypeOfInvolvementType     map[string]string `json:"base types of involvement types"`     // The base type of each involvement type
		RelationTypeOfInvolvementType map[string]string `json:"relation types of involvement types"` // The relation type of each involvement type

		// For relation types
		RelationTypes                     map[string]bool             `json:"relation types"`                         // The relation types
		InvolvementTypesOfRelationType    map[string]map[string]bool  `json:"involvement types of relation types"`    // The involvement types of each relation type
		AlternativeReadingsOfRelationType map[string]map[string]bool  `json:"alternative readings of relation types"` // The alternative readings of each relation type
		PrimaryReadingOfRelationType      map[string]string           `json:"primary readings of relation types"`     // The primary reading of each relation type
		ReadingDefinition                 map[string]TRelationReading `json:"reading definition"`                     // The definition of each relation type reading
	}
)

/*
 *
 * Functionality related to the CDM model
 *
 */

// Cleaning the model
func (m *TCDMModel) Clean() {
	// Resetting all fields
	m.ModelName = ""
	m.ConcreteIndividualTypes = map[string]bool{}
	m.QualityTypes = map[string]bool{}
	m.RelationTypes = map[string]bool{}
	m.InvolvementTypes = map[string]bool{}
	m.TypeName = map[string]string{}
	m.DomainOfQualityType = map[string]string{}
	m.BaseTypeOfInvolvementType = map[string]string{}
	m.RelationTypeOfInvolvementType = map[string]string{}
	m.InvolvementTypesOfRelationType = map[string]map[string]bool{}
	m.AlternativeReadingsOfRelationType = map[string]map[string]bool{}
	m.PrimaryReadingOfRelationType = map[string]string{}
	m.ReadingDefinition = map[string]TRelationReading{}
}

// Generating a new element ID
func (m *TCDMModel) NewElementID() string {
	return generics.GetTimestamp()
}

func (m *TCDMModel) SetModelName(name string) {
	m.ModelName = name
}

func (m *TCDMModel) AddConcreteIndividualType(name string) string {
	id := m.NewElementID()
	m.ConcreteIndividualTypes[id] = true
	m.TypeName[id] = name

	return id
}

func (m *TCDMModel) AddQualityType(name, domain string) string {
	id := m.NewElementID()
	m.QualityTypes[id] = true
	m.TypeName[id] = name
	m.DomainOfQualityType[id] = domain

	return id
}

func (m *TCDMModel) AddInvolvementType(name string, base string) string {
	id := m.NewElementID()
	m.InvolvementTypes[id] = true
	m.TypeName[id] = name
	m.BaseTypeOfInvolvementType[id] = base

	return id
}

func (m *TCDMModel) AddRelationType(name string, involvementTypes ...string) string {
	id := m.NewElementID()
	m.RelationTypes[id] = true
	m.TypeName[id] = name

	m.InvolvementTypesOfRelationType[id] = map[string]bool{}
	for _, involvementType := range involvementTypes {
		m.RelationTypeOfInvolvementType[involvementType] = id
		m.InvolvementTypesOfRelationType[id][involvementType] = true
	}

	m.AlternativeReadingsOfRelationType[id] = map[string]bool{}

	return id
}

func (m *TCDMModel) AddRelationTypeReading(relationType string, stringsAndInvolvementTypes ...string) string {
	reading := TRelationReading{}

	isReadingString := true
	for _, element := range stringsAndInvolvementTypes {
		if isReadingString {
			reading.ReadingElements = append(reading.ReadingElements, element)
		} else {
			reading.InvolvementTypes = append(reading.InvolvementTypes, element)
		}
		isReadingString = !isReadingString
	}

	readingID := m.NewElementID()
	m.AlternativeReadingsOfRelationType[relationType][readingID] = true
	m.ReadingDefinition[readingID] = reading

	if m.PrimaryReadingOfRelationType[relationType] == "" {
		m.PrimaryReadingOfRelationType[relationType] = readingID
	}

	return readingID
	// Does require a check to see if all InvolvementTypesss of the relation have been used ... and used only once
	// But ... as this is only "Hello World" for now, so we won't do so yet.
}

/*
 *
 * Initialisation and creation
 *
 */

func CreateCDMModel() TCDMModel {
	CDMModel := TCDMModel{}
	CDMModel.Clean()

	return CDMModel
}

/*
 *
 * Posting models to the artefactBus
 *
 */

func CreateCDMPoster(ModellingBusConnector connect.TModellingBusConnector, modelID string) TCDMModel {
	CDMPosterModel := CreateCDMModel()

	// Note: One ModellingBusConnector can be used for different artefacts with different json versions.
	CDMPosterModel.ModellingBusArtefactPoster = connect.CreateModellingBusArtefactConnector(ModellingBusConnector, ModelJSONVersion)
	CDMPosterModel.ModellingBusArtefactPoster.PrepareForPosting(modelID)

	return CDMPosterModel
}

func (m *TCDMModel) PostState() {
	m.ModellingBusArtefactPoster.PostJSONArtefactState(json.Marshal(m))
}

func (m *TCDMModel) PostUpdate() {
	m.ModellingBusArtefactPoster.PostJSONArtefactUpdate(json.Marshal(m))
}

func (m *TCDMModel) PostConsidering() {
	m.ModellingBusArtefactPoster.PostJSONArtefactConsidering(json.Marshal(m))
}

/*
 *
 * Reading models from the artefactBus
 *
 */

// Note: One ModellingBusConnector can be used for different models of different kinds.
func CreateCDMListener(ModellingBusConnector connect.TModellingBusConnector) connect.TModellingBusArtefactConnector {
	ModellingBusCDMModelListener := connect.CreateModellingBusArtefactConnector(ModellingBusConnector, ModelJSONVersion)

	return ModellingBusCDMModelListener
}

func (m *TCDMModel) GetStateFromBus(artefactBus connect.TModellingBusArtefactConnector) bool {
	m.Clean()
	err := json.Unmarshal(artefactBus.CurrentContent, m)

	return err == nil
}

func (m *TCDMModel) GetUpdatedFromBus(artefactBus connect.TModellingBusArtefactConnector) bool {
	m.Clean()
	err := json.Unmarshal(artefactBus.UpdatedContent, m)

	return err == nil
}

func (m *TCDMModel) GetConsideredFromBus(artefactBus connect.TModellingBusArtefactConnector) bool {
	m.Clean()
	err := json.Unmarshal(artefactBus.ConsideredContent, m)

	return err == nil
}
