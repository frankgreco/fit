package main

import "errors"

type File struct {
	Header  *FileHeader  `json:"header"`
	Records []DataRecord `json:"records"`
}

type FileHeader struct {
	Size            uint8  `json:"size"`
	ProtocolVersion uint8  `json:"protocol_version"`
	ProfileVersion  uint16 `json:"profile_version"`
	DataSize        uint32 `json:"data_size"`
	DataType        string `json:"data_type"`
	CRC             uint16 `json:"crc"`
}

type DataRecord struct {
	Header            *DataRecordHeader  `json:"header"`
	DefinitionMessage *DefinitionMessage `json:"definition_message,omitempty"`
	DataMessage       *DataMessage       `json:"data_message,omitempty"`
}

type DataRecordHeader struct {
	Type             DataRecordHeaderType  `json:"type"`
	LocalMessageType uint8                 `json:"local_message_type"`
	MessageType      DataRecordMessageType `json:"message_type,omitempty"`
	DeveloperData    bool                  `json:"developer_data"`
	TimeOffset       uint8                 `json:"time_offset"`
}

type FieldDefinitions []FieldDefinition

type DefinitionMessage struct {
	Architecture      uint8             `json:"architecture"`
	GlobalMessageType GlobalMessageType `json:"global_message_type"`
	NumFields         uint8             `json:"num_fields"`
	Fields            FieldDefinitions  `json:"fields"`
}

type FieldDefinition struct {
	Type               DataRecordFieldType `json:"type"`
	Number             uint8               `json:"number"`
	Size               uint8               `json:"size"`
	BaseType           *BaseType           `json:"base_type,omitempty"`
	DeveloperDataIndex uint8               `json:"developer_data_index"`
}

type BaseType struct {
	Number        uint8 `json:"number"`
	EndianAbility uint8 `json:"endian_ability"`
}

type DataMessage struct {
	NormalFields    []uint64 `json:"normal_fields"`
	DeveloperFields [][]byte `json:"developer_fields"`
}

type DataRecordHeaderType int

type DataRecordMessageType int

type DataRecordFieldType int

type GlobalMessageType int

const (
	DataRecordHeaderType_Normal DataRecordHeaderType = iota
	DataRecordHeaderType_CompressedTimestamp

	DataRecordMessageType_Definition DataRecordMessageType = iota
	DataRecordMessageType_Data

	DataRecordFieldType_Developer DataRecordFieldType = iota
	DataRecordFieldType_Normal

	GlobalMessageType_FileID GlobalMessageType = iota
	GlobalMessageType_Capabilities
	GlobalMessageType_Unknown
)

var (
	GlobalMessageNumber_Types = map[uint16]GlobalMessageType{
		0: GlobalMessageType_FileID,
		1: GlobalMessageType_Capabilities,
	}

	GlobalMessageType_Names = map[GlobalMessageType]string{
		GlobalMessageType_FileID:       "FILE_ID",
		GlobalMessageType_Capabilities: "CAPABILITIES",
	}

	ErrorTypeNotDefined  = errors.New("type not defined")
	ErrorMalformedBuffer = errors.New("malformed buffer")

	ContextKeyDataRecordHeader        = "DATA_RECORD_HEADER"
	ContextKeyDataRecordFieldType     = "DATA_RECORD_FIELD_TYPE"
	ContextKeyCurrentDefinitionRecord = "CURRENT_DEFINITION_RECORD"
)
