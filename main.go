package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	MinimumHeaderSize = 12
	MaximumeaderSize  = 14
)

func (r *DataRecord) GetMessageLength() (int64, error) {
	switch r.Header.LocalMessageType {
	case 0:
		return 1, nil
	}
	return 0, nil
}

func (t *DataRecordMessageType) MarshalJSON() ([]byte, error) {
	if t == nil {
		return nil, errors.New("data record message type must be defined")
	}

	switch *t {
	case DataRecordMessageType_Definition:
		return json.Marshal("DEFINITION")
	case DataRecordMessageType_Data:
		return json.Marshal("DATA")
	default:
		return nil, fmt.Errorf("unknown data record message type %d", *t)
	}
}

func (t *DataRecordHeaderType) MarshalJSON() ([]byte, error) {
	if t == nil {
		return nil, errors.New("data record header type must be defined")
	}

	switch *t {
	case DataRecordHeaderType_Normal:
		return json.Marshal("NORMAL")
	case DataRecordHeaderType_CompressedTimestamp:
		return json.Marshal("COMPRESSED_TIMESTAMP")
	default:
		return nil, fmt.Errorf("unknown data record header type %d", *t)
	}
}

func (h *DataRecordHeader) Unmarshal(data []byte) error {
	if data == nil || len(data) != 1 {
		return errors.New("a data record header must be exactly 1 byte")
	}
	exploded := ExplodeByte(uint8(data[0]))
	if len(exploded) != 8 {
		return errors.New("bit represention of a byte must be of length 8")
	}

	isNormalHeader, err := strconv.ParseInt(string(exploded[0:1]), 2, 8)
	if err != nil {
		return err
	}
	if isNormalHeader == 0 {
		h.Type = DataRecordHeaderType_Normal
		localMessageType, err := strconv.ParseInt(string(exploded[4:8]), 2, 8)
		if err != nil {
			return err
		}
		h.LocalMessageType = uint8(localMessageType)
		messageType, err := strconv.ParseInt(string(exploded[1:2]), 2, 8)
		if err != nil {
			return err
		}
		if messageType == 1 {
			h.MessageType = DataRecordMessageType_Definition
		} else if messageType == 0 {
			h.MessageType = DataRecordMessageType_Data
		}
		developerData, err := strconv.ParseInt(string(exploded[2:3]), 2, 8)
		if err != nil {
			return err
		}
		if developerData == 1 {
			h.DeveloperData = true
		}
	} else if isNormalHeader == 1 {
		h.Type = DataRecordHeaderType_CompressedTimestamp
		h.MessageType = DataRecordMessageType_Data
		timeOffset, err := strconv.ParseInt(string(exploded[3:8]), 2, 8)
		if err != nil {
			return err
		}
		h.TimeOffset = uint8(timeOffset)
		localMessageType, err := strconv.ParseInt(string(exploded[1:3]), 2, 8)
		if err != nil {
			return err
		}
		h.LocalMessageType = uint8(localMessageType)
	}

	return nil
}

func (m *DefinitionMessage) DataMessageSize() int64 {
	size := int64(0)

	for _, feildDefinition := range m.Fields {
		size += int64(feildDefinition.Size)
	}

	return size
}

func (m *DataMessage) Unmarshal(def *DefinitionMessage, data []byte) error {
	if m == nil || def == nil || int64(len(data)) != def.DataMessageSize() {
		return ErrorMalformedBuffer
	}

	m.NormalFields = []uint64{}
	m.DeveloperFields = [][]byte{}

	offset := 0
	for i := 0; i < int(def.NumFields); i++ {

		// Is this is a developer field?
		if def.Fields[i].BaseType == nil {
			m.DeveloperFields = append(m.DeveloperFields, data[offset:offset+int(def.Fields[i].Size)])
			offset += int(def.Fields[i].Size)
			continue
		}

		switch def.Fields[i].BaseType.Number {
		case 0, 1, 2, 7, 10, 13:
			m.NormalFields = append(m.NormalFields, uint64(data[offset]))
		case 3, 4, 11:
			m.NormalFields = append(m.NormalFields, uint64(binary.BigEndian.Uint16(data[offset:offset+int(def.Fields[i].Size)])))
		case 5, 6, 8, 12:
			m.NormalFields = append(m.NormalFields, uint64(binary.BigEndian.Uint32(data[offset:offset+int(def.Fields[i].Size)])))
		case 9, 14, 15, 16:
			m.NormalFields = append(m.NormalFields, binary.BigEndian.Uint64(data[offset:offset+int(def.Fields[i].Size)]))
		default:
			return errors.New("unknown base type number")
		}
		offset += int(def.Fields[i].Size)
	}

	return nil
}

func (h *FileHeader) Unmarshal(data []byte) error {
	if data == nil || len(data) < 1 {
		return errors.New("data must be defined")
	}

	h.Size = data[0]
	if h.Size != MinimumHeaderSize && h.Size != MaximumeaderSize {
		return fmt.Errorf("valid header sizes are %d and %d", MinimumHeaderSize, MaximumeaderSize)
	}
	if uint8(len(data)) != h.Size {
		return fmt.Errorf("data must be at least of size %d", h.Size)
	}

	h.ProtocolVersion = data[1]
	h.ProfileVersion = binary.BigEndian.Uint16(data[2:4])
	h.DataSize = binary.BigEndian.Uint32(data[4:8])
	h.DataType = string(data[8:12])
	if h.Size == MaximumeaderSize {
		h.CRC = binary.LittleEndian.Uint16(data[12:14])
	}
	return nil
}

func ExplodeByte(data uint8) string {
	var b strings.Builder

	for i := 7; i >= 0; i-- {
		if (data & (1 << uint(i))) != 0 {
			b.WriteString("1")
		} else {
			b.WriteString("0")
		}
	}

	return b.String()
}

func (dm *DataMessage) ReadAndUnmarshal(ctx context.Context, r io.Reader) (int, error) {
	if dm == nil {
		return 0, errors.New("definition message is nil")
	}

	def := ctx.Value(ContextKeyCurrentDefinitionRecord).(DataRecord)
	if def.DefinitionMessage == nil {
		return 0, errors.New("definition message must preceed a data message")
	}

	var totalBytesRead int

	buf := make([]byte, def.DefinitionMessage.DataMessageSize())
	n, err := r.Read(buf)
	totalBytesRead += n
	if err != nil {
		return totalBytesRead, err
	}

	if err := dm.Unmarshal(def.DefinitionMessage, buf); err != nil {
		return totalBytesRead, err
	}

	return totalBytesRead, nil
}

func (f *FieldDefinitions) Unmarshal(ctx context.Context, data []byte) error {
	// Each field is exactly 3 bytes.
	numFields := len(data) / 3
	newFields := make([]FieldDefinition, numFields)
	isDeveloperField := false

	if h, ok := ctx.Value(ContextKeyDataRecordFieldType).(DataRecordFieldType); ok && h == DataRecordFieldType_Developer {
		isDeveloperField = true
	}

	for i := 0; i < (numFields * 3); i += 3 {
		newFields[i/3] = FieldDefinition{
			Number: uint8(data[i]),
			Size:   uint8(data[i+1]),
		}

		if isDeveloperField {
			newFields[i/3].DeveloperDataIndex = uint8(data[i+2])
			continue
		}

		baseTypeExploded := ExplodeByte(uint8(data[i+2]))

		baseTypeNumber, err := strconv.ParseInt(string(baseTypeExploded[3:8]), 2, 8)
		if err != nil {
			return err
		}

		endianAbility, err := strconv.ParseInt(string(baseTypeExploded[0:1]), 2, 8)
		if err != nil {
			return err
		}

		newFields[i/3].BaseType = &BaseType{
			Number:        uint8(baseTypeNumber),
			EndianAbility: uint8(endianAbility),
		}
	}

	*f = append(*f, newFields...)

	return nil
}

func (dm *DefinitionMessage) ReadAndUnmarshal(ctx context.Context, r io.Reader) (int, error) {
	var totalBytesRead int

	if dm == nil {
		return 0, ErrorTypeNotDefined
	}

	fixedContentBuffer := make([]byte, 5)
	fixedContentBytesRead, err := r.Read(fixedContentBuffer)
	totalBytesRead += fixedContentBytesRead
	if err != nil {
		return totalBytesRead, err
	}

	// This does not yet include developer field definitions.
	dm.NumFields = fixedContentBuffer[4]
	dm.Architecture = fixedContentBuffer[1]

	if t, ok := GlobalMessageNumber_Types[binary.BigEndian.Uint16(fixedContentBuffer[2:4])]; !ok {
		dm.GlobalMessageType = GlobalMessageType_Unknown
	} else {
		dm.GlobalMessageType = t
	}

	normalFieldDefinitionsBuffer := make([]byte, int64(dm.NumFields)*3)
	normalFieldDefinitionsBytesRead, err := r.Read(normalFieldDefinitionsBuffer)
	totalBytesRead += normalFieldDefinitionsBytesRead
	if err != nil {
		return totalBytesRead, err
	}

	if err := dm.Fields.Unmarshal(context.WithValue(ctx, ContextKeyDataRecordFieldType, DataRecordFieldType_Normal), normalFieldDefinitionsBuffer); err != nil {
		return totalBytesRead, err
	}

	if h, ok := ctx.Value(ContextKeyDataRecordHeader).(*DataRecordHeader); ok && h.DeveloperData {
		numDeveloperFieldsBuffer := make([]byte, 1)
		numDeveloperFieldsBytesRead, err := r.Read(numDeveloperFieldsBuffer)
		totalBytesRead += numDeveloperFieldsBytesRead
		if err != nil {
			return totalBytesRead, err
		}
		numDeveloperFields := numDeveloperFieldsBuffer[0]
		dm.NumFields += uint8(numDeveloperFields)

		developerFieldDefinitionsBuffer := make([]byte, int64(numDeveloperFields)*3)
		developerFieldDefinitionsBytesRead, err := r.Read(developerFieldDefinitionsBuffer)
		totalBytesRead += developerFieldDefinitionsBytesRead
		if err != nil {
			return totalBytesRead, err
		}

		if err := dm.Fields.Unmarshal(context.WithValue(ctx, ContextKeyDataRecordFieldType, DataRecordFieldType_Developer), developerFieldDefinitionsBuffer); err != nil {
			return totalBytesRead, err
		}
	}

	return totalBytesRead, nil
}

func (dr *DataRecord) ReadAndUnmarshal(ctx context.Context, r io.Reader) (int, error) {
	if dr == nil {
		return 0, ErrorTypeNotDefined
	}

	var totalBytesRead int

	dr.Header = new(DataRecordHeader)

	// no matter what the message type, the header size is the same
	buf := make([]byte, 1)
	n, err := r.Read(buf)
	totalBytesRead += n
	if err != nil {
		return totalBytesRead, err
	}

	if err := dr.Header.Unmarshal(buf); err != nil {
		return totalBytesRead, err
	}
	switch dr.Header.MessageType {
	case DataRecordMessageType_Definition:
		dr.DefinitionMessage = new(DefinitionMessage)
		n, err := dr.DefinitionMessage.ReadAndUnmarshal(context.WithValue(ctx, ContextKeyDataRecordHeader, (*dr).Header), r)
		totalBytesRead += n
		if err != nil {
			return totalBytesRead, err
		}
	case DataRecordMessageType_Data:
		dr.DataMessage = new(DataMessage)
		n, err := dr.DataMessage.ReadAndUnmarshal(ctx, r)
		totalBytesRead += n
		if err != nil {
			return totalBytesRead, err
		}
	default:
		return totalBytesRead, errors.New("unkown message type")
	}

	return totalBytesRead, nil
}

func (h *FileHeader) ReadAndUnmarshal(r io.Reader) (int, error) {
	buf := make([]byte, 14)
	n, err := r.Read(buf)
	if err != nil {
		return n, err
	}
	return n, h.Unmarshal(buf)
}

func (f *File) ReadAndUnmarshal(r io.Reader) (int, error) {
	if f == nil {
		return 0, ErrorTypeNotDefined
	}

	if f.Header == nil {
		f.Header = new(FileHeader)
	}

	if f.Records == nil {
		f.Records = []DataRecord{}
	}

	var totalBytesRead int

	n, err := f.Header.ReadAndUnmarshal(r)
	totalBytesRead += n
	if err != nil {
		return totalBytesRead, err
	}

	recordsProcessed := 0

	var currentDefinitionMessage DataRecord

	dataRecordsBytesLeftToProcess := int(f.Header.DataSize)

	for dataRecordsBytesLeftToProcess > 0 {
		dr := new(DataRecord)

		ctx := context.Background()
		if recordsProcessed > 0 {
			ctx = context.WithValue(ctx, ContextKeyCurrentDefinitionRecord, currentDefinitionMessage)
		}

		n, err := dr.ReadAndUnmarshal(ctx, r)
		totalBytesRead += n
		dataRecordsBytesLeftToProcess -= n
		if err != nil {
			return totalBytesRead, err
		}

		if dr.Header.MessageType == DataRecordMessageType_Definition {
			currentDefinitionMessage = *dr
		}

		f.Records = append(f.Records, *dr)
		recordsProcessed++
	}
	return totalBytesRead, nil
}

func main() {
	rawFile, err := os.Open(os.Getenv("FILE_LOCATION"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer rawFile.Close()

	fitFile := new(File)
	if _, err := fitFile.ReadAndUnmarshal(rawFile); err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	data, err := json.Marshal(fitFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println(string(data))
}
