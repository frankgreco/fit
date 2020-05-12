package main

type DataRecordHeaderType int
type DataRecordMessageType int

type File struct {
	Header      *FileHeader  `json:"header,omitempty"`
	DataRecords []DataRecord `json:"data_records"`
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
	DeveloperData    bool                  `json:"developer_data,omitempty"`
	TimeOffset       uint8                 `json:"time_offset"`
}

type DefinitionMessage struct {
	FixedContent    *DefinitionMessageFixedContent    `json:"fixed_content"`
	VariableContent *DefinitionMessageVariableContent `json:"variable_content"`
}

type DefinitionMessageFixedContent struct {
	Architecture        uint8  `json:"architecture"`
	GlobalMessageNumber uint16 `json:"global_message_number"`
	NumFields           uint8  `json:"num_fields"`
}

type DefinitionMessageVariableContent struct {
	Fields []FieldDefinition `json:"fields"`
}

type FieldDefinition struct {
	Number   uint8    `json:"number"`
	Size     uint8    `json:"size"`
	BaseType BaseType `json:"base_type"`
}

type BaseType struct {
	Number        uint8 `json:"number"`
	EndianAbility uint8 `json:"endian_ability"`
}

type DataMessage struct {
	Fields []uint64 `json:"fields"`
}

const (
	DataRecordHeaderType_Normal DataRecordHeaderType = iota
	DataRecordHeaderType_CompressedTimestamp

	DataRecordMessageType_Definition DataRecordMessageType = iota
	DataRecordMessageType_Data
)
